package payment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/toVersus/otel-demo/pkg/utils"
)

type data struct {
	Amount int `json:"amount" validate:"required"`
}

func (s *Server) transferAmount(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.UserIDFromContext("/payments/transfer/id/", r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("userID: %s", userID)

	var data data
	if err := utils.ReadBody(w, r, &data); err != nil {
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	// send the request to user service
	url := fmt.Sprintf("%s/users/%s", s.userUrl, userID)
	resp, err := utils.SendRequest(r.Context(), http.MethodPut, url, payload)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("payment failed. got response: %s", b))
		return
	}

	utils.WriteResponse(w, http.StatusOK, data)
}
