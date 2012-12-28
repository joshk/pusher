package pusher

import (
    "crypto/md5"
    "crypto/hmac"
    "crypto/sha256"
    "io"
    "strconv"
    "encoding/hex"
    "encoding/json"
    "strings"
    "net/url"
    "net/http"
    "time"
    "bytes"
)


var Endpoint = "api.pusherapp.com"

const AuthVersion = "1.0"


type Client struct {
    appid, key, secret string
    secure             bool
}

type Payload struct {
    Name        string `json:"name"`
    Channels  []string `json:"channels"`
    Data        string `json:"data"`
}


func NewClient(appid, key, secret string, secure bool) *Client {
    return &Client{appid, key, secret, secure}
}


func (c *Client) Publish(data, event string, channels ...string) error {
    timestamp  := c.stringTimestamp()
    
    content, err := c.jsonifyData(data, event, channels)
    if err != nil {
        return err
    }
    
    md5Content := c.md5(content)
    
    signature  := c.signature(timestamp, md5Content)
    
    query := c.encodedQuery(timestamp, md5Content, signature)

    err = c.post(content, c.publishUrl(), query)

    return err
}


func (c *Client) md5(content string) string {
    hash := md5.New()
    io.WriteString(hash, content)
    return hex.EncodeToString(hash.Sum(nil))
}


func (c *Client) hmacSha256(content string) string {
    hash := hmac.New(sha256.New, []byte(c.secret))
    io.WriteString(hash, content)
    return hex.EncodeToString(hash.Sum(nil))
}


func (c *Client) encodedQuery(timestamp, md5Content, signature string) string {
    query := make(url.Values)
    query.Set("auth_key",       c.key)
    query.Set("auth_timestamp", timestamp)
    query.Set("auth_version",   AuthVersion)
    query.Set("body_md5",       md5Content)
    query.Set("auth_signature", signature)
    return query.Encode()
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

    if c.secure {
        postUrl.Scheme   = "https"
    } else {
        postUrl.Scheme   = "http"
    }

    postUrl.RawQuery = query

    resp, err := http.Post(postUrl.String(), "application/json", buffer)
    if err != nil {
        return err
    }
    
    defer resp.Body.Close()
  
    return nil
}


func (c *Client) publishPath() string {
    return "apps/" + c.appid + "/events"
}


func (c *Client) publishUrl() string {
    return "http://" + Endpoint + "/" + c.publishPath()
}


func (c *Client) signature(timestamp, md5Content string) string {
    authParts := strings.Join([]string{"auth_key=" + c.key, "auth_timestamp=" + timestamp, "auth_version=" + AuthVersion, "body_md5=" + md5Content}, "&")
    signatureContent := strings.Join([]string{"POST", "/" + c.publishPath(), authParts}, "\n")
    return c.hmacSha256(signatureContent)
}


func (c *Client) stringTimestamp() string {
    t := time.Now()
    return strconv.FormatInt(t.Unix(), 10)
}


