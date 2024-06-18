package service

import (
	"context"
	"fmt"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/ljinf/im_server_standalone/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, error)
	// 用户信息
	GetProfile(ctx context.Context, userId int64) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx context.Context, userId int64, req *v1.UpdateProfileRequest) error
	// 更新注册表
	UpdateRegisterInfo(ctx context.Context, userId int64, req *v1.UpdateRegisterInfoRequest) error
}

type userService struct {
	userRepo repository.UserRepository
	*Service
}

func NewUserService(service *Service, userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
		Service:  service,
	}
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// check username
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("GetByEmail", req))
		return v1.ErrInternalServerError
	}

	if user != nil {
		return v1.ErrEmailAlreadyUse
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(fmt.Sprintf("GenerateFromPassword %v", err))
		return v1.ErrGenerateFromPassword
	}
	// Generate user ID
	userId, err := s.sid.GenUint64()
	if err != nil {
		s.logger.Error(fmt.Sprintf("GenerateUserID %v", err))
		return v1.ErrGenerateUserID
	}

	account := &model.AccountInfo{
		UserId:   int64(userId),
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err = s.userRepo.CreateRegister(ctx, account); err != nil {
		s.logger.Error(err.Error(), zap.Any("accountInfo", account))
		return v1.ErrInternalServerError
	}
	return nil
}

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (string, error) {
	info, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || info == nil {
		return "", v1.ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(info.Password), []byte(req.Password))
	if err != nil {
		return "", v1.ErrPasswordFailed
	}

	token, err := s.jwt.GenToken(info.UserId, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) UpdateRegisterInfo(ctx context.Context, userId int64, req *v1.UpdateRegisterInfoRequest) error {
	info, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return v1.ErrUnauthorized
	}

	if info != nil {
		return v1.ErrEmailAlreadyUse
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Error(fmt.Sprintf("GenerateFromPassword %v", err))
			return v1.ErrGenerateFromPassword
		}
		req.Password = string(hashedPassword)
	}

	registerInfo := &model.Register{
		UserId:   userId,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}

	if err = s.userRepo.UpdateRegister(ctx, registerInfo); err != nil {
		s.logger.Error(err.Error(), zap.Any("registerInfo", registerInfo))
		return v1.ErrInternalServerError
	}
	return nil
}

func (s *userService) GetProfile(ctx context.Context, userId int64) (*v1.GetProfileResponseData, error) {
	user, err := s.userRepo.GetAccountInfoByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &v1.GetProfileResponseData{
		UserId:   user.UserId,
		NickName: user.NickName,
		Phone:    user.Phone,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userId int64, req *v1.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	user.Avatar = req.Avatar
	user.NickName = req.NickName

	if user.Gender == 3 {
		if req.Gender == 1 || req.Gender == 2 {
			user.Gender = req.Gender
		}
	}

	if err = s.userRepo.UpdateUserInfo(ctx, user); err != nil {
		return err
	}

	return nil
}
