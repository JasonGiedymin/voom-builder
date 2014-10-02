// Supervisor
// Uses `suture`, a golang supervisor tree impl.
//
// - [x] create supervisor
// - [x] add services that impl the service interface
// - [x] super.serve()
// - [x] add comm channel that services report back on

package main

import (
    "log"
    // "time"

    "github.com/thejerf/suture"
)

type Super struct {
    ref          *suture.Supervisor
    QueuePending chan string
    QueueDone    chan string
}

func (s *Super) New(name string) {
    if s.ref == nil {
        s.ref = suture.NewSimple("Supervisor")
        s.QueuePending = make(chan string)
        s.QueueDone = make(chan string)
    } else {
        log.Fatal("Tried to construct a new supervisor over an existing one. For safety reasons, this command is not allowed.")
    }
}

func (s *Super) Add(service suture.Service) {
    s.ref.Add(service)
}

func (s *Super) Run() {
    go s.ref.ServeBackground()
}

func (s *Super) Stop() {
    s.ref.Stop()
}

func (s *Super) Listen() {
    go func() {
        // for jobDone := range s.QueueDone {
        //     log.Println("** Job was done: ", jobDone)
        // }
        for {
            select {
            case jobDone := <-s.QueueDone:
                log.Println("** Job was done: ", jobDone)
            }
        }
    }()
}

func (s *Super) SendWork() {
    s.QueuePending <- "Fake Work"
}

type Worker struct {
    // counter      int
    queuePending chan string
    queueDone    chan string
    //next         chan string
    stop chan bool
}

func NewWorker(queuePending chan string, queueDone chan string) *Worker {
    return &Worker{queuePending, queueDone, make(chan bool)}
}

func (i *Worker) Work() string {
    log.Println("...working...")
    // time.Sleep(time.Millisecond * 500)
    return "some work!"
}

func (i *Worker) Stop() {
    log.Println("Stopping the service")
    i.stop <- true
}

func (i *Worker) Serve() {
    for {
        select {
        case <-i.queuePending:
            log.Println("Pending work found...")
            i.queueDone <- i.Work()
            //log.Println("pending queue found: ", msg)
            //log.Println("Current: ", i.counter)
            //i.counter = 0
            // case i.queueDone <- i.Work():
            // log.Println("Doing work...")
            // log.Println("Current: ", i.counter)
            // i.counter += 1
            // if i.counter == 10 {
            //     log.Println("counter: ", i.counter)
            //     close(i.next)
            //     i.Stop()
            //     return
            // }
        }
    }
}

func main() {
    log.Println("Started")

    var super Super
    super.New("Supervisor2")
    service := NewWorker(super.QueuePending, super.QueueDone)

    super.Add(service)
    go super.Listen()
    super.Run()
    super.SendWork()
    super.SendWork()

    // func() {
    //     time.Sleep(time.Millisecond * 1000)
    //     super.Stop() // this stops the world
    // }()

    // We sync here just to guarantee the output of "Stopping the service"
    //<-service.stop
}
