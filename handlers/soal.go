package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"path/filepath"

	"get-pdf/utils"

	"github.com/gorilla/mux"
)

func removeFontFamilyFromHTML(html string) string {
	return strings.ReplaceAll(html, "font-family:", "")
}

func getSoal(mapel string) (map[string]any, map[string]any, error) {
	soalPath := fmt.Sprintf("148/soal/soal/%s", mapel)
	soalDoc, err := utils.GetDocument(soalPath)
	if err != nil {
		return nil, nil, err
	}

	fields, ok := soalDoc["fields"].(map[string]any)
	if !ok {
		return nil, nil, fmt.Errorf("invalid soal document format")
	}

	parsed := utils.ParseFirestoreFields(fields)

	if soalList, ok := parsed["soal"].([]any); ok {
		for _, item := range soalList {
			if soalMap, ok := item.(map[string]any); ok {
				if s, ok := soalMap["soal"].(string); ok {
					soalMap["soal"] = removeFontFamilyFromHTML(s)
				}
				if opts, ok := soalMap["options"].([]any); ok {
					for _, opt := range opts {
						if optMap, ok := opt.(map[string]any); ok {
							if ans, ok := optMap["answer"].(string); ok {
								optMap["answer"] = removeFontFamilyFromHTML(ans)
							}
						}
					}
				}
			}
		}
	}

	name, _ := soalDoc["name"].(string)
	id := mapel
	if name != "" {
		id = name[strings.LastIndex(name, "/")+1:]
	}
	parsed["id"] = id

	kunciPath := fmt.Sprintf("148/datakunci/datakunci/%s", mapel)
	kunciDoc, err := utils.GetDocument(kunciPath)
	if err != nil {
		return parsed, map[string]any{}, nil
	}

	kunciFields, ok := kunciDoc["fields"].(map[string]any)
	if !ok {
		return parsed, map[string]any{}, nil
	}

	datakunci := utils.ParseFirestoreFields(kunciFields)

	return parsed, datakunci, nil
}

func GetSoalJsonHandler(w http.ResponseWriter, r *http.Request) {
	mapel := mux.Vars(r)["mapel"]

	datasoal, datakunci, err := getSoal(mapel)
	if err != nil {
		http.Error(w, "Soal tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"datasoal":  datasoal,
		"datakunci": datakunci,
	})
}

func GetOriginalSoalHandler(w http.ResponseWriter, r *http.Request) {
	mapel := mux.Vars(r)["mapel"]

	datasoal, datakunci, err := getSoal(mapel)
	if err != nil {
		http.Error(w, "Soal tidak ditemukan", http.StatusNotFound)
		return
	}

	html, err := utils.RenderTemplateToString(filepath.Join("views", "template.html"), datasoal, datakunci)
	if err != nil {
		http.Error(w, "Gagal render template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func GetCustomSoalHandler(w http.ResponseWriter, r *http.Request) {
	mapel := mux.Vars(r)["mapel"]

	var input struct {
		Datakunci map[string]any `json:"datakunci"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Datakunci["kunci"] == nil {
		http.Error(w, "Format datakunci tidak valid", http.StatusBadRequest)
		return
	}

	datasoal, _, err := getSoal(mapel)
	if err != nil {
		http.Error(w, "Soal tidak ditemukan", http.StatusNotFound)
		return
	}

	html, err := utils.RenderTemplateToString(filepath.Join("views", "template.html"), datasoal, input.Datakunci)
	if err != nil {
		http.Error(w, "Gagal render template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
