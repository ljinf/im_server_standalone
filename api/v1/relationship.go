package v1

type ApplyFriendshipRequest struct {
	UserId      string `json:"user_id"`                                  //用户id 拥有者  被申请人
	TargetId    string `json:"target_id" binding:"required" example:"1"` //用户id 对方    申请人
	Remark      string `json:"remark"`                                   //对方的别名备注
	Description string `json:"description"`                              //申请描述
	Status      int    `json:"status"`                                   //状态 1申请中 2通过 3被拒绝
}

type RelationshipRequest struct {
	UserId           string `json:"user_id"`                                          //用户id 拥有者
	TargetId         string `json:"target_id" binding:"required" example:"1"`         //用户id 对方
	Remark           string `json:"remark"`                                           //验证信息
	RelationshipType int    `json:"relationship_type" binding:"required" example:"1"` //关系类型  1好友 2关注
	Status           int    `json:"status"`                                           //状态 1正常 2拉黑 3删除
	Extra            string `json:"extra"`                                            //其他信息
}

type RelationshipListReq struct {
	RelationshipType int `json:"relationship_type" binding:"required" example:"1"` //关系类型  1好友 2关注
	PageNum          int `json:"page_num"`
	PageSize         int `json:"page_size"`
}

type RelationshipRespData struct {
	UserId           string `json:"user_id"` //好友ID
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	NickName         string `json:"nick_name"`                                        //昵称
	Avatar           string `json:"avatar"`                                           //头像
	Gender           int    `json:"gender"`                                           //性别
	Remark           string `json:"remark"`                                           //别名
	RelationshipType int    `json:"relationship_type" binding:"required" example:"1"` //关系类型  1好友 2关注
	Status           int    `json:"status"`                                           //状态 1正常 2拉黑 3删除
	Extra            string `json:"extra"`                                            //其他信息
	CreatedAt        int64  `json:"created_at"`
}
