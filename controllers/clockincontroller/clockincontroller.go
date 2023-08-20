package clockincontroller

import (
	"attendance/models"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/carlogit/phash"
	"gorm.io/gorm"
)

func Lockin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	id_employee_val := r.FormValue("id_employee")
	id_employee, err := strconv.Atoi(id_employee_val)
	if err != nil {
		id_employee = 0
	}

	currentTime := time.Now()
	formattedTime := currentTime.Format("YYYY-MM-DD HH:MM:SS")

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the uploaded file
	filePath := "./login_uploads/" + handler.Filename
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

	var user models.Employee
	if err := models.DB.Where("Id = ?", id_employee).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			http.Error(w, "Error employee not found", http.StatusInternalServerError)
			return
		default:
			http.Error(w, "Something wrong with the process", http.StatusInternalServerError)
			return
		}
	}

	fileSourcePath := user.Photo

	if CheckIsPhotoSimilar(fileSourcePath, filePath) {
		var userInput = models.Lockon{
			Id_employee:    id_employee,
			Datetime_logon: formattedTime,
			Image:          filePath,
		}

		if err := models.DB.Create(&userInput).Error; err != nil {
			http.Error(w, "Error while saving lockin record", http.StatusInternalServerError)
			return
		}

		user.Islogin = 1
		result := models.DB.Save(&user)
		if result.Error != nil {
			http.Error(w, "Error updating employee status", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK) // Set the response status code to 200
		fmt.Fprintln(w, "Employee lockin successfully!")

	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Println("Photo is not match. Please upload your profile")
	}
}

func CheckIsPhotoSimilar(fileSourceUrl string, fileUploadUrl string) bool {

	a := hash(fileSourceUrl)
	b := hash(fileUploadUrl)
	distance := phash.GetDistance(a, b)

	if distance <= 20 {
		return true
	} else {
		return false
	}
}

func hash(filename string) string {
	img, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer img.Close()

	ahash, err := phash.GetHash(img)
	if err != nil {
		log.Fatal(err)
	}
	return ahash
}
