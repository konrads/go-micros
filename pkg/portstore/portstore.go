package portstore

import (
	"log"

	"golang.org/x/net/context"
)

type PortStore struct{}

func (s *PortStore) PersistPorts(PortStore_PersistPortsServer) error {
	log.Printf("...persisting: ...")
	return nil
}

func (s *PortStore) GetPort(ctx context.Context, req *PortReq) (*OptionalPortResp, error) {
	log.Printf("...getting: %v", *req)
	return &OptionalPortResp{}, nil
}
