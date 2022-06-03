package handler

import (
	"context"
	"github.com/core-go/search"
	sv "github.com/core-go/service"
	"net/http"
	"reflect"

	. "go-service/internal/user/domain"
	. "go-service/internal/user/service"
)

func NewUserHandler(find func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), service UserService, status sv.StatusConfig, logError func(context.Context, string), validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), action *sv.ActionConfig) *HttpUserHandler {
	filterType := reflect.TypeOf(UserFilter{})
	modelType := reflect.TypeOf(User{})
	params := sv.CreateParams(modelType, &status, logError, validate, action)
	searchHandler := search.NewSearchHandler(find, modelType, filterType, logError, params.Log)
	return &HttpUserHandler{service: service, SearchHandler: searchHandler, Params: params}
}

type HttpUserHandler struct {
	service UserService
	*search.SearchHandler
	*sv.Params
}

func (h *HttpUserHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		res, err := h.service.Load(r.Context(), id)
		sv.RespondModel(w, r, res, err, h.Error, nil)
	}
}
func (h *HttpUserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := sv.Decode(w, r, &user)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Create) {
			res, er3 := h.service.Create(r.Context(), &user)
			sv.AfterCreated(w, r, &user, res, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Create)
		}
	}
}
func (h *HttpUserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := sv.DecodeAndCheckId(w, r, &user, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Update) {
			res, er3 := h.service.Update(r.Context(), &user)
			sv.HandleResult(w, r, &user, res, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Update)
		}
	}
}
func (h *HttpUserHandler) Patch(w http.ResponseWriter, r *http.Request) {
	var user User
	r, json, er1 := sv.BuildMapAndCheckId(w, r, &user, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Patch) {
			res, er3 := h.service.Patch(r.Context(), json)
			sv.HandleResult(w, r, json, res, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Patch)
		}
	}
}
func (h *HttpUserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		res, err := h.service.Delete(r.Context(), id)
		sv.HandleDelete(w, r, res, err, h.Error, h.Log, h.Resource, h.Action.Delete)
	}
}
