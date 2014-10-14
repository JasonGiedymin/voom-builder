package actors

import (
    "github.com/JasonGiedymin/voom-builder/stats"

    "fmt"
    "log"
    "os"
    "sync"
    // "time"

    "github.com/nu7hatch/gouuid"
    "github.com/thejerf/suture"
)

const ()

type SupportSupervisor struct {
    AbstractSupervisor
    JobsDone chan string
}

func (s *SupportSupervisor) Name() string {
    return "Supervisor (Support)"
}

func (s *SupportSupervisor) Serve() {
    var wg sync.WaitGroup
    wg.Add(1)

    fmt.Println("**** (SS) Serving")

    func() { // don't put in goroutine or else it will exit
        for {
            select {
            case doing := <-s.JobsDone:
                s.baseStats.Success(1)
                s.Log("worker is complete with %s", doing)
            }
        }
    }()

    wg.Wait()
}

func (s *SupportSupervisor) Run() {
    fmt.Println("Running")
    go s.ref.ServeBackground()
    // s.ref.Serve()
    go s.Serve()
}

func (s *SupportSupervisor) Work(work string) {
}

func NewSupportSupervisor(supervisorSpec SupervisorSpec) *SupportSupervisor {
    var supervisor SupportSupervisor
    supervisorUUID, err := uuid.NewV4()
    if err != nil {
        log.Fatalf("Error: Could not generate uuid! - %v", err)
        // fatal will exit here
    }

    name := fmt.Sprintf("Supervisor [support] (%s)", supervisorUUID)
    log.Printf("Created: %s", name)

    spec := suture.Spec{Log: supervisor.LogServiceFailure}
    superRef := suture.New(supervisorUUID.String(), spec)

    // create a new supervisor, with a null suture reference
    supervisor = SupportSupervisor{
        AbstractSupervisor{
            name,
            superRef,
            map[string]Worker{},
            map[string]Supervisor{},
            supervisorUUID.String(),
            stats.NewSupervisorStats(),
        },
        make(chan string, supervisorSpec.BufferLimit),
    }

    // log.Println("******", supervisor.ref)

    return &supervisor
}

func (s *SupportSupervisor) ServiceTag() ServiceTag {

    hostname, err := os.Hostname()
    if err != nil {
        log.Fatalf("Error: Could not get hostname! - %v", err)
        // fatal will exit here
    }

    return ServiceTag{
        hostname,
        s.supervisorUUID,
        len(s.workers),
    }
}

// func (s *SupportSupervisor) Add(worker Worker) {
//     s.ref.Add(worker)

//     if _, ok := s.workers[worker.Name()]; ok {
//         s.LogFatal("Worker '%s' cannot be added twice for supervision!", worker.Name())
//     } else {
//         s.workers[worker.Name()] = worker
//     }
// }

// func (s *SupportSupervisor) AddSupervisor(supervisor Supervisor) {
//     s.ref.Add(supervisor)

//     if _, ok := s.supervisors[supervisor.Name()]; ok {
//         s.LogFatal("Supervisor '%s' cannot be added twice for supervision!", supervisor.Name)
//     } else {
//         s.supervisors[supervisor.Name()] = supervisor
//     }
// }
