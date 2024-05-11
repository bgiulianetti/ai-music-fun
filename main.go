package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	URL                    = "https://api.replicate.com/v1/predictions"
	PORT                   = "8080"
	API_TOKEN              = "{token}"
	AUTH_TOKEN             = "Bearer " + API_TOKEN
	VERSION                = "671ac645ce5e552cc63a54a2bbff63fcf798043055d2dac5fc9e36a837eedcfb"
	MODEL_VERSION          = "stereo-large"
	OUTPUT_FORMAT          = "mp3"
	NORMALIZATION_STRATEGY = "peak"
)

type Input struct {
	Prompt                string `json:"prompt"`
	ModelVersion          string `json:"model_version"`
	OutputFormat          string `json:"output_format"`
	NormalizationStrategy string `json:"normalization_strategy"`
}

type RequestBody struct {
	Version string `json:"version"`
	Input   Input  `json:"input"`
}

type InnerAPIResponse struct {
	ID        string `json:"id"`
	Model     string `json:"model"`
	Version   string `json:"version"`
	Input     Input  `json:"input"`
	Logs      string `json:"logs"`
	Error     string `json:"error"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	URLs      URLs   `json:"urls"`
}

type URLs struct {
	Cancel string `json:"cancel"`
	Get    string `json:"get"`
}

type Response struct {
	Message string `json:"message"`
	Song    string `json:"song"`
	Input   Input  `json:"input"`
	Status  string `json:"status"`
	Created string `json:"created_at"`
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	requestBody.Version = VERSION
	requestBody.Input.ModelVersion = MODEL_VERSION
	requestBody.Input.OutputFormat = OUTPUT_FORMAT
	requestBody.Input.NormalizationStrategy = NORMALIZATION_STRATEGY

	requestBodyJSON, _ := json.Marshal(requestBody)

	log.Println(string(requestBodyJSON))

	apiRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		http.Error(w, "Error preparing API request", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(apiRequestBody))
	if err != nil {
		http.Error(w, "Error creating HTTP request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", AUTH_TOKEN)
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error calling prediction API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading API response", http.StatusInternalServerError)
		return
	}

	log.Printf("Prediction API Response: %s\n", respBody)

	var innerResponse InnerAPIResponse
	err = json.Unmarshal(respBody, &innerResponse)
	if err != nil {
		http.Error(w, "Error parsing API response", http.StatusInternalServerError)
		return
	}

	transformedResponse := Response{
		Message: "Your song is being processed, get the song with the link below",
		Song:    fmt.Sprintf("http://localhost:%s/song/%s", PORT, innerResponse.ID),
		Input:   innerResponse.Input,
		Status:  innerResponse.Status,
		Created: innerResponse.CreatedAt,
	}

	transformedRespJSON, err := json.Marshal(transformedResponse)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(transformedRespJSON)
}

func songHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/song/")
	if id == "" {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", URL, id), nil)
	if err != nil {
		http.Error(w, "Error creating HTTP request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", AUTH_TOKEN)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error calling prediction API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading API response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func main() {
	http.HandleFunc("/song", createHandler)
	http.HandleFunc("/song/", songHandler)

	log.Printf("Server started on port %s", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil))
}
