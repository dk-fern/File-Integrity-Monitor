package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getIPReport(apiKey string, ipAddr string) (*IP, error) {
	// Make GET request to api endpoint
	url := fmt.Sprintf("https://www.virustotal.com/api/v3/ip_addresses/%v", ipAddr)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error in GET request: %v", err)
	}

	// Add headers
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-apikey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("the client encountered an error making the request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the response body: %v", err)
	}

	// Init ip variable to marshall the json data to our struct
	var ip IP

	if err := json.Unmarshal(body, &ip); err != nil {
		return nil, fmt.Errorf("error unmarshaling to json: %v", err)
	}

	return &ip, nil

}

func getDomainReport(apiKey string, domainToSearch string) (*Domain, error) {
	// Make GET request to api endpoint
	url := fmt.Sprintf("https://www.virustotal.com/api/v3/domains/%v", domainToSearch)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error in GET request: %v", err)
	}

	// Add headers
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-apikey", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("the client encountered an error making the request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the response body: %v", err)
	}

	// Init ip variable to marshall the json data to our struct
	var domain Domain

	if err := json.Unmarshal(body, &domain); err != nil {
		return nil, fmt.Errorf("error unmarshaling to json: %v", err)
	}

	return &domain, nil

}
