package main

import(
	"net/http"
	"encoding/json"
	"strconv"
	"time"
)

func SearchHandler(index* Index) http.HandlerFunc {
	return func(w http.ResponseWriter,r *http.Request){

		start := time.Now()

		q:=r.URL.Query().Get("q")
		if q=="" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad request"))
			return
		}

		limit:=10
		offset:=0

		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		if limitStr!="" {
			if v,err := strconv.Atoi(limitStr); err==nil {
				limit=v
			}
		}

		if offsetStr != "" {
			if v, err := strconv.Atoi(offsetStr); err == nil {
				offset = v
			}
		}

		if limit <= 0 {
			limit = 10
		}
		if limit > 50 {
			limit = 50
		}
		if offset < 0 {
			offset = 0
		}

		results,total:=Search(index,q,limit,offset) 

		duration := time.Since(start)
		elapsed := duration.Milliseconds()

		resp := SearchResponse{
			Total: total,
			Results: results,
			TimeMS: elapsed,
		}
		
		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(resp)
		
	}
}