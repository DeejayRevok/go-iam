package requestPasswordReset

import (
	"context"
	"errors"
	"go-uaa/mocks"
	"go-uaa/src/domain/user"
	"go-uaa/src/infrastructure/logging"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	UserRepository              *mocks.UserRepository
	UserPasswordResetRepository *mocks.UserPasswordResetRepository
	Hasher                      *mocks.Hasher
	EventPublisher              *mocks.EventPublisher
	UseCase                     *RequestPasswordResetUseCase
}

func setUp(t *testing.T) testCase {
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	userRepository := mocks.NewUserRepository(t)
	userPasswordResetRepository := mocks.NewUserPasswordResetRepository(t)
	hasher := mocks.NewHasher(t)
	eventPublisher := mocks.NewEventPublisher(t)

	return testCase{
		UserRepository:              userRepository,
		UserPasswordResetRepository: userPasswordResetRepository,
		Hasher:                      hasher,
		EventPublisher:              eventPublisher,
		UseCase:                     NewRequestPasswordResetUseCase(userRepository, userPasswordResetRepository, eventPublisher, hasher, logger),
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
	testCase.UserPasswordResetRepository.AssertNotCalled(t, "Save")
	testCase.Hasher.AssertNotCalled(t, "Hash")
	testCase.EventPublisher.AssertNotCalled(t, "Publish")
}

func TestExecuteFindUserError(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testEmail"
	request := RequestPasswordResetRequest{
		Email: testEmail,
	}
	ctx := context.Background()
	testError := errors.New("Test find error")
	testCase.UserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, testError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as user repository find")
	}
	testCase.UserRepository.AssertCalled(t, "FindByEmail", ctx, testEmail)
	testCase.UserPasswordResetRepository.AssertNotCalled(t, "Save")
	testCase.Hasher.AssertNotCalled(t, "Hash")
	testCase.EventPublisher.AssertNotCalled(t, "Publish")
}

func TestExecuteHashingError(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testEmail"
	request := RequestPasswordResetRequest{
		Email: testEmail,
	}
	ctx := context.Background()
	testUser := user.User{
		Username: "testUser",
		Email:    testEmail,
	}
	testCase.UserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)
	testError := errors.New("Test hash error")
	testCase.Hasher.On("Hash", mock.Anything).Return(nil, testError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as hasher hash error")
	}
	testCase.UserRepository.AssertCalled(t, "FindByEmail", ctx, testEmail)
	testCase.UserPasswordResetRepository.AssertNotCalled(t, "Save")
	testCase.Hasher.AssertCalled(t, "Hash", mock.Anything)
	testCase.EventPublisher.AssertNotCalled(t, "Publish")
}

func TestExecuteUserPasswordResetSaveError(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testEmail"
	request := RequestPasswordResetRequest{
		Email: testEmail,
	}
	ctx := context.Background()
	testUser := user.User{
		ID:       uuid.New(),
		Username: "testUser",
		Email:    testEmail,
	}
	testCase.UserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)
	testResetTokenHash := "testHash"
	testCase.Hasher.On("Hash", mock.Anything).Return(&testResetTokenHash, nil)
	testError := errors.New("Test hash error")
	testCase.UserPasswordResetRepository.On("Save", mock.Anything, mock.Anything).Return(testError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as user password reset repository save error")
	}
	testCase.UserRepository.AssertCalled(t, "FindByEmail", ctx, testEmail)
	testCase.UserPasswordResetRepository.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(reset user.UserPasswordReset) bool {
		return reset.Token == testResetTokenHash && reset.UserID == testUser.ID
	}))
	testCase.Hasher.AssertCalled(t, "Hash", mock.Anything)
	testCase.EventPublisher.AssertNotCalled(t, "Publish")
}

func TestExecuteEventPublishError(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testEmail"
	request := RequestPasswordResetRequest{
		Email: testEmail,
	}
	ctx := context.Background()
	testUser := user.User{
		ID:       uuid.New(),
		Username: "testUser",
		Email:    testEmail,
	}
	testCase.UserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)
	testResetTokenHash := "testHash"
	testCase.Hasher.On("Hash", mock.Anything).Return(&testResetTokenHash, nil)
	testCase.UserPasswordResetRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
	testError := errors.New("Test hash error")
	testCase.EventPublisher.On("Publish", mock.Anything).Return(testError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as event publish error")
	}
	testCase.UserRepository.AssertCalled(t, "FindByEmail", ctx, testEmail)
	testCase.UserPasswordResetRepository.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(reset user.UserPasswordReset) bool {
		return reset.Token == testResetTokenHash && reset.UserID == testUser.ID
	}))
	testCase.Hasher.AssertCalled(t, "Hash", mock.Anything)
	testCase.EventPublisher.AssertCalled(t, "Publish", mock.MatchedBy(func(event user.UserPasswordResetRequestedEvent) bool {
		return event.UserID == testUser.ID.String()
	}))
}

func TestExecuteSuccess(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testEmail"
	request := RequestPasswordResetRequest{
		Email: testEmail,
	}
	ctx := context.Background()
	testUser := user.User{
		ID:       uuid.New(),
		Username: "testUser",
		Email:    testEmail,
	}
	testCase.UserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)
	testResetTokenHash := "testHash"
	testCase.Hasher.On("Hash", mock.Anything).Return(&testResetTokenHash, nil)
	testCase.UserPasswordResetRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
	testCase.EventPublisher.On("Publish", mock.Anything).Return(nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	if response.Content != nil {
		t.Fatal("Expected use case to return an empty response")
	}
	testCase.UserRepository.AssertCalled(t, "FindByEmail", ctx, testEmail)
	testCase.UserPasswordResetRepository.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(reset user.UserPasswordReset) bool {
		return reset.Token == testResetTokenHash && reset.UserID == testUser.ID
	}))
	testCase.Hasher.AssertCalled(t, "Hash", mock.Anything)
	testCase.EventPublisher.AssertCalled(t, "Publish", mock.MatchedBy(func(event user.UserPasswordResetRequestedEvent) bool {
		return event.UserID == testUser.ID.String()
	}))
}
