package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"pthd-notifications/pkg/domain"
	"pthd-notifications/pkg/domain/model"
	domain_mocks "pthd-notifications/tests/mocks/domain"
	"testing"
)

type NotificationHandlerTestSuite struct {
	suite.Suite

	controller *gomock.Controller
	service    *domain_mocks.MockIService

	handler *notificationHandler
}

func (s *NotificationHandlerTestSuite) SetupSuite() {
	s.controller = gomock.NewController(s.T())

	s.service = domain_mocks.NewMockIService(s.controller)

	s.handler = newNotificationHandler(s.service, initializeDecoder(), initializeValidator())
}

func (s *NotificationHandlerTestSuite) Test_Handler_NoParameters_ErrorResponse() {
	req, err := http.NewRequest("GET", "/api/v1/notification", nil)
	if err != nil {
		s.T().Errorf("Error creating a new request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.Handle)
	handler.ServeHTTP(rr, req)

	var response map[string]interface{}
	if decodeErr := json.NewDecoder(rr.Body).Decode(&response); decodeErr != nil {
		s.T().Errorf("Error decoding response: %v", decodeErr)
	}

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Equal(
		s.T(),
		response,
		map[string]interface{}{
			"reason": "validation error",
			"errors": map[string]interface{}{
				"type": "required",
			},
		},
	)
}

func (s *NotificationHandlerTestSuite) Test_Handler_UsersConnectedMissingParams_ErrorResponse() {
	req, err := http.NewRequest("GET", "/api/v1/notification?type=users_connected", nil)
	if err != nil {
		s.T().Errorf("Error creating a new request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.Handle)
	handler.ServeHTTP(rr, req)

	var response map[string]interface{}
	if decodeErr := json.NewDecoder(rr.Body).Decode(&response); decodeErr != nil {
		s.T().Errorf("Error decoding response: %v", decodeErr)
	}

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Equal(
		s.T(),
		response,
		map[string]interface{}{
			"reason": "validation error",
			"errors": map[string]interface{}{
				"discord_id": "required",
				"usernames":  "required",
			},
		},
	)
}

func (s *NotificationHandlerTestSuite) Test_Handler_UsersConnectedNoSettings_ErrorResponse() {
	req, err := http.NewRequest(
		"GET",
		"/api/v1/notification?type=users_connected&discord_id=1&usernames=test",
		nil,
	)
	if err != nil {
		s.T().Errorf("Error creating a new request: %v", err)
	}

	serviceErr := &domain.ErrNoSettings{
		DiscordId: 1,
		Type:      "users_connected",
	}
	s.service.EXPECT().SendNotification(gomock.Any()).Return(serviceErr)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.Handle)
	handler.ServeHTTP(rr, req)

	var response map[string]interface{}
	if decodeErr := json.NewDecoder(rr.Body).Decode(&response); decodeErr != nil {
		s.T().Errorf("Error decoding response: %v", decodeErr)
	}

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Equal(
		s.T(),
		response,
		map[string]interface{}{
			"reason": "send error",
			"errors": map[string]interface{}{
				"parameters": "no settings for such parameters",
			},
		},
	)
}

func (s *NotificationHandlerTestSuite) Test_Handler_UsersConnectedNoMessage_ErrorResponse() {
	req, err := http.NewRequest(
		"GET",
		"/api/v1/notification?type=users_connected&discord_id=1&usernames=test",
		nil,
	)
	if err != nil {
		s.T().Errorf("Error creating a new request: %v", err)
	}

	serviceErr := &model.ErrNoMessage{
		DiscordId: 1,
		Type:      "users_connected",
	}
	s.service.EXPECT().SendNotification(gomock.Any()).Return(serviceErr)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.Handle)
	handler.ServeHTTP(rr, req)

	var response map[string]interface{}
	if decodeErr := json.NewDecoder(rr.Body).Decode(&response); decodeErr != nil {
		s.T().Errorf("Error decoding response: %v", decodeErr)
	}

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Equal(
		s.T(),
		response,
		map[string]interface{}{
			"reason": "send error",
			"errors": map[string]interface{}{
				"message": "no message for such parameters",
			},
		},
	)
}

func (s *NotificationHandlerTestSuite) Test_Handler_UsersConnectedSent_StatusOk() {
	req, err := http.NewRequest(
		"GET",
		"/api/v1/notification?type=users_connected&discord_id=1&usernames=test",
		nil,
	)
	if err != nil {
		s.T().Errorf("Error creating a new request: %v", err)
	}

	s.service.EXPECT().SendNotification(gomock.Any()).Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.Handle)
	handler.ServeHTTP(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
}

func (s *NotificationHandlerTestSuite) Test_Handler_UsersLeftChannelMissingParams_ErrorResponse() {
	req, err := http.NewRequest(
		"GET",
		"/api/v1/notification?type=users_left_channel",
		nil,
	)
	if err != nil {
		s.T().Errorf("Error creating a new request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.Handle)
	handler.ServeHTTP(rr, req)

	var response map[string]interface{}
	if decodeErr := json.NewDecoder(rr.Body).Decode(&response); decodeErr != nil {
		s.T().Errorf("Error decoding response: %v", decodeErr)
	}

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Equal(
		s.T(),
		response,
		map[string]interface{}{
			"reason": "validation error",
			"errors": map[string]interface{}{
				"discord_id": "required",
			},
		},
	)
}

func (s *NotificationHandlerTestSuite) Test_Handler_UsersLeftChannelMissingParams_StatusOk() {
	req, err := http.NewRequest(
		"GET",
		"/api/v1/notification?type=users_left_channel&discord_id=1",
		nil,
	)
	if err != nil {
		s.T().Errorf("Error creating a new request: %v", err)
	}

	s.service.EXPECT().SendNotification(gomock.Any()).Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.Handle)
	handler.ServeHTTP(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
}

func TestNotificationHandlerTestSuite(t *testing.T) {
	testSuite := NotificationHandlerTestSuite{}
	suite.Run(t, &testSuite)
}
