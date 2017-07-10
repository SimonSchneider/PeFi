package main

import (
	"net/http"
)

//var templates = template.Must(template.ParseFiles("static/daily.html"))

type (
	M map[string]interface{}

	Page struct {
		Title string
		Body  []byte
	}
)

//func dailyHandel(w http.ResponseWriter, r *http.Request) {
//as, err := getAccounts()
//if err != nil {
//fmt.Println(err)
//return
//}

//err = templates.ExecuteTemplate(w, "daily.html", M{
//"ExternalAccounts": as,
//"InternalAccounts": as,
//})
//if err != nil {
//fmt.Println(err)
//return
//}
//}

func main() {
	//http.HandleFunc("/daily", dailyHandel)
	http.Handle("/", http.FileServer(http.Dir("web/static/")))

	http.ListenAndServe(":8080", nil)
}
