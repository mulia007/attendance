package registercontroller

import (
	"attendance/models"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name_employee")

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the uploaded file
	filePath := "./profile_upload/" + handler.Filename
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "Error copying file content", http.StatusInternalServerError)
		return
	}

	var userInput = models.Employee{
		Name_employee: name,
		Photo:         filePath,
		Islogin:       0,
	}

	if err := models.DB.Create(&userInput).Error; err != nil {
		http.Error(w, "Error saving employee data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) // Set the response status code to 200
	fmt.Fprintln(w, "Success!")  // Write a response message
}
