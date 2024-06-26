package v1

type MsgReq struct {
	ConversationId int64  `json:"conversation_id"` //会话ID
	UserId         int64  `json:"user_id"`         //发送者ID
	TargetId       int64  `json:"target_id"`       //接收者ID
	Content        string `json:"content"`         //消息文本
	ContentType    int    `json:"content_type"`    //内容类型
	SendTime       int64  `json:"send_time"`       //发送时间
	CreatedAt      int64  `json:"created_at"`
}

type MsgResp struct {
	UserId         int64  `json:"user_id"`         //发送者ID
	MsgId          int64  `json:"msg_id"`          //消息ID
	ConversationId int64  `json:"conversation_id"` //会话ID
	Content        string `json:"content"`         //消息文本
	ContentType    int    `json:"content_type"`    //内容类型  1文本  2图片 3音频文件  4音频文件  5实时语音  6实时视频
	Status         int    `json:"status"`          //消息状态枚举，0可见 1屏蔽 2撤回
	Seq            int64  `json:"seq"`
	SendTime       int64  `json:"send_time"` //发送时间
	CreatedAt      int64  `json:"created_at"`
}
