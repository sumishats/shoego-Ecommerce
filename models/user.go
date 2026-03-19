package models

type SignupDetail struct {
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type SignupDetailResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=6,max=20"`
}

type UserLoginResponse struct {
	ID       uint   `json:"id" gorm:"column:id"`
	Name     string `json:"name" gorm:"column:name"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"password" gorm:"column:password"`
	Phone    string `json:"phone" gorm:"column:phone"`
	Blocked  bool   `json:"blocked" gorm:"column:blocked"`
	IsAdmin  bool   `json:"isadmin" gorm:"column:isadmin"`
}
type UserProfileResponse struct {
	ID           uint              `json:"id"`
	Name         string            `json:"name"`
	Email        string            `json:"email"`
	Phone        string            `json:"phone"`
	ProfileImage string            `json:"profile_image"`
	Addresses    []AddressResponse `json:"addresses"`
}

type EditProfileRequest struct {
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	ProfileImage string `json:"profile_image"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type EmailChangeRequest struct {
	NewEmail string `json:"new_email"`
}

type VerifyEmailChangeRequest struct {
	NewEmail string `json:"new_email"`
	OTP      string `json:"otp"`
}

type AddAddressRequest struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	HouseName string `json:"house_name"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Pincode   string `json:"pincode"`
	IsDefault bool   `json:"is_default"`
}

type EditAddressRequest struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	HouseName string `json:"house_name"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Pincode   string `json:"pincode"`
	IsDefault bool   `json:"is_default"`
}

type AddressResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	HouseName string `json:"house_name"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Pincode   string `json:"pincode"`
	IsDefault bool   `json:"is_default"`
}
