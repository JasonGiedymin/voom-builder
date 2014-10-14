// One criticism I have is not having enough time to make a lot more
// of this code abstract and easier to use. Not nearly enough interfaces.
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
)

func main() {
    // var supervisor actors.RepoSupervisor
    // var supportsupervisor actors.SupportSupervisor

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
            "test", //serviceTag.Hostname,
            "test", //serviceTag.Uuid,
        )

        regSpec := actors.RegistrationWorkerSpec{
            "RegistrationWorker-1",
            &appConfig.Etcd,
            clients.EtcdPair{
                key,
                "test", //serviceTag.Json(), //val (will turn into pair)
            },
            appConfig.Etcd.RegistrationInterval,
            supportsupervisor.JobsDone,
        }

        // add registration worker to supervisor
        supportsupervisor.Add(actors.NewRegistrationWorker(regSpec))
    }

    setupSupportWorkers(setupWorkers()) // arity to show explicit relationship (serviceTag)
    //supervisor.AddSupervisor(supportsupervisor)

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
