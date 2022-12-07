package getAuthenticatedUser

import (
	"context"
	"errors"
	"go-iam/mocks"
	"go-iam/src/domain/auth/accessToken"
	"go-iam/src/domain/user"
	"go-iam/src/infrastructure/logging"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	UserRepo *mocks.UserRepository
	UseCase  *GetAuthenticatedUserUseCase
}

func setUp(t *testing.T) testCase {
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	userRepositoryMock := mocks.NewUserRepository(t)
	return testCase{
		UserRepo: userRepositoryMock,
		UseCase:  NewGetAuthenticatedUserUseCase(userRepositoryMock, logger),
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
	testCase.UserRepo.AssertNotCalled(t, "FindByEmail")
}

func TestExecuteFindError(t *testing.T) {
	testCase := setUp(t)
	findError := errors.New("Test find user error")
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, findError)
	testUsername := "TestUser"
	testToken := accessToken.AccessToken{
		Sub: testUsername,
	}
	request := GetAuthenticatedUserRequest{
		Token: &testToken,
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != findError {
		t.Fatal("Expected use case to return user repository find error")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, testUsername)
}

func TestExecuteSuccess(t *testing.T) {
	testCase := setUp(t)
	testUsername := "TestUser"
	testUser := user.User{
		Username: testUsername,
	}
	testCase.UserRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(&testUser, nil)
	testToken := accessToken.AccessToken{
		Sub: testUsername,
	}
	request := GetAuthenticatedUserRequest{
		Token: &testToken,
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	responseUser := *response.Content.(*user.User)
	if !reflect.DeepEqual(responseUser, testUser) {
		t.Fatal("Expected use case ro return same user as the repository")
	}
	testCase.UserRepo.AssertCalled(t, "FindByEmail", ctx, testUsername)
}
