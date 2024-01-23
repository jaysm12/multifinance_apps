package installment

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/jaysm12/multifinance-apps/internal/handler/utilhttp"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CreateInstallmentRequest struct {
	CreditOptionID uint    `json:"credit_option_id"`
	OtrAmount      float64 `json:"otr_amount"`
	AssetName      string  `json:"asset_name"`
}

type CreateInstallmentPublishPayload struct {
	UserID         uint    `json:"user_id"`
	CreditOptionID uint    `json:"credit_option_id"`
	OtrAmount      float64 `json:"otr_amount"`
	AssetName      string  `json:"asset_name"`
}

const (
	queueCreateInstallment = "create_installment"
)

func (h *InstallmentHandler) CreateInstallment(w http.ResponseWriter, r *http.Request) {
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

	var body CreateInstallmentRequest
	data, err := io.ReadAll(r.Body)
	if err != nil {
		code = http.StatusBadRequest
		err = fmt.Errorf("bad request")
		return
	}

	err = json.Unmarshal(data, &body)
	if err != nil {
		code = http.StatusBadRequest
		err = fmt.Errorf("bad request")
		return
	}
	payload := CreateInstallmentPublishPayload{
		UserID:         r.Context().Value("id").(uint),
		CreditOptionID: body.CreditOptionID,
		OtrAmount:      body.OtrAmount,
		AssetName:      body.AssetName,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		code = http.StatusInternalServerError
		err = fmt.Errorf("failed to marshal json data, err : %v", err)
		return
	}

	err = h.rabbitmqClient.Channel.PublishWithContext(
		ctx,
		"",
		queueCreateInstallment,
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
