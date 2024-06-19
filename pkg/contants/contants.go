package contants

const (
	//申请状态
	ApplyFriendshipStatusApplying = 1 //申请中
	ApplyFriendshipStatusPending  = 2 //待处理
	ApplyFriendshipStatusApproved = 3 //通过
	ApplyFriendshipStatusRejected = 4 //拒绝
	ApplyFriendshipStatusExpired  = 5 //过期

	//关系类型
	RelationshipTypeFriend = 1 //好友关系
	RelationshipTypeFollow = 2 //关注

	//关系状态
	RelationshipStatusNormal = 1 //正常
	RelationshipStatusBlock  = 2 //拉黑
	RelationshipStatusDel    = 3 //删除

	MsgTypeNotify  = 1 //通知消息
	MsgTypeCommand = 2 //指令消息
	MsgTypeChat    = 3 //普通聊天消息
)
