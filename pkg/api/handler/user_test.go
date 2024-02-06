package handler

import (
	"bytes"
	"clean/pkg/usecase/mock"
	"clean/pkg/utils/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_UserSignup(t *testing.T) {
	testCase := map[string]struct {
		input         models.SignupDetail
		buildStub     func(useCaseMock *mock.MockUserUseCase, signupData models.SignupDetail)
		checkResponse func(t *testing.T, responseRecorder *httptest.ResponseRecorder)
	}{
		"Valid Signup": {
			input: models.SignupDetail{
				Firstname:    "preya",
				Lastname:     "v",
				Email:        "preya@gmail.com",
				Password:     "1234",
				Phone:        "7902689612",
				ReferralCode: "659823",
			},
			buildStub: func(useCaseMock *mock.MockUserUseCase, signupData models.SignupDetail) {
				err := validator.New().Struct(signupData)
				if err != nil {
					fmt.Println("validation failed")
				}
				useCaseMock.EXPECT().UserSignup(signupData).Times(1).Return(&models.TokenUser{
					Users: models.SignupDetailResponse{
						Id:        1,
						Firstname: "preya",
						Lastname:  "v",
						Email:     "preya@gmail.com",
						Phone:     "7902689612",
					},
					AccessToken:  "adfsae.thjjshahfiurhf.ajherkuefeu",
					RefreshToken: "fkdgker.jrijigsiejggj.rlisjgjisg3",
				}, nil)
			}, 
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, responseRecorder.Code)
			},
		},
		"user couldnot sign up": {
			input: models.SignupDetail{
				Firstname:    "preya",
				Lastname:     "v",
				Email:        "preya@gmail.com",
				Password:     "1234",
				Phone:        "7902689612",
				ReferralCode: "659823",
			},
			buildStub: func(useCaseMock *mock.MockUserUseCase, signupData models.SignupDetail) {
				// copying signupData to domain.user for pass to Mock usecase
				err := validator.New().Struct(signupData)
				if err != nil {
					fmt.Println("validation failed")  
				}

				useCaseMock.EXPECT().UserSignup(signupData).Times(1).Return(&models.TokenUser{}, errors.New("cannot sign up"))
			},
			checkResponse: func(t *testing.T, responseRecorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

			},
		},
	}
	for testName, test := range testCase {
		test := test
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockUseCase := mock.NewMockUserUseCase(ctrl)
			test.buildStub(mockUseCase, test.input)
			userHandler := NewUserHandler(mockUseCase)

			server := gin.Default()
			server.POST("/signup", userHandler.UserSignUp)

			jsonData, err := json.Marshal(test.input)
			assert.NoError(t, err)
			body := bytes.NewBuffer(jsonData)

			mockRequest, err := http.NewRequest(http.MethodPost, "/signup", body)
			assert.NoError(t, err)
			responseRecorder := httptest.NewRecorder()
			server.ServeHTTP(responseRecorder, mockRequest)

			test.checkResponse(t, responseRecorder)

		})
	}
}
