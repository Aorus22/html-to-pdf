package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type FirestoreResponse struct {
	Document struct {
		Fields struct {
			ID struct {
				StringValue string `json:"stringValue"`
			} `json:"id"`
		} `json:"fields"`
		UpdateTime string `json:"updateTime"`
	} `json:"document"`
}

type ListMapelResponse struct {
	MapelList []string `json:"mapelList"`
}

func ListMapelHandler(w http.ResponseWriter, r *http.Request) {
	mapelList, err := getMapelList()
	if err != nil {
		log.Printf("Error saat mengambil soal: %v", err)
		http.Error(w, "Terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListMapelResponse{MapelList: mapelList})
}

func getMapelList() ([]string, error) {
	const firestoreURL = "https://firestore.googleapis.com/v1/projects/cbt-01-a21ba/databases/(default)/documents/148/soal:runQuery"

	queryBody := map[string]any{
		"structuredQuery": map[string]any{
			"from": []map[string]string{
				{"collectionId": "soal"},
			},
			"select": map[string][]map[string]string{
				"fields": {
					{"fieldPath": "id"},
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(queryBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", firestoreURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []FirestoreResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	now := time.Now()
	minDate := now.AddDate(0, 0, -30)
	var mapelList []string

	for _, item := range result {
		doc := item.Document
		if doc.UpdateTime == "" {
			continue
		}

		updateTime, err := time.Parse(time.RFC3339Nano, doc.UpdateTime)
		if err != nil {
			continue
		}

		if updateTime.Before(minDate) || updateTime.After(now) {
			continue
		}

		mapelList = append(mapelList, doc.Fields.ID.StringValue)
	}

	return mapelList, nil
}
