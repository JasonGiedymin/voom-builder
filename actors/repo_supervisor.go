package actors

import (
    // "github.com/JasonGiedymin/voom-builder/clients"
    "github.com/JasonGiedymin/voom-builder/common"
    // "github.com/JasonGiedymin/voom-builder/config"
    "github.com/JasonGiedymin/voom-builder/stats"

    "encoding/json"
    "fmt"
    "log"
    "os"
    "sync"

    "github.com/nu7hatch/gouuid"
    "github.com/thejerf/suture"
)

const (
    UuidFormat           = "Supervisor - %s - amuxbit.com" // used to generate the uuid
    SupervisorNameFormat = "(%s,%s)"                       // used as readable format, hostname and full uuid pair which are unique
    SupervisorFormat     = "%s (success: %d, fail: %d, stats: %v)"
)

type ServiceTag struct {
    Hostname string
    Uuid     string
    Workers  int
}

func (s *ServiceTag) Json() string {
    data, err := json.Marshal(s)

    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Json data: %s", string(data))

    return string(data)
}

func (s ServiceTag) String() string {
    return fmt.Sprintf(SupervisorNameFormat, s.Hostname, s.Uuid)
}

type SupervisorSpec struct {
    BufferLimit int
}

// == RepoSupervisor ==
type RepoSupervisor struct {
    AbstractSupervisor
    JobsPending chan string
    JobsDone    chan string
    fired       chan bool
}

func NewSupervisor(supervisorSpec SupervisorSpec) *RepoSupervisor {
    var supervisor RepoSupervisor
    supervisorUUID, err := uuid.NewV4()
    if err != nil {
        log.Fatalf("Error: Could not generate uuid! - %v", err)
        // fatal will exit here
    }

    name := fmt.Sprintf("Supervisor [repo] (%s)", supervisorUUID)
    log.Printf("Created: %s", name)

    spec := suture.Spec{Log: supervisor.LogServiceFailure}
    superRef := suture.New(supervisorUUID.String(), spec)

    // create a new supervisor, with a null suture reference
    supervisor = RepoSupervisor{
        AbstractSupervisor{
            name,
            superRef,
            map[string]Worker{},
            map[string]Supervisor{},
            supervisorUUID.String(),
            stats.NewSupervisorStats(),
        },
        make(chan string, supervisorSpec.BufferLimit),
        make(chan string, supervisorSpec.BufferLimit),
        make(chan bool),
    }

    return &supervisor
}

func (s *RepoSupervisor) Name() string {
    return "Supervisor (Repo)"
}

func (s *RepoSupervisor) ServiceTag() ServiceTag {

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

// Upon failure, this method will be called
// func (s *RepoSupervisor) LogServiceFailure(failure string) {
//     s.stats.Withdraw(1, &common.WorkError{failure})
//     log.Println(failure)
// }

func (s *RepoSupervisor) Log(doing string, v ...interface{}) {
    entityName := fmt.Sprintf(
        SupervisorFormat,
        s.Name(),
        s.baseStats.SuccessCount(),
        s.baseStats.Errors(),
        s.baseStats.Snapshot().RateMean(),
    )
    log.Printf(common.LogEntryf(entityName, doing, v...))
}

func (s *RepoSupervisor) LogFatal(doing string, v ...interface{}) {
    entityName := fmt.Sprintf(SupervisorFormat, s.ref.Name)
    log.Fatalf(common.LogEntryf(entityName, doing, v...))
}

func (s *RepoSupervisor) Serve() {
    var wg sync.WaitGroup
    wg.Add(1)

    fmt.Println("**** (RS) Serving")
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

func (s *RepoSupervisor) Run() {
    go s.ref.ServeBackground()
    // s.ref.Serve()
    go s.Serve()
}

func (s *RepoSupervisor) Work(work string) {
    // like akka, need a router
    s.Log("Work added to queue: %s", work)
    s.JobsPending <- work
    // s.baseStats.Reserve(1)
}

// func (s *RepoSupervisor) AddSupervisor(supervisor Supervisor) {
//     s.ref.Add(supervisor)

//     if _, ok := s.supervisors[supervisor.Name()]; ok {
//         s.LogFatal("Supervisor '%s' cannot be added twice for supervision!", supervisor.Name)
//     } else {
//         s.supervisors[supervisor.Name()] = supervisor
//     }
// }
