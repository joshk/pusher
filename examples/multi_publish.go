package main

import (
    "fmt"
    "github.com/timonv/pusher"
    "time"
)

func main() {
    workers := 100
    messages := make(chan string)
    done := make(chan bool)

    client := pusher.NewClient(appId, key, secret, false)

    for i := 0; i < workers; i++ {
        go func() {
            for data := range messages {
                err := client.Publish(data, "test", "test")
                if err != nil {
                    fmt.Printf("E %s\n", err)
                } else {
                    fmt.Println(".")
                }
            }
        }()
    }

    go func() {
        for i := 0; i < 5000; i++ {
            messages <- "test"
        }
        done <- true
    }()

    select {
    case <-done:
        fmt.Println("Done :-)")
    case <-time.After(1 * time.Minute):
        fmt.Println("Timeout :-(")
    }

    fmt.Println("")
}
