package requests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

var port = 8080

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	params := map[string]string{}
	for k, v := range query {
		params[k] = v[0]
	}
	body, _ := json.Marshal(params)
	w.Write(body)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		req, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("invalid body: " + err.Error()))
			return
		}
		w.Write(req)
		return
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("parse body: " + err.Error()))
			return
		}
		data := map[string]interface{}{}
		for k, v := range r.PostForm {
			data[k] = v[0]
		}
		body, _ := json.Marshal(data)
		w.Write(body)
		return
	}
}

func timeoutHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	w.Write([]byte("OK"))
}

func headerHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := json.Marshal(r.Header)
	w.Write(bytes)
}

func cooclerHandler(w http.ResponseWriter, r *http.Request) {
	r.Cookies()
	w.Write([]byte("OK"))
}

func TestMain(m *testing.M) {
	http.HandleFunc("/", handler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/timeout", timeoutHandler)
	http.HandleFunc("/header", headerHandler)
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			panic(err)
		}
	}()
	code := m.Run()
	os.Exit(code)
}
