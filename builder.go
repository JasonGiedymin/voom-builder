// One criticism I have is not having enough time to make a lot more
// of this code abstract and easier to use.
// Also there seems to be a bug with the way supervisor trees are created
// and share channels. There is almost a high possiblity of blocking in
// the top level supervisor, making it impossible for it to manage
// other supervisors. This impl is probably flawed for the long term and
// as of this moment used only for service restarts.
//
// Recommend that it be re-written in possibly Rust.
//
package main

import (
    "fmt"
    "runtime"
    "sync"
    //"log"
    "time"

    "github.com/JasonGiedymin/voom-builder/actors"
    "github.com/JasonGiedymin/voom-builder/clients"
    "github.com/JasonGiedymin/voom-builder/config"
)

const (
    CONFIG_FILE = "config.yml"
    VERSION     = "v0.1.2"
)

func banner() {
    fmt.Printf("\nVoom Builder %s\n\n", VERSION)
}

func main() {
    banner()

    appConfig, _ := config.ReadConfig(CONFIG_FILE)
    workers := appConfig.Workers

    runtime.GOMAXPROCS(appConfig.Workers)

    supervisorSpec := actors.SupervisorSpec{
        BufferLimit: workers * 5,
    }
    supervisor := actors.NewSupervisor(supervisorSpec)
    supportsupervisor := actors.NewSupportSupervisor(supervisorSpec)

    setupWorkers := func() actors.ServiceTag {
        for i := 0; i < workers; i++ {
            spec := actors.WorkerSpec{
                fmt.Sprintf("worker-%d", i),
                supervisor.JobsPending,
                supervisor.JobsDone,
                time.Millisecond * 5000,
            }

            // add general workers to supervisor
            supervisor.Add(actors.NewWorker(spec))
        }

        return supervisor.ServiceTag()
    }

    setupSupportWorkers := func(serviceTag actors.ServiceTag) {

        // Supervisor setup is complete, register it with etcd via registration worker
        key := fmt.Sprintf(
            "%s/%s/%s",
            appConfig.Etcd.Paths.Supervisors,
            serviceTag.Hostname,
            serviceTag.Uuid,
        )

        regSpec := actors.RegistrationWorkerSpec{
            "RegistrationWorker-1",
            &appConfig.Etcd,
            clients.EtcdPair{
                key,
                serviceTag.Json(),
            },
            appConfig.Etcd.RegistrationInterval,
            supportsupervisor.JobsDone,
        }

        // add registration worker to supervisor
        supportsupervisor.Add(actors.NewRegistrationWorker(regSpec))
    }

    setupSupportWorkers(setupWorkers()) // arity to show explicit relationship (serviceTag)
    // supervisor.AddSupervisor(supportsupervisor)

    // every two seconds, send work to supervisor
    // this can be done in another worker
    count := 0
    go func() {
        for {
            time.Sleep(time.Millisecond * 100)
            count += 1
            if count <= appConfig.WorkLimit {
                supervisor.Work(fmt.Sprintf("Piling %d rocks", count))
            } else {
                return
            }
        }
    }()

    var wg sync.WaitGroup
    wg.Add(1)

    supervisor.Run()
    supportsupervisor.Run()

    wg.Wait()
}
