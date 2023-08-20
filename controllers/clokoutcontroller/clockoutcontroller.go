package clockoutcontroller

import (
	"attendance/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/carlogit/phash"
	"gorm.io/gorm"
)

func Lockout(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the uploaded file
	filePath := "./logout_uploads/" + handler.Filename
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

	id_employee_val := r.FormValue("id_employee")
	id_employee, err := strconv.Atoi(id_employee_val)
	if err != nil {
		id_employee = 0
	}

	var user models.Employee
	if err := models.DB.Where("Id = ? and Islogin = ?", id_employee, 1).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			w.WriteHeader(http.StatusOK) // Set the response status code to 200
			fmt.Fprintln(w, "Employee locked out!")
			return
		default:
			http.Error(w, "Something wrong with the process", http.StatusInternalServerError)
			return
		}
	}

	fileSourcePath := user.Photo

	if CheckIsPhotoSimilar(fileSourcePath, filePath) {

		currentTime := time.Now()
		lockoutTime := currentTime.Format("YYYY-MM-DD HH:MM:SS")

		var userInput = models.Lockout{
			Datetime_lockout: lockoutTime,
			Id_employee:      id_employee,
			Image:            filePath,
		}

		if err := models.DB.Create(&userInput).Error; err != nil {
			http.Error(w, "Error create lockout record", http.StatusInternalServerError)
			return
		}

		user.Islogin = 0
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
