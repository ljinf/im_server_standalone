package v1

type CodeRequest struct {
	Account     string `json:"account" binding:"required" example:"123456"`      //账号
	AccountType int    `json:"account_type" binding:"required" example:"123456"` //账号类型 1-email 2-手机号 3-微信
}

type RegisterRequest struct {
	Account     string `json:"account" binding:"required" example:"123456"`      //账号
	AccountType int    `json:"account_type" binding:"required" example:"123456"` //账号类型 1-email 2-手机号 3-微信
	Password    string `json:"password"`
	Code        string `json:"code" binding:"required" example:"123456"` //验证码
}

type LoginRequest struct {
	Account     string `json:"account" binding:"required" example:"123456"`      //账号
	AccountType int    `json:"account_type" binding:"required" example:"123456"` //账号类型 1-email 2-手机号 3-微信
	Password    string `json:"password"`
	Code        string `json:"code" binding:"required" example:"123456"` //验证码
}
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}

type UpdateRegisterInfoRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	NickName      string `json:"nick_name"`      //昵称
	Avatar        string `json:"avatar"`         //头像
	Gender        int    `json:"gender"`         //性别
	Background    string `json:"background"`     //个人背景图
	SelfSignature string `json:"self_signature"` //个性签名
}
type GetProfileResponseData struct {
	UserId        string `json:"user_id"`
	Account       string `json:"account"`        //账号
	AccountType   int    `json:"account_type"`   //账号类型 1-email 2-手机号 3-微信
	NickName      string `json:"nick_name"`      //昵称
	Avatar        string `json:"avatar"`         //头像
	Background    string `json:"background"`     //个人背景图
	SelfSignature string `json:"self_signature"` //个性签名
	Gender        int    `json:"gender"`         //性别
	InitInfo      bool   `json:"init_info"`      //是否初始化基本信息
	Status        int    `json:"status"`         //用户状态    1:正常 2:封禁  3:注销
	AccessToken   string `json:"access_token"`
}
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}

type SearchUserRequest struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
}
