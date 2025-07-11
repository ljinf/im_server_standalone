package service

import (
	"context"
	"fmt"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/cache"
	"github.com/ljinf/im_server_standalone/internal/model"
	"github.com/ljinf/im_server_standalone/internal/repository"
	"github.com/ljinf/im_server_standalone/pkg/contants"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService interface {
	SendCode(ctx context.Context, req *v1.CodeRequest) error
	Register(ctx context.Context, req *v1.RegisterRequest) (*v1.GetProfileResponseData, error)
	Login(ctx context.Context, req *v1.LoginRequest) (*v1.GetProfileResponseData, error)
	// 用户信息
	GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error)
	GetProfileByAccount(ctx context.Context, account string, accountType int) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) (*v1.GetProfileResponseData, error)
	// 更新注册表
	UpdateRegisterInfo(ctx context.Context, userId string, req *v1.UpdateRegisterInfoRequest) error
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

func (s *userService) SendCode(ctx context.Context, req *v1.CodeRequest) error {

	code := GenerateRandomNumberString(6)

	if err := cache.SetVerificationCodeCache(s.cache.Redis(ctx), req.Account, code); err != nil {
		s.logger.Error(err.Error(), zap.Any("SetVerificationCodeCacheErr", req.Account))
		return v1.ErrInternalServerError
	}

	//todo 发送验证码

	return nil
}

func (s *userService) checkCode(ctx context.Context, account, code string) bool {

	codeCache, err := cache.GetVerificationCodeCache(s.cache.Redis(ctx), account)
	if err != nil {
		s.logger.Error(err.Error(), zap.String("GetVerificationCodeCacheErr", ""))
		return false
	}

	if codeCache == code {
		return true
	}

	return false
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.GetProfileResponseData, error) {
	//验证code
	if s.checkCode(ctx, req.Account, req.Code) {
		return nil, v1.ErrVerificationCodeInvalid
	}

	// check username
	registerInfo, err := s.userRepo.GetByAccountType(ctx, req.Account, req.AccountType)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("GetByAccountType", req))
		return nil, v1.ErrInternalServerError
	}

	if registerInfo != nil {
		return nil, v1.ErrAccountAlreadyUse
	}

	return s.doRegister(ctx, req)
}

func (s *userService) doRegister(ctx context.Context, req *v1.RegisterRequest) (*v1.GetProfileResponseData, error) {

	var (
		password string
		err      error
	)

	if req.AccountType == contants.AccountTypeEmail {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Error(fmt.Sprintf("GenerateFromPassword %v", err))
			return nil, v1.ErrGenerateFromPassword
		}
		password = string(hashedPassword)
	}

	// Generate user ID
	userId, err := s.sid.GenUint64()
	if err != nil {
		s.logger.Error(fmt.Sprintf("GenerateUserID %v", err))
		return nil, v1.ErrGenerateUserID
	}

	account := &model.AccountInfo{
		UserId:      fmt.Sprintf("%v", userId),
		Account:     req.Account,
		AccountType: req.AccountType,
		Password:    password,
	}

	if err = s.userRepo.CreateRegister(ctx, account); err != nil {
		s.logger.Error(err.Error(), zap.Any("accountInfo", account))
		return nil, v1.ErrInternalServerError
	}

	token, err := s.jwt.GenToken(fmt.Sprintf("%v", userId), time.Now().Add(time.Hour*24*90))
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("GenToken", userId))
		return nil, v1.ErrInternalServerError
	}

	return &v1.GetProfileResponseData{
		UserId:      fmt.Sprintf("%v", userId),
		Account:     req.Account,
		AccountType: req.AccountType,
		InitInfo:    false,
		AccessToken: token,
	}, nil
}

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.GetProfileResponseData, error) {

	//验证code
	if !s.checkCode(ctx, req.Account, req.Code) {
		return nil, v1.ErrVerificationCodeInvalid
	}

	info, err := s.userRepo.GetByAccountType(ctx, req.Account, req.AccountType)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, v1.ErrInternalServerError
	}

	//账号不存在
	if info == nil {
		return s.doRegister(ctx, &v1.RegisterRequest{
			Account:     req.Account,
			AccountType: req.AccountType,
			Password:    req.Password,
		})
	}

	if info.AccountType == contants.AccountTypeEmail {
		err = bcrypt.CompareHashAndPassword([]byte(info.Password), []byte(req.Password))
		if err != nil {
			return nil, v1.ErrPasswordFailed
		}
	}

	userInfo, err := s.userRepo.GetByID(ctx, info.UserId)

	token, err := s.jwt.GenToken(info.UserId, time.Now().Add(time.Hour*24*90))
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("GenToken", userInfo.UserId))
		return nil, v1.ErrInternalServerError
	}

	return &v1.GetProfileResponseData{
		UserId:        userInfo.UserId,
		Account:       req.Account,
		AccountType:   req.AccountType,
		NickName:      userInfo.NickName,
		Avatar:        userInfo.Avatar,
		SelfSignature: userInfo.SelfSignature,
		Background:    userInfo.Background,
		Gender:        userInfo.Gender,
		InitInfo:      userInfo.InitInfo,
		Status:        userInfo.Status,
		AccessToken:   token,
	}, nil
}

func (s *userService) UpdateRegisterInfo(ctx context.Context, userId string, req *v1.UpdateRegisterInfoRequest) error {
	/*info, err := s.userRepo.GetByEmail(ctx, req.Email)
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
	}*/
	return nil
}

func (s *userService) GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error) {
	user, err := s.userRepo.GetAccountInfoByID(ctx, userId)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("GetAccountInfoByID", userId))
		return nil, v1.ErrInternalServerError
	}

	return &v1.GetProfileResponseData{
		UserId:        user.UserId,
		Account:       user.Account,
		AccountType:   user.AccountType,
		NickName:      user.NickName,
		Avatar:        user.Avatar,
		Background:    user.Background,
		SelfSignature: user.SelfSignature,
		Gender:        user.Gender,
		InitInfo:      user.InitInfo,
		Status:        user.Status,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) (*v1.GetProfileResponseData, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	user.Avatar = req.Avatar
	user.NickName = req.NickName
	user.SelfSignature = req.SelfSignature
	user.Background = req.Background

	if user.Gender == 3 {
		if req.Gender == 1 || req.Gender == 2 {
			user.Gender = req.Gender
		}
	}

	if err = s.userRepo.UpdateUserInfo(ctx, user); err != nil {
		return nil, err
	}

	return &v1.GetProfileResponseData{
		UserId:        user.UserId,
		NickName:      user.NickName,
		Avatar:        user.Avatar,
		Background:    user.Background,
		SelfSignature: user.SelfSignature,
		Gender:        user.Gender,
		Status:        user.Status,
	}, nil
}

func (s *userService) GetProfileByAccount(ctx context.Context, account string, accountType int) (*v1.GetProfileResponseData, error) {
	user, err := s.userRepo.GetAccountInfoByEmail(ctx, account, contants.AccountTypeEmail)
	if err != nil {
		return nil, err
	}

	return &v1.GetProfileResponseData{
		UserId:        user.UserId,
		Account:       user.Account,
		AccountType:   user.AccountType,
		NickName:      user.NickName,
		Avatar:        user.Avatar,
		Background:    user.Background,
		SelfSignature: user.SelfSignature,
		Gender:        user.Gender,
		InitInfo:      user.InitInfo,
		Status:        user.Status,
	}, nil
}
