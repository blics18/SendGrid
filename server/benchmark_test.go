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
<<<<<<< HEAD
	client.PopulateDB(10, 100, 5)
	cfg := client.GetEnv()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		client.Populate(cfg)
		b.StopTimer()
		client.Clear(cfg)
		b.StartTimer()
=======
	for n := 0; n < b.N; n++ {
		client.Clear()
>>>>>>> master
	}
}

func BenchmarkPopulate(b *testing.B) {
<<<<<<< HEAD
	client.PopulateDB(10, 100, 5)
	cfg := client.GetEnv()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		client.Clear(cfg)
		b.StopTimer()
		client.Populate(cfg)
		b.StartTimer()
	}
	
}

func benchmarkCheck(i int, email []string, b *testing.B) {
	db, _ := client.PopulateDB(10, 100, 5)
	cfg := client.GetEnv()
	client.Populate(cfg)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := client.Check(cfg, i, email)
		if err != nil{
			b.Error(err)
			b.FailNow()
		}
	}
	// client.DropTables(5, db)
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

// func BenchmarkCheck1000(b *testing.B) {
// 	users := client.MakeRandomUsers(1000, 100)
// 	for n := 0; n < len(users); n++ {
// 		benchmarkCheck(*users[n].UserID, users[n].Email, b)
// 	}
// }
=======
	for n := 0; n < b.N; n++ {
		client.Populate()
	}
}

func BenchmarkCheck(b *testing.B) {
	for n := 0; n < b.N; n++ {
		client.Check()
	}
}
>>>>>>> master
