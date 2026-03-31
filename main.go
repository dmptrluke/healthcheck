package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: healthcheck http <host:port> <path> | file-age <path> <max-seconds>")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "http":
		if len(os.Args) != 4 {
			fmt.Fprintln(os.Stderr, "usage: healthcheck http <host:port> <path>")
			os.Exit(1)
		}
		httpCheck(os.Args[2], os.Args[3])
	case "file-age":
		if len(os.Args) != 4 {
			fmt.Fprintln(os.Stderr, "usage: healthcheck file-age <path> <max-seconds>")
			os.Exit(1)
		}
		fileAgeCheck(os.Args[2], os.Args[3])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func httpCheck(host, path string) {
	client := &http.Client{Timeout: 4 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://%s%s", host, path))
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "healthcheck failed: status %d\n", resp.StatusCode)
		os.Exit(1)
	}
}

func fileAgeCheck(path, maxAgeStr string) {
	maxAge, err := strconv.ParseFloat(maxAgeStr, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid max-age: %v\n", err)
		os.Exit(1)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck failed: %v\n", err)
		os.Exit(1)
	}
	ts, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid timestamp: %v\n", err)
		os.Exit(1)
	}
	age := math.Abs(float64(time.Now().Unix()) - ts)
	if age > maxAge {
		fmt.Fprintf(os.Stderr, "healthcheck failed: age %.0fs exceeds %.0fs\n", age, maxAge)
		os.Exit(1)
	}
}