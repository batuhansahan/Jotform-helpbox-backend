package main

import (
	"encoding/json"
	"fmt"
    "log"
    "net/http"
    "strings"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
    "github.com/PuerkitoBio/goquery"
)


func main() {
    var router = mux.NewRouter()
    router.Use(commonMiddleware)
	router.HandleFunc("/healthcheck", healthCheck).Methods("GET")
    router.HandleFunc("/headings", handleHeadings).Methods("GET")

	headersOk := handlers.AllowedHeaders([]string{"Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	fmt.Println("Running server!")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
func commonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        next.ServeHTTP(w, r)
    })
}

func handleHeadings(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	searchQuery := vars.Get("query")
    searchQuery = strings.Replace(searchQuery, " ", "+", -1)
    var headings []string
    var links []string
    out := map[string]interface{}{}
    // var summary []string
    // var links []string
    

    response, err := http.Get("https://www.jotform.com/help/keyword_search.php?rpp=0&search="+searchQuery)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    // Create a goquery document from the HTTP response
    document, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        log.Fatal("Error loading HTTP response body. ", err)
    }

    document.Find(".chapterTitle").Each(func(i int, s *goquery.Selection) {
        // class, _ := s.Attr("class")
        href, _ := s.Attr("href")
        headings = append(headings,s.Text())
        links = append(links,href)
        // fmt.Println(s.Text())
        
    })

    out["headings"] = headings;
    out["links"] = links;

    
	json.NewEncoder(w).Encode(out)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}

//https://www.jotform.com/help/keyword_search.php?rpp=0&search=conditions