package grpc

import (
	"context"
	"example.com/api"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type NonFungibleTokenRepo struct{
	data map[uint64]*api.NonFungibleToken
	api.UnimplementedNonFungibleTokenServiceServer
	mu *sync.RWMutex
}

func (c *NonFungibleTokenRepo) All(ctx context.Context, req *api.Empty) (*api.NonFungibleTokens, error){
	c.mu.RLock()
	defer c.mu.RUnlock()
	nonFungibleTokens := &api.NonFungibleTokens{
		NonFungibleTokens: make([]*api.NonFungibleToken, 0, len(c.data)),
	}
	for _, nonFungibleToken := range c.data{
		nonFungibleTokens.NonFungibleTokens = append(nonFungibleTokens.NonFungibleTokens, nonFungibleToken)
	}
	return nonFungibleTokens, nil
}

func (c *NonFungibleTokenRepo) ByID(ctx context.Context, req *api.Id) (*api.NonFungibleToken, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	id := req.Id
	fmt.Println(id)
	nonFungibleToken, ok := c.data[id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("nonFungibleToken with id %d does not exist", req.Id))
	}
	return nonFungibleToken, nil
}
func (c*NonFungibleTokenRepo) Create(ctx context.Context, req *api.NonFungibleToken) (*api.Empty, error){
	c.mu.Lock()
	defer c.mu.Unlock()
	id := req.Id
	if _, exist := c.data[id]; exist {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprintf("[error] non fungible token with id %d exists", req.Id))
	}
	c.data[req.Id] = req
	return &api.Empty{}, nil
}

func (c *NonFungibleTokenRepo) Update (ctx context.Context, req *api.NonFungibleToken) (*api.Empty, error){
	c.mu.Lock()
	defer c.mu.Unlock()
	id := req.Id
	if _, exist := c.data[id]; !exist {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("[error] non fungible token with id %d does not exist", req.Id))
	}
	c.data[req.Id] = req
	return &api.Empty{}, nil
}
func (c *NonFungibleTokenRepo) Delete(ctx context.Context, req *api.Id) (*api.Empty, error){
	c.mu.Lock()
	defer c.mu.Unlock()
	id := req.Id
	if _, exist := c.data[id]; !exist {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("[error] non fungible token with id %d does not exist", req.Id))
	}
	delete(c.data, req.Id)
	return &api.Empty{}, nil
}