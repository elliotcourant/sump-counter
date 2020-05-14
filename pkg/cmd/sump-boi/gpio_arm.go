// +build arm

package main

import (
	"fmt"
	"github.com/elliotcourant/sump-counter/pkg/ioserver"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	_ GPIOCloser = &grpcGPIOBase{}
)

type grpcGPIOBase struct {
	conn   *grpc.ClientConn
	client ioserver.GPIOServiceClient
}

func NewGPIO() (GPIOCloser, error) {
	grpcConn, err := grpc.Dial(fmt.Sprintf("unix:%s", gpioSocketPath), grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial gpio service")
	}

	gpioClient := ioserver.NewGPIOServiceClient(grpcConn)

	return &grpcGPIOBase{
		conn:   grpcConn,
		client: gpioClient,
	}, nil
}

func (g *grpcGPIOBase) Client() ioserver.GPIOServiceClient {
	return g.client
}

func (g *grpcGPIOBase) Close() error {
	return g.conn.Close()
}
