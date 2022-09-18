package createUser

import (
	"context"
	"errors"
	"go-uaa/mocks"
	"go-uaa/src/domain/role"
	"go-uaa/src/domain/user"
	"go-uaa/src/infrastructure/logging"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	UserRepo       *mocks.UserRepository
	PasswordHasher *mocks.Hasher
	RoleRepo       *mocks.RoleRepository
	EventPublisher *mocks.EventPublisher
	EmailValidator *mocks.EmailValidator
	UseCase        *CreateUserUseCase
}

func setUp(t *testing.T) testCase {
	userRepositoryMock := mocks.NewUserRepository(t)
	passwordHasherMock := mocks.NewHasher(t)
	roleRepositoryMock := mocks.NewRoleRepository(t)
	eventPublisherMock := mocks.NewEventPublisher(t)
	emailValidatorMock := mocks.NewEmailValidator(t)
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	return testCase{
		UserRepo:       userRepositoryMock,
		PasswordHasher: passwordHasherMock,
		RoleRepo:       roleRepositoryMock,
		EventPublisher: eventPublisherMock,
		EmailValidator: emailValidatorMock,
		UseCase:        NewCreateUserUseCase(userRepositoryMock, passwordHasherMock, roleRepositoryMock, eventPublisherMock, emailValidatorMock, logger),
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
	testCase.RoleRepo.AssertNotCalled(t, "FindByIDs")
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
		Roles:     make([]string, 0),
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
	testCase.RoleRepo.AssertNotCalled(t, "FindByIDs")
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
		Roles:     make([]string, 0),
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
	testCase.RoleRepo.AssertNotCalled(t, "FindByIDs")
	testCase.UserRepo.AssertNotCalled(t, "Save")
	testCase.EventPublisher.AssertNotCalled(t, "Publish")
}

func TestExecuteFindRolesError(t *testing.T) {
	testCase := setUp(t)
	testEmail := "testValidEmail"
	testCase.EmailValidator.On("Validate", mock.Anything).Return(nil)
	testPassword := "testPassword"
	testPasswordHash := "testPasswordHash"
	testCase.PasswordHasher.On("Hash", mock.Anything).Return(&testPasswordHash, nil)
	testUUID1, _ := uuid.NewUUID()
	testUUID2, _ := uuid.NewUUID()
	testUUID1Str := testUUID1.String()
	testUUID2Str := testUUID2.String()
	roleIDs := []uuid.UUID{testUUID1, testUUID2}
	roleIDsStr := []string{testUUID1Str, testUUID2Str}
	findError := errors.New("Test find roles error")
	testCase.RoleRepo.On("FindByIDs", mock.Anything, mock.Anything).Return(nil, findError)
	request := &CreateUserRequest{
		Username:  "Test",
		Email:     testEmail,
		Password:  testPassword,
		Roles:     roleIDsStr,
		Superuser: false,
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != findError {
		t.Fatal("Expected use case to return role repository find error")
	}
	testCase.EmailValidator.AssertCalled(t, "Validate", testEmail)
	testCase.PasswordHasher.AssertCalled(t, "Hash", testPassword)
	testCase.RoleRepo.AssertCalled(t, "FindByIDs", ctx, roleIDs)
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
	testUUID1, _ := uuid.NewUUID()
	testUUID2, _ := uuid.NewUUID()
	testUUID1Str := testUUID1.String()
	testUUID2Str := testUUID2.String()
	roleIDs := []uuid.UUID{testUUID1, testUUID2}
	roleIDsStr := []string{testUUID1Str, testUUID2Str}
	testRole := role.Role{
		ID:   testUUID1,
		Name: "Test role",
	}
	roles := []role.Role{testRole, testRole}
	testCase.RoleRepo.On("FindByIDs", mock.Anything, mock.Anything).Return(roles, nil)
	saveError := errors.New("Test find roles error")
	testCase.UserRepo.On("Save", mock.Anything, mock.Anything).Return(saveError)
	testUsername := "Test user name"
	testIsSuperuser := false
	request := &CreateUserRequest{
		Username:  testUsername,
		Email:     testEmail,
		Password:  testPassword,
		Roles:     roleIDsStr,
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
	testCase.RoleRepo.AssertCalled(t, "FindByIDs", ctx, roleIDs)
	testCase.UserRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Username == testUsername && reflect.DeepEqual(user.Roles, roles) && user.Email == testEmail && user.Password == testPasswordHash && user.Superuser == testIsSuperuser
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
	testUUID1, _ := uuid.NewUUID()
	testUUID2, _ := uuid.NewUUID()
	testUUID1Str := testUUID1.String()
	testUUID2Str := testUUID2.String()
	roleIDs := []uuid.UUID{testUUID1, testUUID2}
	roleIDsStr := []string{testUUID1Str, testUUID2Str}
	testRole := role.Role{
		ID:   testUUID1,
		Name: "Test role",
	}
	roles := []role.Role{testRole, testRole}
	testCase.RoleRepo.On("FindByIDs", mock.Anything, mock.Anything).Return(roles, nil)
	testCase.UserRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	testCase.EventPublisher.On("Publish", mock.Anything).Return(nil)
	testUsername := "Test user name"
	testIsSuperuser := false
	request := &CreateUserRequest{
		Username:  testUsername,
		Email:     testEmail,
		Password:  testPassword,
		Roles:     roleIDsStr,
		Superuser: testIsSuperuser,
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	testCase.EmailValidator.AssertCalled(t, "Validate", testEmail)
	testCase.PasswordHasher.AssertCalled(t, "Hash", testPassword)
	testCase.RoleRepo.AssertCalled(t, "FindByIDs", ctx, roleIDs)
	testCase.UserRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(user user.User) bool {
		return user.Username == testUsername && reflect.DeepEqual(user.Roles, roles) && user.Email == testEmail && user.Password == testPasswordHash && user.Superuser == testIsSuperuser
	}))
	testCase.EventPublisher.AssertCalled(t, "Publish", mock.MatchedBy(func(event user.UserCreatedEvent) bool {
		return event.Username == testUsername && event.Email == testEmail
	}))
}
