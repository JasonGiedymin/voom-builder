package experiments

import (
    "log"
    // "net/http"
    "fmt"
    "runtime"
    "time"
)

const (
    workers        = 2                // number of Poller goroutines to launch
    pullWait       = 1 * time.Second  // minimum pull relief wait timing
    poisonPillTime = 10 * time.Second // poison pill in 10 seconds
    // statusInterval = 10 * time.Second // how often to log status to stdout
    // errTimeout = 10 * time.Second // back-off timeout on error
    buildqueue = "http://buildqueue.amuxbit.com"
)

type GitRepo struct {
    protocol string
    location string
    fqString string
}

type Work struct {
    repo     GitRepo
    output   string
    state    string
    canceled bool
}

func (w *Work) Build(id int) error {
    if !w.canceled {
        mockBuildTime := 5 * time.Second
        log.Println(" -> building...")
        time.Sleep(mockBuildTime)
        w.state = fmt.Sprintf("OK - done by worker %d", id)
    } else {
        log.Println(" -> Work poisoned, stopping.")
        w.state = fmt.Sprintf("CANCELED - by worker %d", id)
    }

    return nil
}

type Queue struct {
    location string
    canceled bool
}

// Mocked now
// When `Pull()` is called, it will `Sleep` for `pullWait` and then
// send `Work` to the `pending` channel.
func (q *Queue) Pull(pending chan<- *Work) {
    if !q.canceled {
        time.Sleep(pullWait)
        sampleWork := &Work{
            repo: GitRepo{fqString: "git://test.repo.git"},
        }
        workToDo := workers + 1
        for i := 0; i < workToDo; i++ {
            pending <- sampleWork
        }
        log.Println("*Populated queue")
    } else {
        log.Println("*Queue poisoned!")
        return
    }
}

func (q *Queue) Poison(pending chan<- *Work) {
    time.Sleep(poisonPillTime)
    sampleWork := &Work{
        canceled: true,
    }
    workToDo := workers + 1
    for i := 0; i < workToDo; i++ {
        pending <- sampleWork
    }
    q.canceled = true
    log.Println("*Poisoned queue!")
}

func NewWorker(id int, inQueue chan *Work, outQueue chan<- *Work) {
    for w, ok := range inQueue { // grap from the inQueue
        if ok { // channel not closed
            log.Printf("Worker %d ready, Queue length: %d", id, len(inQueue))
            if !w.canceled {
                _ = w.Build(id) // do work!
                outQueue <- w   // send work to the outQueue when done
            } else {
                log.Printf("Worker %d done (found poison). Queue length was: %d", id, len(inQueue))
                if len(inQueue) == 0 {
                    close(outQueue)
                    close(inQueue)
                    return
                }
            }
        }
    }
}

func main() {
    runtime.GOMAXPROCS(workers)
    // Create our input and output channels.
    workQueue := Queue{location: buildqueue}
    pending, complete, done := make(chan *Work, 5), make(chan *Work), make(chan struct{})
    defer close(done)
    //status := StateMonitor(statusInterval)

    go workQueue.Pull(pending)   // pull once to populate queue
    go workQueue.Poison(pending) // method will wait poison time, then poison the work queue

    for i := 0; i < workers; i++ { // create workers
        go NewWorker(i, pending, complete) // create and immediately start building
    }

    for w := range complete { // when worker is complete print the state
        log.Println(w.state)
        if !w.canceled {
            go workQueue.Pull(pending) // pull from queue after each worker is done
        } else {
            log.Println("Found poisoned!")
            return
        }
    }
}
