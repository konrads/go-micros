package starstore

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

func ToProtobuff(s *model.Star) *OptionalStarResp_Star {
	return &OptionalStarResp_Star{
		Id:                s.Id,
		Name:              s.Name,
		Alias:             s.Alias,
		Constellation:     s.Constellation,
		Coordinates:       s.Coordinates,
		Distance:          s.Distance,
		ApparentMagnitude: s.ApparentMagnitude,
	}
}

func (s *OptionalStarResp_Star) ToModel() *model.Star {
	return &model.Star{
		Id:                s.Id,
		Name:              s.Name,
		Alias:             s.Alias,
		Constellation:     s.Constellation,
		Coordinates:       s.Coordinates,
		Distance:          s.Distance,
		ApparentMagnitude: s.ApparentMagnitude,
	}
}

// Server
type StarStoreServerImpl struct {
	address    string
	db         *db.DB
	grpcServer *grpc.Server
}

// part of StarStoreServer interface
func (ss *StarStoreServerImpl) PersistStars(stream StarStore_PersistStarsServer) error {
	for {
		starBuff, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Grpc EOF, for now...")
			break
		} else if err != nil {
			log.Printf("Failed to fetch a star due to %v", err)
			break
		} else {
			star := starBuff.ToModel()
			log.Printf("Persisting star... %v", *star)
			_, err := (*ss.db).SaveAll([]model.Star{*star})
			if err != nil {
				log.Fatalf("Failed due to %v", err)
				return err
			}
		}
	}
	return nil
}

// part of StarStoreServer interface
func (ss *StarStoreServerImpl) GetStar(ctx context.Context, starReq *StarReq) (*OptionalStarResp, error) {
	log.Printf("...getting: %v", starReq)
	star, err := (*ss.db).Get(starReq.StarId)
	if err != nil {
		return nil, err
	} else if star == nil {
		return &OptionalStarResp{Resp: nil}, nil
	} else {
		protoStar := ToProtobuff(star)
		protoOptStar := OptionalStarResp{Resp: protoStar}
		return &protoOptStar, err
	}
}

func RunStarServer(address string, db *db.DB) error {
	s := &StarStoreServerImpl{db: db, grpcServer: grpc.NewServer(), address: address}
	RegisterStarStoreServer(s.grpcServer, s)
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Printf("Star store listen error: %v", err)
		return err
	} else {
		return s.grpcServer.Serve(listener)
	}
}

// Client
type StarStoreClientImpl struct {
	conn *grpc.ClientConn
}

func NewStarClient(grpcUri string) (*StarStoreClientImpl, error) {
	conn, err := grpc.Dial(grpcUri, grpc.WithBlock(), grpc.WithInsecure()) // awaits the connection, no transport security (eg. TLS/SSL)
	if err != nil {
		return nil, err
	} else {
		return &StarStoreClientImpl{conn}, nil
	}
}

type processor func(model.Star) error
type cleanup func()

func (ss *StarStoreClientImpl) GetStarPersistor() (processor, cleanup, error) {
	client := NewStarStoreClient(ss.conn)
	ctx, ctxCancel := context.WithTimeout(context.Background(), 100*time.Second)
	if stream, err := client.PersistStars(ctx); err != nil {
		return nil, nil, err
	} else {
		processor := func(p model.Star) error {
			log.Printf("gRPC send: %v", p)
			return stream.Send(ToProtobuff(&p))
		}
		cleanup := func() {
			log.Printf("gRPC cleanup...")
			stream.CloseSend()
			_ = ctxCancel // not running ctxCancel() due to status switching to Canceled prematurely, as per: https://github.com/grpc/grpc-go/issues/1099
		}
		return processor, cleanup, nil
	}
}

func (ss *StarStoreClientImpl) GetStar(id string) (*model.Star, error) {
	client := NewStarStoreClient(ss.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	optProtoStar, err := client.GetStar(ctx, &StarReq{StarId: id})
	if err != nil {
		return nil, err
	} else {
		protoStar := optProtoStar.GetResp()
		if protoStar == nil {
			return nil, fmt.Errorf("Failed to find the star for id %v", id)
		} else {
			return protoStar.ToModel(), nil
		}
	}
}
