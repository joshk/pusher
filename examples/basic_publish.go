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
            fmt.Printf("Error %s\n", err)
        } else {
            fmt.Println("Message Published!")
        }
        done <- true
    }()

    select {
    case <-done:
        fmt.Println("Done :-)")
    case <-time.After(1 * time.Minute):
        fmt.Println("Timeout :-(")
    }
}
