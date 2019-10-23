package paramfileserver

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

type ParamFileServer struct {
	Folder string
}

func (pfs ParamFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := path.Join(pfs.Folder, r.URL.RawQuery)
	log.Printf("Attempting to serve `%s`...\n", filename)
	file, err := os.Open(filename)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusNotFound)
	}
	io.Copy(w, file)
}
