package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func Fetch(url string) (string, error) {
    //Creating a client with a timeout so workers don't get stuck forever
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

    //Creating a request so we can set headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

    // User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

    //Check for non-200 status codes (like 403 Forbidden or 404)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}