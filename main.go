package main

import (
	"log"
	"net/http"
)

func main() {
	path:="index.gob"
	idx,err:=LoadIndex(path)         //If LoadIndex() succeeds, crawling never happens.
	if err!=nil {
		idx,err=BuildIndexFromScratch()
		if err==nil {
			_=idx.SaveIndex(path)
		}
	}

	handler:=SearchHandler(idx)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/search",handler)

	log.Fatal(http.ListenAndServe(":8080",nil))
}
