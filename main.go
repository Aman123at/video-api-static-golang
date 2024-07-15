package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Video struct {
	Id          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	UploadedBy  *User  `json:"uploadedBy"`
}

type User struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email"`
}

func (v *Video) IsEmpty() bool {
	return v.Title == ""
}

var allVideos = []Video{}

func main() {

	// seeding data
	allVideos = append(allVideos, Video{
		Id:          uuid.New().String(),
		Title:       "Video 1",
		Description: "Description 1",
		Duration:    100,
		UploadedBy:  &User{Id: uuid.New().String(), Name: "John Doe", Email: "johndoe@gmail.com"},
	}, Video{
		Id:          uuid.New().String(),
		Title:       "Video 2",
		Description: "Description 2",
		Duration:    200,
		UploadedBy:  &User{Id: uuid.New().String(), Name: "Jane", Email: "jane@gmail.com"}})

	fmt.Println("Welcome to youtube api")

	// initiate router
	r := mux.NewRouter()

	r.HandleFunc("/", showWelcome).Methods("GET")

	r.HandleFunc("/videos", getAllVideos).Methods("GET")

	r.HandleFunc("/video", addVideo).Methods("POST")

	r.HandleFunc("/video/update/{id}", updateVideo).Methods("PUT")

	r.HandleFunc("/video/delete/{id}", deleteVideo).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":4000", r))
}

func showWelcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to youtube api in golang</h1>"))
}

func getAllVideos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allVideos)
}

func addVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}
	var video Video
	_ = json.NewDecoder(r.Body).Decode(&video)

	if video.IsEmpty() {
		json.NewEncoder(w).Encode("No data inside request body")
	}
	allVideos = append(allVideos, video)
	json.NewEncoder(w).Encode(video)
}

func updateVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}
	var video Video
	_ = json.NewDecoder(r.Body).Decode(&video)

	for index, item := range allVideos {
		if item.Id == params["id"] {
			allVideos = append(allVideos[:index], allVideos[index+1:]...)
			allVideos = append(allVideos, video)
			break
		}
	}
	json.NewEncoder(w).Encode(allVideos)
}

func deleteVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range allVideos {
		if item.Id == params["id"] {
			allVideos = append(allVideos[:index], allVideos[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(allVideos)
}
