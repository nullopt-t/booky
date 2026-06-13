package user

import (
	"booky-backend/internal/middleware"
	"booky-backend/internal/model"
	"booky-backend/pkg/api"
	"booky-backend/pkg/api/security"
	"booky-backend/pkg/config"
	"booky-backend/pkg/log"
	"booky-backend/pkg/utils/jwt"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	ResendOTPLimit = 5
	VerifyOTPLimit = 3
)

type GetAllUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
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

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,max=999999"`
}

type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type UserURIParam struct {
	UserID uuid.UUID `json:"id" validate:"required,uuid"`
}

type UserResponse struct {
	ID             uuid.UUID  `json:"id,omitempty"`
	Role           string     `json:"role,omitempty"`
	Email          string     `json:"email,omitempty"`
	Phone          *string    `json:"phone,omitempty"`
	Status         string     `json:"status,omitempty"`
	SuspendedUntil *time.Time `json:"suspended_until,omitempty"`
	LockedUntil    *time.Time `json:"locked_unitl,omitzero"`
	CreatedAt      time.Time  `json:"created_at,omitzero"`
	UpdatedAt      time.Time  `json:"updated_at,omitzero"`
}

func ToUserResponse(user model.User) UserResponse {
	return UserResponse{
		ID:             user.ID,
		Role:           string(user.Role),
		Email:          user.Email,
		Phone:          user.Phone,
		Status:         string(user.Status),
		SuspendedUntil: user.SuspendedUntil,
		LockedUntil:    user.LockedUntil,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func ToUserListResponse(users []model.User) GetAllUsersResponse {
	var list GetAllUsersResponse
	list.Users = make([]UserResponse, 0, len(users))
	for _, user := range users {
		list.Users = append(list.Users, ToUserResponse(user))
	}
	return list
}

type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Handler struct {
	userService *UserService
	authService *AuthService
	limiter     *security.RateLimiter
	secrets     *config.Secrets
	logger      log.Logger
}

func NewHandler(
	userService *UserService,
	authService *AuthService,
	limiter *security.RateLimiter,
	secrets *config.Secrets,
	logger log.Logger,
) *Handler {
	return &Handler{
		userService: userService,
		authService: authService,
		limiter:     limiter,
		secrets:     secrets,
		logger:      logger,
	}
}

func (h *Handler) handleValidationError(c *gin.Context, err error) {
	if ve, ok := errors.AsType[validator.ValidationErrors](err); ok && ve != nil {
		fieldErrors := make([]api.FieldError, 0, len(ve))
		for _, e := range ve {
			fieldErrors = append(fieldErrors, api.FieldError{
				Field: e.Field(),
				Tags:  e.Tag(),
			})
		}
		c.Error(security.NewSecureError(
			http.StatusBadRequest,
			security.CodeValidation,
			"bad request data",
			err,
		).WithFields(fieldErrors))
		return
	}
	c.Error(
		security.NewSecureError(
			http.StatusBadRequest,
			security.CodeValidation,
			"bad request data",
			err,
		),
	)
}

func (h *Handler) UserRegister(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleValidationError(c, err)
		return
	}

	err := h.authService.Register(
		c.Request.Context(),
		req,
	)
	if err != nil {
		h.logger.Error(
			"user register",
			log.Meta{
				"Error": err,
			},
		)
	}

	c.JSON(
		http.StatusCreated,
		api.SuccessResponse{
			Message: "If the email is valid, you will receive a verification email.",
		},
	)
}

func (h *Handler) UserLogin(c *gin.Context) {
	var req LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleValidationError(c, err)
		return
	}

	tokens, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Data: map[string]string{
				"access_token":  tokens.AccessToken,
				"refresh_token": tokens.RefreshToken,
			},
		},
	)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	var uri UserURIParam
	if err := c.ShouldBindJSON(&uri); err != nil {
		h.handleValidationError(c, err)
		return
	}

	user, err := h.userService.GetUserByID(
		c.Request.Context(),
		uri.UserID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Data: UserResponse{
				ID:          user.ID,
				Email:       user.Email,
				LockedUntil: user.LockedUntil,
				CreatedAt:   user.CreatedAt,
				UpdatedAt:   user.UpdatedAt,
			},
		},
	)
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	var q api.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		h.handleValidationError(c, err)
		return
	}

	if q.PageSize == 0 {
		q.PageSize = 5
	}

	if q.Page == 0 {
		q.Page = 1
	}

	users, page, err := h.userService.GetAllUsers(
		c.Request.Context(),
		q,
	)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Data: ToUserListResponse(users),
			Meta: page,
		},
	)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	var uri UserURIParam
	if err := c.ShouldBindUri(&uri); err != nil {
		h.handleValidationError(c, err)
		return
	}

	err := h.userService.DeleteUserByID(
		c.Request.Context(),
		uri.UserID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Message: "User deleted successfully",
		},
	)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleValidationError(c, err)
		return
	}

	claims, err := jwt.VerifyToken(
		req.RefreshToken,
		h.secrets.JwtRefreshTokenSecretKey,
	)
	if err != nil {
		c.Error(
			security.NewSecureError(
				http.StatusUnauthorized,
				security.CodeAuth,
				"invalid or expired refresh token",
				err,
			),
		)
		return
	}

	accessToken, err := jwt.CreateToken(
		claims.Subject,
		h.secrets.JwtAccessTokenSecretKey,
		jwt.AccessTokenTTL,
		jwt.AccessTokenType,
	)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Data: RefreshTokenResponse{
				AccessToken: accessToken,
			},
		},
	)
}

// func (h *Handler) ForgetPassword(c *gin.Context) {
// 	var req ForgetPasswordRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		h.handleValidationError(c, err)
// 		return
// 	}

// 	/// check user exists by email
// 	user, err := h.userService.GetUserByEmail(
// 		c.Request.Context(),
// 		req.Email,
// 	)
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	var resetToken string
// 	if user.ResetTokenExpireAt != nil && user.ResetTokenExpireAt.After(time.Now()) {
// 		resetToken = *user.ResetToken
// 	} else {
// 		subjectStr, err := json.Marshal(
// 			token.UserSubject{
// 				UserID:   user.ID,
// 				UserRole: user.Role,
// 			},
// 		)
// 		if err != nil {
// 			c.Error(err)
// 			return
// 		}

// 		resetToken, err = jwt.CreateToken(
// 			string(subjectStr),
// 			h.secrets.JwtResetPassTokenSecretKey,
// 			jwt.ResetPassTokenTTL,
// 			jwt.ResetPassTokenType,
// 		)
// 		if err != nil {
// 			c.Error(err)
// 			return
// 		}

// 		err = h.userService.SetResetToken(
// 			c.Request.Context(),
// 			user.ID,
// 			resetToken,
// 		)
// 		if err != nil {
// 			c.JSON(
// 				http.StatusInternalServerError,
// 				api.ErrorResponse{
// 					Code:    "INTERNAL_ERROR",
// 					Message: err.Error(),
// 				},
// 			)
// 			return
// 		}
// 	}

// 	c.JSON(
// 		http.StatusOK,
// 		api.SuccessResponse{
// 			Message: "email sent successfully",
// 		},
// 	)
// }

// func (h *Handler) ResetPassword(c *gin.Context) {
// 	var req VerifyResetTokenRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		h.handleValidationError(c, err)
// 		return
// 	}

// 	claims, err := jwt.VerifyToken(
// 		req.Token,
// 		h.secrets.JwtResetPassTokenSecretKey,
// 	)
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	var subject token.UserSubject
// 	if err := json.Unmarshal(
// 		[]byte(claims.Subject),
// 		&subject,
// 	); err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	if err := h.userService.CheckPasswordResetToken(
// 		c.Request.Context(),
// 		subject.UserID,
// 		req.Token,
// 	); err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	if claims.ExpiresAt.Before(time.Now()) {
// 		c.Error(
// 			security.NewSecureError(
// 				http.StatusBadRequest,
// 				security.CodeValidation,
// 				"expired reset token",
// 				nil,
// 			),
// 		)
// 		return
// 	}

// 	if err := h.userService.CheckPassword(
// 		c.Request.Context(),
// 		subject.UserID,
// 		req.OldPassword,
// 	); err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	if err := h.userService.UpdatePassword(
// 		c.Request.Context(),
// 		subject.UserID,
// 		req.NewPassword,
// 	); err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	c.JSON(
// 		http.StatusOK,
// 		api.SuccessResponse{
// 			Message: "password updated successfully",
// 		},
// 	)
// }

func (h *Handler) GetMe(c *gin.Context) {
	u, err := middleware.GetUserWithContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	user, err := h.userService.GetUserByID(
		c.Request.Context(),
		u.UserID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Data: UserResponse{
				ID:             user.ID,
				Role:           string(user.Role),
				Email:          user.Email,
				SuspendedUntil: user.SuspendedUntil,
				LockedUntil:    user.LockedUntil,
				CreatedAt:      user.CreatedAt,
				UpdatedAt:      user.UpdatedAt,
			},
		},
	)
}

func (h *Handler) VerifyOTP(c *gin.Context) {
	ctx := c.Request.Context()
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleValidationError(c, err)
		return
	}

	key := "otp:verify:" + req.Email
	allowed, err := h.limiter.Allow(
		ctx,
		key,
		VerifyOTPLimit,
		10*time.Minute,
	)
	if err != nil {
		c.Error(err)
		return
	}

	if !allowed {
		c.Error(
			security.NewSecureError(
				http.StatusTooManyRequests,
				"RATE_LIMIT_EXCEEDED",
				"too many verification attempts",
				nil,
			),
		)
		return
	}

	if err := h.authService.VerifyOTP(ctx, req); err != nil {
		c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Message: "OK",
		},
	)
}

func (h *Handler) ResendOTP(c *gin.Context) {
	ctx := c.Request.Context()

	var req ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleValidationError(c, err)
		return
	}

	key := "otp:verify:" + req.Email
	allowed, err := h.limiter.Allow(
		ctx,
		key,
		ResendOTPLimit,
		1*time.Hour,
	)
	if err != nil {
		c.Error(err)
		return
	}

	if !allowed {
		c.Error(
			security.NewSecureError(
				http.StatusTooManyRequests,
				"RATE_LIMIT_EXCEEDED",
				"too many otp resend requests",
				nil,
			),
		)
		return
	}

	// we should not return an error here, as the email may not exist
	// if the email does not exist, we will simply not send an OTP
	err = h.authService.SendEmailOTP(ctx, req.Email)
	if err != nil {
		h.logger.Error(
			"send otp email",
			log.Meta{
				"Error": err,
			},
		)
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Message: "go check your email !",
		},
	)
}
