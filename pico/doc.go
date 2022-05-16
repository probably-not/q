// Package pico contains the most minimalistic version of the implementation of this queue.
// It enables a single producer and single consumer to work concurrently on the same slice of jobs,
// with the indices of jobs pushed to the queue and jobs popped from the queue fully managed by the `pico.Q`.
// The slice of jobs itself is managed by an outside source, however access to this slice of jobs should be fully
// managed by the `pico.Q` type.
package pico
