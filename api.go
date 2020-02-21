package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sdomino/scribble"
	"github.com/google/uuid"
)

// ============================================================================
// ============================================================================

// Data model
type Snippet struct {
	Id string
	Name string
	Description string
	Tags string
	Notes string
	Reference string
	Related string
	Files string
	Language string
	Public bool
}

func returnAllSnippets(w http.ResponseWriter, r * http.Request) {
	fmt.Println("Endpoint Hit: returnAll " + DATABASE_SNIPPET)
	records, err:= db.ReadAll(DATABASE_SNIPPET)
	if err != nil {
		fmt.Println("Error", err)
	}
	items:= [] Snippet {}
	for _, f:= range records {
		snippetFound:= Snippet {}
		if err:= json.Unmarshal([] byte(f), & snippetFound);err != nil {
			fmt.Println("Error", err)
		}
		items = append(items, snippetFound)
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func returnSingleSnippet(w http.ResponseWriter, r * http.Request) {
	fmt.Println("Endpoint Hit: return one " + DATABASE_SNIPPET)
	vars:= mux.Vars(r)
	key:= vars["id"]

	// Loop over all of our Articles
	// if the article.Id equals the key we pass in
	// return the article encoded as JSON
	item:= Snippet {}
	if err:= db.Read(DATABASE_SNIPPET, key, & item);
	err != nil {
		fmt.Println("Error", err)
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func createNewSnippet(w http.ResponseWriter, r * http.Request) {
	fmt.Println("Endpoint Hit: create new " + DATABASE_SNIPPET)
	b, err:= ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var item Snippet
	err = json.Unmarshal(b, & item)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	id, err:= uuid.NewUUID()
	if err != nil {
		// handle error
		fmt.Println("Error", err)
	}
	// fmt.Printf(id.String())
	item.Id = id.String()
	db.Write(DATABASE_SNIPPET, item.Id, item)

	output, err:= json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

func updateSnippet(w http.ResponseWriter, r * http.Request) {
	fmt.Println("Endpoint Hit: update " + DATABASE_SNIPPET)
	b, err:= ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var item Snippet
	err = json.Unmarshal(b, & item)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	db.Write(DATABASE_SNIPPET, item.Id, item)

	output, err:= json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

func deleteSnippet(w http.ResponseWriter, r * http.Request) {
	fmt.Println("Endpoint Hit: delete one " + DATABASE_SNIPPET)
	vars:= mux.Vars(r)
	id:= vars["id"]
	fmt.Println("Removing: " + string(DATABASE_SNIPPET) + "/" + string(id))
	if err:= db.Delete(DATABASE_SNIPPET, string(id));
	err != nil {
		fmt.Println("Error", err)
	}
	fmt.Fprintf(w, "%+v", "Removed: " + string(DATABASE_SNIPPET) + "/" + string(id))
}

// ============================================================================
// ============================================================================

func handleRequests() {
	router:= mux.NewRouter().StrictSlash(true)

	// NOTE: Ordering is important here!
	SECTION_NAME:= "snippet"
	router.HandleFunc(API_PREFIX + SECTION_NAME + "/", returnAllSnippets).Methods("GET")
	router.HandleFunc(API_PREFIX + SECTION_NAME + "", createNewSnippet).Methods("POST")
	router.HandleFunc(API_PREFIX + SECTION_NAME + "/{id}", updateSnippet).Methods("PUT")
	router.HandleFunc(API_PREFIX + SECTION_NAME + "/{id}", deleteSnippet).Methods("DELETE")
	router.HandleFunc(API_PREFIX + SECTION_NAME + "/{id}", returnSingleSnippet)

	// All static assets
	router.
	PathPrefix(STATIC_DIR).
	Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("." + STATIC_DIR))))

	// All static pages
	router.
	PathPrefix("/").
	Handler(http.StripPrefix("/", http.FileServer(http.Dir("." + PAGES_DIR))))

	log.Fatal(http.ListenAndServe(":" + PORT, router))
}

const (
	PAGES_DIR = "/pages/"
	STATIC_DIR = "/static/"
	PORT = "10000"
	DATA = "./DATA"
	API_PREFIX = "/api/v1/"
	DATABASE_SNIPPET = "snippet"
)

var db, err = scribble.New(DATA, nil)

func main() {

	fmt.Println("Rest API v1.0 - Dan's Mux Router test")
	handleRequests()
}