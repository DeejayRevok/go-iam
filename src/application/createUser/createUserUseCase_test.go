package createUser

import (
	"context"
	"errors"
	"go-iam/mocks"
	"go-iam/src/domain/user"
	"go-iam/src/infrastructure/logging"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	UserRepo       *mocks.UserRepository
	PasswordHasher *mocks.Hasher
	EventPublisher *mocks.EventPublisher
	EmailValidator *mocks.EmailValidator
	UseCase        *CreateUserUseCase
}

func setUp(t *testing.T) testCase {
	userRepositoryMock := mocks.NewUserRepository(t)
	passwordHasherMock := mocks.NewHasher(t)
	eventPublisherMock := mocks.NewEventPublisher(t)
	emailValidatorMock := mocks.NewEmailValidator(t)
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	return testCase{
		UserRepo:       userRepositoryMock,
		PasswordHasher: passwordHasherMock,
		EventPublisher: eventPublisherMock,
		EmailValidator: emailValidatorMock,
		UseCase:        NewCreateUserUseCase(userRepositoryMock, passwordHasherMock, eventPublisherMock, emailValidatorMock, logger),
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
	testCase.EmailValidator.AssertNotCalled(t, "Validate")
	testCase.PasswordHasher.AssertNotCalled(t, "Hash")
	testCase.UserRepo.AssertNotCalled(t, "Save")
	testCase.EventPublisher.AssertNotCalled(t, "Publish")
}

func TestExecuteEmailNotValid(t *testing.T) {
	testCase := setUp(t)
	validationError := errors.New("Test email validation error")
	testEmail := "testWrongEmail"
	testCase.EmailValidator.On("Validate", mock.Anything).Return(validationError)
	request := &CreateUserRequest{
		Username:  "Test",
		Email:     testEmail,
		Password:  "Test",
		Superuser: false,
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != validationError {
		t.Fatal("Expected use case to return email validation error")
	}
	testCase.EmailValidator.AssertCalled(t, "Validate", testEmail)
	testCase.PasswordHasher.AssertNotCalled(t, "Hash")
	testCase.UserRepo.AssertNotCalled(t, "Save")
	testCase.EventPublisher.AssertNotCalled(t, "Publish")
}

func TestExecutePasswordHashError(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testValidEmail"
	testCase.EmailValidator.On("Validate", mock.Anything).Return(nil)
	testPassword := "testPassword"
	hashError := errors.New("Test password hash error")
	testCase.PasswordHasher.On("Hash", mock.Anything).Return(nil, hashError)
	request := &CreateUserRequest{
		Username:  "Test",
		Email:     testEmail,
		Password:  testPassword,
		Superuser: false,
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != hashError {
		t.Fatal("Expected use case to return email validation error")
	}
	testCase.EmailValidator.AssertCalled(t, "Validate", testEmail)
	testCase.PasswordHasher.AssertCalled(t, "Hash", testPassword)
	testCase.UserRepo.AssertNotCalled(t, "Save")
	testCase.EventPublisher.AssertNotCalled(t, "Publish")
}

func TestExecuteSaveError(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testValidEmail"
	testCase.EmailValidator.On("Validate", mock.Anything).Return(nil)
	testPassword := "testPassword"
	testPasswordHash := "testPasswordHash"
	testCase.PasswordHasher.On("Hash", mock.Anything).Return(&testPasswordHash, nil)
	saveError := errors.New("Test save error")
	testCase.UserRepo.On("Save", mock.Anything, mock.Anything).Return(saveError)
	testUsername := "Test user name"
	testIsSuperuser := false
	request := &CreateUserRequest{
		Username:  testUsername,
		Email:     testEmail,
		Password:  testPassword,
		Superuser: testIsSuperuser,
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != saveError {
		t.Fatal("Expected use case to return user repository save error")
	}
	testCase.EmailValidator.AssertCalled(t, "Validate", testEmail)
	testCase.PasswordHasher.AssertCalled(t, "Hash", testPassword)
	testCase.UserRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Username == testUsername && user.Email == testEmail && user.Password == testPasswordHash && user.Superuser == testIsSuperuser
	}))
	testCase.EventPublisher.AssertNotCalled(t, "Publish")
}

func TestExecuteSuccess(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testValidEmail"
	testCase.EmailValidator.On("Validate", mock.Anything).Return(nil)
	testPassword := "testPassword"
	testPasswordHash := "testPasswordHash"
	testCase.PasswordHasher.On("Hash", mock.Anything).Return(&testPasswordHash, nil)
	testCase.UserRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	testCase.EventPublisher.On("Publish", mock.Anything).Return(nil)
	testUsername := "Test user name"
	testIsSuperuser := false
	request := &CreateUserRequest{
		Username:  testUsername,
		Email:     testEmail,
		Password:  testPassword,
		Superuser: testIsSuperuser,
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	testCase.EmailValidator.AssertCalled(t, "Validate", testEmail)
	testCase.PasswordHasher.AssertCalled(t, "Hash", testPassword)
	testCase.UserRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Username == testUsername && user.Email == testEmail && user.Password == testPasswordHash && user.Superuser == testIsSuperuser
	}))
	testCase.EventPublisher.AssertCalled(t, "Publish", mock.MatchedBy(func(event user.UserCreatedEvent) bool {
		return event.Username == testUsername && event.Email == testEmail
	}))
}
