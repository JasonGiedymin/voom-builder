package main

import (
    "github.com/coreos/go-etcd/etcd"

    "errors"
    "log"
)

const (
    ConfigPath = "/app/config"
)

var (
    ErrorNoEtcdClient      = errors.New("None or invalid etcd client given.")
    ErrorNoEtcdConnection  = errors.New("Could not contact etcd.")
    ErrorNoEtcdConfigFound = errors.New("Could not find config in etcd.")
)

type BuildWorker struct {
    etcdClient *etcd.Client
    etcdConfig string
}

func NewBuildWorker(client *etcd.Client) *BuildWorker {
    return &BuildWorker{etcdClient: client}
}

func (b BuildWorker) check() error {
    if b.etcdClient == nil {
        return ErrorNoEtcdClient
    }

    if cluster := b.etcdClient.GetCluster(); len(cluster) == 0 {
        return ErrorNoEtcdConnection
    }

    if _, err := b.etcdClient.CreateDir("ClientStatus", 1); err != nil {
        return ErrorNoEtcdConnection
    }

    return nil
}

// Queries etcd for runtime config such as the Queue endpoint
func (b BuildWorker) EtcdConfig() error {
    // standard check
    err := b.check()
    if err != nil {
        return err
    }

    if r, err := b.etcdClient.Get(ConfigPath, true, true); err != nil {
        // log.Println(ErrorNoEtcdConfigFound.Error() + " - " + err.Error())
        return err
    } else {
        // if len(r.Node.Nodes) == 0 {
        //     return ErrorNoEtcdConfigFound
        // }
        // for _, n := range r.Node.Nodes {
        //     log.Printf("k:%s, v:%s, i:%d", n.Key, n.Value, n.CreatedIndex)
        // }
        log.Printf("k:%s, v:%s, i:%d", r.Node.Key, r.Node.Value, r.Node.CreatedIndex)
    }

    return nil
}
