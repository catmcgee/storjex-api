// TO DO
// Get access grants working
// Deploy
package main

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type server struct{}

const (
    myAccessGrant = "1cCYAHiNyogwPRkZbkpo8C6txEeJ4cktWZoHzwAKTwNzGeAKhXyJZojfb8jSGGFSAc5M3zGgswJQXVEkrJ8na6nj4aeRRdtwvjgGFC5GPFV8NGo9NGsKn8XjsRADBE7RziJKd5XkDNxyQ8HWh1maxtonsLqd5zNKnvh1nFPHnntTqyhXVi6NwubuzeX78ranRrnn1Wxt5vuqcKMxdFeNF8RWhRzRD2HskLm54GgTcRcbG5vnga2Dfx34vxgvs19s1j2TR5BPoHiS"
    bucket = "storjex"
)

func main() {
    r := mux.NewRouter()
    api := r.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("", uploadFile).Methods(http.MethodPost)
    api.HandleFunc("", deleteFile).Methods(http.MethodDelete)
    api.HandleFunc("/file/{passphrase}", downloadFile).Methods(http.MethodGet)
	// api.HandleFunc("/fileInfo", getFileInfo).Methods(http.MethodGet)
    log.Fatal(http.ListenAndServe(":10000", r))
}

func uploadFile(w http.ResponseWriter, r * http.Request) {
    r.ParseMultipartForm(10 << 20)
	file, _, err:= r.FormFile("file")
    if err != nil {
        message := fmt.Sprint("Error parsing form: ", err)
        w.Write([]byte(message)); 
    } 
	name := r.FormValue("name")
    numberOfDownloads, err := strconv.Atoi(r.FormValue("numberOfDownloads"))
    if err != nil {
        message := fmt.Sprint("Error parsing form: ", err)
        w.Write([]byte(message)); 
    }

    defer file.Close()
    var fileByte[] byte
    fileByte = make([] byte, 1000000)
    fileByte, byteErr:= ioutil.ReadAll(file)
    if byteErr != nil {
        message := fmt.Sprint("Error converting file: ", err)
        w.Write([]byte(message)); 
    }

    passphrase, adminPassphrase, objectKey := UploadData(context.Background(),
        myAccessGrant, bucket, name, fileByte, numberOfDownloads)

	w.WriteHeader(http.StatusOK)
	message := fmt.Sprintf(`{"passphrase": "%s", "adminPassphrase": "%s", bucket: "%s", key: "%s"}`, passphrase, adminPassphrase, bucket, objectKey)
    w.Write([]byte(message))
}


func downloadFile(w http.ResponseWriter, r * http.Request) {
	pathParams := mux.Vars(r)
    passphrase := pathParams["passphrase"]
	        // pass this passphrase to download data
    fileByte, err := DownloadData(passphrase)
    if err != nil {
        message := fmt.Sprint("Error: ", err)
        w.Write([]byte(message)); 
    } else {
	encodedString := b64.StdEncoding.EncodeToString([]byte(fileByte))
	message := fmt.Sprintf(`{"image": "%s"}`, encodedString)
	w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(message)); 
}}

func deleteFile(w http.ResponseWriter, r * http.Request) {
	r.ParseMultipartForm(10 << 20)
	passphrase := r.FormValue("adminPassphrase")
	response := HandleDelete(passphrase) 
    w.Write([]byte(response))
}


