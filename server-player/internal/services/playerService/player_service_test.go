package playerService

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
	mocksvc "github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/services/playerService/mocks"
)

type PlayerServiceSuite struct {
	suite.Suite
	ctx      context.Context
	ctrl     *gomock.Controller
	storage  *mocksvc.MockPlayerStorage
	producer *mocksvc.MockQuizResultProducer
	svc      *PlayerService
}

func (s *PlayerServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.ctrl = gomock.NewController(s.T())
	s.storage = mocksvc.NewMockPlayerStorage(s.ctrl)
	s.producer = mocksvc.NewMockQuizResultProducer(s.ctrl)
	s.svc = NewPlayerService(s.ctx, s.storage, s.producer)
}

func (s *PlayerServiceSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *PlayerServiceSuite) TestGetAllTopics_Success() {
	expected := []*models.ActiveTopics{{}, {}}
	s.storage.EXPECT().GetAllTopics(gomock.Any()).Return(expected, nil)

	got, err := s.svc.GetAllTopics(s.ctx)

	s.Require().NoError(err)
	s.Equal(expected, got)
}

func (s *PlayerServiceSuite) TestGetAllTopics_Error() {
	wantErr := errors.New("db error")

	s.storage.EXPECT().GetAllTopics(gomock.Any()).Return([]*models.ActiveTopics(nil), wantErr)

	got, err := s.svc.GetAllTopics(s.ctx)

	s.ErrorIs(err, wantErr)
	s.Nil(got)
}

func (s *PlayerServiceSuite) TestGetTenQuestionsByTopic_Success() {
	topicName := "topic"
	expectedQuestions := []*models.Question{{}, {}}

	s.storage.EXPECT().GetTopicIdByName(gomock.Any(), topicName).Return(42, nil)
	s.storage.EXPECT().GetTenQuestionsByTopicID(gomock.Any(), 42).Return(expectedQuestions, nil)

	got, tid, err := s.svc.GetTenQuestionsByTopic(s.ctx, topicName)

	s.Require().NoError(err)
	s.Equal(42, tid)
	s.Equal(expectedQuestions, got)
}

func (s *PlayerServiceSuite) TestGetTenQuestionsByTopic_TopicIdError() {
	topicName := "topic"
	wantErr := errors.New("not found")

	s.storage.EXPECT().GetTopicIdByName(gomock.Any(), topicName).Return(-1, wantErr)

	got, tid, err := s.svc.GetTenQuestionsByTopic(s.ctx, topicName)

	s.ErrorIs(err, wantErr)
	s.Nil(got)
	s.Equal(-1, tid)
}

func (s *PlayerServiceSuite) TestGetTenQuestionsByTopic_QuestionsError() {
	topicName := "topic"
	wantErr := errors.New("db error")

	s.storage.EXPECT().GetTopicIdByName(gomock.Any(), topicName).Return(7, nil)
	s.storage.EXPECT().GetTenQuestionsByTopicID(gomock.Any(), 7).Return([]*models.Question(nil), wantErr)

	got, tid, err := s.svc.GetTenQuestionsByTopic(s.ctx, topicName)

	s.ErrorIs(err, wantErr)
	s.Nil(got)
	s.Equal(-1, tid)
}

func (s *PlayerServiceSuite) TestSetResultsByQuiz_Success() {
	req := &models.QuizRequest{}
	s.producer.EXPECT().SendQuizResult(gomock.Any(), req).Return(nil)

	err := s.svc.SetResultsByQuiz(s.ctx, req)

	s.NoError(err)
}

func (s *PlayerServiceSuite) TestSetResultsByQuiz_ProducerError() {
	req := &models.QuizRequest{}
	wantErr := errors.New("kafka error")

	s.producer.EXPECT().SendQuizResult(gomock.Any(), req).Return(wantErr)

	err := s.svc.SetResultsByQuiz(s.ctx, req)

	s.ErrorIs(err, wantErr)
}

func TestPlayerServiceSuite(t *testing.T) {
	suite.Run(t, new(PlayerServiceSuite))
}
