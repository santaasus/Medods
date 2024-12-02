package auth

import "time"

// LoginUser is a struct that contains the request body for the login user
type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserData struct {
	UserName  string `json:"user_name" example:"UserName" gorm:"unique"`
	Email     string `json:"email" example:"some@mail.com" gorm:"unique"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Status    bool   `json:"status" example:"1"`
	Role      string `json:"role" example:"admin"`
	ID        int    `json:"id" example:"123"`
}

type SecurityData struct {
	Warning                   string    `json:"warning" example:"Seems your ip address was changed. If you didn't visit pornhub.com right now, open this instruction: lalala.com/safe_this_world's_eyes_from_you_bdsm_content"`
	JWTRefreshToken           string    `json:"jwt_refresh_token"`
	JWTAccessToken            string    `json:"jwt_access_token"`
	ExpirationAccessDateTime  time.Time `json:"expiration_access_date_time" example:"2023-02-02T21:03:53.196419-06:00"`
	ExpirationRefreshDateTime time.Time `json:"expiration_refresh_date_time" example:"2023-02-03T06:53:53.196419-06:00"`
}
