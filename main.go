package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	myip, err := httpGet("http://readmyip.appspot.com/")
	if err != nil {
		log.Fatalf("Failed to: %v", err)
	}

	until := time.Now().Add(10 * time.Minute)
	for {
		if time.Now().After(until) {
			log.Fatalf("Timeout waiting change IP")
		}
		s, err := httpGet("http://readmyip.appspot.com/")
		if err != nil {
			log.Fatalf("Failed to: %v", err)
		}
		if s != myip {
			os.Exit(0)
		}
		time.Sleep(10 * time.Second)
	}
}

func httpGet(s string) (string, error) {
	c := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   500 * time.Millisecond,
				KeepAlive: 30 * time.Second,
				DualStack: false,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 5 * time.Second,
			MaxIdleConns:          10,
			MaxIdleConnsPerHost:   http.DefaultMaxIdleConnsPerHost,
			IdleConnTimeout:       20 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get(s)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf := make([]byte, 1024)
	_, err = resp.Body.Read(buf)
	if err == io.EOF {
		err = nil
	}
	return string(buf), err
}
