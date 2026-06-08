package otp

import "booky-backend/pkg/utils"

var OTPLength map[string]int = map[string]int{
	"login": 6,
	"reset": 6,
}

type OTPGenerator interface {
	GenerateOTP(
		otpType string,
	) (string, error)
}

type OTPGeneratorImpl struct{}

func NewGenerator() OTPGenerator {
	return &OTPGeneratorImpl{}
}

func (g *OTPGeneratorImpl) GenerateOTP(otpType string) (string, error) {
	if _, ok := OTPLength[otpType]; !ok {
		otpType = "login"
	}

	otp, err := utils.GenerateOTP(OTPLength[otpType])
	if err != nil {
		return "", err
	}

	return otp, nil
}
