package model

type ValidateUser struct {
	Email    string
	Password string
}

type JwtToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int32  `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int32  `json:"refresh_expires_in"`
}

type AuthTokenRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
