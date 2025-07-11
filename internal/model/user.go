package model

import (
	"gorm.io/gorm"
	"time"
)

// 注册表
type Register struct {
	Id          int64          `json:"id" gorm:"primarykey"`
	UserId      string         `json:"user_id"`
	Account     string         `json:"account"` //账号
	Password    string         `json:"password"`
	AccountType int            `json:"account_type"` //账号类型 1-email 2-手机号 3-微信
	Channel     string         `json:"channel"`      //渠道
	AppVersion  string         `json:"app_version"`  //app版本
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (r *Register) TableName() string {
	return "register"
}

// 用户信息表
type UserInfo struct {
	Id            int64          `json:"id" gorm:"primarykey"`
	UserId        string         `json:"user_id"`
	NickName      string         `json:"nick_name"`      //昵称
	Avatar        string         `json:"avatar"`         //头像
	Background    string         `json:"background"`     //个人背景图
	SelfSignature string         `json:"self_signature"` //个性签名
	Gender        int            `json:"gender"`         //性别
	InitInfo      bool           `json:"init_info"`      //是否初始化基本信息
	Status        int            `json:"status"`         //用户状态    1:正常 2:封禁  3:注销
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *UserInfo) TableName() string {
	return "user_info"
}

type AccountInfo struct {
	UserId      string `json:"user_id"`
	Account     string `json:"account"`      //账号
	AccountType int    `json:"account_type"` //账号类型 1-email 2-手机号 3-微信
	Password    string `json:"-"`
	Salt        string `json:"-"`

	NickName      string `json:"nick_name"`      //昵称
	Avatar        string `json:"avatar"`         //头像
	Background    string `json:"background"`     //个人背景图
	SelfSignature string `json:"self_signature"` //个性签名
	Gender        int    `json:"gender"`         //性别
	InitInfo      bool   `json:"init_info"`      //是否初始化基本信息
	Status        int    `json:"status"`         //用户状态    1:正常 2:封禁  3:注销
}
