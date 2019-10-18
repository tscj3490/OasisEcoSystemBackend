package server

import (
	"com/merkinsio/oasis-api/config"
	"com/merkinsio/oasis-api/services"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

//UploadFile Uploads a file to the server and returns its data
func UploadFile(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	uploadedFile, headers, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, file, err := services.CreateFile(headers.Filename, uploadedFile)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, *file)
}

//DownloadFile Downloads the file
func DownloadFile(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)

	// fileName := vars["fileName"]
	fileID := vars["fileId"]

	filePath := fmt.Sprintf("%s/%s", config.Config.GetString("uploads.path"), fileID)
	f, err := os.Open(filePath)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "File not found")
	} else {
		defer f.Close()

		io.Copy(w, f)
	}
}
