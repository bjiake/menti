package handler

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"menti/pkg/db"
	"menti/pkg/domain/account"
	"net/http"
	"strconv"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var login account.Login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		log.Error("Ошибка декодирования данных", err)
		return
	}
	result, err := h.service.Login(ctx, login)
	if err != nil {
		switch err.Error() {
		case db.ErrNotExist.Error():
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
			return
		default:
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		ID string `json:"id"`
	}{ID: strconv.FormatInt(result, 10)}
	json.NewEncoder(w).Encode(response)
	return
}
