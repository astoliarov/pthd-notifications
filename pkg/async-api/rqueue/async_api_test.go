package rqueue

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"pthd-notifications/pkg/domain/model"
	"pthd-notifications/tests/mocks"
	"strings"
	"testing"
)

type RedisAsyncAPITestSuite struct {
	suite.Suite

	controller    *gomock.Controller
	executorMock  *mocks.MockIExecutor
	connectorMock *mocks.MockIRedisConnector

	api *RedisAsyncAPI
}

func (s *RedisAsyncAPITestSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())

	s.executorMock = mocks.NewMockIExecutor(s.controller)
	s.connectorMock = mocks.NewMockIRedisConnector(s.controller)

	s.api = NewRedisAsyncAPI(
		s.executorMock,
		s.connectorMock,
	)
}

func (s *RedisAsyncAPITestSuite) Test_executeRead_NoMessage_NoErr() {
	s.connectorMock.EXPECT().ReadFromQueue(gomock.Any()).Return("", "", nil)

	readErr := s.api.executeRead(context.TODO())

	assert.Nil(s.T(), readErr)
}

func (s *RedisAsyncAPITestSuite) Test_executeRead_connectorErr_Err() {
	testErr := errors.New("test error")

	s.connectorMock.EXPECT().ReadFromQueue(gomock.Any()).Return("", "", testErr)

	readErr := s.api.executeRead(context.TODO())

	assert.Equal(s.T(), testErr, readErr)
}

func (s *RedisAsyncAPITestSuite) Test_executeRead_notExpectedMessage_NotSentNoErr() {
	data := `{"test": "test"}`

	s.connectorMock.EXPECT().ReadFromQueue(gomock.Any()).Return("", data, nil)

	readErr := s.api.executeRead(context.TODO())

	assert.Nil(s.T(), readErr)
}

func (s *RedisAsyncAPITestSuite) Test_executeRead_correctMessage_Sent() {
	data := `{"version":1,"type":"users_connected","data":{"usernames":["ticoncerto"]},"channel_id":1,"happened_at":"2023-09-12T17:44:00.418879Z"}`

	notificationContext := &model.UsersConnectedNotificationContext{
		Id:          1,
		Names:       []string{"ticoncerto"},
		NamesJoined: strings.Join([]string{"ticoncerto"}, ","),
	}

	s.connectorMock.EXPECT().ReadFromQueue(gomock.Any()).Return("", data, nil)
	s.executorMock.EXPECT().SendNotification(notificationContext).Return(nil)

	readErr := s.api.executeRead(context.TODO())

	assert.Nil(s.T(), readErr)
}

func (s *RedisAsyncAPITestSuite) Test_executeRead_correctMessageSendErr_NoErr() {
	data := `{"version":1,"type":"users_connected","data":{"usernames":["ticoncerto"]},"channel_id":1,"happened_at":"2023-09-12T17:44:00.418879Z"}`

	notificationContext := &model.UsersConnectedNotificationContext{
		Id:          1,
		Names:       []string{"ticoncerto"},
		NamesJoined: strings.Join([]string{"ticoncerto"}, ","),
	}

	testErr := errors.New("test error")

	s.connectorMock.EXPECT().ReadFromQueue(gomock.Any()).Return("", data, nil)
	s.executorMock.EXPECT().SendNotification(notificationContext).Return(testErr)

	readErr := s.api.executeRead(context.TODO())

	assert.Nil(s.T(), readErr)
}

func TestRedisAsyncAPI(t *testing.T) {

	testSuite := RedisAsyncAPITestSuite{}
	suite.Run(t, &testSuite)
}
