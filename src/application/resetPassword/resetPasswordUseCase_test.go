package resetPassword

import (
	"context"
	"errors"
	"fmt"
	"go-iam/mocks"
	"go-iam/src/domain/user"
	"go-iam/src/infrastructure/logging"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	UserRepository              *mocks.UserRepository
	UserPasswordResetRepository *mocks.UserPasswordResetRepository
	Hasher                      *mocks.Hasher
	HashComparator              *mocks.HashComparator
	UseCase                     *ResetPasswordUseCase
}

func setUp(t *testing.T) testCase {
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	userRepository := mocks.NewUserRepository(t)
	userPasswordResetRepository := mocks.NewUserPasswordResetRepository(t)
	hasher := mocks.NewHasher(t)
	hashComparator := mocks.NewHashComparator(t)

	return testCase{
		UserRepository:              userRepository,
		UserPasswordResetRepository: userPasswordResetRepository,
		Hasher:                      hasher,
		HashComparator:              hashComparator,
		UseCase:                     NewResetPasswordUseCase(userRepository, userPasswordResetRepository, hashComparator, hasher, logger),
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
	testCase.UserPasswordResetRepository.AssertNotCalled(t, "FindByUserID")
	testCase.Hasher.AssertNotCalled(t, "Hash")
	testCase.HashComparator.AssertNotCalled(t, "Compare")
}

func TestExecuteFindUserError(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testEmail"
	request := ResetPasswordRequest{
		UserEmail: testEmail,
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
	testCase.HashComparator.AssertNotCalled(t, "Compare")
}

func TestExecuteFindPasswordResetError(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testEmail"
	request := ResetPasswordRequest{
		UserEmail: testEmail,
	}
	ctx := context.Background()
	testUser := user.User{
		ID:       uuid.New(),
		Username: "testUser",
		Email:    testEmail,
	}
	testError := errors.New("Test find password reset error")
	testCase.UserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)
	testCase.UserPasswordResetRepository.On("FindByUserID", mock.Anything, mock.Anything).Return(nil, testError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as user password reset repository find")
	}
	testCase.UserRepository.AssertCalled(t, "FindByEmail", ctx, testEmail)
	testCase.UserPasswordResetRepository.AssertCalled(t, "FindByUserID", ctx, testUser.ID)
	testCase.Hasher.AssertNotCalled(t, "Hash")
	testCase.HashComparator.AssertNotCalled(t, "Compare")
}

func TestExecuteResetTokenExpired(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testEmail"
	testResetToken := "testResetToken"
	request := ResetPasswordRequest{
		UserEmail:  testEmail,
		ResetToken: testResetToken,
	}
	ctx := context.Background()
	testUser := user.User{
		ID:       uuid.New(),
		Username: "testUser",
		Email:    testEmail,
	}
	testUserPasswordReset := user.UserPasswordReset{
		Token:      "testResetToken",
		Expiration: time.Now().Add(-1 * time.Hour),
		UserID:     testUser.ID,
	}
	testCase.UserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)
	testCase.UserPasswordResetRepository.On("FindByUserID", mock.Anything, mock.Anything).Return(&testUserPasswordReset, nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	expectedError := fmt.Sprintf("reset token %s is expired", testResetToken)
	if response.Err.Error() != expectedError {
		t.Fatal("Expected use case to return reset token expired error")
	}
	testCase.UserRepository.AssertCalled(t, "FindByEmail", ctx, testEmail)
	testCase.UserPasswordResetRepository.AssertCalled(t, "FindByUserID", ctx, testUser.ID)
	testCase.Hasher.AssertNotCalled(t, "Hash")
	testCase.HashComparator.AssertNotCalled(t, "Compare")
}

func TestExecuteResetTokenComparisionFailed(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testEmail"
	testResetToken := "testResetToken"
	request := ResetPasswordRequest{
		UserEmail:  testEmail,
		ResetToken: testResetToken,
	}
	ctx := context.Background()
	testUser := user.User{
		ID:       uuid.New(),
		Username: "testUser",
		Email:    testEmail,
	}
	testUserPasswordReset := user.UserPasswordReset{
		Token:      "testResetToken",
		Expiration: time.Now().Add(1 * time.Hour),
		UserID:     testUser.ID,
	}
	testCase.UserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)
	testCase.UserPasswordResetRepository.On("FindByUserID", mock.Anything, mock.Anything).Return(&testUserPasswordReset, nil)
	testError := errors.New("Test hash comparision error")
	testCase.HashComparator.On("Compare", mock.Anything, mock.Anything).Return(testError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as hash comparision one")
	}
	testCase.UserRepository.AssertCalled(t, "FindByEmail", ctx, testEmail)
	testCase.UserPasswordResetRepository.AssertCalled(t, "FindByUserID", ctx, testUser.ID)
	testCase.Hasher.AssertNotCalled(t, "Hash")
	testCase.HashComparator.AssertCalled(t, "Compare", request.ResetToken, testUserPasswordReset.Token)
}

func TestExecutePasswordHashError(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testEmail"
	testResetToken := "testResetToken"
	testNewPassword := "testNewPassword"
	request := ResetPasswordRequest{
		UserEmail:   testEmail,
		ResetToken:  testResetToken,
		NewPassword: testNewPassword,
	}
	ctx := context.Background()
	testUser := user.User{
		ID:       uuid.New(),
		Username: "testUser",
		Email:    testEmail,
	}
	testUserPasswordReset := user.UserPasswordReset{
		Token:      "testResetToken",
		Expiration: time.Now().Add(1 * time.Hour),
		UserID:     testUser.ID,
	}
	testCase.UserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)
	testCase.UserPasswordResetRepository.On("FindByUserID", mock.Anything, mock.Anything).Return(&testUserPasswordReset, nil)
	testCase.HashComparator.On("Compare", mock.Anything, mock.Anything).Return(nil)
	testError := errors.New("Test hashing error")
	testCase.Hasher.On("Hash", mock.Anything).Return(nil, testError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as hashing one")
	}
	testCase.UserRepository.AssertCalled(t, "FindByEmail", ctx, testEmail)
	testCase.UserPasswordResetRepository.AssertCalled(t, "FindByUserID", ctx, testUser.ID)
	testCase.Hasher.AssertCalled(t, "Hash", testNewPassword)
	testCase.HashComparator.AssertCalled(t, "Compare", request.ResetToken, testUserPasswordReset.Token)
}

func TestExecuteDeleteResetError(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testEmail"
	testResetToken := "testResetToken"
	testNewPassword := "testNewPassword"
	testPasswordHash := "testPasswordHash"
	request := ResetPasswordRequest{
		UserEmail:   testEmail,
		ResetToken:  testResetToken,
		NewPassword: testNewPassword,
	}
	ctx := context.Background()
	testUser := user.User{
		ID:       uuid.New(),
		Username: "testUser",
		Email:    testEmail,
	}
	testUserPasswordReset := user.UserPasswordReset{
		Token:      "testResetToken",
		Expiration: time.Now().Add(1 * time.Hour),
		UserID:     testUser.ID,
	}
	testCase.UserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)
	testCase.UserPasswordResetRepository.On("FindByUserID", mock.Anything, mock.Anything).Return(&testUserPasswordReset, nil)
	testCase.HashComparator.On("Compare", mock.Anything, mock.Anything).Return(nil)
	testCase.Hasher.On("Hash", mock.Anything).Return(&testPasswordHash, nil)
	testCase.UserRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
	testError := errors.New("Test password reset delete error")
	testCase.UserPasswordResetRepository.On("Delete", mock.Anything, mock.Anything).Return(testError)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != testError {
		t.Fatal("Expected use case to return same error as reset repository deletion one")
	}
	testCase.UserRepository.AssertCalled(t, "FindByEmail", ctx, testEmail)
	testCase.UserRepository.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Email == testEmail && user.ID == testUser.ID && user.Password == testPasswordHash
	}))
	testCase.UserPasswordResetRepository.AssertCalled(t, "FindByUserID", ctx, testUser.ID)
	testCase.Hasher.AssertCalled(t, "Hash", testNewPassword)
	testCase.HashComparator.AssertCalled(t, "Compare", request.ResetToken, testUserPasswordReset.Token)
	testCase.UserPasswordResetRepository.AssertCalled(t, "Delete", ctx, testUserPasswordReset)
}

func TestExecuteSuccess(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testEmail"
	testResetToken := "testResetToken"
	testNewPassword := "testNewPassword"
	testPasswordHash := "testPasswordHash"
	request := ResetPasswordRequest{
		UserEmail:   testEmail,
		ResetToken:  testResetToken,
		NewPassword: testNewPassword,
	}
	ctx := context.Background()
	testUser := user.User{
		ID:       uuid.New(),
		Username: "testUser",
		Email:    testEmail,
	}
	testUserPasswordReset := user.UserPasswordReset{
		Token:      "testResetToken",
		Expiration: time.Now().Add(1 * time.Hour),
		UserID:     testUser.ID,
	}
	testCase.UserRepository.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)
	testCase.UserPasswordResetRepository.On("FindByUserID", mock.Anything, mock.Anything).Return(&testUserPasswordReset, nil)
	testCase.HashComparator.On("Compare", mock.Anything, mock.Anything).Return(nil)
	testCase.Hasher.On("Hash", mock.Anything).Return(&testPasswordHash, nil)
	testCase.UserRepository.On("Save", mock.Anything, mock.Anything).Return(nil)
	testCase.UserPasswordResetRepository.On("Delete", mock.Anything, mock.Anything).Return(nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	if response.Content != nil {
		t.Fatal("Expected use case to return empty response")
	}
	testCase.UserRepository.AssertCalled(t, "FindByEmail", ctx, testEmail)
	testCase.UserRepository.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Email == testEmail && user.ID == testUser.ID && user.Password == testPasswordHash
	}))
	testCase.UserPasswordResetRepository.AssertCalled(t, "FindByUserID", ctx, testUser.ID)
	testCase.Hasher.AssertCalled(t, "Hash", testNewPassword)
	testCase.HashComparator.AssertCalled(t, "Compare", request.ResetToken, testUserPasswordReset.Token)
	testCase.UserPasswordResetRepository.AssertCalled(t, "Delete", ctx, testUserPasswordReset)
}
