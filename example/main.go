package main

import (
	"fmt"
	"time"

	"github.com/pro200/go-workerpool"
)

type Job struct {
	Num int
}

type handler struct{}

func (h *handler) Process(job any) error {
	time.Sleep(time.Millisecond * 300)

	j, ok := job.(Job)
	if !ok {
		return fmt.Errorf("invalid job type")
	}
	fmt.Println("num->", j.Num)
	return nil
}

type producer struct {
	cnt int
}

func (p *producer) Next() (any, error) {
	p.cnt++
	if p.cnt > 10 {
		return nil, nil
	}
	return Job{Num: p.cnt}, nil
}

func (p *producer) Close() error {
	return nil
}

func main() {
	wp := workerpool.NewWorkerPool(3, &handler{})
	results := wp.Run(&producer{})

	for r := range results {
		fmt.Println("num:", r.Job.(Job).Num, "err:", r.Err)
	}
}
