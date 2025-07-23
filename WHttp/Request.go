package main

type Request struct {
	Method  string            // HTTP method (GET, POST, etc.)
	URL     string            // Request URL
	Headers map[string]string // Request headers
}
