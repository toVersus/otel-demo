package order

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/toVersus/otel-demo/pkg/datastore"
	"github.com/toVersus/otel-demo/pkg/utils"
)

const (
	tracerName = "order"
)

type data struct {
	ID          int64  `json:"id"`
	UserID      int    `json:"user_id" validate:"required"`
	ProductName string `json:"product_name" validate:"required"`
	Price       int    `json:"price" validate:"required"`
}

type user struct {
	ID       int64  `json:"id"`
	UserName string `json:"user_name"`
	Account  string `json:"account"`
	Amount   int
}

func (s *Server) createOrder(w http.ResponseWriter, r *http.Request) {
	var request data
	if err := utils.ReadBody(w, r, &request); err != nil {
		return
	}

	// get user details from user service
	url := fmt.Sprintf("%s/users/%d", s.userUrl, request.UserID)
	userResponse, err := utils.SendRequest(r.Context(), http.MethodGet, url, nil)
	if err != nil {
		log.Printf("%v", err)
		utils.WriteResponse(w, http.StatusInternalServerError, err)
		return
	}

	b, err := ioutil.ReadAll(userResponse.Body)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	defer userResponse.Body.Close()

	if userResponse.StatusCode != http.StatusOK {
		utils.WriteErrorResponse(w, userResponse.StatusCode, fmt.Errorf("payment failed. got response: %s", b))
		return
	}

	var user user
	if err := json.Unmarshal(b, &user); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	// basic check for the user balance
	if user.Amount < request.Price {
		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, fmt.Errorf("insufficient balance. add %d more amount to account", request.Price-user.Amount))
		return
	}

	ctx, insertSpan := otel.Tracer(tracerName).Start(r.Context(), "insert order")
	id, err := s.db.InsertOne(ctx, datastore.InsertParams{
		Query: `insert into ORDERS(ACCOUNT, PRODUCT_NAME, PRICE, ORDER_STATUS) VALUES (?,?,?,?)`,
		Vars:  []interface{}{user.Account, request.ProductName, request.Price, "SUCCESS"},
	})
	if err != nil {
		msg := "insert order error"
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", msg, err))
		insertSpan.SetStatus(codes.Error, msg)
		insertSpan.RecordError(err)
		insertSpan.End()
		return
	}
	insertSpan.End()

	ctx, updateSpan := otel.Tracer(tracerName).Start(r.Context(), "update user amount")
	// update the pending amount in user table
	if err := s.db.UpdateOne(ctx, datastore.UpdateParams{
		Query: `update USERS set AMOUNT = AMOUNT - ? where ID = ?`,
		Vars:  []interface{}{request.Price, user.ID},
	}); err != nil {
		msg := "update user amount error"
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		updateSpan.SetStatus(codes.Error, msg)
		updateSpan.RecordError(err)
		updateSpan.End()
		return
	}
	updateSpan.End()

	// send response
	response := request
	response.ID = id
	utils.WriteResponse(w, http.StatusCreated, response)
}
