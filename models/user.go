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

// type TokenUser struct {
// 	Users        SignupDetailResponse
// 	AccessToken  string
// 	RefreshToken string
// }

type UserLogin struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=6,max=20"`
}

//	type UserLoginResponse struct {
//		Id        uint   `json:"id"`
//		UserId    uint   `json:"user_id"`
//		Name      string `json:"firstname"`
//		Email     string `json:"email"`
//		Phone     string `json:"phone"`
//		Password  string `json:"password"`
//	}
type UserLoginResponse struct {
	ID       uint   `json:"id" gorm:"column:id"`
	Name     string `json:"name" gorm:"column:name"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"password" gorm:"column:password"`
	Phone    string `json:"phone" gorm:"column:phone"`
	Blocked  bool   `json:"blocked" gorm:"column:blocked"`
	IsAdmin  bool   `json:"isadmin" gorm:"column:isadmin"`
}
