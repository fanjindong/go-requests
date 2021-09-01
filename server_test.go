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
	contentType := r.Header.Get("content-Type")
	switch contentType {
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
	default:
		req, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("invalid body: " + err.Error()))
			return
		}
		w.Write(req)
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

func uploadFile(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("file Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	fileName := r.FormValue("fileField")
	if fileName == "" {
		fileName = "file"
	}
	// FormFile returns the first file for the given key `fileName`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile(fileName)
	if err != nil {
		//fmt.Println("Error Retrieving the file", fileName, r.MultipartForm.file, r.MultipartForm.Value)
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded file: %+v\n", handler.Filename)
	//fmt.Printf("file Size: %+v\n", handler.Size)
	//fmt.Printf("MIME Header: %+v\n", handler.Header)
	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(fileBytes)
}

func TestMain(m *testing.M) {
	http.HandleFunc("/", handler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/timeout", timeoutHandler)
	http.HandleFunc("/header", headerHandler)
	http.HandleFunc("/upload", uploadFile)
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			panic(err)
		}
	}()
	code := m.Run()
	os.Exit(code)
}
