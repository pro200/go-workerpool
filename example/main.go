package main

import (
	"fmt"

	"github.com/pro200/go-workerpool"
)

type Job struct {
	ID int
}

type Result struct {
	ID      int
	Message string
}

// JobHandler 구현
type Handler struct{}

func (h *Handler) Process(job Job) (Result, error) {
	// 실제 작업 로직 작성
	return Result{ID: job.ID, Message: "Success"}, nil
}

// JobProducer 구현
type Producer struct {
	count int
	max   int
}

func (p *Producer) Next() (Job, bool, error) {
	if p.count >= p.max {
		return Job{}, false, nil
	}
	p.count++
	return Job{ID: p.count}, true, nil
}

func (p *Producer) Close() error {
	return nil
}

func main() {
	handler := &Handler{}
	producer := &Producer{max: 10}

	// 워커 5개로 풀 생성
	pool := workerpool.NewWorkerPool[Job, Result](3, handler)

	// 실행 및 결과 수신
	results := pool.Run(producer)

	for r := range results {
		if r.Err != nil {
			fmt.Printf("Error: %v\n", r.Err)
			continue
		}
		fmt.Printf("Result: %v\n", r.Value)
	}
}
