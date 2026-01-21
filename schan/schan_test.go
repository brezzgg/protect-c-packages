package schan

import (
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"unbuffered", 0},
		{"small", 10},
		{"large", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := New[int](tt.size)
			if sc == nil {
				t.Fatal("New returned nil")
			}
			if sc.IsClosed() {
				t.Error("channel should not be closed")
			}
			sc.Close()
		})
	}
}

func TestSendReceive(t *testing.T) {
	sc := New[int](10)
	defer sc.Close()

	values := []int{0, 42, -100}
	for _, val := range values {
		if !sc.Send(val) {
			t.Errorf("Send(%d) failed", val)
		}
		got, ok := sc.Receive()
		if !ok || got != val {
			t.Errorf("got (%d, %v), want (%d, true)", got, ok, val)
		}
	}
}

func TestSendAfterClose(t *testing.T) {
	sc := New[int](10)
	sc.Close()

	if sc.Send(42) {
		t.Error("Send to closed channel should return false")
	}
	if !sc.IsClosed() {
		t.Error("IsClosed() should return true")
	}
}

func TestReceiveFromClosed(t *testing.T) {
	sc := New[int](10)
	sc.Send(1)
	sc.Send(2)
	sc.Close()

	if val, ok := sc.Receive(); !ok || val != 1 {
		t.Errorf("got (%d, %v), want (1, true)", val, ok)
	}
	if val, ok := sc.Receive(); !ok || val != 2 {
		t.Errorf("got (%d, %v), want (2, true)", val, ok)
	}
	if _, ok := sc.Receive(); ok {
		t.Error("Receive from empty closed channel should return ok=false")
	}
}

func TestMultipleClose(t *testing.T) {
	sc := New[int](10)
	for range 10 {
		sc.Close()
	}
	if !sc.IsClosed() {
		t.Error("channel should be closed")
	}
}

func TestConcurrentSendReceive(t *testing.T) {
	sc := New[int](100)
	var wg sync.WaitGroup

	const writers, readers, messages = 5, 3, 100

	sent := make(map[int]bool)
	received := make(map[int]bool)
	var mu sync.Mutex

	for i := range writers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range messages {
				val := id*1000 + j
				if sc.Send(val) {
					mu.Lock()
					sent[val] = true
					mu.Unlock()
				}
			}
		}(i)
	}

	for range readers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				val, ok := sc.Receive()
				if !ok {
					return
				}
				mu.Lock()
				received[val] = true
				mu.Unlock()
			}
		}()
	}

	time.Sleep(100 * time.Millisecond)
	sc.Close()
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()

	if len(sent) == 0 {
		t.Error("no messages sent")
	}
	for val := range sent {
		if !received[val] {
			t.Errorf("value %d sent but not received", val)
		}
	}
}

func TestConcurrentSendClose(t *testing.T) {
	sc := New[int](100)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range sc.Chan() {
		}
	}()

	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range 100 {
				sc.Send(id*100 + j)
			}
		}(i)
	}

	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			sc.Close()
		}()
	}

	wg.Wait()

	if !sc.IsClosed() {
		t.Error("channel should be closed")
	}
	if sc.Send(999) {
		t.Error("Send to closed channel should return false")
	}
}

func TestChanRange(t *testing.T) {
	sc := New[int](10)

	go func() {
		for i := range 5 {
			sc.Send(i)
		}
		sc.Close()
	}()

	received := []int{}
	for val := range sc.Chan() {
		received = append(received, val)
	}

	if len(received) != 5 {
		t.Errorf("received %d values, want 5", len(received))
	}
	for i, val := range received {
		if val != i {
			t.Errorf("received[%d] = %d, want %d", i, val, i)
		}
	}
}

func TestGenerics(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		sc := New[string](10)
		defer sc.Close()
		sc.Send("hello")
		if val, ok := sc.Receive(); !ok || val != "hello" {
			t.Errorf("got (%s, %v), want (hello, true)", val, ok)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type data struct {
			id   int
			name string
		}
		sc := New[data](10)
		defer sc.Close()

		d := data{id: 1, name: "test"}
		sc.Send(d)
		if val, ok := sc.Receive(); !ok || val != d {
			t.Errorf("got (%+v, %v), want (%+v, true)", val, ok, d)
		}
	})
}

func TestBufferFull(t *testing.T) {
	const size = 3
	sc := New[int](size)
	defer sc.Close()

	for i := range size {
		if !sc.Send(i) {
			t.Errorf("Send(%d) failed", i)
		}
	}

	done := make(chan bool)
	go func() {
		sc.Send(999)
		done <- true
	}()

	select {
	case <-done:
		t.Error("Send should block on full buffer")
	case <-time.After(50 * time.Millisecond):
	}

	sc.Receive()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Error("Send should complete after buffer has space")
	}
}

func TestZeroValues(t *testing.T) {
	sc := New[int](10)
	defer sc.Close()
	sc.Send(0)
	if val, ok := sc.Receive(); !ok || val != 0 {
		t.Errorf("got (%d, %v), want (0, true)", val, ok)
	}
}

func BenchmarkSend(b *testing.B) {
	sc := New[int](b.N)
	go func() {
		for range sc.Chan() {
		}
	}()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sc.Send(i)
	}
	b.StopTimer()
	sc.Close()
}

func BenchmarkReceive(b *testing.B) {
	sc := New[int](b.N)
	go func() {
		for i := 0; i < b.N; i++ {
			sc.Send(i)
		}
		sc.Close()
	}()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sc.Receive()
	}
}

func BenchmarkConcurrent(b *testing.B) {
	sc := New[int](10000)
	go func() {
		for range sc.Chan() {
		}
	}()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sc.Send(i)
			i++
		}
	})
	b.StopTimer()
	sc.Close()
}

func BenchmarkNative(b *testing.B) {
	ch := make(chan int, b.N)
	go func() {
		for range ch {
		}
	}()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch <- i
	}
	b.StopTimer()
	close(ch)
}
