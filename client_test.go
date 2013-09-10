package pusher

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "net/url"
    "reflect"
    "testing"
)

func setupTestServer(handler http.Handler) (server *httptest.Server) {
    server = httptest.NewServer(handler)
    url, _ := url.Parse(server.URL)
    // FIXME: This is not thread-safe at all, and will not handle running tests in parallel.
    //   The endpoint field should probably be moved to the Client struct.
    Endpoint = url.Host
    return
}

func verifyRequest(t *testing.T, prefix string, req *http.Request, method, path string) (payload Payload) {
    if method != req.Method {
        t.Errorf("%s: Expected method %s, got %s", prefix, method, req.Method)
    }
    if path != req.URL.Path {
        t.Errorf("%s: Expected path '%s', got '%s'", prefix, path, req.URL.Path)
    }

    err := json.NewDecoder(req.Body).Decode(&payload)
    if err != nil {
        fmt.Println("Got error:", err)
    }

    return
}

func stringSlicesEqual(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}

func TestPublish(t *testing.T) {
    server := setupTestServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
        w.WriteHeader(200)
        fmt.Fprintf(w, "{}")

        payload := verifyRequest(t, "Publish()", request, "POST", "/apps/1/events")

        if payload.Name != "event" {
            t.Errorf("Publish(): Expected body[name] = \"event\", got %q", payload.Name)
        }
        if !reflect.DeepEqual(payload.Channels, []string{"mychannel", "c2"}) {
            t.Errorf("Publish(): Expected body[channels] = [mychannel c2], got %+v", payload.Channels)
        }
    }))
    defer server.Close()

    client := NewClient("1", "key", "secret", false)
    err := client.Publish("data", "event", "mychannel", "c2")

    if err != nil {
        t.Errorf("Publish(): %v", err)
    }
}
