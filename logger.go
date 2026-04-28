package main

import(
	"strings"
	"log"
)
type Logger struct {

}
func (l *Logger) DetectUnseenTerms(tokens []string,dfmap map[string]int) {
	terms:=[]string{}
	unseen:=make(map[string]struct{})
	for _,token:= range tokens {
		_,ok:=dfmap[token]
		if !ok {
			unseen[token]=struct{}{}
		}
	}
	if len(unseen)==0 {
		return
	}
	for key :=range unseen {
		terms=append(terms,key)
	}
	unseenStr := strings.Join(terms, ", ")
	log.Printf("[WARN] unseen query terms: %s",unseenStr)

}