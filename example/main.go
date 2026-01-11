package main

import (
	"fmt"

	"github.com/pro200/go-workerpool"
)

type Job struct {
	Num int
}

type handler struct{}

func (h *handler) Process(job any) (any, error) {
	j, ok := job.(Job)
	if !ok {
		return nil, fmt.Errorf("invalid job type")
	}
	fmt.Println("process->", j.Num)
	return j.Num, nil
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
		if r.Err != nil {
			fmt.Println("error->", r.Value.(int), r.Err)
			continue
		}
		fmt.Println("done->", r.Value.(int))
	}
}
