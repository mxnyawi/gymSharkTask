package api

import (
	"encoding/json"
	"net/http"

	"github.com/mxnyawi/gymSharkTask/internal/model"
)

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	total := order.CalculateTotal()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(total)
}
