package user

type UpdateUserRequest struct {
	Email      *string `json:"email"`
	IsInactive *bool   `json:"is_inactive"`
}

type GetAllUsersResponse struct {
	Users []UserResponse `json:"users"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Total int            `json:"total"`
}

type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

type RefreshTokenRequest struct {
	Refresh_token string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email"`
}

type VerifyResetTokenRequest struct {
	Token       string `json:"token"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type VerifyEmailOTPRequest struct {
	Otp string `json:"otp"`
}
