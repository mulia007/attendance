package main

import (
	"attendance/controllers/clockincontroller"
	clockoutcontroller "attendance/controllers/clokoutcontroller"
	"attendance/controllers/registercontroller"
	"net/http"
)

func main() {
	http.HandleFunc("/", clockincontroller.Lockin)
	http.HandleFunc("/Register", registercontroller.Insert)
	http.HandleFunc("/Lockin", clockincontroller.Lockin)
	http.HandleFunc("/Lockout", clockoutcontroller.Lockout)

	http.ListenAndServe(":8080", nil)
}
