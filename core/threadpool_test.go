package disel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestThreadPool(t *testing.T) {
	var wg sync.WaitGroup
	pool := NewThreadPool(1, &wg)
	pool.Add(func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Job Id-1")
	})
	pool.Add(func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Job Id-2")
	})
	pool.Add(func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Job Id-3")
	})
	pool.Add(func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Job Id-4")
	})
	pool.Add(func() {
		time.Sleep(4 * time.Second)
		fmt.Println("Job Id-5")
	})
	pool.Add(func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Job Id-6")
	})
	pool.Wait()
}
