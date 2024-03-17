//go:build e2e
// +build e2e

package test

import (
	"log"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestHealthEndpoint(t *testing.T) {
	log.Println("Running E2E tests for health check endpoint")

	client := resty.New()
	resp, err := client.R().Get("http://localhost:8080/api/health")
}
