package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numJobs   = 10
	numWorker = 4
)

var (
	wg sync.WaitGroup
)

type Job struct {
	Id         int
	Complexity int
	Result     int
}

func (j Job) String() string {
	return fmt.Sprint(j.Id)
}

func worker(id int, jobs <-chan Job) {
	for job := range jobs {
		fmt.Printf("Worker %d started work with job %v \n", id, job)
		time.Sleep(time.Duration(job.Complexity) * time.Second)
		fmt.Printf("Worker %d finished work with job %v \n", id, job)
	}
	wg.Done()
}

func main() {
	wg.Add(numWorker)
	jobs := make(chan Job, numJobs)

	for i := 1; i <= numWorker; i++ {
		go worker(i, jobs)
	}

	for i := 1; i <= numJobs; i++ {
		job := Job{i, rand.Intn(3), rand.Intn(10)}
		jobs <- job
	}

	close(jobs)

	wg.Wait()

	fmt.Println("All workers finished their job")

}
