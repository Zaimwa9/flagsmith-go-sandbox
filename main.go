package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/Flagsmith/flagsmith-go-client/v4"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

var client *flagsmith.Client

type Flag struct {
	Name    string      `json:"name"`
	Enabled bool        `json:"enabled"`
	Value   interface{} `json:"value"`
}

type CustomRoundTripper struct{}

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
	rt := CustomRoundTripper{}

	customHttpClient := &http.Client{
		Transport: rt,
	}

	fmt.Println(customHttpClient.Transport)

	restyClient := resty.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	restyClient.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		logger.Info("HTTP request",
			slog.String("method", req.Method),
			slog.String("url", req.URL),
			slog.Any("headers", req.Header),
		)
		return nil
	})

	restyClient.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		logger.Info("HTTP response",
			slog.String("url", resp.Request.URL),
			slog.Int("status", resp.StatusCode()),
			slog.Int("bytes", len(resp.Body())),
			slog.Duration("elapsed", resp.Time()),
		)
		return nil
	})

	apiKey := os.Getenv("FLAGSMITH_API_KEY")
	// baseURL := os.Getenv("FLAGSMITH_BASE_URL")
	if apiKey == "" {
		log.Fatal("FLAGSMITH_API_KEY not set")
	}

	client = flagsmith.NewClient(apiKey, flagsmith.WithRestyClient(restyClient))

	http.HandleFunc("/", handler)
	log.Println("Listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func (ct CustomRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	log.Printf("url: %s | method: %s", resp.Request.URL.String(), resp.Request.Method)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()

	log.Println("number of bytes in tripper:", len(data))
	log.Printf("value in tripper: %s\n", string(data))

	resp.Body = io.NopCloser(bytes.NewReader(data))

	return resp, nil
}
