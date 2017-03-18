package oss

import (
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

func NewClient(endpoint, accessKeyId, accessKeySecret string) *Client {
	return &Client{
		Endpoint:        endpoint,
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
}

func (c *Client) NewBucket(bucket string) *Bucket {
	return &Bucket{
		Name:   bucket,
		client: c,
	}
}

func (c *Client) NewTicket(verb, bucket, object string) *Ticket {
	return &Ticket{
		client:    c,
		Verb:      verb,
		Bucket:    bucket,
		Object:    object,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Header:    make(http.Header),
		Query:     make(url.Values),
	}
}
