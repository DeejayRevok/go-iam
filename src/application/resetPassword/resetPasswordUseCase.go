package resetPassword

import (
	"fmt"
	"go-uaa/src/domain/hash"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/user"
	"time"

	"go.uber.org/zap"
)

type ResetPasswordUseCase struct {
	userRepository              user.UserRepository
	userPasswordResetRepository user.UserPasswordResetRepository
	hashComparator              hash.HashComparator
	hasher                      hash.Hasher
	logger                      *zap.Logger
}

func (useCase *ResetPasswordUseCase) Execute(request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*ResetPasswordRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(fmt.Sprintf("Starting password reset with token %s", validatedRequest.ResetToken))
	defer useCase.logger.Info(fmt.Sprintf("Finished password reset with token %s", validatedRequest.ResetToken))

	user, err := useCase.userRepository.FindByEmail(validatedRequest.UserEmail)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	passwordReset, err := useCase.userPasswordResetRepository.FindByUserID(user.ID)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	err = useCase.validateResetToken(validatedRequest.ResetToken, passwordReset)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	err = useCase.storeNewPassword(validatedRequest.NewPassword, user)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	err = useCase.userPasswordResetRepository.Delete(*passwordReset)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	return internals.EmptyUseCaseResponse()
}

func (useCase *ResetPasswordUseCase) validateResetToken(resetToken string, passwordReset *user.UserPasswordReset) error {
	if passwordReset.Expiration.Before(time.Now()) {
		return fmt.Errorf("Reset token %s is expired", resetToken)
	}
	err := useCase.hashComparator.Compare(resetToken, passwordReset.Token)
	if err != nil {
		return err
	}
	return nil
}

func (useCase *ResetPasswordUseCase) storeNewPassword(newPassword string, user *user.User) error {
	newPasswordHash, err := useCase.hasher.Hash(newPassword)
	if err != nil {
		return err
	}
	user.Password = *newPasswordHash
	return useCase.userRepository.Save(*user)
}

func (*ResetPasswordUseCase) RequiredPermissions() []string {
	return make([]string, 0)
}

func NewResetPasswordUseCase(userRepository user.UserRepository, userPasswordResetRepository user.UserPasswordResetRepository, hashComparator hash.HashComparator, hasher hash.Hasher, logger *zap.Logger) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		userRepository:              userRepository,
		userPasswordResetRepository: userPasswordResetRepository,
		hashComparator:              hashComparator,
		hasher:                      hasher,
		logger:                      logger,
	}
}
