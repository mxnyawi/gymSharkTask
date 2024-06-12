package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/mxnyawi/gymSharkTask/internal/db"
	"github.com/mxnyawi/gymSharkTask/internal/db/mocks"
	"github.com/stretchr/testify/mock"
)

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		contentType    string
		body           string
		mockDBManager  func() *mocks.MockDBManager
		expectedStatus int
	}{
		{
			name:           "Method not allowed",
			method:         http.MethodGet,
			contentType:    "application/json",
			body:           `{"username": "test", "password": "test"}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Content-Type header is not application/json",
			method:         http.MethodPost,
			contentType:    "text/plain",
			body:           `{"username": "test", "password": "test"}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "Invalid request body",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           `{"username": "test", "password":`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Username and password are required",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           `{"username": "", "password": ""}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Could not get database credentials",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"username": "test", "password": "test"}`,
			mockDBManager: func() *mocks.MockDBManager {
				m := &mocks.MockDBManager{}
				m.On("GetDBCreds").Return("", "", "", "", errors.New("test error"))
				return m
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "Could not store user",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"username": "test", "password": "test"}`,
			mockDBManager: func() *mocks.MockDBManager {
				m := &mocks.MockDBManager{}
				m.On("GetDBCreds").Return("bucket", "scope", "collection", "password", nil)
				m.On("WriteDocument", "bucket", "scope", "collection", "test", mock.Anything).Return(errors.New("test error"))
				return m
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "User created",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"username": "test", "password": "test"}`,
			mockDBManager: func() *mocks.MockDBManager {
				m := &mocks.MockDBManager{}
				m.On("GetDBCreds").Return("bucket", "scope", "collection", "password", nil)
				m.On("WriteDocument", "bucket", "scope", "collection", "test", mock.Anything).Return(nil)
				return m
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/registerUser", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()

			dbManager := tt.mockDBManager()

			RegisterHandler(rr, req, dbManager)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		contentType    string
		body           string
		mockDBManager  func() *mocks.MockDBManager
		expectedStatus int
	}{
		{
			name:           "Method not allowed",
			method:         http.MethodGet,
			contentType:    "application/json",
			body:           `{"username": "test", "password": "test"}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Content-Type header is not application/json",
			method:         http.MethodPost,
			contentType:    "text/plain",
			body:           `{"username": "test", "password": "test"}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "Invalid request body",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           `{"username": "test", "password":`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Username and password are required",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           `{"username": "", "password": ""}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Could not get database credentials",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"username": "test", "password": "test"}`,
			mockDBManager: func() *mocks.MockDBManager {
				m := &mocks.MockDBManager{}
				m.On("GetDBCreds").Return("", "", "", "", errors.New("test error"))
				return m
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "User authenticated",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"username": "test", "password": "test"}`,
			mockDBManager: func() *mocks.MockDBManager {
				m := &mocks.MockDBManager{}
				m.On("GetDBCreds").Return("bucket", "scope", "collection", "password", nil)
				hash, _ := argon2id.CreateHash("test", argon2id.DefaultParams)
				m.On("GetUser", "bucket", "scope", "collection", "test").Return(&db.User{Username: "test", Password: hash}, nil)
				return m
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/loginUser", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()

			dbManager := tt.mockDBManager()

			LoginHandler(rr, req, dbManager)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestPostOrderHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		contentType    string
		body           string
		mockDBManager  func() *mocks.MockDBManager
		expectedStatus int
	}{
		{
			name:           "Method not allowed",
			method:         http.MethodGet,
			contentType:    "application/json",
			body:           `{"orderAmount": 10, "packageSizes": [5, 5]}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Content-Type header is not application/json",
			method:         http.MethodPost,
			contentType:    "text/plain",
			body:           `{"orderAmount": 10, "packageSizes": [5, 5]}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "Invalid request body",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           `{"orderAmount": 10, "packageSizes": [5,`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Could not get database credentials",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"orderAmount": 10, "packageSizes": [5, 5]}`,
			mockDBManager: func() *mocks.MockDBManager {
				m := &mocks.MockDBManager{}
				m.On("GetDBCreds").Return("", "", "", "", errors.New("test error"))
				return m
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "Order created",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"orderAmount": 12, "packageSizes": [5, 10, 15, 20, 25]}`,
			mockDBManager: func() *mocks.MockDBManager {
				m := &mocks.MockDBManager{}
				m.On("GetDBCreds").Return("bucket", "scope", "collection", "document", nil)
				m.On("GetDocument", "bucket", "scope", "collection", "document").Return(&db.DocumentHistory{}, nil)
				m.On("WriteDocument", "bucket", "scope", "collection", "document", mock.AnythingOfType("*db.DocumentHistory")).Return(nil)
				return m
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/order", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()

			dbManager := tt.mockDBManager()

			PostOrderHandler(rr, req, dbManager)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestCreateAdminUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		contentType    string
		body           string
		mockDBManager  func() *mocks.MockDBManager
		expectedStatus int
	}{
		{
			name:           "Method not allowed",
			method:         http.MethodGet,
			contentType:    "application/json",
			body:           `{"username": "test", "password": "test"}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Content-Type header is not application/json",
			method:         http.MethodPost,
			contentType:    "text/plain",
			body:           `{"username": "test", "password": "test"}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "Invalid request body",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           `{"username": "test", "password":`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Username and password are required",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           `{"username": "", "password": ""}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Could not create admin user",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"username": "test", "password": "test"}`,
			mockDBManager: func() *mocks.MockDBManager {
				m := &mocks.MockDBManager{}
				m.On("GetDBCreds").Return("bucket", "scope", "collection", "document", nil)
				m.On("CreateAdminUser", "test", "test").Return(errors.New("test error"))
				return m
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "Admin user created",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"username": "test", "password": "test"}`,
			mockDBManager: func() *mocks.MockDBManager {
				m := &mocks.MockDBManager{}
				m.On("GetDBCreds").Return("bucket", "scope", "collection", "document", nil)
				m.On("CreateAdminUser", "test", "test").Return(nil)
				return m
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/createAdminUser", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()

			dbManager := tt.mockDBManager()

			CreateAdminUserHandler(rr, req, dbManager)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestSetDocumentHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		contentType    string
		body           string
		mockDBManager  func() *mocks.MockDBManager
		expectedStatus int
	}{
		{
			name:           "Method not allowed",
			method:         http.MethodGet,
			contentType:    "application/json",
			body:           `{"orderAmount": 10, "packageSizes": [5, 5]}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Content-Type header is not application/json",
			method:         http.MethodPost,
			contentType:    "text/plain",
			body:           `{"orderAmount": 10, "packageSizes": [5, 5]}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "Invalid request body",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           `{"orderAmount": 10, "packageSizes": [5,`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Could not get database credentials",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"orderAmount": 10, "packageSizes": [5, 5]}`,
			mockDBManager: func() *mocks.MockDBManager {
				m := &mocks.MockDBManager{}
				m.On("GetDBCreds").Return("", "", "", "", errors.New("test error"))
				return m
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "Document created",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"orderAmount": 12, "packageSizes": [5, 10, 15, 20, 25]}`,
			mockDBManager: func() *mocks.MockDBManager {
				m := &mocks.MockDBManager{}
				m.On("GetDBCreds").Return("bucket", "scope", "collection", "document", nil)
				m.On("GetDocument", "bucket", "scope", "collection", "document").Return(nil, errors.New("test error"))
				m.On("WriteDocument", "bucket", "scope", "collection", "document", db.Document{}).Return(nil)
				return m
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/setDocument", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()

			dbManager := tt.mockDBManager()

			SetDocumentHandler(rr, req, dbManager)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestGetDocumentHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		contentType    string
		body           string
		mockDBManager  func() *mocks.MockDBManager
		expectedStatus int
	}{
		{
			name:           "Method not allowed",
			method:         http.MethodPost,
			contentType:    "application/json",
			body:           `{"documentID": "test"}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Content-Type header is not application/json",
			method:         http.MethodGet,
			contentType:    "text/plain",
			body:           `{"documentID": "test"}`,
			mockDBManager:  func() *mocks.MockDBManager { return &mocks.MockDBManager{} },
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:        "Could not get database credentials",
			method:      http.MethodGet,
			contentType: "application/json",
			body:        `{"documentID": "test"}`,
			mockDBManager: func() *mocks.MockDBManager {
				m := &mocks.MockDBManager{}
				m.On("GetDBCreds").Return("", "", "", "", errors.New("test error"))
				return m
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "Document retrieved",
			method:      http.MethodGet,
			contentType: "application/json",
			body:        `{"documentID": "test"}`,
			mockDBManager: func() *mocks.MockDBManager {
				m := &mocks.MockDBManager{}
				m.On("GetDBCreds").Return("bucket", "scope", "collection", "document", nil)
				m.On("GetDocument", "bucket", "scope", "collection", "document").Return(&db.DocumentHistory{}, nil)
				return m
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/getDocument", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()

			dbManager := tt.mockDBManager()

			GetDocumentHandler(rr, req, dbManager)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
