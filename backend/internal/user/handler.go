package user

import (
	"booky-backend/internal/middleware"
	"booky-backend/internal/model"
	"booky-backend/internal/shared/token"
	"booky-backend/pkg/api"
	"booky-backend/pkg/api/security"
	"booky-backend/pkg/config"
	"booky-backend/pkg/utils/jwt"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	emailOTPPrefix = "email"
	phoneOTPPrefix = "phone"
)

type GetAllUsersResponse struct {
	Users []UserResponse `json:"users"`
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
type UserURIParam struct {
	UserID uuid.UUID `json:"id" validate:"required,uuid"`
}

type UserResponse struct {
	ID              uuid.UUID  `json:"id"`
	Role            string     `json:"role"`
	Email           *string    `json:"email,omitempty"`
	IsEmailVerified bool       `json:"is_email_verified"`
	Phone           *string    `json:"phone,omitempty"`
	IsPhoneVerified bool       `json:"is_phone_verified"`
	IsInactive      bool       `json:"is_inactive"`
	LockedUntil     *time.Time `json:"locked_unitl,omitzero"`
	CreatedAt       time.Time  `json:"created_at,omitzero"`
	UpdatedAt       time.Time  `json:"updated_at,omitzero"`
}

func ToUserResponse(user model.User) UserResponse {
	return UserResponse{
		ID:              user.ID,
		Role:            string(user.Role),
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		Phone:           user.Phone,
		IsPhoneVerified: user.IsPhoneVerified,
		IsInactive:      user.IsInactive,
		LockedUntil:     user.LockedUntil,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
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
	Email    *string `json:"email" binding:"omitempty,email,required_without=Phone,excluded_with=Phone"`
	Phone    *string `json:"phone" binding:"omitempty,e164,required_without=Email,excluded_with=Email"`
	Password string  `json:"password" binding:"required,min=8"`
}

type RegisterUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginUserRequest struct {
	Email    *string `json:"email" binding:"omitempty,email"`
	Phone    *string `json:"phone" binding:"omitempty,e164,excluded_with=Email"`
	Password string  `json:"password" binding:"required"`
}

type Handler struct {
	userService UserService
	otpService  OTPService
	config      *config.Config
}

func NewHandler(
	userService UserService,
	otpService OTPService,
	config *config.Config,
) *Handler {
	return &Handler{
		userService: userService,
		otpService:  otpService,
		config:      config,
	}
}

func (h *Handler) UserRegister(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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
		return
	}

	userID, err := h.userService.CreateUser(
		c.Request.Context(),
		req,
	)
	if err != nil {
		c.Error(err)
		return
	}

	user, err := h.userService.GetUserByID(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.otpService.SendOTP(
		c.Request.Context(),
		userID,
		"register",
	)
	if err != nil {
		c.Error(err)
		return
	}

	subject, err := json.Marshal(
		token.UserSubject{
			UserID:          user.ID,
			UserRole:        user.Role,
			IsEmailVerified: user.IsEmailVerified,
		},
	)
	if err != nil {
		c.Error(err)
		return
	}

	token, err := jwt.CreateToken(
		string(subject),
		h.config.JwtSecretKey,
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
			Data: RegisterUserResponse{
				AccessToken: token,
			},
		},
	)
}

func (h *Handler) UserLogin(c *gin.Context) {
	var req LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok && ve != nil {
			fieldErrors := make([]api.FieldError, 0, len(ve))
			for _, e := range ve {
				fieldErrors = append(fieldErrors, api.FieldError{
					Field: e.Field(),
					Tags:  e.Tag(),
				})
			}
			c.Error(
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
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
		return
	}

	if (req.Email == nil) &&
		(req.Phone == nil) {
		c.Error(
			security.NewSecureError(
				http.StatusBadRequest,
				security.CodeValidation,
				"bad request data",
				errors.New("email or phone is required"),
			),
		)
		return
	}

	var user *model.User
	var err error
	if req.Email != nil {
		user, err = h.userService.GetUserByEmail(
			c.Request.Context(),
			*req.Email,
		)
	} else {
		user, err = h.userService.GetUserByPhone(
			c.Request.Context(),
			*req.Phone,
		)
	}
	if err != nil {
		c.Error(err)
		return
	}

	err = h.userService.CheckPassword(
		c.Request.Context(),
		user.ID,
		req.Password,
	)
	if err != nil {
		c.Error(err)
		return
	}

	subject, err := json.Marshal(
		token.UserSubject{
			UserID:          user.ID,
			UserRole:        user.Role,
			IsEmailVerified: user.IsEmailVerified,
		},
	)
	if err != nil {
		c.Error(err)
		return
	}

	accessToken, err := jwt.CreateToken(
		string(subject),
		h.config.JwtSecretKey,
		jwt.AccessTokenTTL,
		jwt.AccessTokenType,
	)
	if err != nil {
		c.Error(
			security.NewSecureError(
				http.StatusInternalServerError,
				security.CodeInternal,
				"internal error, please try again",
				err,
			),
		)
		return
	}

	refreshToken, err := jwt.CreateToken(
		string(subject),
		h.config.JwtSecretKey,
		jwt.RefreshTokenTTL,
		jwt.RefreshTokenType,
	)
	if err != nil {
		c.Error(
			security.NewSecureError(
				http.StatusInternalServerError,
				security.CodeInternal,
				"internal error, please try again",
				err,
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Data: RegisterUserResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		},
	)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	var uri UserURIParam
	if err := c.ShouldBindJSON(&uri); err != nil {
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok && ve != nil {
			fieldErrors := make([]api.FieldError, 0, len(ve))
			for _, e := range ve {
				fieldErrors = append(fieldErrors, api.FieldError{
					Field: e.Field(),
					Tags:  e.Tag(),
				})
			}
			c.Error(
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
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
				ID:              user.ID,
				Email:           user.Email,
				Phone:           user.Phone,
				IsEmailVerified: user.IsEmailVerified,
				IsPhoneVerified: user.IsPhoneVerified,
				IsInactive:      user.IsInactive,
				LockedUntil:     user.LockedUntil,
				CreatedAt:       user.CreatedAt,
				UpdatedAt:       user.UpdatedAt,
			},
		},
	)
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	var q api.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok && ve != nil {
			fieldErrors := make([]api.FieldError, 0, len(ve))
			for _, e := range ve {
				fieldErrors = append(fieldErrors, api.FieldError{
					Field: e.Field(),
					Tags:  e.Tag(),
				})
			}
			c.Error(
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
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
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok && ve != nil {
			fieldErrors := make([]api.FieldError, 0, len(ve))
			for _, e := range ve {
				fieldErrors = append(fieldErrors, api.FieldError{
					Field: e.Field(),
					Tags:  e.Tag(),
				})
			}
			c.Error(
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
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
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok && ve != nil {
			fieldErrors := make([]api.FieldError, 0, len(ve))
			for _, e := range ve {
				fieldErrors = append(fieldErrors, api.FieldError{
					Field: e.Field(),
					Tags:  e.Tag(),
				})
			}
			c.Error(
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
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
		return
	}

	claims, err := jwt.VerifyToken(
		req.Refresh_token,
		h.config.JwtSecretKey,
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
		h.config.JwtSecretKey,
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

func (h *Handler) ForgetPassword(c *gin.Context) {
	var req ForgetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok && ve != nil {
			fieldErrors := make([]api.FieldError, 0, len(ve))
			for _, e := range ve {
				fieldErrors = append(fieldErrors, api.FieldError{
					Field: e.Field(),
					Tags:  e.Tag(),
				})
			}
			c.Error(
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
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
		return
	}

	/// check user exists by email
	user, err := h.userService.GetUserByEmail(
		c.Request.Context(),
		req.Email,
	)
	if err != nil {
		c.Error(err)
		return
	}

	var resetToken string
	if user.ResetTokenExpireAt != nil && user.ResetTokenExpireAt.After(time.Now()) {
		resetToken = *user.ResetToken
	} else {
		subjectStr, err := json.Marshal(
			token.UserSubject{
				UserID:   user.ID,
				UserRole: user.Role,
			},
		)
		if err != nil {
			c.Error(err)
			return
		}

		resetToken, err = jwt.CreateToken(
			string(subjectStr),
			h.config.JwtSecretKey,
			jwt.ResetPassTokenTTL,
			jwt.ResetPassTokenType,
		)
		if err != nil {
			c.Error(err)
			return
		}

		err = h.userService.SetResetToken(
			c.Request.Context(),
			user.ID,
			resetToken,
		)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				api.ErrorResponse{
					Code:    "INTERNAL_ERROR",
					Message: err.Error(),
				},
			)
			return
		}
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Message: "email sent successfully",
		},
	)
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var req VerifyResetTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok && ve != nil {
			fieldErrors := make([]api.FieldError, 0, len(ve))
			for _, e := range ve {
				fieldErrors = append(fieldErrors, api.FieldError{
					Field: e.Field(),
					Tags:  e.Tag(),
				})
			}
			c.Error(
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
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
		return
	}

	claims, err := jwt.VerifyToken(
		req.Token,
		h.config.JwtSecretKey,
	)
	if err != nil {
		c.Error(err)
		return
	}

	var subject token.UserSubject
	if err := json.Unmarshal(
		[]byte(claims.Subject),
		&subject,
	); err != nil {
		c.Error(err)
		return
	}

	if err := h.userService.CheckPasswordResetToken(
		c.Request.Context(),
		subject.UserID,
		req.Token,
	); err != nil {
		c.Error(err)
		return
	}

	if claims.ExpiresAt.Before(time.Now()) {
		c.Error(
			security.NewSecureError(
				http.StatusBadRequest,
				security.CodeValidation,
				"expired reset token",
				nil,
			),
		)
		return
	}

	if err := h.userService.CheckPassword(
		c.Request.Context(),
		subject.UserID,
		req.OldPassword,
	); err != nil {
		c.Error(err)
		return
	}

	if err := h.userService.UpdatePassword(
		c.Request.Context(),
		subject.UserID,
		req.NewPassword,
	); err != nil {
		c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Message: "password updated successfully",
		},
	)
}

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
				ID:              user.ID,
				Role:            string(user.Role),
				Email:           user.Email,
				IsEmailVerified: user.IsEmailVerified,
				IsInactive:      user.IsInactive,
				LockedUntil:     user.LockedUntil,
				CreatedAt:       user.CreatedAt,
				UpdatedAt:       user.UpdatedAt,
			},
		},
	)
}

func (h *Handler) VerifyEmailOTP(c *gin.Context) {
	var req VerifyEmailOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok && ve != nil {
			fieldErrors := make([]api.FieldError, 0, len(ve))
			for _, e := range ve {
				fieldErrors = append(fieldErrors, api.FieldError{
					Field: e.Field(),
					Tags:  e.Tag(),
				})
			}
			c.Error(
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
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
		return
	}

	u, err := middleware.GetUserWithContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	if u.IsEmailVerified {
		c.Error(
			security.NewSecureError(
				http.StatusBadRequest,
				security.CodeValidation,
				"email already verified",
				nil,
			),
		)
		return
	}

	if err := h.otpService.VerifyOTP(
		c.Request.Context(),
		u.UserID,
		"register",
		req.Otp,
	); err != nil {
		c.Error(err)
		return
	}

	if err := h.userService.VerifyUserEmail(
		c.Request.Context(),
		u.UserID,
	); err != nil {
		c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Message: "email verified successfully",
		},
	)
}

func (h *Handler) ResendEmailOTP(c *gin.Context) {
	u, err := middleware.GetUserWithContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	if u.IsEmailVerified {
		c.Error(
			security.NewSecureError(
				http.StatusBadRequest,
				security.CodeValidation,
				"email already verified",
				nil,
			),
		)
		return
	}

	err = h.otpService.SendOTP(
		c.Request.Context(),
		u.UserID,
		"register",
	)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Message: "email sent successfully",
		},
	)
}
