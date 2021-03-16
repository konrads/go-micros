package portstore

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/konrads/go-micros/pkg/db"
	"github.com/konrads/go-micros/pkg/model"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func ToProtobuff(p *model.Port) *OptionalPortResp_Port {
	return &OptionalPortResp_Port{
		Id:          p.Id,
		Name:        p.Name,
		Coordinates: p.Coordinates[:],
		City:        p.City,
		Province:    p.Province,
		Country:     p.Country,
		Alias:       p.Alias,
		Regions:     p.Regions,
		Timezone:    p.Timezone,
		Unlocs:      p.Unlocs,
		Code:        p.Code,
	}
}

func (p *OptionalPortResp_Port) ToModel() *model.Port {
	return &model.Port{
		Id:          p.Id,
		Name:        p.Name,
		Coordinates: p.Coordinates,
		City:        p.City,
		Province:    p.Province,
		Country:     p.Country,
		Alias:       p.Alias,
		Regions:     p.Regions,
		Timezone:    p.Timezone,
		Unlocs:      p.Unlocs,
		Code:        p.Code,
	}
}

// Server
type PortStoreServerImpl struct {
	address    string
	db         *db.DB
	grpcServer *grpc.Server
}

// part of PortStoreServer interface
func (ps *PortStoreServerImpl) PersistPorts(stream PortStore_PersistPortsServer) error {
	for {
		portBuff, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Grpc EOF, for now...")
			break
		} else if err != nil {
			log.Printf("Failed to fetch a port due to %v", err)
			break
		} else {
			port := portBuff.ToModel()
			log.Printf("Persisting port... %v", port)
			_, dbErr := ps.db.SaveAll([]model.Port{*port})
			if dbErr != nil {
				return dbErr
			}
		}
	}
	return nil
}

// part of PortStoreServer interface
func (ps *PortStoreServerImpl) GetPort(ctx context.Context, portReq *PortReq) (*OptionalPortResp, error) {
	log.Printf("...getting: %v", portReq)
	port, err := ps.db.Get(portReq.PortId)
	if err != nil {
		return nil, err
	}
	protoPort := ToProtobuff(port)
	protoMaybePort := OptionalPortResp{Resp: protoPort}
	return &protoMaybePort, err
}

func RunPortServer(address string) error {
	s := &PortStoreServerImpl{db: db.New(), grpcServer: grpc.NewServer(), address: address}
	RegisterPortStoreServer(s.grpcServer, s)
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Printf("Port store listen error: %v", err)
		return err
	} else {
		return s.grpcServer.Serve(listener)
	}
}

// Client
type PortStoreClientImpl struct {
	conn *grpc.ClientConn
}

func NewPortClient(serverAddress string) (*PortStoreClientImpl, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())    // FIXME: revisit
	opts = append(opts, grpc.WithInsecure()) // FIXME: revisit
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		return nil, err
	} else {
		return &PortStoreClientImpl{conn}, nil
	}
}

func (ps *PortStoreClientImpl) GetPortPersistor() (func(port model.Port) error /*processor*/, func() /*cleanup*/, error) {
	client := NewPortStoreClient(ps.conn)
	// FIXME: understand context...
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	if stream, err := client.PersistPorts(ctx); err != nil {
		return nil, nil, err
	} else {
		processor := func(p model.Port) error {
			log.Printf("Sending via grpc... %v", p)
			return stream.Send(ToProtobuff(&p))
		}
		cleanup := func() {
			log.Printf("Grpc cleanup...")
			defer stream.CloseSend()
		}
		return processor, cleanup, nil
	}
}

func (ps *PortStoreClientImpl) GetPort(id string) (*model.Port, error) {
	client := NewPortStoreClient(ps.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	optProtoPort, err := client.GetPort(ctx, &PortReq{PortId: id})
	if err != nil {
		return nil, err
	} else {
		protoPort := optProtoPort.GetResp()
		if protoPort == nil {
			return nil, fmt.Errorf("Failed to find the port for id %v", id)
		} else {
			return protoPort.ToModel(), nil
		}
	}
}
