package sensetive

import (
	"github.com/importcjj/sensitive"
	"log"
	"sync"
)

var (
	Filter     *sensitive.Filter
	filterOnce sync.Once
)

const WordDictPath = "./sensetive_doc/sensitiveWord.txt"

func InitFilter() {
	filterOnce.Do(func() {
		Filter = sensitive.New()
		err := Filter.LoadWordDict(WordDictPath)
		if err != nil {
			log.Println("InitFilter Fail,Err=" + err.Error())
		}
	})
}

func GetFilter() *sensitive.Filter {
	InitFilter()
	return Filter
}