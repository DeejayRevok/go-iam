package createUser

import (
	"fmt"
	"go-uaa/src/domain/events"
	"go-uaa/src/domain/hash"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/role"
	"go-uaa/src/domain/user"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreateUserUseCase struct {
	userRepository user.UserRepository
	passwordHasher hash.Hasher
	roleRepository role.RoleRepository
	eventPublisher events.EventPublisher
	emailValidator user.EmailValidator
	logger         *zap.Logger
}

func (useCase *CreateUserUseCase) Execute(request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*CreateUserRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(fmt.Sprintf("Starting user creation for username %s", validatedRequest.Username))
	defer useCase.logger.Info(fmt.Sprintf("Finished user creation for username %s", validatedRequest.Username))

	err := useCase.emailValidator.Validate(validatedRequest.Email)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	passwordHash, err := useCase.passwordHasher.Hash(validatedRequest.Password)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	userRoles, err := useCase.findRoles(validatedRequest.Roles)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	user := user.User{
		ID:        uuid.New(),
		Username:  validatedRequest.Username,
		Email:     validatedRequest.Email,
		Password:  *passwordHash,
		Roles:     userRoles,
		Superuser: validatedRequest.Superuser,
	}
	if err = useCase.userRepository.Save(user); err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	useCase.publishEvent(&user)
	return internals.EmptyUseCaseResponse()
}

func (useCase *CreateUserUseCase) findRoles(roleIDs []string) ([]role.Role, error) {
	var roleUUIDs []uuid.UUID
	for _, roleID := range roleIDs {
		roleUUID, err := uuid.Parse(roleID)
		if err != nil {
			return nil, err
		}
		roleUUIDs = append(roleUUIDs, roleUUID)
	}
	roles, err := useCase.roleRepository.FindByIDs(roleUUIDs)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (useCase *CreateUserUseCase) publishEvent(userCreated *user.User) error {
	event := user.UserCreatedEvent{
		ID:       userCreated.ID.String(),
		Username: userCreated.Username,
		Email:    userCreated.Email,
	}
	return useCase.eventPublisher.Publish(event)
}

func (*CreateUserUseCase) RequiredPermissions() []string {
	return []string{}
}

func NewCreateUserUseCase(userRepository user.UserRepository, passwordHasher hash.Hasher, roleRepository role.RoleRepository, eventPublisher events.EventPublisher, emailValidator user.EmailValidator, logger *zap.Logger) *CreateUserUseCase {
	useCase := CreateUserUseCase{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
		roleRepository: roleRepository,
		eventPublisher: eventPublisher,
		emailValidator: emailValidator,
		logger:         logger,
	}
	return &useCase
}
