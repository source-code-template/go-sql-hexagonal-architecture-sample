package handler

import (
	"context"
	"encoding/json"
	"github.com/core-go/search"
	sv "github.com/core-go/service"
	"github.com/gorilla/mux"
	"go-service/internal/user/domain"
	"go-service/internal/user/service"
	"net/http"
	"reflect"
)

func NewUserHandler(find func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), service service.UserService, logError func(context.Context, string)) *HttpUserHandler {
	filterType := reflect.TypeOf(domain.UserFilter{})
	modelType := reflect.TypeOf(domain.User{})
	searchHandler := search.NewSearchHandler(find, modelType, filterType, logError, nil)
	return &HttpUserHandler{service: service, SearchHandler: searchHandler}
}

type HttpUserHandler struct {
	service service.UserService
	*search.SearchHandler
}

func (h *HttpUserHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}

	user, err := h.service.Load(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusOK, user)
}
func (h *HttpUserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	er1 := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}

	res, er2 := h.service.Create(r.Context(), &user)
	if er2 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusCreated, res)
}
func (h *HttpUserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	er1 := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		http.Error(w, "Id not match", http.StatusBadRequest)
		return
	}

	res, er2 := h.service.Update(r.Context(), &user)
	if er2 != nil {
		http.Error(w, er2.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusOK, res)
}
func (h *HttpUserHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}

	var user domain.User
	userType := reflect.TypeOf(user)
	_, jsonMap, _ := sv.BuildMapField(userType)
	body, er1 := sv.BuildMapAndStruct(r, &user)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	if len(user.Id) == 0 {
		user.Id = id
	} else if id != user.Id {
		http.Error(w, "Id not match", http.StatusBadRequest)
		return
	}
	json, er2 := sv.BodyToJsonMap(r, user, body, []string{"id"}, jsonMap)
	if er2 != nil {
		http.Error(w, er2.Error(), http.StatusInternalServerError)
		return
	}

	res, er3 := h.service.Patch(r.Context(), json)
	if er3 != nil {
		http.Error(w, er3.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusOK, res)
}
func (h *HttpUserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	res, err := h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusOK, res)
}

func JSON(w http.ResponseWriter, code int, res interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(res)
}
