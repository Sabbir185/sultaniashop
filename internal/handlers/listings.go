package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
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
	db     *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *listingHandler {
	return &listingHandler{
		db:     db,
		logger: logger,
	}
}

func (h *listingHandler) ListAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rows, err := h.db.QueryContext(ctx,
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
	ctx := r.Context()
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}
	_, err := h.db.ExecContext(ctx, `DELETE FROM listing WHERE id = $1`, id)
	if err != nil {
		h.logger.Error("delete failed", "id", id, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "deleted successfully",
	})

}
