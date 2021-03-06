package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blics18/SendGrid/client"
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
	userID := 1
	userIDEmail := "1|eFeOnJkMqw@aol.com"
	email := "eFeOnJkMqw@aol.com"
	bf := NewBloomFilter(1000)
	bf.Filter.Add([]byte(userIDEmail))

	user := client.User{
		UserID: &userID,
		Email:  []string{email},
	}

	userJSON, err := json.MarshalIndent(user, "", " ")

	req, err := http.NewRequest("GET", "/checkBF", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(bf.checkBF)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var statStruct client.Stats
	statStruct.Miss = 1
	statStruct.Total = 1
	
	structToStr, _ := json.MarshalIndent(statStruct, "", " ")
	
	if rr.Body.String() != string(structToStr) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(structToStr))
	}
}

func TestUnsubClear(t *testing.T) {
	userID := 3
	userIDEmail := "3|hzSfMqs@aol.com"
	email := "hzSfMqs@aol.com"
	bf := NewBloomFilter(1000)

	bf.Filter.Add([]byte(userIDEmail))

	user := client.User{
		UserID: &userID,
		Email:  []string{email},
	}

	userJSON, err := json.MarshalIndent(user, "", " ")

	checkReq, err := http.NewRequest("GET", "/checkBF", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(bf.checkBF)

	handler.ServeHTTP(rr, checkReq)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var statStruct client.Stats
	statStruct.Miss = 1
	statStruct.Total = 
	1
	structToStr, _ := json.MarshalIndent(statStruct, "", " ")

	if rr.Body.String() != string(structToStr) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(structToStr))
	}

	clearReq, err := http.NewRequest("GET", "/clearBF", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}

	clearRR := httptest.NewRecorder()
	clearHandler := http.HandlerFunc(bf.clearBF)

	clearHandler.ServeHTTP(clearRR, clearReq)

	if status := clearRR.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	clearExpected := "Successfully Cleared Bloom Filter"
	if clearRR.Body.String() != clearExpected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			clearRR.Body.String(), clearExpected)
	}

	req, err := http.NewRequest("GET", "/checkBF", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}

	checkRR := httptest.NewRecorder()
	checkHandler := http.HandlerFunc(bf.checkBF)

	checkHandler.ServeHTTP(checkRR, req)

	if status := checkRR.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	statStruct.Hits = 1
	statStruct.Miss = 0
	statStruct.Total = 1
	structToStr, _ = json.MarshalIndent(statStruct, "", " ")

	if checkRR.Body.String() != string(structToStr) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			checkRR.Body.String(), string(structToStr))
	}
}