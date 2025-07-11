package contants

const (
	AccountTypeEmail  = 1
	AccountTypePhone  = 2
	AccountTypeWechat = 3
)

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

	ConversationTypeC2C   = 0 //单聊
	ConversationTypeGroup = 1 //群聊

	MsgTypeNotify  = 1 //通知消息
	MsgTypeCommand = 2 //指令消息
	MsgTypeChat    = 3 //普通聊天消息

	ChatSayHello = "从此我们是好友关系啦！"

	MsgContentTypeTxt   = 1 //文字
	MsgContentTypeImg   = 2 //语音
	MsgContentTypeVideo = 3 //视频

	DefaultMsgListSize = 100
)

const (
	DirectLessThan    = "lt"
	DirectGreaterThan = "gt"
)

// social
const (

	//时刻状态 1-审核中 2-正常 3-违规 4-删除
	MomentStatusPending   = 1 //审核中
	MomentStatusNormal    = 2 //正常
	MomentStatusViolation = 3 //违规
	MomentStatusDelete    = 4 //删除

	//点赞
	MomentLikeNormal = 1 //正常
	MomentLikeCancel = 2 //取消

	//评论  //状态，1-审核中 2-正常 3-违规 4-删除
	MomentCommentStatusPending   = 1 //审核中
	MomentCommentStatusNormal    = 2 //正常
	MomentCommentStatusViolation = 3 //违规
	MomentCommentStatusDelete    = 4 //删除
)
