package user

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	userService "github.com/jaysm12/multifinance-apps/internal/service/user"
	"github.com/jaysm12/multifinance-apps/internal/service/user/mock"
)

func TestUserHandler_CreateUserKycHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	m := mock.NewMockUserServiceMethod(mockCtrl)
	defer mockCtrl.Finish()

	type args struct {
		userID uint
		body   string
	}

	type want struct {
		code int
		body string
	}

	tests := []struct {
		name        string
		args        args
		want        want
		mockFunc    func()
		mockContext func() (context.Context, func())
	}{
		{
			name: "success flow",
			args: args{
				body:   `{"nik":"1234567890123456","legal_name":"John Doe","birth_date":"1990-05-15","birth_address":"123 Main Street, Cityville","salary_amount":"50000","photo_id_url":"http://example.com/photo_id.jpg","photo_selfie_url":"http://example.com/photo_selfie.jpg"}`,
				userID: 1,
			},
			mockFunc: func() {
				m.EXPECT().CreateUserKyc(userService.CreateUserKycRequest{
					UserId:         uint(1),
					NIK:            "1234567890123456",
					LegalName:      "John Doe",
					BirthDate:      "1990-05-15",
					BirthAddress:   "123 Main Street, Cityville",
					SalaryAmount:   "50000",
					PhotoIDUrl:     "http://example.com/photo_id.jpg",
					PhotoSelfieUrl: "http://example.com/photo_selfie.jpg",
				}).Return(nil)
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 201,
				body: `{"code":201,"message":"success"}`,
			},
		},
		{
			name: "error on service",
			args: args{
				body:   `{"nik":"1234567890123456","legal_name":"John Doe","birth_date":"1990-05-15","birth_address":"123 Main Street, Cityville","salary_amount":"50000","photo_id_url":"http://example.com/photo_id.jpg","photo_selfie_url":"http://example.com/photo_selfie.jpg"}`,
				userID: 1,
			},
			mockFunc: func() {
				m.EXPECT().CreateUserKyc(userService.CreateUserKycRequest{
					UserId:         uint(1),
					NIK:            "1234567890123456",
					LegalName:      "John Doe",
					BirthDate:      "1990-05-15",
					BirthAddress:   "123 Main Street, Cityville",
					SalaryAmount:   "50000",
					PhotoIDUrl:     "http://example.com/photo_id.jpg",
					PhotoSelfieUrl: "http://example.com/photo_selfie.jpg",
				}).Return(fmt.Errorf("some error"))
			},
			mockContext: func() (context.Context, func()) {
				return context.Background(), func() {}
			},
			want: want{
				code: 500,
				body: `{"code":500,"message":"some error"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			defer mockCtrl.Finish()
			handler := &UserHandler{
				service:      m,
				timeoutInSec: 100,
			}
			r := httptest.NewRequest(http.MethodGet, "/user/kyc", strings.NewReader(tt.args.body))
			ctx, cancel := tt.mockContext()
			defer cancel()
			r = r.WithContext(ctx)
			if tt.args.userID > 0 {
				r = r.WithContext(context.WithValue(r.Context(), "id", tt.args.userID))
			}
			w := httptest.NewRecorder()
			handler.CreateUserKyc(w, r)
			result := w.Result()
			resBody, err := ioutil.ReadAll(result.Body)

			if err != nil {
				t.Fatalf("Error read body err = %v\n", err)
			}

			if string(resBody) != tt.want.body {
				t.Fatalf("GetStatHandler body got =%s, want %s \n", string(resBody), tt.want.body)
			}

			if result.StatusCode != tt.want.code {
				t.Fatalf("GetStatHandler status code got =%d, want %d \n", result.StatusCode, tt.want.code)
			}
		})
	}
}
