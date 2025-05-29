package model

import "encoding/json"

// 用户发布的时刻
type CommunityMoment struct {
	Id             int64           `json:"id"`
	UserId         string          `json:"user_id"`         //发送者ID
	MomentId       string          `json:"moment_id"`       //时刻ID
	Content        string          `json:"content"`         //描述内容
	Attachment     json.RawMessage `json:"attachment"`      //图片/音频/视频的url集合
	AttachmentType int             `json:"attachment_type"` //类型 1-图片  2-音频  3-视频文件
	Public         int             `json:"public"`          //可见范围 1-公共  2-私密
	Status         int             `json:"status"`          //状态，1-审核中 2-正常 3-违规
	IsDel          int             `json:"is_del"`          //删除状态，1-正常 2-删除
	CreatedAt      int64           `json:"created_at"`
}

func (c *CommunityMoment) TableName() string {
	return "community_moment_list"
}

// 时刻历史点赞评论计数
type MomentCount struct {
	Id              int64  `json:"id"`
	MomentId        string `json:"moment_id"`         //时刻ID
	LikeCount       int    `json:"like_count"`        //点赞数
	LikeCancelCount int    `json:"like_cancel_count"` //点赞取消数
	CommentCount    int    `json:"comment_count"`     //评论数
}

func (c *MomentCount) TableName() string {
	return "moment_count_list"
}

// 时刻被点赞记录
type MomentLike struct {
	Id        int64  `json:"id"`
	MomentId  string `json:"moment_id"` //时刻ID
	UserId    string `json:"user_id"`
	Status    int    `json:"status"` //状态，1-正常 2-取消
	CreatedAt int64  `json:"created_at"`
}

func (c *MomentLike) TableName() string {
	return "moment_like_list"
}

// 时刻的评论记录
type MomentComment struct {
	Id             int64  `json:"id"`
	CommentId      string `json:"comment_id"`       //评论ID
	ParentId       string `json:"parent_id"`        //父评论ID  顶级评论
	MomentId       string `json:"moment_id"`        //时刻ID
	UserId         string `json:"user_id"`          //用户ID
	ReplyId        string `json:"reply_id"`         //回复 用户ID
	ReplyCommentId string `json:"reply_comment_id"` //回复评论ID
	Content        string `json:"content"`          //评论内容
	Status         int    `json:"status"`           //状态，1-审核中 2-正常 3-违规
	IsDel          int    `json:"is_del"`           //删除状态，1-正常 2-删除
	CreatedAt      int64  `json:"created_at"`
}

func (c *MomentComment) TableName() string {
	return "moment_comment_list"
}

type MomentCommentResp struct {
	Id        int64  `json:"id"`
	CommentId string `json:"comment_id"` //评论ID
	ParentId  string `json:"parent_id"`  //父评论ID
	MomentId  string `json:"moment_id"`  //时刻ID
	UserId    string `json:"user_id"`    //用户ID
	ReplyId   string `json:"reply_id"`   //回复 用户ID
	Content   string `json:"content"`    //评论内容
	Status    int    `json:"status"`     //状态，1-审核中 2-正常 3-违规 4-删除
	CreatedAt int64  `json:"created_at"`

	LikeCount       int `json:"like_count"`        //点赞数
	LikeCancelCount int `json:"like_cancel_count"` //点赞取消数
	CommentCount    int `json:"comment_count"`     //评论回复数
}

// 评论被点赞的记录
type MomentCommentLike struct {
	Id        int64  `json:"id"`
	CommentId string `json:"comment_id"` //评论ID
	UserId    string `json:"user_id"`
	Status    int    `json:"status"` //状态，1-正常 2-取消
	CreatedAt int64  `json:"created_at"`
}

func (c *MomentCommentLike) TableName() string {
	return "moment_comment_like_list"
}

// 评论被点赞回复的计数
type MomentCommentCount struct {
	Id              int64  `json:"id"`
	CommentId       string `json:"comment_id"`        //评论ID
	LikeCount       int    `json:"like_count"`        //点赞数
	LikeCancelCount int    `json:"like_cancel_count"` //点赞取消数
	CommentCount    int    `json:"comment_count"`     //评论数
}

func (c *MomentCommentCount) TableName() string {
	return "moment_comment_count"
}

type CommunityMomentResp struct {
	Id             int64           `json:"id"`
	UserId         string          `json:"user_id"`         //发送者ID
	MomentId       string          `json:"moment_id"`       //时刻ID
	Content        string          `json:"content"`         //描述内容
	Attachment     json.RawMessage `json:"attachment"`      //图片/音频/视频的url集合
	AttachmentType int             `json:"attachment_type"` //类型 1-图片  2-音频  3-视频文件
	Public         int             `json:"public"`          //可见范围 1-公共  2-私密
	Status         int             `json:"status"`          //状态，1-审核中 2-正常 3-违规 4-删除
	CreatedAt      int64           `json:"created_at"`

	LikeCount       int `json:"like_count"`        //点赞数
	LikeCancelCount int `json:"like_cancel_count"` //点赞取消数
	LikeStatus      int `json:"like_status"`       //点赞状态
	CommentCount    int `json:"comment_count"`     //评论数
}
