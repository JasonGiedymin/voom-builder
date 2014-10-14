package actors

import (
    "fmt"
    "log"
    "os"

    "github.com/JasonGiedymin/voom-builder/common"
    "github.com/JasonGiedymin/voom-builder/stats"

    "github.com/thejerf/suture"
)

const (
    AbstractSupervisorFormat = "Supervisor - %s"
)

type Supervisor interface {
    Name() string
    ServiceTag() ServiceTag
    LogServiceFailure(failure string)
    Log(doing string, v ...interface{})
    LogFatal(doing string, v ...interface{})
    Serve() // useful when you need supervisor to do light work, call this in 'Run()''
    Add(worker Worker)
    AddSupervisor(supervisor Supervisor)
    Run()
    Stop()
    Work(work string)
}

type AbstractSupervisor struct {
    name           string
    ref            *suture.Supervisor
    workers        map[string]Worker
    supervisors    map[string]Supervisor
    supervisorUUID string
    baseStats      *stats.SupervisorStats
}

func (s *AbstractSupervisor) Name() string {
    if s.name == "" {
        return "AbstractSupervisor"
    } else {
        return s.name
    }
}

func (s *AbstractSupervisor) Log(doing string, v ...interface{}) {
    entityName := fmt.Sprintf(
        AbstractSupervisorFormat,
        s.Name(),
    )
    log.Printf(common.LogEntryf(entityName, doing, v...))
}

func (s *AbstractSupervisor) LogFatal(doing string, v ...interface{}) {
    entityName := fmt.Sprintf(
        AbstractSupervisorFormat,
        s.Name(),
    )
    log.Fatalf(common.LogEntryf(entityName, doing, v...))
}

// Upon failure, this method will be called
func (s *AbstractSupervisor) LogServiceFailure(failure string) {
    s.baseStats.Failure(1, &common.WorkError{failure})
    log.Println(failure)
}

func (s *AbstractSupervisor) Add(worker Worker) {
    s.ref.Add(worker)

    if _, ok := s.workers[worker.Name()]; ok {
        s.LogFatal("Worker '%s' cannot be added twice for supervision!", worker.Name())
    } else {
        s.workers[worker.Name()] = worker
    }
}

func (s *AbstractSupervisor) AddSupervisor(supervisor Supervisor) {
    s.ref.Add(supervisor)

    if _, ok := s.supervisors[supervisor.Name()]; ok {
        s.LogFatal("Supervisor '%s' cannot be added twice for supervision!", supervisor.Name)
    } else {
        s.supervisors[supervisor.Name()] = supervisor
    }
}

func (s *AbstractSupervisor) ServiceTag() ServiceTag {

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

func (s *AbstractSupervisor) Stop() {
    s.Log("Stopping")
}

func (s *AbstractSupervisor) Serve() {
}

func (s *AbstractSupervisor) Run() {
    s.ref.ServeBackground()
    s.Serve()
}
