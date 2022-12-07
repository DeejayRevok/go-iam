package sendPasswordResetToken

import (
	"context"
	"errors"
	"go-iam/mocks"
	"go-iam/src/domain/user"
	"go-iam/src/infrastructure/logging"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	UserRepository           *mocks.UserRepository
	PasswordResetTokenSender *mocks.PasswordResetTokenSender
	UseCase                  *SendPasswordResetTokenUseCase
}

func setUp(t *testing.T) testCase {
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	userRepository := mocks.NewUserRepository(t)
	resetTokenSender := mocks.NewPasswordResetTokenSender(t)

	return testCase{
		UserRepository:           userRepository,
		PasswordResetTokenSender: resetTokenSender,
		UseCase:                  NewSendPasswordResetTokenUseCase(userRepository, resetTokenSender, logger),
	}
}

func TestExecuteWrongRequest(t *testing.T) {
	testCase := setUp(t)
	request := "wrongRequest"
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	testCase.UserRepository.AssertNotCalled(t, "FindByEmail")
	testCase.PasswordResetTokenSender.AssertNotCalled(t, "Send")
}

func TestExecuteFindUserError(t *testing.T) {
	testCase := setUp(t)
	testReceiverID := uuid.New()
	testResetToken := "testResetToken"
	request := SendPasswordResetTokenRequest{
		UserID:     testReceiverID.String(),
		ResetToken: testResetToken,
	}
	ctx := context.Background()
	testError := errors.New("Test find error")
	testCase.UserRepository.On("FindByID", mock.Anything, mock.Anything).Return(nil, testError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as user repository find")
	}
	testCase.UserRepository.AssertCalled(t, "FindByID", ctx, testReceiverID)
	testCase.PasswordResetTokenSender.AssertNotCalled(t, "Send")
}

func TestExecuteSendError(t *testing.T) {
	testCase := setUp(t)
	testReceiverID := uuid.New()
	testResetToken := "testResetToken"
	request := SendPasswordResetTokenRequest{
		UserID:     testReceiverID.String(),
		ResetToken: testResetToken,
	}
	ctx := context.Background()
	testUser := user.User{
		ID:       testReceiverID,
		Username: "testUser",
	}
	testCase.UserRepository.On("FindByID", mock.Anything, mock.Anything).Return(&testUser, nil)
	testError := errors.New("Test send error")
	testCase.PasswordResetTokenSender.On("Send", mock.Anything, mock.Anything).Return(testError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as sending error")
	}
	testCase.UserRepository.AssertCalled(t, "FindByID", ctx, testReceiverID)
	testCase.PasswordResetTokenSender.AssertCalled(t, "Send", testResetToken, &testUser)
}

func TestExecuteSuccess(t *testing.T) {
	testCase := setUp(t)
	testReceiverID := uuid.New()
	testResetToken := "testResetToken"
	request := SendPasswordResetTokenRequest{
		UserID:     testReceiverID.String(),
		ResetToken: testResetToken,
	}
	ctx := context.Background()
	testUser := user.User{
		ID:       testReceiverID,
		Username: "testUser",
	}
	testCase.UserRepository.On("FindByID", mock.Anything, mock.Anything).Return(&testUser, nil)
	testCase.PasswordResetTokenSender.On("Send", mock.Anything, mock.Anything).Return(nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	if response.Content != nil {
		t.Fatal("Expected use case to return empty response")
	}
	testCase.UserRepository.AssertCalled(t, "FindByID", ctx, testReceiverID)
	testCase.PasswordResetTokenSender.AssertCalled(t, "Send", testResetToken, &testUser)
}
