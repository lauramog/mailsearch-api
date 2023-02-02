package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	client "github.com/zinclabs/sdk-go-zincsearch"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	addr := flag.String("addr", ":8081", "HTTP network address")
	flag.Parse()
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Get("/search", searchHandler)

	err := http.ListenAndServe(*addr, router)
	if err != nil {
		log.Fatal("can not start the server", err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	term := r.URL.Query().Get("term")
	if strings.TrimSpace(term) == "" {
		http.Error(w, "missing search term", http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(context.Background(), client.ContextBasicAuth, client.BasicAuth{
		UserName: os.Getenv("UserName"),
		Password: os.Getenv("Password"),
	})
	configuration := client.NewConfiguration()
	configuration.Servers = client.ServerConfigurations{
		client.ServerConfiguration{
			URL: "http://localhost:4080",
		},
	}
	apiClient := client.NewAPIClient(configuration)

	index := "inbox"
	query := *client.NewV1ZincQuery()
	query.SetSearchType("match")
	params := *client.NewV1QueryParams()
	params.SetTerm(term)
	query.SetQuery(params)
	resp, _, err := apiClient.Search.SearchV1(ctx, index).Query(query).Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var sources []map[string]interface{}
	for _, hit := range resp.GetHits().Hits {
		sources = append(sources, hit.GetSource())
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	if err = enc.Encode(sources); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
