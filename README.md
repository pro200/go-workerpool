# Go Worker Pool

이 프로젝트는 Go 언어의 고루틴(Goroutine)과 채널(Channel)을 활용하여 구현된 효율적인 비동기 작업 처리 라이브러리입니다. 인터페이스 기반의 설계를 통해 어떤 종류의 작업이든 유연하게 처리할 수 있으며, 워커의 수를 제한하여 시스템 리소스를 안정적으로 관리합니다.

## 주요 특징

- **동시성 제어**: 고정된 수의 워커를 사용하여 과도한 고루틴 생성을 방지합니다.
- **백프레셔(Backpressure) 지원**: 워커가 처리할 수 있는 속도에 맞춰 작업 공급을 자동으로 조절합니다.
- **인터페이스 기반**: `JobHandler`와 `JobProducer` 인터페이스를 통해 비즈니스 로직과 데이터 공급원을 자유롭게 정의할 수 있습니다.
- **안전한 종료**: 모든 작업이 처리된 후 채널을 닫고 리소스를 정리하는 Graceful Shutdown을 보장합니다.

## 시작하기

### 설치

```bash
go get github.com/pro200/go-workerpool
```

### 사용 방법
#### 1. Handler 구현하기
`JobHandler`인터페이스를 구현하여 각 작업이 어떻게 처리될지 정의합니다.
```go
type MyHandler struct{}

func (h *MyHandler) Process(job any) error {
    // 실제 처리 로직을 여기에 작성하세요
    fmt.Printf("Processing: %v\n", job)
    return nil
}
```

#### 2. Producer 구현하기
인터페이스를 구현하여 처리할 데이터를 공급하는 로직을 정의합니다. `JobProducer`
```go
type MyProducer struct {
    count int
}

func (p *MyProducer) Next() (any, error) {
    if p.count >= 10 {
        return nil, nil // 더 이상 작업이 없으면 nil 반환
    }
    p.count++
    return p.count, nil
}

func (p *MyProducer) Close() error {
    return nil
}
```

#### 3. WorkerPool 실행하기
워커 풀을 생성하고 실행하여 결과를 수신합니다.
```go
func main() {
    // 3개의 워커를 가진 풀 생성
    wp := workerpool.NewWorkerPool(3, &MyHandler{})

    // 풀 실행 및 결과 채널 수신
    results := wp.Run(&MyProducer{})

    // 결과 처리
    for res := range results {
        if res.Err != nil {
            fmt.Printf("Error processing job %v: %v\n", res.Job, res.Err)
            continue
        }
        fmt.Printf("Result for job %v is success\n", res.Job)
    }
}
```

## 구조체 및 인터페이스 설명
### `JobHandler`
작업의 비즈니스 로직을 담당합니다.
- `Process(job any) error`: 작업을 인자로 받아 실행하며, 에러 발생 시 이를 반환합니다. 

### `JobProducer`
작업 데이터의 공급을 담당합니다.
- `Next() (any, error)`: 다음 작업을 반환합니다. 첫 번째 반환값이 `nil`이면 모든 작업이 끝난 것으로 간주합니다. 
- `Close() error`: 데이터 공급이 끝난 후 리소스를 정리합니다. 

### `Result`
작업의 처리 결과를 담고 있습니다.
- `Job any`: 처리된 작업 객체 
- `Err error`: 발생한 에러 (성공 시 `nil`) 
