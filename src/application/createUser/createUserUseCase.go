package createUser

import (
	"context"
	"fmt"
	"go-iam/src/domain/events"
	"go-iam/src/domain/hash"
	"go-iam/src/domain/internals"
	"go-iam/src/domain/user"

	"github.com/google/uuid"
)

type CreateUserUseCase struct {
	userRepository user.UserRepository
	passwordHasher hash.Hasher
	eventPublisher events.EventPublisher
	emailValidator user.EmailValidator
	logger         internals.Logger
}

func (useCase *CreateUserUseCase) Execute(ctx context.Context, request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*CreateUserRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(ctx, fmt.Sprintf("Starting user creation for username %s", validatedRequest.Username))
	defer useCase.logger.Info(ctx, fmt.Sprintf("Finished user creation for username %s", validatedRequest.Username))

	err := useCase.emailValidator.Validate(validatedRequest.Email)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	passwordHash, err := useCase.passwordHasher.Hash(validatedRequest.Password)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	user := user.User{
		ID:        uuid.New(),
		Username:  validatedRequest.Username,
		Email:     validatedRequest.Email,
		Password:  *passwordHash,
		Superuser: validatedRequest.Superuser,
	}
	if err = useCase.userRepository.Save(ctx, user); err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	if err = useCase.publishEvent(&user); err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	return internals.EmptyUseCaseResponse()
}

func (useCase *CreateUserUseCase) publishEvent(userCreated *user.User) error {
	event := user.UserCreatedEvent{
		ID:        userCreated.ID.String(),
		Username:  userCreated.Username,
		Email:     userCreated.Email,
		Superuser: userCreated.Superuser,
	}
	return useCase.eventPublisher.Publish(&event)
}

func NewCreateUserUseCase(userRepository user.UserRepository, passwordHasher hash.Hasher, eventPublisher events.EventPublisher, emailValidator user.EmailValidator, logger internals.Logger) *CreateUserUseCase {
	useCase := CreateUserUseCase{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
		eventPublisher: eventPublisher,
		emailValidator: emailValidator,
		logger:         logger,
	}
	return &useCase
}
