package main

import (
	"flag"
	"log"
	"net/http"
	"mime"

	"get-pdf/handlers"
	"get-pdf/utils"
	"github.com/gorilla/mux"
)

func main() {
	downloadOnly := flag.Bool("download-only", false, "Download Rod Browser")
	flag.Parse()

	if *downloadOnly {
		log.Println("Downloading Rod Browser...")
		html := "<html><body><h1>Dummy HTML to trigger Rod browser download</h1></body></html>"
		_, err := utils.GeneratePDF(html)
		if err != nil {
			log.Fatalf("Failed to download browser: %v", err)
		}
		log.Println("Browser downloaded successfully")
		return
	}

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./views"))
	r.PathPrefix("/views/").Handler(http.StripPrefix("/views/", fs))

	mime.AddExtensionType(".ttf", "font/ttf")
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	r.HandleFunc("/listmapel", handlers.ListMapelHandler).Methods("GET")
	r.HandleFunc("/pdf/{mapel}", handlers.GetOriginalSoalPdfHandler).Methods("GET")
	r.HandleFunc("/pdf/{mapel}", handlers.GetCustomSoalPdfHandler).Methods("POST")
	r.HandleFunc("/json/{mapel}", handlers.GetSoalJsonHandler).Methods("GET")
	r.HandleFunc("/soal/{mapel}", handlers.GetOriginalSoalHandler).Methods("GET")
	r.HandleFunc("/soal/{mapel}", handlers.GetCustomSoalHandler).Methods("POST")

	port := ":3000"
	log.Printf("Server is running on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
