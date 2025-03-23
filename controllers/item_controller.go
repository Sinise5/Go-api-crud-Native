package controllers

import (
	"encoding/json"
	"myapp/config"
	"myapp/models"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	json.NewDecoder(r.Body).Decode(&item)

	_, err := config.DB.Exec("INSERT INTO items (title, body) VALUES ($1, $2)", item.Title, item.Body)
	if err != nil {
		http.Error(w, "Error creating item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item created"})
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, title, body FROM items")
	if err != nil {
		http.Error(w, "Error fetching items", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		rows.Scan(&item.ID, &item.Title, &item.Body)
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(items)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, err := config.DB.Exec("DELETE FROM items WHERE id = $1", params["id"])
	if err != nil {
		http.Error(w, "Error deleting item", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted"})
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var item models.Item
	json.NewDecoder(r.Body).Decode(&item)

	_, err := config.DB.Exec("UPDATE items SET title = $1, body = $2 WHERE id = $3", item.Title, item.Body, params["id"])
	if err != nil {
		http.Error(w, "Error updating item", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Item updated successfully"})
}
