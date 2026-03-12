package models


type VerifyOTP struct {
    Email    string `json:"email"`
    OTP      string `json:"otp"`
    Name     string `json:"name"`     // For signup
    Password string `json:"password"` 
    Phone    string `json:"phone"`    
}

type ResendOTP struct {
	Email string `json:"email"`
}
type TokenUser struct {
	User         SignupDetailResponse `json:"user"`
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
}

type ForgotPassword struct{
	Email string `json:"email"`
}

type ResetPassword struct {
	Email       string `json:"email" binding:"required,email"`
	OTP         string `json:"otp" binding:"required"`
	NewPassword string `json:"newpassword" binding:"required,min=6"`
}