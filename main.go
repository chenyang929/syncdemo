package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syncdemo/gracehttp"
)

var bs []byte
var graceful = flag.Bool("graceful", false, "listen on fd open 3 (internal use only)")

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()
	var err error
	f, err := os.Open("words.txt")
	if err != nil {
		return
	}
	bs, err = ioutil.ReadAll(f)
	if err != nil {
		return
	}
}

func getRouter() http.Handler {
	http.Handle("/", http.HandlerFunc(sayhelloName))
	return http.DefaultServeMux
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	w.Write(bs)
}

func main() {
	router := getRouter()
	host := ":9090"
	//graceServer, err := grace.New(host, router, *graceful)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//graceServer.Start()
	err := gracehttp.ListenAndServe(host, router)
	if err != nil {
		log.Println(err)
	}
}