package pusher

import (
  "testing"
  "fmt"
  "time"
  //"github.com/timonv/pusher"
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

