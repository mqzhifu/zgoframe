package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	fmt.Println(11)
	ssHH()
	for {

	}
}

func ssHH() {
	hostPort := "127.0.0.1:10000"
	http.HandleFunc("/", wh)
	err := http.ListenAndServe(hostPort, nil)
	if err != nil {
		fmt.Println("ListenAndServe err:", err)
	}
}

func wh(w http.ResponseWriter, r *http.Request) {

}
