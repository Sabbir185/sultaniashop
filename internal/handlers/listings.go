package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

type listingHandler struct {
	db *sql.DB
}

func NewListingHandler(db *sql.DB) *listingHandler {
	return &listingHandler{
		db: db,
	}
}

func (h *listingHandler) ListAllProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(
		`SELECT id, title, description, price, city, created_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT 100`)
	if err != nil {
		log.Printf("Query: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := []listing{}
	for rows.Next() {
		var list listing
		if err := rows.Scan(&list.ID, &list.Title, &list.Description, &list.Price, &list.City, &list.CreatedAt); err != nil {
			log.Printf("Scan: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		products = append(products, list)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Err: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)

}

func (h *listingHandler) DeleteList(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}
	_, err := h.db.Exec(`DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		log.Printf("Delete: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "deleted successfully",
	})

}
