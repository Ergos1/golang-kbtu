package grpc

import (
	"context"
	"example.com/api"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type UserRepo struct{
	data map[uint64]*api.User
	api.UnimplementedUserServiceServer
	mu *sync.RWMutex
}

func (c *UserRepo) All(ctx context.Context, req *api.Empty) (*api.Users, error){
	c.mu.RLock()
	defer c.mu.RUnlock()
	users := &api.Users{
		Users: make([]*api.User, 0, len(c.data)),
	}
	for _, user := range c.data{
		users.Users = append(users.Users, user)
	}
	return users, nil
}
func (c *UserRepo) ByID(ctx context.Context, req *api.Id) (*api.User, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	id := req.Id
	fmt.Println(id)
	user, ok := c.data[id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("user with id %d does not exist", req.Id))
	}
	return user, nil
}
func (c*UserRepo) Create(ctx context.Context, req *api.User) (*api.Empty, error){
	c.mu.Lock()
	defer c.mu.Unlock()
	id := req.Id
	if _, exist := c.data[id]; exist {
		return &api.Empty{}, status.Errorf(codes.AlreadyExists, fmt.Sprintf("[error] user with id %d exist", req.Id))
	}
	c.data[req.Id] = req
	return &api.Empty{}, nil
}

func (c *UserRepo) Update (ctx context.Context, req *api.User) (*api.Empty, error){
	c.mu.Lock()
	defer c.mu.Unlock()
	id := req.Id
	if _, exist := c.data[id]; !exist {
		return &api.Empty{}, status.Errorf(codes.NotFound, fmt.Sprintf("[error] user with id %d does not exist", req.Id))
	}
	c.data[req.Id] = req
	return &api.Empty{}, nil
}
func (c *UserRepo) Delete(ctx context.Context, req *api.Id) (*api.Empty, error){
	c.mu.Lock()
	defer c.mu.Unlock()
	id := req.Id
	if _, exist := c.data[id]; !exist {
		return &api.Empty{}, status.Errorf(codes.NotFound, fmt.Sprintf("[error] user with id %d does not exist", req.Id))
	}
	delete(c.data, req.Id)
	return &api.Empty{}, nil
}