package createRole

import (
	"context"
	"errors"
	"go-uaa/mocks"
	"go-uaa/src/domain/permission"
	"go-uaa/src/domain/role"
	"go-uaa/src/infrastructure/logging"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.elastic.co/apm/v2"
)

type testCase struct {
	PermissionRepo *mocks.PermissionRepository
	RoleRepo       *mocks.RoleRepository
	UseCase        *CreateRoleUseCase
}

func setUp(t *testing.T) testCase {
	tracer := apm.DefaultTracer()
	logger := logging.NewZapTracedLogger(tracer)
	permissionRepoMock := mocks.NewPermissionRepository(t)
	roleRepoMock := mocks.NewRoleRepository(t)
	return testCase{
		PermissionRepo: permissionRepoMock,
		RoleRepo:       roleRepoMock,
		UseCase:        NewCreateRoleUseCase(roleRepoMock, permissionRepoMock, logger),
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
	testCase.RoleRepo.AssertNotCalled(t, "Save")
	testCase.PermissionRepo.AssertNotCalled(t, "FindByNames")
}

func TestExecuteRoleFindPermissionsError(t *testing.T) {
	testCase := setUp(t)
	findError := errors.New("Test find error")
	permissionName := "Test permission"
	testCase.PermissionRepo.On("FindByNames", mock.Anything, mock.Anything).Return(nil, findError)
	permissionsNames := []string{permissionName}
	roleName := "Test role"
	request := CreateRoleRequest{
		Name:        roleName,
		Permissions: permissionsNames,
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != findError {
		t.Fatal("Error expected to be the same as the permissions repository returned error")
	}
	testCase.PermissionRepo.AssertCalled(t, "FindByNames", ctx, permissionsNames)
	testCase.RoleRepo.AssertNotCalled(t, "Save")
}

func TestExecuteRoleSaveError(t *testing.T) {
	testCase := setUp(t)
	saveError := errors.New("Test save error")
	testCase.RoleRepo.On("Save", mock.Anything, mock.Anything).Return(saveError)
	permissionName := "Test permission"
	permissions := []permission.Permission{{Name: permissionName}}
	testCase.PermissionRepo.On("FindByNames", mock.Anything, mock.Anything).Return(permissions, nil)
	permissionsNames := []string{permissionName}
	roleName := "Test role"
	request := CreateRoleRequest{
		Name:        roleName,
		Permissions: permissionsNames,
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != saveError {
		t.Fatal("Error expected to be the same as the role repository returned error")
	}
	testCase.PermissionRepo.AssertCalled(t, "FindByNames", ctx, permissionsNames)
	testCase.RoleRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(role role.Role) bool {
		return role.Name == roleName && reflect.DeepEqual(role.Permissions, permissions)
	}))
}

func TestExecuteRoleSaveSuccess(t *testing.T) {
	testCase := setUp(t)
	testCase.RoleRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	permissionName := "Test permission"
	permissions := []permission.Permission{{Name: permissionName}}
	testCase.PermissionRepo.On("FindByNames", mock.Anything, mock.Anything).Return(permissions, nil)
	permissionsNames := []string{permissionName}
	roleName := "Test role"
	request := CreateRoleRequest{
		Name:        roleName,
		Permissions: permissionsNames,
	}
	ctx := context.Background()

	response := testCase.UseCase.Execute(ctx, &request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	testCase.PermissionRepo.AssertCalled(t, "FindByNames", ctx, permissionsNames)
	testCase.RoleRepo.AssertCalled(t, "Save", ctx, mock.MatchedBy(func(role role.Role) bool {
		return role.Name == roleName && reflect.DeepEqual(role.Permissions, permissions)
	}))
}
