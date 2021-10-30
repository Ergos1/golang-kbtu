package grpc

import (
	"context"
	"example.com/api"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type WalletRepo struct{
	data map[uint64]*api.Wallet
	api.UnimplementedWalletServiceServer
	mu *sync.RWMutex
}

func (c *WalletRepo) All(ctx context.Context, req *api.Empty) (*api.Wallets, error){
	c.mu.RLock()
	defer c.mu.RUnlock()
	wallets := &api.Wallets{
		Wallets: make([]*api.Wallet, 0, len(c.data)),
	}
	for _, wallet := range c.data{
		wallets.Wallets = append(wallets.Wallets, wallet)
	}
	return wallets, nil
}
func (c *WalletRepo) ByID(ctx context.Context, req *api.Id) (*api.Wallet, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	id := req.Id
	fmt.Println(id)
	wallet, ok := c.data[id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("wallet with id %d does not exist", req.Id))
	}
	return wallet, nil
}
func (c*WalletRepo) Create(ctx context.Context, req *api.Wallet) (*api.Empty, error){
	c.mu.Lock()
	defer c.mu.Unlock()
	id := req.Id
	if _, exist := c.data[id]; exist {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprintf("[error] wallet with id %d exists", req.Id))
	}
	c.data[req.Id] = req
	return &api.Empty{}, nil
}

func (c *WalletRepo) Update (ctx context.Context, req *api.Wallet) (*api.Empty, error){
	c.mu.Lock()
	defer c.mu.Unlock()
	id := req.Id
	if _, exist := c.data[id]; !exist {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("[error] wallet with id %d does not exist", req.Id))
	}
	c.data[req.Id] = req
	return &api.Empty{}, nil
}
func (c *WalletRepo) Delete(ctx context.Context, req *api.Id) (*api.Empty, error){
	c.mu.Lock()
	defer c.mu.Unlock()
	id := req.Id
	if _, exist := c.data[id]; !exist {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("[error] wallet with id %d does not exists", req.Id))
	}
	delete(c.data, req.Id)
	return &api.Empty{}, nil
}