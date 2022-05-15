package api

import (
	"encoding/json"
	"net/http"

	"github.com/AntonyIS/go-loco/app"
	js "github.com/AntonyIS/go-loco/serializer/json"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

type LocomotiveHandler interface {
	Get(http.ResponseWriter, *http.Request)
	// Post(http.ResponseWriter, *http.Request)
	// Update(http.ResponseWriter, *http.Request)
	// Delete(http.ResponseWriter, *http.Request)
}

type handler struct {
	locomotiveService app.LocomotiveService
}

func NewHandler(locomotiveService app.LocomotiveService) LocomotiveHandler {
	return &handler{locomotiveService: locomotiveService}
}

// func setResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
// 	w.Header().Set("Content-Type", contentType)
// 	w.WriteHeader(statusCode)

// 	_, err := w.Write(body)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

func (h *handler) serializer(contentType string) app.LocomotiveSerializer {
	return &js.Locomotive{}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	loco_id := chi.URLParam(r, "loco_id")
	loco, err := h.locomotiveService.GetLoco(loco_id)
	if err != nil {
		if errors.Cause(err) == app.ErrorLocomotiveNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loco)
}
