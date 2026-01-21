// Package schan provide a channel with safe write
package schan

import (
    "sync"
    "sync/atomic"
)

type SafeChannel[T any] struct {
    ch     chan T
    closed atomic.Bool
    once   sync.Once
}

func New[T any](bufferSize int) *SafeChannel[T] {
    return &SafeChannel[T]{
        ch: make(chan T, bufferSize),
    }
}

func (sc *SafeChannel[T]) Send(value T) (ok bool) {
    if sc.closed.Load() {
        return false
    }
    
    defer func() {
        if recover() != nil {
            ok = false
        }
    }()
    
    sc.ch <- value
    return true
}

func (sc *SafeChannel[T]) Receive() (T, bool) {
    val, ok := <-sc.ch
    return val, ok
}

func (sc *SafeChannel[T]) Close() {
    sc.once.Do(func() {
        sc.closed.Store(true)
        close(sc.ch)
    })
}

func (sc *SafeChannel[T]) IsClosed() bool {
    return sc.closed.Load()
}

func (sc *SafeChannel[T]) Chan() <-chan T {
    return sc.ch
}
