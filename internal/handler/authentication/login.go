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

// LoginUserRequest is list request parameter for Login Api
type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginUserResponse is list response parameter for Login Api
type LoginUserResponse struct {
	Token string `json:"token"`
}

// LoginUserHandler is func handler for login
func (h *AuthenticationHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println("[LoginUserHandler]-Error Marshal Response :", err)
			code = http.StatusInternalServerError
			data = []byte(`{"code":500,"message":"Internal Server Error"}`)
		}
		utilhttp.WriteResponse(w, data, code)
	}()

	var body LoginUserRequest
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
	var token string
	go func(ctx context.Context) {
		token, err = h.service.Login(
			authentication.LoginServiceRequest{
				Username: body.Username,
				Password: body.Password,
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
			if err == user.ErrUserNameNotExists || err == user.ErrPasswordIsIncorrect || strings.Contains(err.Error(), "not found") {
				code = http.StatusNotFound
				err = fmt.Errorf("invalid Username or Password")
			} else {
				code = http.StatusInternalServerError
			}
			return
		}
	}

	response = mapResponseLogin(token)
}

func mapResponseLogin(result string) utilhttp.StandardResponse {
	var res utilhttp.StandardResponse
	data := LoginUserResponse{
		Token: result,
	}
	res.Data = data
	return res
}
