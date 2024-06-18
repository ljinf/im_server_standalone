package v1

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// more biz errors
	ErrEmailAlreadyUse      = newError(1001, "The email is already in use.")
	ErrGenerateFromPassword = newError(1002, "密码加密异常")
	ErrGenerateUserID       = newError(1003, "创建用户ID失败")
	ErrPasswordFailed       = newError(1004, "账号密码错误")

	// 申请关系
	ErrAddApplyFriendshipFailed = newError(2001, "申请失败")
	ErrCreateRelationshipFailed = newError(2002, "添加好友失败")
)
