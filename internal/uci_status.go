package internal

import (
	"sync"
	"sync/atomic"
)

type UciStatus struct {
	state int32
	stop  chan struct{}
	once  sync.Once
}

const (
	StatusStopped int32 = iota
	StatusRunning
)

func NewStatus() *UciStatus {
	return &UciStatus{
		stop: make(chan struct{}),
	}
}

func (status *UciStatus) Stop() {
	atomic.StoreInt32(&status.state, StatusStopped)
	status.once.Do(func() {
		close(status.stop)
	})
}

func (status *UciStatus) Done() <-chan struct{} {
	return status.stop
}

func (status *UciStatus) Start() {
	atomic.StoreInt32(&status.state, StatusRunning)
}

func (status *UciStatus) Get() int32 {
	return atomic.LoadInt32(&status.state)
}

func (status *UciStatus) Set(state int32) {
	atomic.StoreInt32(&status.state, state)
}

func (status *UciStatus) IsStopped() bool {
	return atomic.LoadInt32(&status.state) == StatusStopped
}
