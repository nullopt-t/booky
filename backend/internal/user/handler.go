package user

import (
	"booky-backend/internal/middleware"
	"booky-backend/internal/model"
	"booky-backend/internal/shared/token"
	"booky-backend/pkg/api"
	"booky-backend/pkg/api/security"
	"booky-backend/pkg/config"
	"booky-backend/pkg/logger"
	"booky-backend/pkg/utils/jwt"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

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
	service UserService
	config  *config.Config
}

func NewHandler(
	service UserService,
	config *config.Config,
) *Handler {
	return &Handler{
		service: service,
		config:  config,
	}
}

func (h *Handler) handleError(
	c *gin.Context,
	err error,
) {
	if serr, ok := errors.AsType[*security.SecureError](err); ok {
		logger.Log(
			logger.ERROR,
			serr.LogMessage(),
			logger.LMeta{
				"fields": serr.Fields,
			},
		)

		c.JSON(
			serr.Status,
			api.ErrorResponse{
				Code:    serr.Code,
				Message: serr.Error(),
				Details: serr.Fields,
			},
		)
		return
	}

	logger.Log(
		logger.ERROR,
		err.Error(),
	)

	c.JSON(
		http.StatusInternalServerError,
		api.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "internal error please try again later",
		},
	)
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
			h.handleError(
				c,
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
			return
		}
		h.handleError(
			c,
			security.NewSecureError(
				http.StatusBadRequest,
				security.CodeValidation,
				"bad request data",
				err,
			),
		)
		return
	}

	userID, err := h.service.CreateUser(
		c.Request.Context(),
		req,
	)
	if err != nil {
		h.handleError(c, err)
		return
	}

	user, err := h.service.GetUserByID(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		h.handleError(c, err)
		return
	}

	subject, err := json.Marshal(
		token.UserSubject{
			UserID:   user.ID,
			UserRole: user.Role,
		},
	)
	if err != nil {
		h.handleError(c, err)
		return
	}

	token, err := jwt.CreateToken(
		string(subject),
		h.config.JwtSecretKey,
		jwt.AccessTokenTTL,
		jwt.AccessTokenType,
	)
	if err != nil {
		h.handleError(c, err)
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
			h.handleError(
				c,
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
			return
		}
		h.handleError(
			c,
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
		h.handleError(
			c,
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
		user, err = h.service.GetUserByEmail(
			c.Request.Context(),
			*req.Email,
		)
	} else {
		user, err = h.service.GetUserByPhone(
			c.Request.Context(),
			*req.Phone,
		)
	}
	if err != nil {
		h.handleError(c, err)
		return
	}

	logger.Log(
		logger.DEBUG,
		"fetched user",
		logger.LMeta{
			"user": user,
		},
	)

	err = h.service.CheckPassword(
		c.Request.Context(),
		user.ID,
		req.Password,
	)
	if err != nil {
		h.handleError(c, err)
		return
	}

	subject, err := json.Marshal(
		token.UserSubject{
			UserID:   user.ID,
			UserRole: user.Role,
		},
	)
	if err != nil {
		h.handleError(c, err)
		return
	}

	accessToken, err := jwt.CreateToken(
		string(subject),
		h.config.JwtSecretKey,
		jwt.AccessTokenTTL,
		jwt.AccessTokenType,
	)
	if err != nil {
		h.handleError(
			c,
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
		h.handleError(
			c,
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
			h.handleError(
				c,
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
			return
		}
		h.handleError(
			c,
			security.NewSecureError(
				http.StatusBadRequest,
				security.CodeValidation,
				"bad request data",
				err,
			),
		)
		return
	}

	user, err := h.service.GetUserByID(
		c.Request.Context(),
		uri.UserID,
	)
	if err != nil {
		h.handleError(c, err)
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
		h.handleError(
			c,
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

	users, page, err := h.service.GetAllUsers(
		c.Request.Context(),
		q,
	)
	if err != nil {
		h.handleError(c, err)
		return
	}

	data := make([]UserResponse, 0, len(users))
	for _, user := range users {
		data = append(data, UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			IsInactive: user.IsInactive,
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		})
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Data: data,
			Meta: page,
		},
	)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	var uri UserURIParam
	if err := c.ShouldBindUri(&uri); err != nil {
		h.handleError(
			c,
			security.NewSecureError(
				http.StatusBadRequest,
				security.CodeValidation,
				"bad request data",
				err,
			),
		)
		return
	}

	err := h.service.DeleteUserByID(
		c.Request.Context(),
		uri.UserID,
	)
	if err != nil {
		h.handleError(c, err)
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
			h.handleError(
				c,
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
			return
		}
		h.handleError(
			c,
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
		h.handleError(
			c,
			security.NewSecureError(
				http.StatusUnauthorized,
				security.CodeAuth,
				"invalid refresh token",
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
		h.handleError(c, err)
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
			h.handleError(
				c,
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
			return
		}
		h.handleError(
			c,
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
	user, err := h.service.GetUserByEmail(
		c.Request.Context(),
		req.Email,
	)
	if err != nil {
		h.handleError(c, err)
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
			h.handleError(c, err)
			return
		}

		resetToken, err = jwt.CreateToken(
			string(subjectStr),
			h.config.JwtSecretKey,
			jwt.ResetPassTokenTTL,
			jwt.ResetPassTokenType,
		)
		if err != nil {
			h.handleError(c, err)
			return
		}

		err = h.service.SetResetToken(
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

	logger.Log(
		logger.INFO,
		"token created",
		logger.LMeta{
			"hash":  resetToken,
			"email": user.Email,
			"id":    user.ID,
		},
	)

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
			h.handleError(
				c,
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
			return
		}
		h.handleError(
			c,
			security.NewSecureError(
				http.StatusBadRequest,
				security.CodeValidation,
				"bad request data",
				err,
			),
		)
		return
	}

	logger.Log(
		logger.DEBUG,
		"verify password reset token",
		logger.LMeta{
			"token":        req.Token,
			"old_password": req.OldPassword,
			"new_password": req.NewPassword,
		},
	)

	claims, err := jwt.VerifyToken(
		req.Token,
		h.config.JwtSecretKey,
	)
	if err != nil {
		h.handleError(c, err)
		return
	}

	var subject token.UserSubject
	if err := json.Unmarshal(
		[]byte(claims.Subject),
		&subject,
	); err != nil {
		h.handleError(c, err)
		return
	}

	if err := h.service.CheckPasswordResetToken(
		c.Request.Context(),
		subject.UserID,
		req.Token,
	); err != nil {
		h.handleError(c, err)
		return
	}

	if claims.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest,
			api.ErrorResponse{
				Code:    "EXPIRED_TOKEN",
				Message: "expired reset token",
			},
		)
		return
	}

	if err := h.service.CheckPassword(
		c.Request.Context(),
		subject.UserID,
		req.OldPassword,
	); err != nil {
		h.handleError(c, err)
		return
	}

	if err := h.service.UpdatePassword(
		c.Request.Context(),
		subject.UserID,
		req.NewPassword,
	); err != nil {
		h.handleError(c, err)
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
		h.handleError(c, err)
		return
	}

	user, err := h.service.GetUserByID(
		c.Request.Context(),
		u.UserID,
	)
	if err != nil {
		h.handleError(c, err)
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
			h.handleError(
				c,
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
			return
		}
		h.handleError(
			c,
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
		h.handleError(c, err)
		return
	}

	if err := h.service.VerifyEmailOTP(
		c.Request.Context(),
		u.UserID,
		req.Otp,
	); err != nil {
		h.handleError(c, err)
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
		h.handleError(c, err)
		return
	}

	err = h.service.ResendEmailOTP(
		c.Request.Context(),
		u.UserID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Message: "email sent successfully",
		},
	)
}

func (h *Handler) VerifyPhoneOTP(c *gin.Context) {
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
			h.handleError(
				c,
				security.NewSecureError(
					http.StatusBadRequest,
					security.CodeValidation,
					"bad request data",
					err,
				).WithFields(fieldErrors),
			)
			return
		}
		h.handleError(
			c,
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
		h.handleError(c, err)
		return
	}

	if err := h.service.VerifyPhoneOTP(
		c.Request.Context(),
		u.UserID,
		req.Otp,
	); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Message: "email verified successfully",
		},
	)
}

func (h *Handler) ResendPhoneOTP(c *gin.Context) {
	u, err := middleware.GetUserWithContext(c)
	if err != nil {
		h.handleError(c, err)
		return
	}

	err = h.service.ResendPhoneOTP(
		c.Request.Context(),
		u.UserID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		api.SuccessResponse{
			Message: "email sent successfully",
		},
	)
}
