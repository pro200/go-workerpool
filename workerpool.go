package workerpool

import (
	"sync"
)

type JobHandler interface {
	Process(job any) error
}

type JobProducer interface {
	Next() (any, error) // job이 nil일경우 더이상 job이 없음
	Close() error
}

type Result struct {
	Job any
	Err error
}

type WorkerPool struct {
	maxWorkers int
	jobs       chan any
	results    chan Result
	wg         sync.WaitGroup
	handler    JobHandler
}

func NewWorkerPool(maxWorkers int, handler JobHandler) *WorkerPool {
	return &WorkerPool{
		maxWorkers: maxWorkers,
		jobs:       make(chan any, maxWorkers),
		results:    make(chan Result),
		handler:    handler,
	}
}

func (p *WorkerPool) Run(producer JobProducer) <-chan Result {
	go func() {
		defer producer.Close()
		p.startWorkers()
		p.produceJobs(producer)
		p.waitAndClose()
	}()
	return p.results
}

func (p *WorkerPool) startWorkers() {
	for i := 0; i < p.maxWorkers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *WorkerPool) worker() {
	defer p.wg.Done()
	for job := range p.jobs {
		err := p.handler.Process(job)
		p.results <- Result{Job: job, Err: err}
	}
}

func (p *WorkerPool) produceJobs(producer JobProducer) {
	for {
		job, err := producer.Next()
		if job == nil {
			break
		}
		if err != nil {
			p.results <- Result{Job: job, Err: err}
			break
		}
		// backpressure 발생 지점
		// worker가 처리 끝날 때까지 자동 대기
		p.jobs <- job
	}
}

func (p *WorkerPool) waitAndClose() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}
