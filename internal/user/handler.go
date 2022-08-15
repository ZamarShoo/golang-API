package user

import (
	"go_advanced/internal/handlers"
	"go_advanced/pkg/logging"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var _ handlers.Handler = &handler{}

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)

type handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(usersURL, h.GetList)
	router.GET(userURL, h.GeUserByUuid)
	router.POST(usersURL, h.Createuser)
	router.PUT(usersURL, h.UpdateUser)
	router.PATCH(usersURL, h.PartiallyUpdateUser)
	router.DELETE(usersURL, h.Deleteuser)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("list of all users"))
	w.WriteHeader(200)
}

func (h *handler) GeUserByUuid(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("get one user"))
	w.WriteHeader(200)
}

func (h *handler) Createuser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("create user"))
	w.WriteHeader(201)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("update user"))
	w.WriteHeader(204)
}

func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("part of update user"))
	w.WriteHeader(204)
}

func (h *handler) Deleteuser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("delete user"))
	w.WriteHeader(204)
}
