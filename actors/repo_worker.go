package actors

import (
    "fmt"
    "log"
    "time"

    "github.com/JasonGiedymin/voom-builder/common"
)

type WorkerSpec struct {
    Name      string
    Pending   chan string
    Done      chan string
    YieldTime time.Duration
}

type RepoWorker struct {
    name        string
    yieldTime   time.Duration
    jobStream   chan string
    jobsDone    chan string
    currentWork string
}

func (w RepoWorker) Log(doing string, v ...interface{}) {
    entityName := fmt.Sprintf("RepoWorker - %s", w.Name())
    log.Printf(common.LogEntryf(entityName, doing, v...))
}

func (w *RepoWorker) Work(doing string) {
    time.Sleep(w.yieldTime)
    w.Log("finished %s", doing)
    w.jobsDone <- doing
}

func (w *RepoWorker) Serve() {
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

func (w *RepoWorker) Send(work string) {
    w.jobStream <- work
}

func (w *RepoWorker) Stop() {
    w.Log("Stopping!")
}

func (w *RepoWorker) Name() string {
    return w.name
}

func NewWorker(spec WorkerSpec) *RepoWorker {
    return &RepoWorker{
        spec.Name,
        spec.YieldTime,
        spec.Pending,
        spec.Done,
        "_WORKER_CREATED", // default init job
    }
}
