package main

import (
    "fmt"
    "github.com/timonv/pusher"
    "time"
)

func main() {
    client := pusher.NewClient(appId, key, secret, false)

    done := make(chan bool)

    go func() {
        err := client.Publish("test", "test", "test")
        if err != nil {
            fmt.Printf("E %s\n", err)
        } else {
            fmt.Print(".")
        }
        done <- true
    }()

    select {
    case <-done:
        fmt.Println("\ndone")
    case <-time.After(1 * time.Minute):
        fmt.Println("timeout")
    }

    fmt.Println("")
}
