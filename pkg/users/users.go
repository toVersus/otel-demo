package users

import (
	"fmt"
	"net/http"

	"github.com/toVersus/otel-demo/pkg/datastore"
	"github.com/toVersus/otel-demo/pkg/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const (
	tracerName = "users"
)

type user struct {
	ID       int64  `json:"id" validate:"-"`
	UserName string `json:"user_name" validate:"required"`
	Account  string `json:"account" validate:"required"`
	Amount   int
}

type paymentData struct {
	Amount int `json:"amount" validate:"required"`
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var u user
	if err := utils.ReadBody(w, r, &u); err != nil {
		return
	}

	ctx, span := otel.Tracer(tracerName).Start(r.Context(), "create user")
	defer span.End()
	id, err := s.db.InsertOne(ctx, datastore.InsertParams{
		Query: `INSERT INTO USERS(USER_NAME, ACCOUNT) VALUES (?, ?)`,
		Vars:  []interface{}{u.UserName, u.Account},
	})
	if err != nil {
		msg := "create user error"
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", msg, err))
		span.SetStatus(codes.Error, msg)
		span.RecordError(err)
		return
	}

	u.ID = id
	utils.WriteResponse(w, http.StatusCreated, u)
}

func (s *Server) manageUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getUser(w, r)
	case http.MethodPut:
		s.updateUser(w, r)
	default:
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("%s method not implemented", r.Method))
	}
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.UserIDFromContext("/users/", r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
	}
	var u user

	ctx, span := otel.Tracer(tracerName).Start(r.Context(), "get user")
	defer span.End()
	if err := s.db.SelectOne(ctx, datastore.SelectParams{
		Query:   `select ID, USER_NAME, ACCOUNT, AMOUNT from USERS where ID = ?`,
		Filters: []interface{}{userID},
		Result:  []interface{}{&u.ID, &u.UserName, &u.Account, &u.Amount},
	}); err != nil {
		msg := "get user error"
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", msg, err))
		span.SetStatus(codes.Error, msg)
		span.RecordError(err)
		return
	}

	utils.WriteResponse(w, http.StatusOK, u)
}

func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.UserIDFromContext("/users/", r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
	}
	var data paymentData
	if err := utils.ReadBody(w, r, &data); err != nil {
		return
	}

	ctx, span := otel.Tracer(tracerName).Start(r.Context(), "update user amount")
	defer span.End()
	if err := s.db.UpdateOne(ctx, datastore.UpdateParams{
		Query: `update USERS set AMOUNT = AMOUNT + ? where ID = ?`,
		Vars:  []interface{}{data.Amount, userID},
	}); err != nil {
		msg := "update user amount error"
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", msg, err))
		span.SetStatus(codes.Error, msg)
		span.RecordError(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
