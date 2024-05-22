package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/emeraldls/platnova-task/internal/generator"
	"github.com/emeraldls/platnova-task/internal/types"
	"github.com/joho/godotenv"
	"github.com/unidoc/unipdf/v3/creator"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	slog.Info("File Upload Endpoint Hit")

	r.ParseMultipartForm(3 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		slog.Error("error retrieving file", "err", err)
		return
	}
	defer file.Close()

	fb, err := io.ReadAll(file)
	if err != nil {
		slog.Error("error reading file", "err", err)
		return
	}

	if handler.Filename[len(handler.Filename)-5:] != ".json" {
		http.Error(w, "file must be a json file", http.StatusBadRequest)
		return
	}

	var accStmt types.AccountStatement
	err = json.Unmarshal(fb, &accStmt)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing json: %v", err), http.StatusInternalServerError)
		return
	}

	apiKey := os.Getenv("apiKey")

	c := types.NewClient(creator.New(), apiKey)
	fn, err := generator.GenerateAccountStatementPDF(*c, accStmt)
	if err != nil {
		http.Error(w, fmt.Sprintf("error generating document: %v", err), http.StatusInternalServerError)
		return
	}

	pdfFile, err := os.Open(fn + ".pdf")
	if err != nil {
		http.Error(w, fmt.Sprintf("error opening generated document: %v", err), http.StatusInternalServerError)
		return
	}
	defer pdfFile.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fn+".pdf")
	w.Header().Set("Content-Type", "application/pdf")

	_, err = io.Copy(w, pdfFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("error sending file: %v", err), http.StatusInternalServerError)
		return
	}

	go func() {
		time.Sleep(5 * time.Second)
		err = os.Remove(fn + ".pdf")
		if err != nil {
			slog.Error("error deleting file", "err", err)
		}
	}()

}

func SetupRoutes() {

	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("error loading .env file", "err", err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("POST /upload", uploadFile)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "2222"

	}

	slog.Info("Starting server", "address", fmt.Sprintf("0.0.0.0:%s", PORT))

	err = http.ListenAndServe(":"+PORT, nil)

	if err != nil {
		slog.Error("error starting server", "err", err)
		return
	}

}
