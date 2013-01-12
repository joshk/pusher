package pusher

import (
  "testing"
  "fmt"
  "time"
)

// Integration level tests
func TestSimplePublish(t *testing.T) {
	client := NewClient("34420", "87bdfd3a6320e83b9289", "f25dfe88fb26ebf75139", false)

	done := make(chan bool)

	go func() {
		err := client.Publish("test", "test", "test")
		if err != nil {
			t.Errorf("Error %s\n", err)
    }
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("Done :-)")
	case <-time.After(1 * time.Minute):
		t.Errorf("Timeout :-(")
	}
}

func TestMultiPublish(t *testing.T) {
	workers := 10
	messageCount := 100
	messages := make(chan string)
	done := make(chan bool)

	client := NewClient("34420", "87bdfd3a6320e83b9289", "f25dfe88fb26ebf75139", false)

	for i := 0; i < workers; i++ {
		go func() {
			for data := range messages {
				err := client.Publish(data, "test", "test")
				if err != nil {
          t.Errorf("Error: ", err)
				} else {
					fmt.Print(".")
				}
			}
		}()
	}

	go func() {
		for i := 0; i < messageCount; i++ {
			messages <- "test"
		}
		done <- true
		close(messages)
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("\nDone :-)")
	case <-time.After(1 * time.Minute):
		fmt.Println("\nTimeout :-(")
	}

	fmt.Println("")
}
