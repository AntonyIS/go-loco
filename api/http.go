package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/AntonyIS/go-loco/app"
	js "github.com/AntonyIS/go-loco/serializer/json"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

type LocomotiveHandler interface {
	GetLoco(http.ResponseWriter, *http.Request)
	GetAllLoco(http.ResponseWriter, *http.Request)
	PostLoco(http.ResponseWriter, *http.Request)
	UpdateLoco(http.ResponseWriter, *http.Request)
	DeleteLoco(http.ResponseWriter, *http.Request)
}

type handler struct {
	locomotiveService app.LocomotiveService
}

func NewHandler(locomotiveService app.LocomotiveService) LocomotiveHandler {
	return &handler{locomotiveService: locomotiveService}
}

func setResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) serializer(contentType string) app.LocomotiveSerializer {
	return &js.Locomotive{}
}

func (h *handler) GetLoco(w http.ResponseWriter, r *http.Request) {
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

func (h *handler) GetAllLoco(w http.ResponseWriter, r *http.Request) {

	loco, err := h.locomotiveService.GetAllLoco()

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

func (h *handler) PostLoco(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	loco, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	loco, err = h.locomotiveService.CreateLoco(loco)

	if err != nil {
		if errors.Cause(err) == app.ErrorInvalidLocomotive {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responseBody, err := h.serializer(contentType).Encode(loco)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	setResponse(w, contentType, responseBody, http.StatusCreated)

}

func (h *handler) UpdateLoco(w http.ResponseWriter, r *http.Request) {
	loco := &app.Locomotive{}
	loco_id := chi.URLParam(r, "loco_id")
	json.NewDecoder(r.Body).Decode(&loco)
	loco.LocoID = loco_id
	results, err := h.locomotiveService.UpdateLoco(loco)
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
	json.NewEncoder(w).Encode(results)

}

func (h *handler) DeleteLoco(w http.ResponseWriter, r *http.Request) {
	loco_id := chi.URLParam(r, "loco_id")
	err := h.locomotiveService.DeleteLoco(loco_id)
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

}
