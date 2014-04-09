//+build skip

package main

import (
	"encoding/json"
	"fmt"
	bleed "github.com/quiffman/Heartbleed/bleed"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var PAYLOAD = []byte("heartbleed.filippo.io")

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://filippo.io/Heartbleed", http.StatusFound)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hi there.")
}

type result struct {
	Code  int    `json:"code"`
	Data  string `json:"data"`
	Error string `json:"error"`
}

func bleedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	host := r.URL.Path[len("/bleed/"):]
	u, err := url.Parse(host)
	if err == nil && u.Host != "" {
		host = u.Host
	}

	tgt := bleed.Target{
		HostIp: string(host),
	}
	data, err := bleed.Heartbleed(&tgt, PAYLOAD)
	var rc int
	var errS string
	if err == bleed.Safe {
		rc = 1
		data = []byte("")
		log.Printf("%v - SAFE", host)
	} else if err != nil {
		rc = 2
		data = []byte("")
		errS = err.Error()
		log.Printf("%v - ERROR", host)
	} else {
		rc = 0
		log.Printf("%v - VULNERABLE", host)
	}
	res := result{rc, string(data), errS}
	j, err := json.Marshal(res)
	if err != nil {
		log.Println("ERROR", err)
	} else {
		w.Write(j)
	}
}

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/bleed/", bleedHandler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
