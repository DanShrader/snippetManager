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


func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: homePage")
    // fmt.Fprintf(w, "Welcome to the HomePage!")
    // http.FileServer(http.Dir("./"))
    // http.StripPrefix("/", http.FileServer(http.Dir("./pages/index.html")))
    // http.ServeFile(w, r, r.URL.Path[1:])
    http.ServeFile(w, r,"./pages/index.html")
    // Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))
}


func handleRequests() {
        
    router := mux.NewRouter().StrictSlash(true)
    
    // All Statis assets
    router.
      PathPrefix(STATIC_DIR).
      Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))
    
    
    
    
    
    router.HandleFunc("/", homePage)
    router.HandleFunc("/articles", returnAllArticles)
    // NOTE: Ordering is important here!
    router.HandleFunc("/article", createNewArticle).Methods("POST")
    router.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
    router.HandleFunc("/article/{id}", returnSingleArticle)
    log.Fatal(http.ListenAndServe(":" + PORT, router))
    
}

const (
    STATIC_DIR = "/static/"
    PORT       = "10000"
    DATA       = "./"
)

var db, err  = scribble.New(DATA, nil)

func main() {
  
    fmt.Println("Rest API v1.0 - Dan's Mux Router test")
    handleRequests()
}