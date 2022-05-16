package nano

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestPop(t *testing.T) {
	testCases := []struct {
		desc                string
		queue               Q
		expectedAllowedPops int
		expectedIdx         int
		expectedIsEmpty     bool
		queueSizeFactor     int
	}{
		{
			desc:            "Zero value of queue is empty",
			queue:           NewQ(),
			expectedIdx:     -1,
			expectedIsEmpty: true,
			queueSizeFactor: 6,
		},
		{
			desc: "Pops are allowed as many times as there are jobs in the queue, and when completed return empty",
			queue: func() Q {
				q := NewQ()
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
			queue: func() Q {
				q := NewQ()
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
		desc               string
		queue              Q
		expectedIsFull     bool
		nextExpectedIsFull bool
		expectedIdx        int
		nextExpectedIdx    int
		queueSizeFactor    int
	}{
		{
			desc:               "Zero value of queue allows pushing",
			queue:              NewQ(),
			expectedIsFull:     false,
			nextExpectedIsFull: false,
			expectedIdx:        0,
			nextExpectedIdx:    1,
			queueSizeFactor:    6,
		},
		{
			desc:               "Queue with random amount of jobs less than size factor allows pushing",
			queue:              Q(randomAmountOfJobs),
			expectedIsFull:     false,
			nextExpectedIsFull: false,
			expectedIdx:        randomAmountOfJobs,
			nextExpectedIdx:    randomAmountOfJobs + 1,
			queueSizeFactor:    6,
		},
		{
			desc:               "Queue at the size factor is marked as full and cannot be pushed to",
			queue:              Q(63),
			expectedIsFull:     true,
			nextExpectedIsFull: false,
			expectedIdx:        -1,
			nextExpectedIdx:    0,
			queueSizeFactor:    6,
		},
		{
			desc:               "Overflow is protected against",
			queue:              Q(4294967295),
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
	q := NewQ()
	const queueSizeFactor = 6
	const availableSlots = 1 << queueSizeFactor
	jobs := [availableSlots]int{}

	var wg sync.WaitGroup
	wg.Add(2) // Add producer and consumer goroutines

	// Producer
	producedSum := 0
	go func() {
		defer wg.Done()
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

	sum := 0
	// Consumer
	go func() {
		defer wg.Done()
		emptyAttempts := 0

		for {
			slot, savepoint, isEmpty := q.Pop(queueSizeFactor)
			if isEmpty {
				emptyAttempts++
				if emptyAttempts > 1000 {
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
