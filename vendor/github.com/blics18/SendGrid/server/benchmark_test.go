package main

import (
	// "bytes"
	// "encoding/json"
	// "fmt"
	// "net/http"
	// "net/http/httptest"
	"testing"

	"github.com/blics18/SendGrid/client"
	// "github.com/rcrowley/go-metrics"
	// "github.com/blics18/SendGrid/client"
	// "github.com/stretchr/testify/assert"
	// "github.com/willf/bloom"
)

func BenchmarkClear(b *testing.B) {
	for n := 0; n < b.N; n++ {
		client.Clear()
	}
}

func BenchmarkPopulate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		client.Populate()
	}
}

func BenchmarkCheck(b *testing.B) {
	for n := 0; n < b.N; n++ {
		client.Check()
	}
}