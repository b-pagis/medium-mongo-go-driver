package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"localhost/medium-mongo-go-driver/databases"
	"localhost/medium-mongo-go-driver/handlers"
	"localhost/medium-mongo-go-driver/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userDBMock struct {
	isError  bool
	objectID primitive.ObjectID
}

func (u userDBMock) FindOne(context.Context, interface{}) (*models.User, error) {
	if u.isError {
		return nil, errors.New("some error")
	}

	return &models.User{ID: u.objectID, Email: "email", Username: "user.name"}, nil
}

func (u userDBMock) Create(context.Context, *models.User) error {
	if u.isError {
		return errors.New("some error")
	}
	return nil
}
func (u userDBMock) DeleteByUsername(context.Context, string) error {
	if u.isError {
		return errors.New("some error")
	}
	return nil
}

func TestUser2020_GetUser(t *testing.T) {
	expectedObjID := primitive.NewObjectID()
	gin.SetMode(gin.TestMode)
	type fields struct {
		DB databases.UserDatabase
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantCode    int
		wantMessage string
	}{
		{
			name:        "failure",
			fields:      fields{DB: &userDBMock{isError: true, objectID: expectedObjID}},
			wantCode:    500,
			wantMessage: `{"err":"some error"}`,
		},
		{
			name:        "pass",
			fields:      fields{DB: &userDBMock{isError: false, objectID: expectedObjID}},
			wantCode:    200,
			wantMessage: `{"id":"` + expectedObjID.Hex() + `","username":"user.name","email":"email"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, engine := gin.CreateTestContext(w)
			basePath := engine.Group("")
			handler := handlers.User2020{
				DB: tt.fields.DB,
			}
			basePath.GET("/apiv2/users/:username", handler.GetUser)

			ctx.Request, _ = http.NewRequest("GET", "/apiv2/users/123", nil)
			engine.ServeHTTP(w, ctx.Request)

			if w.Code != tt.wantCode {
				t.Errorf("GET /apiv2/users/123 HTTP Code = %v, wantHTTPCode %v", w.Code, tt.wantCode)
				return
			}

			responseBody := w.Body.String()

			if responseBody != tt.wantMessage {
				t.Errorf("GET /apiv2/users/123 response = %v, want %v\n len(%v) = len(%v)", responseBody, tt.wantMessage, len(responseBody), len(tt.wantMessage))
			}
		})
	}
}

func TestUser2020_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	type fields struct {
		DB databases.UserDatabase
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantCode    int
		wantMessage string
		request     io.Reader
	}{
		{
			name:        "failure",
			fields:      fields{DB: &userDBMock{isError: true}},
			wantCode:    500,
			wantMessage: `{"err":"some error"}`,
			request:     bytes.NewBuffer([]byte("{\"username\":\"user.name\",\"email\":\"random@example.com\"}")),
		},
		{
			name:        "failure-bind",
			fields:      fields{DB: &userDBMock{isError: true}},
			wantCode:    400,
			wantMessage: `{"err":"EOF"}`,
			request:     bytes.NewBuffer([]byte("")),
		},
		{
			name:        "pass",
			fields:      fields{DB: &userDBMock{isError: false}},
			wantCode:    200,
			wantMessage: `{}`,
			request:     bytes.NewBuffer([]byte("{\"username\":\"user.name\",\"email\":\"random@example.com\"}")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, engine := gin.CreateTestContext(w)
			basePath := engine.Group("")
			handler := handlers.User2020{
				DB: tt.fields.DB,
			}
			basePath.POST("/apiv2/users", handler.CreateUser)

			// because there is no validation, then there is no point in marshaling exact data
			ctx.Request, _ = http.NewRequest("POST", "/apiv2/users", tt.request)
			engine.ServeHTTP(w, ctx.Request)

			if w.Code != tt.wantCode {
				t.Errorf("POST /apiv2/users HTTP Code = %v, wantHTTPCode %v", w.Code, tt.wantCode)
				return
			}

			responseBody := w.Body.String()

			if responseBody != tt.wantMessage {
				t.Errorf("POST /apiv2/users response = %v, want %v\n len(%v) = len(%v)", responseBody, tt.wantMessage, len(responseBody), len(tt.wantMessage))
			}
		})
	}
}

func TestUser2020_DeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	type fields struct {
		DB databases.UserDatabase
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantCode    int
		wantMessage string
	}{
		{
			name:        "failure-bind",
			fields:      fields{DB: &userDBMock{isError: true}},
			wantCode:    500,
			wantMessage: `{"err":"some error"}`,
		},
		{
			name:        "pass",
			fields:      fields{DB: &userDBMock{isError: false}},
			wantCode:    200,
			wantMessage: ``,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			ctx, engine := gin.CreateTestContext(w)
			basePath := engine.Group("")
			handler := handlers.User2020{
				DB: tt.fields.DB,
			}
			basePath.DELETE("/apiv2/users/:username", handler.DeleteUser)

			// because there is no validation, then there is no point in marshaling exact data
			ctx.Request, _ = http.NewRequest("DELETE", "/apiv2/users/123", nil)
			engine.ServeHTTP(w, ctx.Request)

			if w.Code != tt.wantCode {
				t.Errorf("DELETE /apiv2/users/123 HTTP Code = %v, wantHTTPCode %v", w.Code, tt.wantCode)
				return
			}

			responseBody := w.Body.String()

			if responseBody != tt.wantMessage {
				t.Errorf("GET /apiv2/users/123 response = %v, want %v\n len(%v) = len(%v)", responseBody, tt.wantMessage, len(responseBody), len(tt.wantMessage))
			}
		})
	}
}
