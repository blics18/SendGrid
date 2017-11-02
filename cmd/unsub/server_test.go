package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sheilatruong96/unsub"
	"github.com/stretchr/testify/assert"
	"github.com/willf/bloom"
)

func TestBloomfilter(t *testing.T) {
	n := uint(1000)
	filter := bloom.New(20*n, 5)
	filter.Add([]byte("hello"))
	exists := filter.Test([]byte("hello"))
	assert.True(t, exists)
}

func TestUnsubCheck(t *testing.T) {
	bloomFilter = createBloomFilter()
	userID := 1

	user := unsub.User{
		UserID: &userID,
		Email:  []string{"hello@gmail.com", "hello2@aim.com"},
	}

	userJSON, err := json.MarshalIndent(user, "", " ")

	req, err := http.NewRequest("GET", "/checkBF", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(checkBF)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `Success`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
