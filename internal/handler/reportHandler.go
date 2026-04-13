package handler

import (
	"encoding/json"
	"net/http"
)

// ===== REPORT =====

// @Summary Get Report
// @Description get protected report data
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /report [get]
func ReportHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"total_orders":  100,
		"total_revenue": 9999.99,
	}

	json.NewEncoder(w).Encode(data)
}
