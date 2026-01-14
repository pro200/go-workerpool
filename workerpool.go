package workerpool

import (
	"sync"
)

type JobHandler[J any, R any] interface {
	Process(job J) (result R)
}

type JobProducer[J any] interface {
	Next() (job J, ok bool) // ok가 false일 경우 더이상 job이 없음
	Close() error
}

type WorkerPool[J any, R any] struct {
	maxWorkers int
	jobs       chan J
	results    chan R
	wg         sync.WaitGroup
	handler    JobHandler[J, R]
}

func NewWorkerPool[J any, R any](maxWorkers int, handler JobHandler[J, R]) *WorkerPool[J, R] {
	return &WorkerPool[J, R]{
		maxWorkers: maxWorkers,
		jobs:       make(chan J, maxWorkers),
		results:    make(chan R),
		handler:    handler,
	}
}

func (p *WorkerPool[J, R]) Run(producer JobProducer[J]) <-chan R {
	go func() {
		p.startWorkers()
		p.produceJobs(producer)
		p.waitAndClose()
	}()
	return p.results
}

func (p *WorkerPool[J, R]) startWorkers() {
	for i := 0; i < p.maxWorkers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *WorkerPool[J, R]) worker() {
	defer p.wg.Done()
	for job := range p.jobs {
		p.results <- p.handler.Process(job)
	}
}

func (p *WorkerPool[J, R]) produceJobs(producer JobProducer[J]) {
	defer producer.Close()
	for {
		job, ok := producer.Next()
		if !ok {
			break
		}

		// worker가 처리 끝날 때까지 자동 대기
		p.jobs <- job
	}
}

func (p *WorkerPool[J, R]) waitAndClose() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}
