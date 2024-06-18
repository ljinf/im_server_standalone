package v1

type ApplyFriendshipRequest struct {
	UserId      int64  `json:"user_id"`                                  //用户id 拥有者
	TargetId    int64  `json:"target_id" binding:"required" example:"1"` //用户id 对方
	Remark      string `json:"remark"`                                   //对方的别名备注
	Description string `json:"description"`                              //申请描述
	Status      int    `json:"status"`                                   //状态 1申请中 2通过 3被拒绝
}

type RelationshipRequest struct {
	UserId           int64  `json:"user_id"`           //用户id 拥有者
	TargetId         int64  `json:"target_id"`         //用户id 对方
	Remark           string `json:"remark"`            //验证信息
	RelationshipType int    `json:"relationship_type"` //关系类型  1好友 2关注
	Status           int    `json:"status"`            //状态 1正常 2拉黑 3删除
	Extra            string `json:"extra"`             //其他信息
}
