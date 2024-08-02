package model

import (
	"gorm.io/gorm"
	"time"
)

// 关系信息表
type RelationshipList struct {
	Id               int64          `json:"id"`
	UserId           int64          `json:"user_id"`           //用户id 拥有者
	TargetId         int64          `json:"target_id"`         //用户id 对方
	Remark           string         `json:"remark"`            //对方的别名备注
	RelationshipType int            `json:"relationship_type"` //关系类型  1好友 2关注
	Status           int            `json:"status"`            //状态 1正常 2拉黑 3删除
	Extra            string         `json:"extra"`             //其他信息
	CreatedAt        time.Time      `json:"-"`
	UpdatedAt        time.Time      `json:"-"`
	DeletedAt        gorm.DeletedAt `json:"-"`
}

func (t *RelationshipList) TableName() string {
	return "relationship_list"
}

// 好友申请记录表
type ApplyFriendshipList struct {
	Id          int64          `json:"id"`
	UserId      int64          `json:"user_id"`     //用户id 拥有者
	TargetId    int64          `json:"target_id"`   //用户id 对方
	Remark      string         `json:"remark"`      //对方的别名备注
	Description string         `json:"description"` //申请描述
	Status      int            `json:"status"`      //状态 1申请中 2待处理 3通过 4被拒绝 5过期
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func (t *ApplyFriendshipList) TableName() string {
	return "apply_friendship_list"
}
