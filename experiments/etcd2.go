package main

import (
    "log"

    "./config"
    "./worker"
)

const (
    file = "config.yml"
)

func main() {
    configFile, err := config.ReadConfig(file)
    if err != nil {
        log.Fatalf("Error reading config file")
        return
    }

    client := config.NewEtcdClient(configFile)
    worker := worker.NewBuildWorker(client)

    if err := worker.EtcdConfig(); err != nil {
        log.Fatal(err)
    } else {
        log.Println("all good")
    }

    log.Println(client.GetCluster())

    client.CreateDir("JobsReady", 0)
    client.CreateDir("JobsInProgress", 0)
    client.CreateDir("JobsDone", 0)
    client.Set("JobsReady/3", "{id:3, repo:'http://github.com/example'}", 0)
    client.Set("JobsReady/1", "{id:1, repo:'http://github.com/example'}", 0)
    client.Set("JobsReady/2", "{id:2, repo:'http://github.com/example'}", 0)

    listAll := func() {
        if r, err := client.Get("JobsReady", true, true); err == nil {
            for _, n := range r.Node.Nodes {
                log.Printf("k:%s, v:%s, i:%d", n.Key, n.Value, n.CreatedIndex)
            }
        }
    }

    firstJob := func() uint64 {
        if r, err := client.Get("JobsReady", true, true); err == nil {
            for _, n := range r.Node.Nodes {
                if n.Key == "/JobsReady/1" {
                    log.Printf("found it, index: %d", n.CreatedIndex)
                    return n.CreatedIndex
                }
                log.Printf("k:%s, v:%s, i:%d", n.Key, n.Value, n.CreatedIndex)
            }
        }

        return 0
    }

    listAll()

    client.CompareAndDelete("JobsReady/1", "{id:1, repo:'http://github.com/example'}", firstJob())

    listAll()
}
