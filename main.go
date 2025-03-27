package main

import (
    "fmt"
    "log"
    "net/http"
    "math/rand"
    "sync"
    
    "github.com/gorilla/mux"
)

var urlStore = struct {
    sync.RWMutex
    data    map[string]string
}{data: make(map[string]string)}

const shortURLLength = 6
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main(){
    router := mux.NewRouter()
    router.HandleFunc("/{shortcode}", redirectHandler).Methods("GET")
    router.HandleFunc("/shorten", shortenURLHandler).Methods("POST")

    fmt.Println("Listening on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func redirectHandler(w http.ResponseWriter, r *http.Request){
    var vars = mux.Vars(r)
    shortURL := vars["shortcode"]

    urlStore.RLock()
    longURL, exists := urlStore.data[shortURL]
    urlStore.RUnlock()
    if !exists {
        http.Error(w,"URL not found.",http.StatusNotFound)
        return
    }

    fmt.Fprintf(w,"%v\n", longURL)
}

func generateShortCode() string{
    shortURL := make([]byte, shortURLLength)
    for i := 0;i < len(shortURL); i++ {
       shortURL[i] = charset[rand.Intn(len(charset))]
    }
    return string(shortURL)
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
    longURL := r.URL.Query().Get("url")

    shortURL := generateShortCode() 
    urlStore.Lock()
    urlStore.data[shortURL] = longURL
    urlStore.Unlock()

    fmt.Fprintf(w, "Short URL: http://localhost:8080/%s\n", shortURL)
}
