package v1

type SendMsgReq struct {
	ConversationId int64  `json:"conversation_id"` //会话ID
	UserId         int64  `json:"user_id"`         //发送者ID
	TargetId       int64  `json:"target_id"`       //接收者ID
	Content        string `json:"content"`         //消息文本
	ContentType    int    `json:"content_type"`    //内容类型
	SendTime       int64  `json:"send_time"`       //发送时间
}

type SendMsgResp struct {
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

type ConversationResp struct {
	ConversationId int64  `json:"conversation_id"` //会话ID
	Type           int    `json:"type"`            //会话类型枚举，0单聊 1群聊
	Avatar         string `json:"avatar"`          //会话头像
	LastReadSeq    int64  `json:"last_read_seq"`   //此会话用户已读的最后一条消息
	NotifyType     int    `json:"notify_type"`     //会话收到消息的提醒类型，0未屏蔽，正常提醒 1屏蔽 2强提醒
	IsTop          int    `json:"is_top"`          //会话是否被置顶展示

	RecentMsg SendMsgResp              `json:"recent_msg"` //此会话最新产生的消息
	UserList  []GetProfileResponseData `json:"user_list"`  //此会话的用户列表
}

type HistoryMsgListReq struct {
	UserId         int64 `json:"user_id"`         //用户ID
	ConversationId int64 `json:"conversation_id"` //会话ID
	Seq            int64 `json:"seq"`             //消息序列号
	PageNum        int   `json:"page_num"`
	PageSize       int   `json:"page_size"`
}

type ReportReadReq struct {
	UserId         int64 `json:"user_id"`         //用户ID
	ConversationId int64 `json:"conversation_id"` //会话ID
	Seq            int64 `json:"seq"`             //消息序列号
}
