package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestBadRoute(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/badroute")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}
func TestAuthorization(t *testing.T) {

	reqBody := map[string]interface{}{
		"retailer":     "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": []map[string]interface{}{
			{"shortDescription": "Mountain Dew 12PK", "price": "6.49"},
			{"shortDescription": "Emils Cheese Pizza", "price": "12.25"},
			{"shortDescription": "Knorr Creamy Chicken", "price": "1.26"},
			{"shortDescription": "Doritos Nacho Cheese", "price": "3.35"},
			{"shortDescription": "Klarbrunn 12-PK 12 FL OZ", "price": "12.00"},
		},
		"total": "35.35",
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post("http://localhost:8080/receipts/process", "application/json", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestProcessAndGetPointsEx1(t *testing.T) {
	// test cases with request body and expected points
	testCases := []struct {
		ReqBody        map[string]interface{}
		ExpectedPoints int
	}{
		{
			ReqBody: map[string]interface{}{
				"retailer":     "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": []map[string]interface{}{
					{"shortDescription": "Mountain Dew 12PK", "price": "6.49"},
					{"shortDescription": "Emils Cheese Pizza", "price": "12.25"},
					{"shortDescription": "Knorr Creamy Chicken", "price": "1.26"},
					{"shortDescription": "Doritos Nacho Cheese", "price": "3.35"},
					{"shortDescription": "Klarbrunn 12-PK 12 FL OZ", "price": "12.00"},
				},
				"total": "35.35",
			},
			ExpectedPoints: 28,
		},
		{
			ReqBody: map[string]interface{}{
				"retailer":     "M&M Corner Market",
				"purchaseDate": "2022-03-20",
				"purchaseTime": "14:33",
				"items": []map[string]interface{}{
					{"shortDescription": "Gatorade", "price": "2.25"},
					{"shortDescription": "Gatorade", "price": "2.25"},
					{"shortDescription": "Gatorade", "price": "2.25"},
					{"shortDescription": "Gatorade", "price": "2.25"},
				},
				"total": "9.00",
			},
			ExpectedPoints: 109,
		},
	}

	// Run tests for each test case
	for _, tc := range testCases {
		err := sendRequestAndGetPoints(tc.ReqBody, tc.ExpectedPoints)
		if err != nil {
			t.Fatal(err)
		}

	}
}

// Helper function to validate the processing and getting points for a receipt
func sendRequestAndGetPoints(reqBody map[string]interface{}, expectedPoints int) error {
	var response struct {
		ID string `json:"id"`
	}
	var pointsResponse struct {
		Points int `json:"points"`
	}

	// Build request to process receipt and get points
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/receipts/process", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return err
	}

	// generate a hash of the API key and set the header
	// simplified for this project to use a predictable key
	hash := sha256.Sum256([]byte("key1"))
	hashedKey := hex.EncodeToString(hash[:])
	req.Header.Set("Authorization", hashedKey)
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	// Build request to get points for receipt
	req, err = http.NewRequest("GET", "http://localhost:8080/receipts/"+response.ID+"/points", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", hashedKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&pointsResponse)
	if err != nil {
		return err
	}

	if pointsResponse.Points != expectedPoints {
		return fmt.Errorf("Expected points %d, got %d", expectedPoints, pointsResponse.Points)
	}
	return nil
}
