package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
	  "github.com/sdomino/scribble"
)
// a fish
type Fish struct{ Name string }

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnAllArticles")
    records, err := db.ReadAll("fish")
  	if err != nil {
  		fmt.Println("Error", err)
  	}
	  fishies := []Fish{}
	  for _, f := range records {
	  	fishFound := Fish{}
		  if err := json.Unmarshal([]byte(f), &fishFound); err != nil {
			  fmt.Println("Error", err)
		  }
		  fishies = append(fishies, fishFound)
	  }
    json.NewEncoder(w).Encode(fishies)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["id"]

    // Loop over all of our Articles
    // if the article.Id equals the key we pass in
    // return the article encoded as JSON
  	onefish := Fish{}
  	if err := db.Read("fish", key, &onefish); err != nil {
  		fmt.Println("Error", err)
  	}
  	
  	json.NewEncoder(w).Encode(onefish)
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Create New")
  	b, err := ioutil.ReadAll(r.Body)
  	defer r.Body.Close()

  	var fish Fish
  	err = json.Unmarshal(b, &fish)
  	if err != nil {
  		http.Error(w, err.Error(), 500)
  		return
  	}

    db.Write("fish", fish.Name, fish)
  
  	output, err := json.Marshal(fish)
  	if err != nil {
  		http.Error(w, err.Error(), 500)
  		return
  	}
  	w.Header().Set("content-type", "application/json")
  	w.Write(output)
}


func deleteArticle(w http.ResponseWriter, r *http.Request) {
    // once again, we will need to parse the path parameters
    vars := mux.Vars(r)
    // we will need to extract the `id` of the article we
    // wish to delete
    id := vars["id"]
    fmt.Println("Removing: " + string("fish") + "/" + string(id))
  	// Delete a fish from the database
  	if err := db.Delete("fish", string(id)); err != nil {
  		fmt.Println("Error", err)
  	}
    fmt.Fprintf(w, "%+v", "Removed: " + string("fish") + "/" + string(id))
}


func handleRequests() {
        
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/articles", returnAllArticles)
    // NOTE: Ordering is important here! This has to be defined before
    // the other `/article` endpoint.
    myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
    myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
    myRouter.HandleFunc("/article/{id}", returnSingleArticle)
    log.Fatal(http.ListenAndServe(":10000", myRouter))
    
}

var db, err  = scribble.New("./", nil)

func main() {
  
    fmt.Println("Rest API v2.0 - Mux Routers")
    handleRequests()
}