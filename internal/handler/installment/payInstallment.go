package installment

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jaysm12/multifinance-apps/internal/handler/utilhttp"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PayInstallmentPublishPayload struct {
	UserID        uint    `json:"user_id"`
	ContractID    string  `json:"contract_id"`
	InstallmentID uint    `json:"installment_id"`
	PaidAmount    float64 `json:"paid_amount"`
}

type PayInstallmentRequest struct {
	InstallmentID uint    `json:"installment_id"`
	PaidAmount    float64 `json:"paid_amount"`
}

const (
	queuePayInstallment = "pay_installment"
)

func (h *InstallmentHandler) PayInstallment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(h.timeoutInSec)*time.Second)
	defer cancel()

	var err error
	var response utilhttp.StandardResponse
	var code int = http.StatusOK

	defer func() {
		response.Code = code
		if err == nil {
			response.Message = "success"
		} else {
			response.Message = err.Error()
		}

		data, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			log.Println("[CreateInstallment]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body PayInstallmentRequest
	data, err := io.ReadAll(r.Body)
	if err != nil {
		code = http.StatusBadRequest
		err = fmt.Errorf("bad Request")
		return
	}

	err = json.Unmarshal(data, &body)
	if err != nil {
		code = http.StatusBadRequest
		err = fmt.Errorf("bad Request")
		return
	}

	vars := mux.Vars(r)
	contractId := vars["contract_id"]

	if body.InstallmentID == 0 {
		code = http.StatusBadRequest
		err = fmt.Errorf("bad Request")
		return
	}

	toPublish := PayInstallmentPublishPayload{
		InstallmentID: body.InstallmentID,
		PaidAmount:    body.PaidAmount,
		UserID:        r.Context().Value("id").(uint),
		ContractID:    contractId,
	}

	jsonData, err := json.Marshal(toPublish)
	if err != nil {
		code = http.StatusInternalServerError
		err = fmt.Errorf("failed to marshal json: %v", err)
		return
	}

	err = h.rabbitmqClient.Channel.PublishWithContext(
		ctx,
		"",
		queuePayInstallment,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		code = http.StatusInternalServerError
		err = fmt.Errorf("failed to publish a message: %v", err)
		return
	}
}
