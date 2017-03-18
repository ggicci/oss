package oss

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Ticket wraps information for accessing to OSS resources.
type Ticket struct {
	client *Client

	Verb      string      `json:"verb"`       // GET, POST, PUT
	Bucket    string      `json:"bucket"`     // which bucket to use
	Object    string      `json:"object"`     // e.g. "/avatar/19910628.png"
	ExpiresAt time.Time   `json:"expires_at"` // expires at
	Header    http.Header `json:"header"`     // headers
	Query     url.Values  `json:"query"`      // sub resource and other queries
	SignedURL string      `json:"url"`        // final signed url
}

func (tk *Ticket) Sign() error {
	c := tk.client

	u, err := url.ParseRequestURI(c.Endpoint)
	if err != nil {
		return fmt.Errorf("invalid oss endpoint: \"%s\", %v", c.Endpoint, err)
	}

	if tk.Bucket != "" {
		u.Host = tk.Bucket + "." + u.Host // 3rd level domain
	}
	if tk.Object != "" && tk.Bucket == "" {
		return errors.New("set object, but bucket not set")
	}

	expires := strconv.FormatInt(tk.ExpiresAt.Unix(), 10)
	// Make signature.
	var bf bytes.Buffer
	bf.WriteString(tk.Verb)
	bf.WriteString("\n")
	bf.WriteString(tk.Header.Get("Content-MD5"))
	bf.WriteString("\n")
	bf.WriteString(tk.Header.Get("Content-Type"))
	bf.WriteString("\n")
	bf.WriteString(expires)
	bf.WriteString("\n")
	canonicalizedOSSHeaders := tk.buildCanonicalizedOSSHeadersString()
	if canonicalizedOSSHeaders != "" {
		bf.WriteString(canonicalizedOSSHeaders)
		bf.WriteString("\n")
	}

	bf.WriteString(tk.buildCanonicalizedResourceString())

	mac := hmac.New(sha1.New, []byte(c.AccessKeySecret))
	mac.Write(bf.Bytes())
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	tk.Query.Set("OSSAccessKeyId", c.AccessKeyId)
	tk.Query.Set("Expires", expires)
	tk.Query.Set("Signature", signature)

	fmt.Printf("%X", bf.Bytes())

	// Assemble the URL.
	u.RawQuery = tk.Query.Encode()
	u.Path = tk.resourcePath()
	println(u.String())
	tk.SignedURL = u.String()
	return nil
}

func (tk *Ticket) buildCanonicalizedOSSHeadersString() string {
	var items []string
	for name, _ := range tk.Header {
		cname := strings.TrimSpace(strings.ToLower(name))
		if !strings.HasPrefix(cname, "x-oss-") {
			continue
		}

		// build item: "key:value"
		item := cname
		value := strings.TrimSpace(tk.Header.Get(name))
		if value != "" {
			item += ":"
			item += value
		}

		items = append(items, item)
	}

	if len(items) == 0 {
		return ""
	}
	sort.Strings(items)
	return strings.Join(items, "\n")
}

// "/object"
func (tk *Ticket) resourcePath() string {
	s := "/"
	s += tk.Object
	return s
}

// "/bucket/object?query"
func (tk *Ticket) buildCanonicalizedResourceString() string {
	s := tk.resourcePath()
	q := tk.Query.Encode()
	if q != "" {
		s += "?"
		s += q
	}
	if tk.Bucket != "" {
		s = "/" + tk.Bucket + s
	}
	return s
}
