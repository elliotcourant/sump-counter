package main

import (
	"github.com/elliotcourant/sump-counter/pkg/ioserver"
)

type GPIOCloser interface {
	Client() ioserver.GPIOServiceClient
	Close() error
}
