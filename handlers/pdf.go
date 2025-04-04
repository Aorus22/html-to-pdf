package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"get-pdf/utils"

	"github.com/gorilla/mux"
)

func GetOriginalSoalPdfHandler(w http.ResponseWriter, r *http.Request) {
	mapel := mux.Vars(r)["mapel"]

	datasoal, datakunci, err := getSoal(mapel)
	if err != nil || datasoal == nil {
		http.Error(w, "Soal tidak ditemukan", http.StatusNotFound)
		return
	}

	html, err := utils.RenderTemplateToString(filepath.Join("views", "template.html"), datasoal, datakunci)
	if err != nil {
		http.Error(w, "Gagal render template", http.StatusInternalServerError)
		return
	}

	pdf, err := utils.GeneratePDF(html)
	if err != nil {
		http.Error(w, "Gagal generate PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, mapel))
	w.Write(pdf)
}

func GetCustomSoalPdfHandler(w http.ResponseWriter, r *http.Request) {
	mapel := mux.Vars(r)["mapel"]

	body, _ := io.ReadAll(r.Body)
	var bodyData struct {
		Datakunci map[string]any `json:"datakunci"`
	}
	if err := json.Unmarshal(body, &bodyData); err != nil || bodyData.Datakunci["kunci"] == nil {
		http.Error(w, "Format datakunci tidak valid", http.StatusBadRequest)
		return
	}

	datasoal, _, err := getSoal(mapel)
	if err != nil || datasoal == nil {
		http.Error(w, "Soal tidak ditemukan", http.StatusNotFound)
		return
	}

	html, err := utils.RenderTemplateToString(filepath.Join("views", "template.html"), datasoal, bodyData.Datakunci)
	if err != nil {
		http.Error(w, "Gagal render template", http.StatusInternalServerError)
		return
	}

	pdf, err := utils.GeneratePDF(html)
	if err != nil {
		http.Error(w, "Gagal generate PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, mapel))
	w.Write(pdf)
}
