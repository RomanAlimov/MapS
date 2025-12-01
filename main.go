package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"zov/models" // ✅ Правильный импорт

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var tmpl *template.Template

func init() {
	var err error
	dsn := "host=localhost user=postgres password=Kroker970z dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Moscow"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	db.AutoMigrate(&models.Tree{})

	tmpl = template.Must(template.ParseFiles("templates/index.html"))
}

func main() {
	log.Println("Сервер запущен на http://localhost:8080")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/trees", treesHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl.Execute(w, nil)
}

func treesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		var trees []models.Tree
		query := db.Model(&models.Tree{})

		if species := r.URL.Query().Get("species"); species != "" {
			query = query.Where("species ILIKE ?", "%"+species+"%")
		}
		if health := r.URL.Query().Get("health"); health != "" {
			query = query.Where("health_status = ?", health)
		}

		query.Find(&trees)
		json.NewEncoder(w).Encode(trees)

	} else if r.Method == http.MethodPost {
		r.ParseForm()
		species := strings.TrimSpace(r.FormValue("species"))
		lat, _ := strconv.ParseFloat(r.FormValue("lat"), 64)
		lng, _ := strconv.ParseFloat(r.FormValue("lng"), 64)
		health := r.FormValue("health_status")
		if health == "" {
			health = "healthy"
		}

		var plantedYear *int
		if yearStr := r.FormValue("planted_year"); yearStr != "" {
			if year, err := strconv.Atoi(yearStr); err == nil {
				plantedYear = &year
			}
		}

		tree := models.Tree{
			Species:      species,
			Lat:          lat,
			Lng:          lng,
			PlantedYear:  plantedYear,
			HealthStatus: health,
		}

		result := db.Create(&tree)
		if result.Error != nil {
			http.Error(w, "Ошибка сохранения", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tree)
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
