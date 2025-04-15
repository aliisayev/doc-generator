package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	// ✅ в
	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

var adminUsername = "admin"
var adminPassword = "1234"

// ===== 📨 Пользователь отправляет анкету =====

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST", http.StatusMethodNotAllowed)
		return
	}

	var data FormData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "JSON decode error", http.StatusBadRequest)
		return
	}

	jsonData, _ := json.Marshal(data)

	var photoPath string
	if data.Photo != "" {
		imgData, _ := base64.StdEncoding.DecodeString(strings.Split(data.Photo, ",")[1])
		photoPath = fmt.Sprintf("uploads/%d.jpg", time.Now().UnixNano())
		_ = os.WriteFile(photoPath, imgData, 0644)
	}

	zipName := fmt.Sprintf("output/contract_%d.zip", time.Now().UnixNano())
	err := generateZIP(data, zipName)
	if err != nil {
		http.Error(w, "ZIP error", 500)
		return
	}

	_, err = DB.Exec("INSERT INTO contracts (json_data, pdf_path, photo_path) VALUES (?, ?, ?)",
		string(jsonData), zipName, photoPath)
	if err != nil {
		http.Error(w, "DB insert error", 500)
		return
	}

	w.Write([]byte("OK"))
}

func generateZIP(data FormData, outPath string) error {
	data.Today = time.Now().Format("02.01.2006")

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	files, _ := os.ReadDir("templates")
	for _, f := range files {
		if !strings.HasPrefix(f.Name(), "contract_") || !strings.HasSuffix(f.Name(), ".html") {
			continue
		}

		tmpl, err := template.New(f.Name()).Funcs(template.FuncMap{
			"indexSafe": func(s []string, i int) string {
				if i >= 0 && i < len(s) {
					return s[i]
				}
				return ""
			},
			"genderWord": func(g string) string {
				if g == "Kişi" {
					return " oğlu"
				}
				return " qızı"
			},
		}).ParseFiles("templates/" + f.Name())
		if err != nil {
			continue
		}

		var filled bytes.Buffer
		_ = tmpl.Execute(&filled, data)

		pdfg, _ := wkhtml.NewPDFGenerator()
		page := wkhtml.NewPageReader(&filled)
		page.EnableLocalFileAccess.Set(true)
		pdfg.AddPage(page)
		pdfg.PageSize.Set(wkhtml.PageSizeA4)
		_ = pdfg.Create()

		writer, _ := zipWriter.Create(strings.TrimSuffix(f.Name(), ".html") + ".pdf")
		_, _ = writer.Write(pdfg.Bytes())
	}

	return nil
}

// ===== 🔐 Админка: Авторизация и проверка =====

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := template.ParseFiles("admin_templates/login.html")
		t.Execute(w, nil)
		return
	}
	r.ParseForm()
	if r.Form.Get("username") == adminUsername && r.Form.Get("password") == adminPassword {
		http.SetCookie(w, &http.Cookie{Name: "admin", Value: "true"})
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

func requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("admin")
	if err != nil || c.Value != "true" {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return false
	}
	return true
}

// ===== 📋 Список договоров =====

func AdminPanelHandler(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}

	rows, _ := DB.Query("SELECT id, created_at, sent, json_data FROM contracts ORDER BY created_at DESC")
	defer rows.Close()

	var contracts []Contract
	for rows.Next() {
		var c Contract
		var jsonData string
		rows.Scan(&c.ID, &c.CreatedAt, &c.Sent, &jsonData)

		// Распарсим FormData из JSON
		var form FormData
		_ = json.Unmarshal([]byte(jsonData), &form)
		c.Form = form

		contracts = append(contracts, c)
	}

	tmpl, _ := template.ParseFiles("admin_templates/panel.html")
	tmpl.Execute(w, contracts)
}

// ===== 🔍 Просмотр договора =====

func AdminViewContractHandler(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	idStr := filepath.Base(r.URL.Path)
	id, _ := strconv.Atoi(idStr)

	var c Contract
	err := DB.QueryRow("SELECT id, created_at, json_data, pdf_path, photo_path, sent, sent_at FROM contracts WHERE id = ?", id).
		Scan(&c.ID, &c.CreatedAt, &c.JSONData, &c.PDFPath, &c.PhotoPath, &c.Sent, &c.SentAt)
	if err != nil {
		http.Error(w, "Müqavilə tapılmadı", 404)
		return
	}

	var form FormData
	if err := json.Unmarshal([]byte(c.JSONData), &form); err != nil {
		http.Error(w, "JSON decode error", 500)
		return
	}

	var positiveAnswers []string
	for i, ans := range form.Answers {
		if ans == "Bəli" {
			positiveAnswers = append(positiveAnswers, fmt.Sprintf("№%d", i+1))
		}
	}

	tmpl, err := template.ParseFiles("admin_templates/contract_view.html")
	if err != nil {
		http.Error(w, "Şablon tapılmadı", 500)
		return
	}

	data := struct {
		Contract
		Form            FormData
		PositiveAnswers []string
	}{
		Contract:        c,
		Form:            form,
		PositiveAnswers: positiveAnswers,
	}

	tmpl.Execute(w, data)
}

// ===== 📤 Отправка в Telegram =====

func AdminSendHandler(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	idStr := filepath.Base(r.URL.Path)
	id, _ := strconv.Atoi(idStr)

	var path, jsonData string
	var sent bool
	var photo string
	DB.QueryRow("SELECT pdf_path, json_data, sent, photo_path FROM contracts WHERE id = ?", id).
		Scan(&path, &jsonData, &sent, &photo)

	force := r.URL.Query().Get("force") == "true"
	if sent && !force {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	var form FormData
	_ = json.Unmarshal([]byte(jsonData), &form)

	_ = generateZIP(form, path)

	if err := sendToTelegram(path); err == nil {
		now := time.Now()
		DB.Exec("UPDATE contracts SET sent = 1, sent_at = ? WHERE id = ?", now, id)
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// ===== ❌ Удаление договора =====

func AdminDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	idStr := filepath.Base(r.URL.Path)
	id, _ := strconv.Atoi(idStr)

	var path, photo string
	DB.QueryRow("SELECT pdf_path, photo_path FROM contracts WHERE id = ?", id).Scan(&path, &photo)

	DB.Exec("DELETE FROM contracts WHERE id = ?", id)

	if path != "" {
		_ = os.Remove(path)
	}
	if photo != "" {
		_ = os.Remove(photo)
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
