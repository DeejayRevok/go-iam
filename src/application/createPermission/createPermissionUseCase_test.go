package createPermission

import (
	"errors"
	"go-uaa/mocks"
	"go-uaa/src/domain/permission"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type testCase struct {
	PermissionRepo *mocks.PermissionRepository
	UseCase        *CreatePermissionUseCase
}

func setUp(t *testing.T) testCase {
	logger, _ := zap.NewDevelopment()
	permissionRepoMock := mocks.NewPermissionRepository(t)
	return testCase{
		PermissionRepo: permissionRepoMock,
		UseCase:        NewCreatePermissionUseCase(permissionRepoMock, logger),
	}
}

func TestExecuteWrongRequest(t *testing.T) {
	testCase := setUp(t)
	request := "wrongRequest"

	response := testCase.UseCase.Execute(request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	testCase.PermissionRepo.AssertNotCalled(t, "Save")
}

func TestExecutePermissionSaveError(t *testing.T) {
	testCase := setUp(t)
	saveError := errors.New("Test save error")
	testCase.PermissionRepo.On("Save", mock.Anything).Return(saveError)
	permissionName := "testPermission"
	request := CreatePermissionRequest{
		Name: permissionName,
	}

	response := testCase.UseCase.Execute(&request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != saveError {
		t.Fatal("Error expected to be the same as the repository returned error")
	}
	expectedSavePermission := permission.Permission{
		Name: permissionName,
	}
	testCase.PermissionRepo.AssertCalled(t, "Save", expectedSavePermission)
}

func TestExecuteSuccess(t *testing.T) {
	testCase := setUp(t)
	testCase.PermissionRepo.On("Save", mock.Anything).Return(nil)
	permissionName := "testPermission"
	request := CreatePermissionRequest{
		Name: permissionName,
	}

	response := testCase.UseCase.Execute(&request)

	if response.Err != nil {
		t.Fatal("Expected use case to not return error")
	}
	expectedSavePermission := permission.Permission{
		Name: permissionName,
	}
	testCase.PermissionRepo.AssertCalled(t, "Save", expectedSavePermission)
}
