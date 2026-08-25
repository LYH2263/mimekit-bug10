package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/LYH2263/go-mimekit"
	"github.com/LYH2263/go-mimekit/workbench"
)

func main() {
	addr := flag.String("addr", ":8221", "listen address")
	flag.Parse()
	p := mimekit.NewPipeline(mimekit.Options{})
	api := &workbench.API{Pipeline: p}
	srv := workbench.New(api)
	log.Printf("mimekit workbench on %s", *addr)
	if err := http.ListenAndServe(*addr, srv.Handler); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
