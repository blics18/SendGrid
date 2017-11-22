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
	client.PopulateDB(10, 100, 5)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		client.Populate()
		b.StopTimer()
		client.Clear()
		b.StartTimer()
	}
}

func BenchmarkPopulate(b *testing.B) {
	client.PopulateDB(10, 100, 5)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		client.Clear()
		b.StopTimer()
		client.Populate()
		b.StartTimer()
	}
}

func benchmarkCheck(i int, email []string, b *testing.B) {
	db := client.PopulateDB(10, 100, 5)
	client.Populate()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		client.Check(i, email)
	}
	client.DropTables(5, db)
	db.Close()
}

func BenchmarkCheck1(b *testing.B) {
	users := client.MakeRandomUsers(1, 100)
	for n := 0; n < len(users); n++ {
		benchmarkCheck(*users[n].UserID, users[n].Email, b)
	}
}

func BenchmarkCheck10(b *testing.B) {
	users := client.MakeRandomUsers(10, 100)
	for n := 0; n < len(users); n++ {
		benchmarkCheck(*users[n].UserID, users[n].Email, b)
	}
}

/*
func BenchmarkCheck1000(b *testing.B) {
	users := client.MakeRandomUsers(1000, 100)
	for n := 0; n < len(users); n++ {
		benchmarkCheck(*users[n].UserID, users[n].Email, b)
	}
}


*/
