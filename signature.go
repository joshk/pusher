package pusher

import (
    "crypto/hmac"
    "crypto/md5"
    "crypto/sha256"
    "encoding/hex"
    "io"
    "net/url"
    "strings"
)

type Signature struct {
    key, secret                                   string
    method, path, timestamp, authVersion, content string
}

func (s *Signature) Sign() string {
    authParts := strings.Join([]string{s.auth_key(), s.auth_timestamp(), s.auth_version(), s.body_md5()}, "&")
    signatureContent := strings.Join([]string{s.method, s.path, authParts}, "\n")
    return s.hmacSha256(signatureContent)
}

func (s *Signature) EncodedQuery() string {
    query := url.Values{}
    query.Set("auth_key", s.key)
    query.Set("auth_timestamp", s.timestamp)
    query.Set("auth_version", s.authVersion)
    query.Set("body_md5", s.md5Content())
    query.Set("auth_signature", s.Sign())
    return query.Encode()
}

func (s *Signature) auth_key() string {
    return "auth_key=" + s.key
}

func (s *Signature) auth_timestamp() string {
    return "auth_timestamp=" + s.timestamp
}

func (s *Signature) auth_version() string {
    return "auth_version=" + s.authVersion
}

func (s *Signature) body_md5() string {
    return "body_md5=" + s.md5Content()
}

func (s *Signature) md5Content() string {
    return s.md5(s.content)
}

func (s *Signature) md5(content string) string {
    hash := md5.New()
    io.WriteString(hash, content)
    return hex.EncodeToString(hash.Sum(nil))
}

func (s *Signature) hmacSha256(content string) string {
    hash := hmac.New(sha256.New, []byte(s.secret))
    io.WriteString(hash, content)
    return hex.EncodeToString(hash.Sum(nil))
}
