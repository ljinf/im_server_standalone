package v1

import "encoding/json"

type AddMomentReq struct {
	MomentId       string          `json:"moment_id"`       //时刻ID
	UserId         string          `json:"user_id"`         //发送者ID
	Content        string          `json:"content"`         //描述内容
	Attachment     json.RawMessage `json:"attachment"`      //图片/音频/视频的url集合
	AttachmentType int             `json:"attachment_type"` //类型 1-图片  2-音频  3-视频文件
	Public         int             `json:"public"`          //可见范围 1-公共  2-私密
	Status         int             `json:"status"`          //状态，1-审核中 2-正常 3-违规 4-删除
}

type GetMomentListReq struct {
	UserId string `json:"user_id"` //当前登录的用户

	WhereUser string `json:"where_user"` //查询条件的用户ID，即是某人的时刻
	Public    int    `json:"public"`     //可见范围 1-公共  2-私密
	Status    int    `json:"status"`     //状态，1-审核中 2-正常 3-违规 4-删除
	CreatedAt int64  `json:"created_at"`
	PageNum   int    `json:"page_num"`
	PageSize  int    `json:"page_size"`
}

type MomentCommentListReq struct {
	UserId    string `json:"user_id"`   //当前登录的用户
	ParentId  string `json:"parent_id"` //父评论ID
	MomentId  string `json:"moment_id"` //时刻ID
	CreatedAt int64  `json:"created_at"`
	PageNum   int    `json:"page_num"`
	PageSize  int    `json:"page_size"`
}

type MomentCommentListResp struct {
	Id              int64  `json:"id"`
	CommentId       string `json:"comment_id"`        //评论ID
	ParentId        string `json:"parent_id"`         //父评论ID
	MomentId        string `json:"moment_id"`         //时刻ID
	UserId          string `json:"user_id"`           //用户ID
	ReplyId         string `json:"reply_id"`          //回复 用户ID
	Content         string `json:"content"`           //评论内容
	LikeCount       int    `json:"like_count"`        //点赞数
	LikeCancelCount int    `json:"like_cancel_count"` //点赞取消数
	CommentCount    int    `json:"comment_count"`     //评论回复数
	LikeStatus      int    `json:"like_status"`       //当前登录用户的点赞状态
	Status          int    `json:"status"`            //状态，1-审核中 2-正常 3-违规 4-删除
	CreatedAt       int64  `json:"created_at"`
}

// 点赞时刻
type AddMomentLikeReq struct {
	MomentId string `json:"moment_id"` //时刻ID
	UserId   string `json:"user_id"`
	Status   int    `json:"status"` //状态，1-正常 2-取消
}

type AddMomentCommentReq struct {
	UserId         string `json:"user_id"`          //用户ID
	ParentId       string `json:"parent_id"`        //父评论ID 顶级评论ID
	MomentId       string `json:"moment_id"`        //时刻ID
	ReplyId        string `json:"reply_id"`         //回复 用户ID
	ReplyCommentId string `json:"reply_comment_id"` //回复评论ID
	Content        string `json:"content"`          //评论内容
}

// 点赞评论
type AddMomentCommentLikeReq struct {
	CommentId string `json:"comment_id"` //评论ID
	UserId    string `json:"user_id"`
	Status    int    `json:"status"` //状态，1-正常 2-取消
}
