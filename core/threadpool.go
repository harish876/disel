package disel

import (
	"fmt"
	"sync"
)

type Work func()

type Worker struct {
	Id       int
	WorkChan chan Work
	Wg       *sync.WaitGroup
}

type ThreadPool struct {
	PoolSize int
	Workers  []Worker
	WorkChan chan Work
	wg       *sync.WaitGroup
}

func NewWorker(id int, workChan chan Work, wg *sync.WaitGroup) Worker {
	return Worker{
		Id:       id,
		WorkChan: workChan,
		Wg:       wg,
	}
}

func (w *Worker) Start() {
	go func() {
		defer w.Wg.Done()
		for work := range w.WorkChan {
			fmt.Printf("Task Picked by Worker: %d\n", w.Id)
			work()
			fmt.Printf("Task Completed by Worker: %d\n", w.Id)
		}
	}()
}

func NewThreadPool(poolSize int, wg *sync.WaitGroup) ThreadPool {
	workers := make([]Worker, 0)
	workChan := make(chan Work, 10)
	for i := 0; i < poolSize; i++ {
		id := i + 1
		worker := NewWorker(id, workChan, wg)
		workers = append(workers, worker)
		worker.Start()
	}
	wg.Add(poolSize)

	return ThreadPool{
		PoolSize: poolSize,
		Workers:  workers,
		WorkChan: workChan,
		wg:       wg,
	}
}

func (t *ThreadPool) Add(work Work) {
	t.WorkChan <- work
}

func (t *ThreadPool) Wait() {
	close(t.WorkChan)
	t.wg.Wait()
}
