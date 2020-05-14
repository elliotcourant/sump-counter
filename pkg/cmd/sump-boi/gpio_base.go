// +build !arm

package main

import (
	"context"
	"github.com/elliotcourant/sump-counter/pkg/ioserver"
	"google.golang.org/grpc"
)

var (
	_ GPIOCloser                 = &mockGPIO{}
	_ ioserver.GPIOServiceClient = &mockClient{}
)

type mockGPIO struct {
	client ioserver.GPIOServiceClient
}

type mockClient struct {
}

func NewGPIO() (GPIOCloser, error) {
	return &mockGPIO{
		client: mockClient{},
	}, nil
}

func (g *mockGPIO) Client() ioserver.GPIOServiceClient {
	return g.client
}

func (g *mockGPIO) Close() error {
	return nil
}

func (m mockClient) WatchPinState(ctx context.Context, in *ioserver.Pin, opts ...grpc.CallOption) (ioserver.GPIOService_WatchPinStateClient, error) {
	panic("implement me")
}

func (m mockClient) UpdatePinState(ctx context.Context, in *ioserver.UpdatePinStateRequest, opts ...grpc.CallOption) (*ioserver.PinState, error) {
	panic("implement me")
}

func (m mockClient) GetPinState(ctx context.Context, in *ioserver.Pin, opts ...grpc.CallOption) (*ioserver.PinState, error) {
	panic("implement me")
}
