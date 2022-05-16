package micro

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPop(t *testing.T) {
	testCases := []struct {
		queue               *Q
		desc                string
		expectedAllowedPops int
		expectedIdx         int
		queueSizeFactor     int
		expectedIsEmpty     bool
	}{
		{
			desc:            "Zero value of queue is empty",
			queue:           NewQ(6),
			expectedIdx:     -1,
			expectedIsEmpty: true,
			queueSizeFactor: 6,
		},
		{
			desc: "Pops are allowed as many times as there are jobs in the queue, and when completed return empty",
			queue: func() *Q {
				q := NewQ(6)
				for i := 0; i < 10; i++ {
					q.PushCommit() // 10 pushes
				}
				return q
			}(),
			expectedAllowedPops: 10,
			expectedIdx:         -1,
			expectedIsEmpty:     true,
			queueSizeFactor:     6,
		},
		{
			desc: "Pops are allowed as many times as there are jobs in the queue, and when not completed return the next index",
			queue: func() *Q {
				q := NewQ(6)
				for i := 0; i < 10; i++ {
					q.PushCommit() // 10 pushes
				}
				return q
			}(),
			expectedAllowedPops: 8,
			expectedIdx:         8,
			expectedIsEmpty:     false,
			queueSizeFactor:     6,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(subT *testing.T) {
			for i := 0; i < tC.expectedAllowedPops; i++ {
				idx, savepoint, isEmpty := tC.queue.Pop(tC.queueSizeFactor)
				if idx == -1 || isEmpty {
					subT.Errorf("unexpected empty queue during allowed pops at pop number %d", i)
					return
				}

				if i != idx {
					subT.Errorf("expected popped job to be %d but got %d", i, idx)
				}

				if !tC.queue.PopCommit(savepoint) {
					subT.Errorf("expected pop commit to pass normally on job %d but got commit failed", i)
				}
			}

			idx, _, isEmpty := tC.queue.Pop(tC.queueSizeFactor)
			if tC.expectedIdx != idx {
				subT.Errorf("expected the returned index to be %d, got %d", tC.expectedIdx, idx)
			}

			if tC.expectedIsEmpty != isEmpty {
				subT.Errorf("expected isEmpty to be %t, got %t", tC.expectedIsEmpty, isEmpty)
			}
		})
	}
}

func TestPush(t *testing.T) {
	randomAmountOfJobs := rand.Intn(62)

	testCases := []struct {
		queue              *Q
		desc               string
		expectedIdx        int
		nextExpectedIdx    int
		queueSizeFactor    int
		expectedIsFull     bool
		nextExpectedIsFull bool
	}{
		{
			desc:               "Zero value of queue allows pushing",
			queue:              NewQ(6),
			expectedIsFull:     false,
			nextExpectedIsFull: false,
			expectedIdx:        0,
			nextExpectedIdx:    1,
			queueSizeFactor:    6,
		},
		{
			desc:               "Queue with random amount of jobs less than size factor allows pushing",
			queue:              &Q{q: uint32(randomAmountOfJobs), queueSizeFactor: 6},
			expectedIsFull:     false,
			nextExpectedIsFull: false,
			expectedIdx:        randomAmountOfJobs,
			nextExpectedIdx:    randomAmountOfJobs + 1,
			queueSizeFactor:    6,
		},
		{
			desc:               "Queue at the size factor is marked as full and cannot be pushed to",
			queue:              &Q{q: uint32(63), queueSizeFactor: 6},
			expectedIsFull:     true,
			nextExpectedIsFull: false,
			expectedIdx:        -1,
			nextExpectedIdx:    0,
			queueSizeFactor:    6,
		},
		{
			desc:               "Overflow is protected against",
			queue:              &Q{q: uint32(4294967295), queueSizeFactor: 6},
			expectedIsFull:     false,
			nextExpectedIsFull: false,
			expectedIdx:        63,
			nextExpectedIdx:    0,
			queueSizeFactor:    6,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(subT *testing.T) {
			idx, isFull := tC.queue.Push(tC.queueSizeFactor)
			if tC.expectedIdx != idx {
				subT.Errorf("expected the returned index to be %d, got %d", tC.expectedIdx, idx)
			}

			if tC.expectedIsFull != isFull {
				subT.Errorf("expected isFull to be %t, got %t", tC.expectedIsFull, isFull)
			}

			tC.queue.PushCommit()

			idx, isFull = tC.queue.Push(tC.queueSizeFactor)
			if tC.nextExpectedIdx != idx {
				subT.Errorf("expected the next returned index to be %d, got %d", tC.nextExpectedIdx, idx)
			}

			if tC.nextExpectedIsFull != isFull {
				subT.Errorf("expected next isFull to be %t, got %t", tC.nextExpectedIsFull, isFull)
			}
		})
	}
}

func TestConcurrentWorkSingleConsumer(t *testing.T) {
	const queueSizeFactor = 6
	q := NewQ(queueSizeFactor)
	const availableSlots = 1 << queueSizeFactor
	jobs := [availableSlots]int{}

	var wg sync.WaitGroup
	wg.Add(2) // Add producer and consumer goroutines

	// Producer
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

	// Consumer
	sum := 0
	go func() {
		defer wg.Done()

		for {
			slot, savepoint, isEmpty := q.Pop(queueSizeFactor)
			if isEmpty {
				if atomic.LoadInt32(&completedProducing) > 0 {
					break
				}
				<-time.After(1 * time.Millisecond) // Allow some sleeping so that it's not a pure busy loop
				continue
			}

			job := jobs[slot]
			if !q.PopCommit(savepoint) {
				continue // Commit failed so we can't run the job
			}
			sum += job
		}
	}()

	wg.Wait()

	if producedSum != sum {
		t.Errorf("expected the sum to be %d but got %d", producedSum, sum)
	}
}

func TestConcurrentWorkMultipleConsumers(t *testing.T) {
	const queueSizeFactor = 6
	q := NewQ(queueSizeFactor)
	const availableSlots = 1 << queueSizeFactor
	// [availableSlots]int64{} is the underlying type of the atomic value.
	atomicJobs := &atomic.Value{}
	atomicJobs.Store([availableSlots]int64{})

	var wg sync.WaitGroup
	wg.Add(1) // Add producer goroutine

	// Producer
	producedSum := int64(0)
	completedProducing := int32(0)
	go func() {
		defer func() {
			atomic.AddInt32(&completedProducing, 1)
			wg.Done()
		}()
		fullAttempts := 0

		for i := int64(0); i < 1000; i++ {
			slot, isFull := q.Push(queueSizeFactor)
			if isFull {
				fullAttempts++
				if fullAttempts > 1000 {
					break
				}
				<-time.After(1 * time.Millisecond) // Allow some sleeping so that it's not a pure busy loop
				continue
			}

			jobs := atomicJobs.Load().([availableSlots]int64)
			newJobs := jobs
			newJobs[slot] = i
			if atomicJobs.CompareAndSwap(jobs, newJobs) {
				producedSum += i
			}
			q.PushCommit()
		}
	}()

	// Consumer
	sum := int64(0)
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				slot, savepoint, isEmpty := q.Pop(queueSizeFactor)
				if isEmpty {
					if atomic.LoadInt32(&completedProducing) > 0 {
						break
					}
					<-time.After(1 * time.Millisecond) // Allow some sleeping so that it's not a pure busy loop
					continue
				}

				jobs := atomicJobs.Load().([availableSlots]int64)
				job := jobs[slot]
				if !q.PopCommit(savepoint) {
					continue // Commit failed so we can't run the job
				}

				atomic.AddInt64(&sum, job)
			}
		}()
	}

	wg.Wait()

	if producedSum < 1000 {
		t.Errorf("expected the produced sum to be greater than 1000 but got %d", producedSum)
	}

	if producedSum != sum {
		t.Errorf("expected the sum to be %d but got %d", producedSum, sum)
	}
}
