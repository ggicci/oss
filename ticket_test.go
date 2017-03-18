package oss_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/ggicci/oss"
)

const (
	endpoint = "https://oss-cn-hangzhou.aliyuncs.com"
	key      = "<key>"
	secret   = "<secret>"
	bucket   = "<bucket>"
	testFile = "./ticket.go"
)

func TestSignTicket(t *testing.T) {
	c := oss.NewClient(endpoint, key, secret)
	tk := c.NewTicket("PUT", bucket, filepath.Base(testFile))
	tk.Header.Set("Content-Type", "text/plain")
	if err := tk.Sign(); err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Logf("ticket: %#v", tk)

	var httpClient http.Client
	f, err := os.Open(testFile)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	// bs := make([]byte, 0, 0)
	var buf = bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	encoder.Encode(tk)
	println(buf.String())

	r, err := http.NewRequest(tk.Verb, tk.SignedURL, f)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	// NB: header must match
	r.Header = tk.Header
	resp, err := httpClient.Do(r)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(errors.New(resp.Status))
		bs, _ := ioutil.ReadAll(resp.Body)
		t.Logf("body: %s", bs)
		t.Fail()
	}
}
