package main

import (
    "log"
    "sync"
    "time"
)

// Needed to use waitgroup or else goroutine would die prematurely
func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    phase1 := make(chan int, 10)
    phase2 := make(chan int, 10)
    done := make(chan bool, 1)

    go func() { // goes on forever
        for { // this for loop is necessary so as to keep reading from the channels
            select { // is none blocking until sleep in the default
            case w, ok := <-phase1:
                if ok {
                    log.Printf("phase1 saw %d", w)
                    if w != 0 {
                        phase2 <- (w + 100)
                    } else {
                        done <- true
                    }
                }
            case w, _ := <-phase2:
                log.Printf("phase2 saw %d", w)
            case _, _ = <-done:
                log.Printf("Saw done msg, closing up phase2")
                close(phase2)
                wg.Done() // releases the wait, allowing exit
                return
            default: // sleep for a second before checking again
                log.Println("Found nothing, sleeping")
                time.Sleep(1 * time.Second)
            }
        }
    }()

    go func() {
        log.Println("Sleeping before populate..")
        time.Sleep(5 * time.Second)
        log.Println("Populating!")
        sendDoneAt := 2
        messagesToSend := 20
        for i := 1; i <= messagesToSend; i++ {
            if i == sendDoneAt {
                phase1 <- 0
            } else {
                phase1 <- i
            }
        }

        // simulates a second round of messages
        time.Sleep(1 * time.Second)

        for i := 21; i <= messagesToSend*3; i++ {
            phase1 <- i
        }

        // time.Sleep(10 * time.Second)
        // phase1 <- 0

    }()

    wg.Wait()
    log.Println("Done.")
}
