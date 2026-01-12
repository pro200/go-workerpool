# Go Worker Pool

Go의 제네릭(Generics)을 활용하여 구현된 간단하고 효율적인 워커 풀(Worker Pool) 라이브러리입니다. 다양한 타입의 작업(Job)과 결과(Result)를 안전하게 처리할 수 있습니다.

## 특징

- **제네릭 지원**: `Job`과 `Result` 타입을 자유롭게 정의하여 사용할 수 있습니다.
- **인터페이스 기반**: `JobHandler`와 `JobProducer` 인터페이스를 구현하여 비즈니스 로직을 분리할 수 있습니다.
- **리소스 관리**: 고루틴의 개수를 제한하여 시스템 자원을 효율적으로 사용하며, 작업 완료 후 채널을 자동으로 닫습니다.

## 설치

```bash
go get github.com/pro200/go-workerpool
```

## 사용 방법
### 1. 핸들러 구현
`workerpool.JobHandler` 인터페이스 구현.
```go
type Job struct {
    ID int
}

type Result struct {
    ID      int
    Message string
}

type Handler struct{}

func (h *Handler) Process(job Job) (Result, error) {
    // 실제 작업 로직 작성
    return Result{ID: job.ID, Message: "Success"}, nil
}
```

### 2. 프로듀서 구현
`workerpool.JobProducer` 인터페이스 구현.
```go
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
```

### 3. 워커 풀 실행
```go
func main() {
    handler := &Handler{}
    producer := &Producer{max: 10}
    
    // 워커 5개로 풀 생성
    pool := workerpool.NewWorkerPool[Job, Result](5, handler)
    
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
```
