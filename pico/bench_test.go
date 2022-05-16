package pico

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// BenchmarkProduce-10    	  507206	      2287 ns/op	       0 B/op	       0 allocs/op
func BenchmarkProduce(b *testing.B) {
	b.StopTimer()

	q := NewQ()
	const queueSizeFactor = 6
	const availableSlots = 1 << queueSizeFactor
	jobs := [availableSlots]int{}
	var wg sync.WaitGroup
	wg.Add(1) // Add consumer goroutine

	producedSum := 0
	completedProducing := int32(0)

	// Start Consumer Outside of the loop since we are benchmarking producing
	sum := 0
	go func() {
		defer wg.Done()

		for {
			slot, isEmpty := q.Pop(queueSizeFactor)
			if isEmpty {
				if atomic.LoadInt32(&completedProducing) > 0 {
					break
				}
				<-time.After(1 * time.Millisecond) // Allow some sleeping so that it's not a pure busy loop
				continue
			}

			job := jobs[slot]
			q.PopCommit()
			sum += job
		}
	}()

	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		// Producer is inside the loop so we can measure producing performance
		for i := 0; i < 1000; i++ {
			slot, isFull := q.Push(queueSizeFactor)
			if isFull {
				continue
			}

			// Simulate job creation allocation
			job := i * 5
			jobs[slot] = job
			q.PushCommit()
			producedSum += job
		}
	}
	atomic.AddInt32(&completedProducing, 1)

	b.StopTimer()
	wg.Wait()

	if producedSum != sum {
		b.Errorf("expected the sum to be %d but got %d", producedSum, sum)
	}
}

// BenchmarkConsume-10    	788615026	         1.353 ns/op	       0 B/op	       0 allocs/op
func BenchmarkConsume(b *testing.B) {
	b.StopTimer()

	q := NewQ()
	const queueSizeFactor = 6
	const availableSlots = 1 << queueSizeFactor
	jobs := [availableSlots]int{}
	var wg sync.WaitGroup
	wg.Add(1) // Add producer goroutine

	sum := 0

	// Start Producer Outside of the loop since we are benchmarking consuming
	producedSum := 0
	completedProducing := int32(0)
	go func() {
		defer func() {
			atomic.AddInt32(&completedProducing, 1)
			wg.Done()
		}()
		fullAttempts := 0

		for i := 0; i < 1000; i++ {
			slot, isFull := q.Push(queueSizeFactor)
			if isFull {
				fullAttempts++
				if fullAttempts > 1000 {
					break
				}
				<-time.After(1 * time.Millisecond) // Allow some sleeping so that it's not a pure busy loop
				continue
			}

			jobs[slot] = i
			producedSum += i
			q.PushCommit()
		}
	}()

	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		// Consumer is inside the loop so we can measure consuming performance
		for {
			slot, isEmpty := q.Pop(queueSizeFactor)
			if isEmpty {
				if atomic.LoadInt32(&completedProducing) > 0 {
					break
				}
				continue
			}

			job := jobs[slot]
			q.PopCommit()
			sum += job
		}
	}

	b.StopTimer()
	wg.Wait()

	if producedSum != sum {
		b.Errorf("expected the sum to be %d but got %d", producedSum, sum)
	}
}
