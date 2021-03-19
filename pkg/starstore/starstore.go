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

const (
	GET_TIMEOUT  = 10 * time.Second
	POST_TIMEOUT = 100 * time.Second
)

func ToProtobuff(s *model.Star) *OptionalStarResp_Star {
	return &OptionalStarResp_Star{
		Id:                s.ID,
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
		ID:                s.Id,
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
			log.Printf("gRPC EOF, for now...")
			break
		}
		if err != nil {
			log.Printf("failed to fetch a star due to %v", err)
			break
		}
		star := starBuff.ToModel()
		log.Printf("persisting star... %v", *star)
		if _, err := (*ss.db).SaveAll([]model.Star{*star}); err != nil {
			log.Fatalf("failed due to %v", err)
			return err
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
	}
	if star == nil {
		return &OptionalStarResp{Resp: nil}, nil
	}
	protoStar := ToProtobuff(star)
	protoOptStar := OptionalStarResp{Resp: protoStar}
	return &protoOptStar, err
}

func RunStarServer(address string, db *db.DB) error {
	s := &StarStoreServerImpl{db: db, grpcServer: grpc.NewServer(), address: address}
	RegisterStarStoreServer(s.grpcServer, s)
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Printf("star store listen error: %v", err)
		return err
	}
	return s.grpcServer.Serve(listener)
}

// Client
type StarStoreClientImpl struct {
	conn *grpc.ClientConn
}

func NewStarClient(grpcUri string) (*StarStoreClientImpl, error) {
	conn, err := grpc.Dial(grpcUri, grpc.WithBlock(), grpc.WithInsecure()) // awaits the connection, no transport security (eg. TLS/SSL)
	if err != nil {
		return nil, err
	}
	return &StarStoreClientImpl{conn}, nil
}

type processor func(model.Star) error
type cleanup func()

func (ss *StarStoreClientImpl) GetStarPersistor() (processor, cleanup, error) {
	client := NewStarStoreClient(ss.conn)
	ctx, ctxCancel := context.WithTimeout(context.Background(), POST_TIMEOUT)
	stream, err := client.PersistStars(ctx)
	if err != nil {
		return nil, nil, err
	}
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

func (ss *StarStoreClientImpl) GetStar(id string) (*model.Star, error) {
	client := NewStarStoreClient(ss.conn)
	ctx, cancel := context.WithTimeout(context.Background(), GET_TIMEOUT)
	defer cancel()
	optProtoStar, err := client.GetStar(ctx, &StarReq{StarId: id})
	if err != nil {
		return nil, err
	}
	protoStar := optProtoStar.GetResp()
	if protoStar == nil {
		return nil, fmt.Errorf("failed to find the star for id %v", id)
	}
	return protoStar.ToModel(), nil
}
