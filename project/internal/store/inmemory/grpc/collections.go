package grpc

import (
	"context"
	"example.com/api"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type CollectionRepo struct{
	data map[uint64]*api.Collection
	api.UnimplementedCollectionServiceServer
	mu *sync.RWMutex
}

func (c *CollectionRepo) All(ctx context.Context, req *api.Empty) (*api.Collections, error){
	c.mu.RLock()
	defer c.mu.RUnlock()
	collections := &api.Collections{
		Collections: make([]*api.Collection, 0, len(c.data)),
	}
	for _, collection := range c.data{
		collections.Collections = append(collections.Collections, collection)
	}
	return collections, nil
}
func (c *CollectionRepo) ByID(ctx context.Context, req *api.Id) (*api.Collection, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	id := req.Id
	collection, ok := c.data[id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("[error] collection with id %d does not exist", req.Id))
	}
	return collection, nil
}
func (c*CollectionRepo) Create(ctx context.Context, req *api.Collection) (*api.Empty, error){
	c.mu.Lock()
	defer c.mu.Unlock()
	id := req.Id // \[ owner id => user ]/
	if _, exist := c.data[id]; exist {
		return &api.Empty{}, status.Errorf(codes.AlreadyExists, fmt.Sprintf("[error] collection with id %d exist", req.Id))
	}
	c.data[req.Id] = req
	return &api.Empty{}, nil
}

func (c *CollectionRepo) Update (ctx context.Context, req *api.Collection) (*api.Empty, error){
	c.mu.Lock()
	defer c.mu.Unlock()
	id := req.Id
	if _, exist := c.data[id]; !exist {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("[error] collection with id %d does not exist", req.Id))
	}
	c.data[req.Id] = req
	return &api.Empty{}, nil
}
func (c *CollectionRepo) Delete(ctx context.Context, req *api.Id) (*api.Empty, error){
	c.mu.Lock()
	defer c.mu.Unlock()
	id := req.Id
	if _, exist := c.data[id]; !exist {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("[error] collection with id %d does not exist", req.Id))
	}
	delete(c.data, req.Id)
	return &api.Empty{}, nil
}