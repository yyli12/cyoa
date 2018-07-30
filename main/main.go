package main

import (
	"net/http"
	"fmt"
	"github.com/gophercises/cyoa"
	"io/ioutil"
	"path/filepath"
	"log"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world! " + r.URL.Path)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	storyServer()
}

func handlerError(e error) {
	if e != nil {
		panic(e)
	}
}

func storyServer() {
	absFilePath, err := filepath.Abs("./cyoa/gopher.json")
	handlerError(err)

	jsonBytes, err := ioutil.ReadFile(absFilePath)
	handlerError(err)

	story, err := cyoa.LoadStory(jsonBytes)
	handlerError(err)

	http.ListenAndServe(":7777", cyoa.StoryHandler(story, defaultMux()))
}