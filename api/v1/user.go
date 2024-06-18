package v1

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
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
	NickName string `json:"nick_name"` //昵称
	Avatar   string `json:"avatar"`    //头像
	Gender   int    `json:"gender"`    //性别
}
type GetProfileResponseData struct {
	UserId   int64  `json:"user_id"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	NickName string `json:"nick_name"` //昵称
	Avatar   string `json:"avatar"`    //头像
	Gender   int    `json:"gender"`    //性别
}
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}
