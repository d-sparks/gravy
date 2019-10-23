package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/d-sparks/ace-of-trades/paramfileserver"
)

var port = flag.Int("port", 8080, "port to serve on")
var folder = flag.String("folder", "./data/mock/alphavantage/", "folder to serve")

func main() {
	flag.Parse()

	pfs := paramfileserver.ParamFileServer{Folder: *folder}
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), pfs); err != nil {
		log.Println(err.Error())
	}
}
