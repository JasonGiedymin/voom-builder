package main

// - [x] supervisor
// - [x] worker
// - [x] suture
// - [x] multi-cpu
// - [x] buffer increase
// - [x] optmize buffer so work is pulled only when needed?
// - [x] cleanup
// - [x] spawn failure test
// - [x] record failures
// - [x] record current state when accepting work
// - [x] introduce concept of job 'claim'
//   - [x] keep track of current claims with channels (as recommended by golang)
//   - [x] keep track of failures and update claims with info
// - [x] use metrics (though this impl really should be a raw channel select
//       or crdt positive-negative counter; atomic will still work yet not
//       as elegant and integration with influxdb is provided)
// - [ ] tests
// - [ ] fleet test deploy
// - [ ] historical data on complete work via influxdb

import (
    "fmt"
    "log"
    "runtime"
    // "sync"
    "time"

    "github.com/rcrowley/go-metrics"
    "github.com/thejerf/suture"
)

const (
    EntitySuper      = "Supervisor"
    EntityWorker     = "Worker"
    workers          = 20
    bufferLimit      = workers * 5
    WorkLimit        = 50
    SupervisorFormat = "Supervisor - %s (claims: %d, stats: %v)"
)

// == Common ==
func LogEntryf(entity string, format string, v ...interface{}) string {
    logFormat := fmt.Sprintf("[%s] - %s", entity, format)
    return fmt.Sprintf(logFormat, v...)
}

// object to help convert this error we get from the supervisor
type WorkError struct {
    FullError string
}

// == Stats ==
type SuperStats struct {
    Claims          metrics.Counter // count of succesful completions
    ClaimCompletion metrics.Meter   // time between completions

    ClaimsErrors          metrics.Counter // count the errors
    ClaimCompletionErrors metrics.Meter   // measure the time between failures
}

func (s *SuperStats) Reserve(shares int64) {
    s.Claims.Inc(shares)
}

func (s *SuperStats) Withdraw(shares int64, err *WorkError) {
    if err == nil {
        s.ClaimCompletion.Mark(1)
    } else {
        s.ClaimsErrors.Inc(1)
        s.ClaimCompletionErrors.Mark(1)
    }
    s.Claims.Dec(shares)
}

// == Supervisor ==
type Supervisor struct {
    jobsPending chan string
    jobsDone    chan string
    fired       chan bool

    ref     *suture.Supervisor
    workers map[string]Worker
    stats   *SuperStats
}

func NewSuperStats() *SuperStats {
    return &SuperStats{
        metrics.NewCounter(),
        metrics.NewMeter(),
        metrics.NewCounter(),
        metrics.NewMeter(),
    }
}

func NewSupervisor() *Supervisor {
    // create a new supervisor, with a null suture reference
    supervisor := &Supervisor{
        make(chan string, bufferLimit),
        make(chan string, bufferLimit),
        make(chan bool),
        nil,
        map[string]Worker{},
        NewSuperStats(),
    }

    // create a new suture reference, using the supervisor method
    // in the spec (requires that I first create a supervisor)
    spec := suture.Spec{Log: supervisor.LogServiceFailure}
    ref := suture.New("Supervisor", spec)

    // assign reference, allowing supervisor to handle failures
    supervisor.ref = ref

    return supervisor
}

func (s *Supervisor) LogServiceFailure(failure string) {
    s.stats.Withdraw(1, &WorkError{failure})
    log.Println(failure)
}

func (s *Supervisor) Log(doing string, v ...interface{}) {
    entityName := fmt.Sprintf(
        SupervisorFormat,
        s.ref.Name,
        s.stats.Claims.Count(),
        s.stats.ClaimCompletion.RateMean(),
    )
    log.Printf(LogEntryf(entityName, doing, v...))
}

func (s *Supervisor) LogFatal(doing string, v ...interface{}) {
    entityName := fmt.Sprintf(SupervisorFormat, s.ref.Name)
    log.Fatalf(LogEntryf(entityName, doing, v...))
}

func (s *Supervisor) Serve() {
    // var wg sync.WaitGroup
    // wg.Add(1)

    func() { // don't put in goroutine or else it will exit
        for {
            select {
            case doing := <-s.jobsDone:
                s.stats.Withdraw(1, nil)
                s.Log("worker is complete with %s", doing)
            }
        }
    }()

    // wg.Wait()
}

func (s *Supervisor) Add(worker Worker) {
    s.ref.Add(&worker)

    if _, ok := s.workers[worker.Name]; ok {
        s.LogFatal("Worker '%s' cannot be added twice for supervision!", worker.Name)
    } else {
        s.workers[worker.Name] = worker
    }
}

func (s *Supervisor) Run() {
    s.ref.ServeBackground()
    s.Serve()
}

func (s *Supervisor) Work(work string) {
    // like akka, need a router
    s.Log("Work added to queue: %s", work)
    s.jobsPending <- work
    s.stats.Reserve(1)
}

// == Worker ==

type Worker struct {
    Name        string
    workTime    int
    jobStream   chan string
    jobsDone    chan string
    currentWork string
}

func (w Worker) Log(doing string, v ...interface{}) {
    entityName := fmt.Sprintf("Worker - %s", w.Name)
    log.Printf(LogEntryf(entityName, doing, v...))
}

func (w *Worker) Work(doing string) {
    time.Sleep(time.Millisecond * 5000)
    w.Log("finished %s", doing)
    w.jobsDone <- doing
}

func (w *Worker) Serve() {
    func() {
        for {
            select {
            case doing := <-w.jobStream:
                w.currentWork = doing
                if doing == "Piling 3 rocks" {
                    // log.Fatal("Force Error!") // force failure!
                    return
                }
                w.Log("Got job, looks like I will be %s", doing)
                w.Work(doing)
            }
        }
    }()
}

func (w *Worker) Send(work string) {
    w.jobStream <- work
}

// func (w *Worker) Serve() {
//     w.Listen()
// }

func (w *Worker) Stop() {
    w.Log("Stopping!")
}

func main() {
    runtime.GOMAXPROCS(workers)

    supervisor := NewSupervisor()

    for i := 0; i < workers; i++ {
        workerName := fmt.Sprintf("worker-%d", i)
        worker := &Worker{
            workerName,
            500,
            supervisor.jobsPending,
            supervisor.jobsDone,
            "_INIT",
        }
        supervisor.Add(*worker)
    }

    // every two seconds, send work to supervisor
    count := 0
    go func() {
        for {
            time.Sleep(time.Millisecond * 100)
            count += 1
            if count < WorkLimit {
                supervisor.Work(fmt.Sprintf("Piling %d rocks", count))
            }
        }
    }()

    supervisor.Run()
}
