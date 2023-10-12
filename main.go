package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// the movie datatype
type Movie struct {
	// each key also represents its json-encoding key.
	// this json-encoding key will be helpful when we want to encode/decode the json.
	ID    string `json:"id"`
	Isbn  string `json:"isbn"`
	Title string `json:"title"`

	// each movie has a director and we are referring to that using pointer.
	Director *Director `json:"director"`
}

// the director datatype
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// since, we are not using the database
// we created a slice for holding the data for us per session.
var movies []Movie

// the main function.
func main() {

	// we are using a package called Gorilla MUX for http routing.
	router := mux.NewRouter()

	// we are adding some dummy movies to our slice
	addDummyData()

	// creating routes and their handlers
	// get all movies
	router.HandleFunc("/movies", getMovies).Methods("GET")
	// get movie by its id
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	// create movie
	router.HandleFunc("/movies", createMovie).Methods("POST")
	// update movie by its id
	router.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	// delete movie by its id
	router.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// the http port to run the server on
	port := ":8080"
	fmt.Printf("Starting server on %s \n", port)
	// start and listen for error if failed to start
	log.Fatal(http.ListenAndServe(port, router))
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	// set the content type in the header.
	w.Header().Set("Content-Type", "application/json")
	// encode the movies from goland struct type to json
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// read the url-params using mux
	// the params are available in a map with string key
	params := mux.Vars(r)

	// find the movie by id and send back in response
	for _, movie := range movies {
		if movie.ID == params["id"] {
			json.NewEncoder(w).Encode(movie)
			break
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// decode the movie from json to Movie struct type
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	// create a random id
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	movies = append(movies, movie)
	// write updated result
	json.NewEncoder(w).Encode(movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			// remove the movie
			movies = append(movies[:index], movies[index+1:]...)
			// add a new movie
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			// use the same ID as old movie
			movie.ID = item.ID
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, movie := range movies {
		if movie.ID == params["id"] {
			// we are using append function to remove a value at found index
			// this line means write whatever we have from index+1 to index, i.e shift
			// hence, removing the value at specified index
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func addDummyData() {
	movies = append(
		movies,
		Movie{
			ID:    "1",
			Isbn:  "564912",
			Title: "Sinbad: The Sailor",
			Director: &Director{
				Firstname: "Avinash",
				Lastname:  "Gautam",
			},
		},
	)

	movies = append(
		movies,
		Movie{
			ID:    "2",
			Isbn:  "12314",
			Title: "Rocky Rani and Kahani",
			Director: &Director{
				Firstname: "Ankita",
				Lastname:  "Gupta",
			},
		},
	)
}
