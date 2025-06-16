package v1

type SendMsgReq struct {
	ClientId       string `json:"client_id"`                                           //客户端的消息唯一识别
	ConversationId string `json:"conversation_id" binding:"required" example:"123456"` //会话ID
	UserId         string `json:"user_id"`                                             //发送者ID
	TargetId       string `json:"target_id" binding:"required" example:"123456"`       //接收者ID
	Content        string `json:"content" binding:"required"`                          //消息文本
	ContentType    int    `json:"content_type" binding:"required"`                     //内容类型
	SendTime       int64  `json:"send_time"`                                           //发送时间
}

type SendMsgResp struct {
	ClientId       string `json:"client_id"`       //客户端的消息唯一识别
	UserId         string `json:"user_id"`         //发送者ID
	MsgId          string `json:"msg_id"`          //消息ID
	ConversationId string `json:"conversation_id"` //会话ID
	Content        string `json:"content"`         //消息文本
	ContentType    int    `json:"content_type"`    //内容类型  1文本  2图片 3音频文件  4音频文件  5实时语音  6实时视频
	Status         int    `json:"status"`          //消息状态枚举，0可见 1屏蔽 2撤回
	Seq            int64  `json:"seq"`             //消息在会话中的序列号，用于保证消息的顺序
	UserSeq        int64  `json:"user_seq"`        //用户消息链序列号
	SendTime       int64  `json:"send_time"`       //发送时间
}

type ConversationResp struct {
	ConversationId string `json:"conversation_id"` //会话ID
	Type           int    `json:"type"`            //会话类型枚举，0单聊 1群聊
	Member         int    `json:"member"`          //与会话相关的用户数量
	Avatar         string `json:"avatar"`          //群组头像
	Announcement   string `json:"announcement"`    //群公告
	CreatedAt      int64  `json:"created_at"`

	UserId      string `json:"user_id"`       //用户ID
	LastReadSeq int64  `json:"last_read_seq"` //此会话用户已读的最后一条消息
	NotifyType  int    `json:"notify_type"`   //会话收到消息的提醒类型，0未屏蔽，正常提醒 1屏蔽 2强提醒
	IsTop       int    `json:"is_top"`        //会话是否被置顶展示 0否 1是
	Version     int    `json:"version"`       //会话版本号
}

type UserConversationResp struct {
	ConversationId string `json:"conversation_id"` //会话ID
	Type           int    `json:"type"`            //会话类型枚举，0单聊 1群聊
	Member         int    `json:"member"`          //与会话相关的用户数量
	Avatar         string `json:"avatar"`          //群组头像
	Announcement   string `json:"announcement"`    //群公告
	CreatedAt      int64  `json:"created_at"`

	UserId      string `json:"user_id"`       //用户ID
	LastReadSeq int64  `json:"last_read_seq"` //此会话用户已读的最后一条消息
	NotifyType  int    `json:"notify_type"`   //会话收到消息的提醒类型，0未屏蔽，正常提醒 1屏蔽 2强提醒
	IsTop       int    `json:"is_top"`        //会话是否被置顶展示 0否 1是
}

type HistoryMsgListReq struct {
	UserId         string `json:"user_id"`         //用户ID
	ConversationId string `json:"conversation_id"` //会话ID
	Seq            int64  `json:"seq"`             //消息序列号
	Limit          int    `json:"limit"`
}

type ReportReadReq struct {
	UserId         string `json:"user_id"`                                             //用户ID
	ConversationId string `json:"conversation_id" binding:"required" example:"123456"` //会话ID
	Seq            int64  `json:"seq"`                                                 //消息序列号
}

type UserConversationListReq struct {
	UserId   string `json:"user_id"` //用户ID
	Version  int64  `json:"version"`
	PageNum  int64  `json:"page_num"`
	PageSize int64  `json:"page_size"`
}

type ConversationUsersReq struct {
	UserId         string `json:"user_id"`                                             //用户ID
	ConversationId string `json:"conversation_id" binding:"required" example:"123456"` //会话ID
}
