package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/jaysm12/multifinance-apps/internal/handler/utilhttp"
	"github.com/jaysm12/multifinance-apps/internal/service/user"
)

type CreateUserKycRequest struct {
	NIK            string `json:"nik"`
	LegalName      string `json:"legal_name"`
	BirthDate      string `json:"birth_date"`
	BirthAddress   string `json:"birth_address"`
	SalaryAmount   string `json:"salary_amount"`
	PhotoIDUrl     string `json:"photo_id_url"`
	PhotoSelfieUrl string `json:"photo_selfie_url"`
}

// CreateUserKyc is service level func to validate and create user kyc
func (h *UserHandler) CreateUserKyc(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[CreateUserKyc]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body CreateUserKycRequest
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

	if len(body.NIK) < 1 {
		code = http.StatusBadRequest
		err = fmt.Errorf("nik is required")
		return
	}

	errChan := make(chan error, 1)
	go func(ctx context.Context) {
		err = h.service.CreateUserKyc(user.CreateUserKycRequest{
			UserId:         r.Context().Value("id").(uint),
			NIK:            body.NIK,
			LegalName:      body.LegalName,
			BirthDate:      body.BirthDate,
			BirthAddress:   body.BirthAddress,
			SalaryAmount:   body.SalaryAmount,
			PhotoIDUrl:     body.PhotoIDUrl,
			PhotoSelfieUrl: body.PhotoSelfieUrl,
		})
		errChan <- err
	}(ctx)

	select {
	case <-ctx.Done():
		code = http.StatusGatewayTimeout
		err = fmt.Errorf("Timeout")
		return
	case err = <-errChan:
		if err != nil {
			code = http.StatusInternalServerError
			return
		}
	}
}
