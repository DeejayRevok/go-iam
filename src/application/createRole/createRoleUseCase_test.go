package createRole

import (
	"errors"
	"go-uaa/mocks"
	"go-uaa/src/domain/permission"
	"go-uaa/src/domain/role"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type testCase struct {
	PermissionRepo *mocks.PermissionRepository
	RoleRepo       *mocks.RoleRepository
	UseCase        *CreateRoleUseCase
}

func setUp(t *testing.T) testCase {
	logger, _ := zap.NewDevelopment()
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

	response := testCase.UseCase.Execute(request)

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
	testCase.PermissionRepo.On("FindByNames", mock.Anything).Return(nil, findError)
	permissionsNames := []string{permissionName}
	roleName := "Test role"
	request := CreateRoleRequest{
		Name:        roleName,
		Permissions: permissionsNames,
	}

	response := testCase.UseCase.Execute(&request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != findError {
		t.Fatal("Error expected to be the same as the permissions repository returned error")
	}
	testCase.PermissionRepo.AssertCalled(t, "FindByNames", permissionsNames)
	testCase.RoleRepo.AssertNotCalled(t, "Save")
}

func TestExecuteRoleSaveError(t *testing.T) {
	testCase := setUp(t)
	saveError := errors.New("Test save error")
	testCase.RoleRepo.On("Save", mock.Anything).Return(saveError)
	permissionName := "Test permission"
	permissions := []permission.Permission{{Name: permissionName}}
	testCase.PermissionRepo.On("FindByNames", mock.Anything).Return(permissions, nil)
	permissionsNames := []string{permissionName}
	roleName := "Test role"
	request := CreateRoleRequest{
		Name:        roleName,
		Permissions: permissionsNames,
	}

	response := testCase.UseCase.Execute(&request)

	if response.Err == nil {
		t.Fatal("Expected use case to return error")
	}
	if response.Err != saveError {
		t.Fatal("Error expected to be the same as the role repository returned error")
	}
	testCase.PermissionRepo.AssertCalled(t, "FindByNames", permissionsNames)
	testCase.RoleRepo.AssertCalled(t, "Save", mock.MatchedBy(func(role role.Role) bool {
		return role.Name == roleName && reflect.DeepEqual(role.Permissions, permissions)
	}))
}

func TestExecuteRoleSaveSuccess(t *testing.T) {
	testCase := setUp(t)
	testCase.RoleRepo.On("Save", mock.Anything).Return(nil)
	permissionName := "Test permission"
	permissions := []permission.Permission{{Name: permissionName}}
	testCase.PermissionRepo.On("FindByNames", mock.Anything).Return(permissions, nil)
	permissionsNames := []string{permissionName}
	roleName := "Test role"
	request := CreateRoleRequest{
		Name:        roleName,
		Permissions: permissionsNames,
	}

	response := testCase.UseCase.Execute(&request)

	if response.Err != nil {
		t.Fatal("Expected use case not to return error")
	}
	testCase.PermissionRepo.AssertCalled(t, "FindByNames", permissionsNames)
	testCase.RoleRepo.AssertCalled(t, "Save", mock.MatchedBy(func(role role.Role) bool {
		return role.Name == roleName && reflect.DeepEqual(role.Permissions, permissions)
	}))
}
