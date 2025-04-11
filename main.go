package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Flagsmith/flagsmith-go-client/v4"
	"github.com/joho/godotenv"
)

var client *flagsmith.Client

type Flag struct {
	Name    string      `json:"name"`
	Enabled bool        `json:"enabled"`
	Value   interface{} `json:"value"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	const exampleFlag = "my_go_flag"
	flags, err := client.GetFlags(r.Context(), nil)
	if err != nil {
		log.Printf("Error getting flags: %v", err)
		http.Error(w, "Error getting flags", http.StatusInternalServerError)
		return
	}
	testFlag, _ := flags.GetFlag(exampleFlag)
	fmt.Printf("Test flag: %v\n", testFlag)
	allFlags := flags.AllFlags()
	flagsMap := make(map[string]Flag)
	for _, f := range allFlags {
		flagsMap[f.FeatureName] = Flag{
			Enabled: f.Enabled,
			Value:   f.Value,
			Name:    f.FeatureName,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(flagsMap); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func main() {
	_ = godotenv.Load()

	apiKey := os.Getenv("FLAGSMITH_API_KEY")
	// baseURL := os.Getenv("FLAGSMITH_BASE_URL")
	if apiKey == "" {
		log.Fatal("FLAGSMITH_API_KEY not set")
	}

	client = flagsmith.NewClient(apiKey)

	http.HandleFunc("/", handler)
	log.Println("Listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
