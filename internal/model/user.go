package model

import (
	"gorm.io/gorm"
	"time"
)

// 注册表
type Register struct {
	Id        int64          `json:"id" gorm:"primarykey"`
	UserId    int64          `json:"user_id"`
	Phone     string         `json:"phone"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Salt      string         `json:"salt"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (r *Register) TableName() string {
	return "register"
}

// 用户信息表
type UserInfo struct {
	Id        int64          `json:"id" gorm:"primarykey"`
	UserId    int64          `json:"user_id"`
	NickName  string         `json:"nick_name"` //昵称
	Avatar    string         `json:"avatar"`    //头像
	Gender    int            `json:"gender"`    //性别
	Status    int            `json:"status"`    //用户状态  0:异常  1:正常
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *UserInfo) TableName() string {
	return "user_info"
}

type AccountInfo struct {
	UserId   int64  `json:"user_id"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	NickName string `json:"nick_name"` //昵称
	Avatar   string `json:"avatar"`    //头像
	Gender   int    `json:"gender"`    //性别
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Status   int    `json:"status"` //用户状态  0:异常  1:正常
}
