package main

import (
    "log"
    "sync"
)

// Needed to use waitgroup or else goroutine would die prematurely
func main() {
    var wg sync.WaitGroup
    wg.Add(3)
    intChan := make(chan int)

    go func() {
        for i := range intChan {
            log.Println(i)
            wg.Done()
        }
    }()

    intChan <- 1
    intChan <- 2
    intChan <- 3

    wg.Wait()

    log.Println("Done.")
}
