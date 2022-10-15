package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Запуск сервера
func runServer(port int) error {
	http.HandleFunc("/", Page)

	log.Println("Starting http server:", fmt.Sprintf("http://127.0.0.1:%v/", port))
	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

// Отдаем html-страницу с параметрами
func ParseAndExecute(w http.ResponseWriter, req *http.Request, str string, params interface{}) {
	t, err := template.ParseFiles(str)
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if err := t.Execute(w, params); err != nil {
		fmt.Printf("Error: %v", err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func Page(w http.ResponseWriter, req *http.Request) {
	var tempId Data
	var ok bool

	temp := dataCache()                              // Получение указателя на кэш
	tempId, ok = temp.Get(req.FormValue("orderUid")) // Получение uid заказа

	// Вывод на страницу
	type Param struct {
		Ok   bool
		Info Data
	}
	ParseAndExecute(w, req, "html/page.html", Param{Ok: ok, Info: tempId})
}
