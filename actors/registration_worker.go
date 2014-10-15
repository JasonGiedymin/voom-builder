package actors

import (
    "fmt"
    "log"
    "time"

    "github.com/JasonGiedymin/voom-builder/clients"
    "github.com/JasonGiedymin/voom-builder/common"
    "github.com/JasonGiedymin/voom-builder/config"
)

type RegistrationWorkerSpec struct {
    Name       string
    EtcdConfig *config.EtcdConfig
    Data       clients.EtcdPair // know about the key path up front
    Interval   int
    Done       chan string
}

type RegistrationWorker struct {
    name string
    // later can have a channel for data, or array, now know it up front
    // path     string // know about the key path up front
    // value    string
    data       clients.EtcdPair
    interval   int
    jobsDone   chan string
    etcdClient *clients.Etcd
}

func NewRegistrationWorker(spec RegistrationWorkerSpec) *RegistrationWorker {
    etcdClient, err := clients.NewEtcdClient(spec.EtcdConfig)
    if err != nil {
        log.Fatal(err)
    }

    etcdClient.Config()

    return &RegistrationWorker{
        spec.Name,
        spec.Data,
        spec.Interval,
        spec.Done,
        etcdClient,
    }
}

func (w RegistrationWorker) Log(msg string) {
    entityName := fmt.Sprintf("RegistrationWorker - %s", w.Name())
    log.Printf(common.LogEntry(entityName, msg))
}

func (w RegistrationWorker) Logf(msg string, v ...interface{}) {
    entityName := fmt.Sprintf("RegistrationWorker - %s", w.Name())
    log.Printf(common.LogEntryf(entityName, msg, v))
}

func (w *RegistrationWorker) Work() {
    if _, err := w.etcdClient.Set(w.data); err != nil {
        w.Logf("Error occured: %s", err.Error())
    } else {
        w.jobsDone <- "registered service with etcd"
    }
}

func (w *RegistrationWorker) Serve() {
    w.Work() // immediately work once

    func() {
        for {
            time.Sleep(time.Second * time.Duration(w.interval))
            w.Work()
        }
    }()
}

func (w *RegistrationWorker) Stop() {
    w.Log("Stopping!")
}

func (w *RegistrationWorker) Name() string {
    return w.name
}
