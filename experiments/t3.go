package main

import (
    "log"
    "sync"
    "time"
)

// Needed to use waitgroup or else goroutine would die prematurely
func main() {
    var wg sync.WaitGroup
    wg.Add(3)
    intChan := make(chan int, 10)
    intChan2 := make(chan int, 10)

    go func() { // goes on forever
        for { // this for loop is necessary so as to keep reading from the channels
            select { // is none blocking until sleep in the default
            case w, _ := <-intChan:
                log.Printf("intChan saw %d", w)
                intChan2 <- (w + 100)
            case w, _ := <-intChan2:
                log.Printf("intChan2 saw %d", w)
                wg.Done()
            default: // sleep for a second before checking again
                log.Println("Found nothing, sleeping")
                time.Sleep(1 * time.Second)
            }
        }
    }()

    go func() {
        log.Println("Sleeping before populate..")
        time.Sleep(2 * time.Second)
        log.Println("Populating!")
        intChan <- 1
        intChan <- 2
        intChan <- 3
    }()

    wg.Wait()
    log.Println("Done.")
}
