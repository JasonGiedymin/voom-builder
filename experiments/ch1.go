package main

func main() {
     // Data is irrelevant
     done := make(chan struct{})

     go func() {
          println("goroutine message")

          // Just send a signal "I'm done"
          close(done)
     }()

     println("main function message")
     <-done
}
