package main

import (
	"bytes"
	"strings"
	"io"
	"os"
	_ "time"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"crypto/rand"
	"mime"
	"path/filepath"
	_ "html/template"
)

type event struct {
	id          string `json:"id"`
	name        string `json:"name"`
	duration    string `json:"duration"`
	size        string `json:"size"`
}

type allEvents []event

var events allEvents
var results []string

const maxUploadSize = 10 * 1024 // 2 MB 
const uploadPath = "./storage"

// func initEvents() {
// 	events = allEvents{
// 		{
// 			id:          "1",
// 			name:    "myfile.wav",
// 			duration:    "300",
// 			size:        "250MB",
// 		},
// 		{
// 			id:          "2",
// 			name:    "bokunofairu.wav",
// 			duration:    "200",
// 			size:        "200MB",
// 		},
// 		{
// 			id:          "3",
// 			name:    "meuarquivo.wav",
// 			duration:    "500",
// 			size:        "700MB",
// 		},
// 	}
// }

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Recorded files")
}

func dlFile(w http.ResponseWriter, r *http.Request) {
	filename, ok := r.URL.Query()["name"]
    
    if !ok || len(filename[0]) < 1 {
        log.Println("Url Param 'name' is missing")
        return
    }

    file := filename[0]

	log.Println("Url Param 'name' is: " + string(file))

	data, err := ioutil.ReadFile(file)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		// log.Fatal(err)
	}
	_ = data
//	http.ServeContent(w, r, file, time.Now(), bytes.NewReader(data))
	http.ServeFile(w, r, file)
}

func ulFile(w http.ResponseWriter, r *http.Request) {

	// if err := r.ParseMultipartForm(maxUploadSize); err != nil {
	// 	fmt.Printf("Could not parse multipart form: %v\n", err)
	// 	// renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
	// 	// panic(err)
	// 	// // return
	// }

	fmt.Println("ParseForm")
	r.ParseForm()
	fmt.Println("PostFormValue")
	fileType := r.PostFormValue("type")
	fmt.Println("FormFile")
	file, fileHeader, err := r.FormFile("uploadFile")
		// if err != nil {
		// 	renderError(w, "INVALID_FILE", http.StatusBadRequest)
		// 	panic(err)
		// 	// return
		// }
	fmt.Println("defer FileClose")
	defer file.Close()

	fileSize := fileHeader.Size
	fmt.Println("File size (bytes): %v\n", fileSize)
	// if fileSize > maxUploadSize {
	// 	renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
	// 	return
	// }

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		panic(err)
		// return
	}

	filetype := http.DetectContentType(fileBytes)
	fmt.Printf("File type: %v\n", filetype)
	if filetype != "audio/vnd.wav" {
		renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
		panic(err)
		// return
	}
	// if filetype != "audio/vnd.wav" && filetype != "image/jpg" &&
	// 	filetype != "image/gif" && filetype != "image/png" &&
	// 	filetype != "application/pdf" {
	// 	renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
	// 	return
	// }

	fileName := randToken(12)
	fileEndings, err := mime.ExtensionsByType(fileType)
	if err != nil {
		renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
		panic(err)
		// return
	}
	newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
	fmt.Printf("FileType: %s, File: %s\n", fileType, newPath)

	
	newFile, err := os.Create(newPath)
	if err != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		panic(err)
		// return
	}

	defer newFile.Close()
	if _, err := newFile.Write(fileBytes); err != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		panic(err)
		// return
	}
	w.Write([]byte("SUCCESS"))
	
	
	// n, err := io.Copy(file, r.Body)
	// if err != nil {
	// 	panic(err)
	// }

	// w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", n)))
}


// https://gist.github.com/ebraminio/576fdfdff425bf3335b51a191a65dbdb
func ulFile3(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request: ",r)
	fmt.Println("file name is :",r.URL.Path[1:])
	fmt.Println("URL:",r.URL)
	fmt.Println("method :", r.Method)
	fmt.Println("FileSize : ", r.ContentLength)
	fmt.Println("Header : ",r.Header)
	fmt.Println("Body : ",r.Body)
	fmt.Println("Host: ",r.Host)
	fmt.Println("TransferEncoding: ",r.TransferEncoding)
	fmt.Println("Trailer: ",r.Trailer)
	fmt.Println("RequestURI: ",r.RequestURI)

	file, err := os.Create("./storage/result")
	if err != nil {
		panic(err)
	}
	n, err := io.Copy(file, r.Body)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", n)))
}

func ulFile2(w http.ResponseWriter, r *http.Request) {

	// r.ParseForm()
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("thefile")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Post form: ",r.PostForm)
	fmt.Println("Form: ",r.Form)
	fmt.Println("Form[data]: ",r.Form.Get("data"))
	fmt.Println("Form[name]: ",r.Form.Get("name"))
	fmt.Println("Form[@]: ",r.Form.Get("@"))
	fmt.Println("Form[file]: ",r.Form.Get("file"))
	fmt.Println("Form[filename]: ",r.Form.Get("filename"))
	fmt.Println("FormValue: ",r.FormValue)
	fmt.Println("FormFile: ",r.FormFile)

	for key, values := range r.Form {
		fmt.Println("Index: ", key)
		fmt.Println("Values: ", values)
		fmt.Println("Value: ", r.Form.Get(key))
	}

	for key, values := range r.PostForm {
		fmt.Println("Index: ", key)
		fmt.Println("Values: ", values)
		fmt.Println("Value: ", r.PostForm.Get(key))
	}

	_ = file

	params := mux.Vars(r)
	fmt.Println("New file name is :", params["fileName"])

	afile, err := os.Create("./storage/result")
	if err != nil {
		panic(err)
	}
	n, err := io.Copy(afile, r.Body)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(fmt.Sprintf("%d bytes are received.\n", n)))

}

func ReceiveFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20) // limit your max input length!
	// r.ParseForm()
    var buf bytes.Buffer
    // in your case file would be fileupload
    file, header, err := r.FormFile("file")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    name := strings.Split(header.Filename, ".")
    fmt.Printf("File name %s\n", name[0])
    // Copy the file data to my buffer
    io.Copy(&buf, file)
    // do something with the contents...
    // I normally have a struct defined and unmarshal into a struct, but this will
    // work as an example
    contents := buf.String()
    fmt.Println(contents)
    // I reset the buffer in case I want to use it again
    // reduces memory allocations in more intense projects
    buf.Reset()
    // do something else
    // etc write header
    return
}

// https://flaviocopes.com/golang-http-post-parameters/
func handleReq(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}
	v := req.Form
	data := req.Form.Get("data")
	name := req.Form.Get("name")
	_ = v
	_ = data
	_ = name
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event filename and duration only in order to update")
	}
	
	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.id == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var updatedEvent event

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event filename and duration only in order to update")
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.id == eventID {
			singleEvent.name = updatedEvent.name
			singleEvent.duration = updatedEvent.duration
			singleEvent.size = updatedEvent.size
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.id == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func main() {
	// initEvents()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/download", dlFile).Methods("GET")
	router.HandleFunc("/post", ulFile3).Methods("POST")
	// router.HandleFunc("/post", ReceiveFile).Methods("POST")
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}