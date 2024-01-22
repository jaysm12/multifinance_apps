package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jaysm12/multifinance-apps/internal/handler/utilhttp"
	"github.com/jaysm12/multifinance-apps/internal/service/authentication"
	"github.com/jaysm12/multifinance-apps/internal/service/user"
)

// RegisterUserRequest is list request parameter for Register Api
type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

// RegisterUserHandler is func handler for Register user
func (h *AuthenticationHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[RegisterUserHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body RegisterUserRequest
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

	// checking valid body
	if len(body.Username) < 1 || len(body.Password) < 1 {
		code = http.StatusBadRequest
		err = fmt.Errorf("invalid parameter request")
		return
	}

	errChan := make(chan error, 1)
	go func(ctx context.Context) {
		err = h.service.Register(authentication.RegisterServiceRequest{
			Username: body.Username,
			Password: body.Password,
			Fullname: body.Fullname,
			Email:    body.Email,
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
			if strings.Contains(err.Error(), user.ErrUserNameAlreadyExists.Error()) {
				code = http.StatusConflict
			} else {
				code = http.StatusInternalServerError
			}
			return
		}
	}

	response = mapResonseRegister()
}

func mapResonseRegister() utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	return res
}
