package main

import (
	"fmt"
	"html"
	"net/http"

	"github.com/willf/bloom"
)

//var n int = uint(1000)
var filter *bloom.BloomFilter

func add(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userid := r.FormValue("userid")
	email := r.FormValue("email")
	fmt.Printf("uid = %s	email=%s\n", userid, email)
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	filter.Add([]byte("hi"))

}

/*func test(t *bloom.BloomFilter) {
	(*filter).Add([]byte("hi"))
}*/
func main() {
	n := uint(1000)
	filter := bloom.New(20*n, 5) // load of 20, 5 keys
	filter.Add([]byte("Love"))
	//	fmt.Println(reflect.TypeOf(filter))
	//	test(&(*filter))
	if filter.Test([]byte("hi")) {
		fmt.Printf("yes")
	} else {
		fmt.Printf("no")
	}
	/*
		i := uint32(100)
		n1 := make([]byte, 4)
		binary.BigEndian.PutUint32(n1, i)
		filter.Add(n1)
		if filter.EstimateFalsePositiveRate(1000) > 0.001 {
			fmt.Printf("yes")
		} else {
			fmt.Printf("no")
		}
	*/
	//	http.HandleFunc("/add", add)
	//	log.Fatal(http.ListenAndServe(":8081", nil))
}
