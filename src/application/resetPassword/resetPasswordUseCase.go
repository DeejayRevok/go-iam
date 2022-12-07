package resetPassword

import (
	"context"
	"fmt"
	"go-iam/src/domain/hash"
	"go-iam/src/domain/internals"
	"go-iam/src/domain/user"
	"time"
)

type ResetPasswordUseCase struct {
	userRepository              user.UserRepository
	userPasswordResetRepository user.UserPasswordResetRepository
	hashComparator              hash.HashComparator
	hasher                      hash.Hasher
	logger                      internals.Logger
}

func (useCase *ResetPasswordUseCase) Execute(ctx context.Context, request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*ResetPasswordRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(ctx, fmt.Sprintf("Starting password reset with token %s", validatedRequest.ResetToken))
	defer useCase.logger.Info(ctx, fmt.Sprintf("Finished password reset with token %s", validatedRequest.ResetToken))

	user, err := useCase.userRepository.FindByEmail(ctx, validatedRequest.UserEmail)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	passwordReset, err := useCase.userPasswordResetRepository.FindByUserID(ctx, user.ID)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	if passwordReset == nil {
		return internals.ErrorUseCaseResponse(fmt.Errorf("password reset not found for %s", validatedRequest.UserEmail))
	}

	err = useCase.validateResetToken(validatedRequest.ResetToken, passwordReset)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	err = useCase.storeNewPassword(ctx, validatedRequest.NewPassword, user)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	err = useCase.userPasswordResetRepository.Delete(ctx, *passwordReset)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	return internals.EmptyUseCaseResponse()
}

func (useCase *ResetPasswordUseCase) validateResetToken(resetToken string, passwordReset *user.UserPasswordReset) error {
	if passwordReset.Expiration.Before(time.Now()) {
		return fmt.Errorf("reset token %s is expired", resetToken)
	}
	err := useCase.hashComparator.Compare(resetToken, passwordReset.Token)
	if err != nil {
		return err
	}
	return nil
}

func (useCase *ResetPasswordUseCase) storeNewPassword(ctx context.Context, newPassword string, user *user.User) error {
	newPasswordHash, err := useCase.hasher.Hash(newPassword)
	if err != nil {
		return err
	}
	user.Password = *newPasswordHash
	return useCase.userRepository.Save(ctx, *user)
}

func NewResetPasswordUseCase(userRepository user.UserRepository, userPasswordResetRepository user.UserPasswordResetRepository, hashComparator hash.HashComparator, hasher hash.Hasher, logger internals.Logger) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		userRepository:              userRepository,
		userPasswordResetRepository: userPasswordResetRepository,
		hashComparator:              hashComparator,
		hasher:                      hasher,
		logger:                      logger,
	}
}
