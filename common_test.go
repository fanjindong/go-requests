package requests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	params := map[string]string{}
	//fmt.Printf("params: %+v", query)
	for k, v := range query {
		params[k] = v[0]
	}
	body, _ := json.Marshal(params)
	w.Write(body)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	data := map[string]interface{}{}
	switch contentType {
	case "application/json":
		req, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("invalid body: " + err.Error()))
			return
		}
		if err = json.Unmarshal(req, &data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("parse body err: " + err.Error()))
			return
		}
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("parse body: " + err.Error()))
			return
		}
		for k, v := range r.PostForm {
			data[k] = v[0]
		}
	}
	body, _ := json.Marshal(data)
	w.Write(body)
}

func timeoutHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	w.Write([]byte("OK"))
}

func TestMain(m *testing.M) {
	http.HandleFunc("/", handler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/timeout", timeoutHandler)
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()
	session = NewSession()
	//time.Sleep(1 * time.Second)

	code := m.Run()
	os.Exit(code)
}
