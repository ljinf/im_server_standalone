package repository

import (
	"context"
	"errors"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/internal/cache"
	"github.com/ljinf/im_server_standalone/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type UserRepository interface {
	//创建注册信息
	CreateRegister(ctx context.Context, req *model.AccountInfo) error
	UpdateRegister(ctx context.Context, req *model.Register) error
	GetByEmail(ctx context.Context, email string) (*model.Register, error)
	GetByPhone(ctx context.Context, phone string) (*model.Register, error)

	GetByID(ctx context.Context, id int64) (*model.UserInfo, error)
	UpdateUserInfo(ctx context.Context, req *model.UserInfo) error

	GetAccountInfoByID(ctx context.Context, userId int64) (*model.AccountInfo, error)
}

func NewUserRepository(r *Repository) UserRepository {
	return &userRepository{
		Repository: r,
	}
}

type userRepository struct {
	*Repository
}

func (r *userRepository) CreateRegister(ctx context.Context, req *model.AccountInfo) error {
	return r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		//注册信息
		registerInfo := model.Register{
			UserId:    req.UserId,
			Phone:     req.Phone,
			Email:     req.Email,
			Password:  req.Password,
			Salt:      req.Salt,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := tx.Create(&registerInfo).Error; err != nil {
			return err
		}

		//创建新用户信息
		userInfo := model.UserInfo{
			UserId:    req.UserId,
			Status:    1,
			CreatedAt: now,
			UpdatedAt: now,
		}
		return tx.Create(&userInfo).Error
	})
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.Register, error) {
	var user model.Register
	if err := r.DB(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*model.Register, error) {
	var user model.Register
	if err := r.DB(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateRegister(ctx context.Context, req *model.Register) error {
	return r.DB(ctx).Where("user_id=?", req.UserId).Updates(req).Error
}

func (r *userRepository) UpdateUserInfo(ctx context.Context, req *model.UserInfo) error {
	return r.DB(ctx).Where("user_id=?", req.UserId).Updates(req).Error
}

func (r *userRepository) GetByID(ctx context.Context, userId int64) (*model.UserInfo, error) {
	var user model.UserInfo
	if err := r.DB(ctx).Where("user_id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAccountInfoByID(ctx context.Context, userId int64) (*model.AccountInfo, error) {

	infoCache, err := cache.GetAccountInfoCache(r.rdb, userId)
	if err != nil {
		r.logger.Error(err.Error(), zap.Any("userId", userId))
	}

	if infoCache != nil {
		return infoCache, nil
	}

	var info model.AccountInfo
	querySql := "SELECT u.`user_id`,u.`nick_name`,u.`avatar`,u.`gender`,r.`email`,r.`phone` " +
		"FROM `user_info` u INNER JOIN `register` r ON u.`user_id`=r.`user_id` WHERE u.`user_id`=?"
	if err := r.DB(ctx).Raw(querySql, userId).Scan(&info).Error; err != nil {
		return nil, err
	}

	if err = cache.SetAccountInfoCache(r.rdb, &info); err != nil {
		r.logger.Error(err.Error(), zap.Any("info", info))
	}

	return &info, nil
}
