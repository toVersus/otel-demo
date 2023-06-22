package order

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

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

	ctx, getUserSpan := otel.Tracer(tracerName).Start(r.Context(), "get user")
	// get user details from user service
	url := fmt.Sprintf("%s/users/%d", s.userUrl, request.UserID)
	userResponse, err := utils.SendRequest(r.Context(), http.MethodGet, url, nil)
	if err != nil {
		log.Printf("%v", err)
		utils.WriteResponse(w, http.StatusInternalServerError, err)
		return
	}
	getUserSpan.End()

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

	ctx, createOrderSpan := otel.Tracer(tracerName).Start(ctx, "create order",
		trace.WithAttributes(attribute.String("user_id", fmt.Sprintf("%d", user.ID))),
	)
	id, err := s.db.InsertOne(ctx, datastore.InsertParams{
		Query: `insert into ORDERS(ACCOUNT, PRODUCT_NAME, PRICE, ORDER_STATUS) VALUES (?,?,?,?)`,
		Vars:  []interface{}{user.Account, request.ProductName, request.Price, "SUCCESS"},
	})
	if err != nil {
		msg := "insert order error"
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("%s: %w", msg, err))
		createOrderSpan.SetStatus(codes.Error, msg)
		createOrderSpan.RecordError(err)
		createOrderSpan.End()
		return
	}
	createOrderSpan.End()

	ctx, updateAmountSpan := otel.Tracer(tracerName).Start(ctx, "update user amount",
		trace.WithAttributes(attribute.String("user_id", fmt.Sprintf("%d", user.ID))),
	)
	// update the pending amount in user table
	if err := s.db.UpdateOne(ctx, datastore.UpdateParams{
		Query: `update USERS set AMOUNT = AMOUNT - ? where ID = ?`,
		Vars:  []interface{}{request.Price, user.ID},
	}); err != nil {
		msg := "update user amount error"
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		updateAmountSpan.SetStatus(codes.Error, msg)
		updateAmountSpan.RecordError(err)
		updateAmountSpan.End()
		return
	}
	updateAmountSpan.End()

	// send response
	response := request
	response.ID = id
	utils.WriteResponse(w, http.StatusCreated, response)
}
