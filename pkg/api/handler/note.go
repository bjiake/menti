package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"menti/pkg/db"
	"menti/pkg/domain/note"
	"net/http"
	"strconv"
)

func (h *Handler) Note(r chi.Router) {
	r.Post("/", post(h))
	r.Get("/", get(h))
}

func post(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookieID, err := r.Cookie("id")
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
			log.Error("Ошибка получения cookie id", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		var request note.Note
		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
			log.Error("Ошибка декодирования запроса", err)
			return
		}

		// вызовите метод NotePost с контекстом и структурой note
		result, err := h.service.NotePost(ctx, cookieID.Value, request)
		if err != nil {
			switch err.Error() {
			case db.ErrAuthorize.Error():
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
				return
			case db.ErrYandexSpeller.Error():
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
				return
			case db.ErrDuplicate.Error():
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusConflict)
				return
			default:
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
		}
		response := struct {
			ID string `json:"id"`
		}{ID: strconv.FormatInt(result, 10)}
		json.NewEncoder(w).Encode(response)
		return
	}
}

func get(h *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookieID, err := r.Cookie("id")
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
			log.Error("Ошибка получения cookie id", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		var result []note.Note

		result, err = h.service.NoteGetAll(ctx, cookieID.Value)
		if err != nil {
			switch err.Error() {
			case db.ErrAuthorize.Error():
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
				return
			default:
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
		}
		response := struct {
			Notes []note.Note `json:"data"`
		}{Notes: result}
		json.NewEncoder(w).Encode(response)
		return
	}
}
