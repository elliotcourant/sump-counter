package ioserver

import (
	"context"
	"google.golang.org/grpc/metadata"
	"testing"
)

var (
	_ GPIOService_WatchPinStateServer = &watchPinStateServerMock{}
	_ GPIOService_WatchPinStateClient = &watchPinStateClientMock{}
)

type (
	watchPinStateServerMock struct {
		t           *testing.T
		ctx         context.Context
		cancelChan  chan struct{}
		sendChannel chan PinState
	}

	watchPinStateClientMock struct {
		t              *testing.T
		ctx            context.Context
		cancelChan     chan struct{}
		receiveChannel chan PinState
	}
)

func NewMockWatchServerClient(t *testing.T) (watchServer GPIOService_WatchPinStateServer, watchClient GPIOService_WatchPinStateClient, cancel context.CancelFunc) {
	channel := make(chan PinState)
	serverCtx, serverCancel := context.WithCancel(context.Background())
	clientCtx, clientCancel := context.WithCancel(context.Background())

	server := &watchPinStateServerMock{
		t:           t,
		ctx:         serverCtx,
		cancelChan:  make(chan struct{}, 1),
		sendChannel: channel,
	}

	client := &watchPinStateClientMock{
		t:              t,
		ctx:            clientCtx,
		cancelChan:     make(chan struct{}, 1),
		receiveChannel: channel,
	}

	return server, client, func() {
		client.cancelChan <- struct{}{}
		serverCancel()
		clientCancel()
	}
}

func (w *watchPinStateServerMock) Send(state *PinState) error {
	w.sendChannel <- *state

	return nil
}

func (w watchPinStateServerMock) SetHeader(md metadata.MD) error {
	panic("implement me")
}

func (w watchPinStateServerMock) SendHeader(md metadata.MD) error {
	panic("implement me")
}

func (w watchPinStateServerMock) SetTrailer(md metadata.MD) {
	panic("implement me")
}

func (w *watchPinStateServerMock) Context() context.Context {
	return w.ctx
}

func (w watchPinStateServerMock) SendMsg(m interface{}) error {
	panic("implement me")
}

func (w watchPinStateServerMock) RecvMsg(m interface{}) error {
	panic("implement me")
}

func (w *watchPinStateClientMock) Recv() (*PinState, error) {
	select {
	case <-w.cancelChan:
		return nil, nil
	case state := <-w.receiveChannel:
		return &state, nil
	}
}

func (w watchPinStateClientMock) Header() (metadata.MD, error) {
	panic("implement me")
}

func (w watchPinStateClientMock) Trailer() metadata.MD {
	panic("implement me")
}

func (w watchPinStateClientMock) CloseSend() error {
	panic("implement me")
}

func (w *watchPinStateClientMock) Context() context.Context {
	return w.ctx
}

func (w watchPinStateClientMock) SendMsg(m interface{}) error {
	panic("implement me")
}

func (w watchPinStateClientMock) RecvMsg(m interface{}) error {
	panic("implement me")
}
