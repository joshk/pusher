package pusher

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
    "time"
)

var Endpoint = "api.pusherapp.com"

const AuthVersion = "1.0"

type Client struct {
    appid, key, secret string
    secure             bool
}

type Payload struct {
    Name     string   `json:"name"`
    Channels []string `json:"channels"`
    Data     string   `json:"data"`
}

func NewClient(appid, key, secret string, secure bool) *Client {
    return &Client{appid, key, secret, secure}
}

func (c *Client) Publish(data, event string, channels ...string) error {
    timestamp := c.stringTimestamp()

    content, err := c.jsonifyData(data, event, channels)
    if err != nil {
        return err
    }

    signature := Signature{c.key, c.secret, "POST", c.publishPath(), timestamp, AuthVersion, content}

    err = c.post(content, c.publishUrl(), signature.EncodedQuery())

    return err
}

func (c *Client) jsonifyData(data, event string, channels []string) (string, error) {
    content := Payload{event, channels, data}
    b, err := json.Marshal(content)
    if err != nil {
        return "", err
    }
    return string(b), nil
}

func (c *Client) post(content string, fullUrl string, query string) error {
    buffer := bytes.NewBuffer([]byte(content))

    postUrl, err := url.Parse(fullUrl)
    if err != nil {
        return err
    }

    postUrl.Scheme = c.scheme()
    postUrl.RawQuery = query

    resp, err := http.Post(postUrl.String(), "application/json", buffer)
    if err != nil {
        return fmt.Errorf("pusher: POST failed: %s", err)
    }

    defer resp.Body.Close()

    if resp.StatusCode == 401 {
        b, _ := ioutil.ReadAll(resp.Body)
        return fmt.Errorf("pusher: POST failed: %s", b)
    }

    return nil
}

func (c *Client) scheme() string {
    if c.secure {
        return "https"
    }
    return "http"
}

func (c *Client) publishPath() string {
    return "/apps/" + c.appid + "/events"
}

func (c *Client) publishUrl() string {
    return "http://" + Endpoint + c.publishPath()
}

func (c *Client) stringTimestamp() string {
    t := time.Now()
    return strconv.FormatInt(t.Unix(), 10)
}
