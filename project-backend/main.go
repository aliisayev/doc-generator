package main

import (
	"log"
	"net/http"
)

func main() {
	InitDB()

	http.HandleFunc("/submit", SubmitHandler)
	http.HandleFunc("/admin/login", AdminLoginHandler)
	http.HandleFunc("/admin", AdminPanelHandler)
	http.HandleFunc("/admin/contracts/", AdminViewContractHandler)
	http.HandleFunc("/admin/send/", AdminSendHandler)

	http.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("frontend"))))
	http.Handle("/output/", http.StripPrefix("/output/", http.FileServer(http.Dir("output"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	log.Println("✅ Сервер запущен на http://localhost:8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal("❌ Ошибка запуска сервера:", err)
	}
}
