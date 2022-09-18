package getUser

import (
	"context"
	"errors"
	"go-uaa/mocks"
	"go-uaa/src/domain/user"
	"go-uaa/src/infrastructure/logging"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	UserRepo *mocks.UserRepository
	UseCase  *GetUserUseCase
}

func setUp(t *testing.T) testCase {
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	userRepositoryMock := mocks.NewUserRepository(t)
	return testCase{
		UserRepo: userRepositoryMock,
		UseCase:  NewGetUserUseCase(userRepositoryMock, logger),
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
	testCase.UserRepo.AssertNotCalled(t, "FindByID")
}

func TestExecuteFindError(t *testing.T) {
	testCase := setUp(t)
	testUserID := uuid.New()
	request := GetUserRequest{
		UserId: testUserID,
	}
	findError := errors.New("Test find user by id error")
	testCase.UserRepo.On("FindByID", mock.Anything).Return(nil, findError)
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != findError {
		t.Fatal("Expected use case to return user repo find error")
	}
	testCase.UserRepo.AssertCalled(t, "FindByID", testUserID)
}

func TestExecuteFindSuccess(t *testing.T) {
	testCase := setUp(t)
	testUserID := uuid.New()
	request := GetUserRequest{
		UserId: testUserID,
	}
	ctx := context.Background()
	testUser := user.User{
		Username: "Test username",
	}
	testCase.UserRepo.On("FindByID", mock.Anything).Return(&testUser, nil)

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	responseUser := *response.Content.(*user.User)
	if !reflect.DeepEqual(responseUser, testUser) {
		t.Fatal("Expected use case ro return same user as the repository")
	}
	testCase.UserRepo.AssertCalled(t, "FindByID", testUserID)
}
