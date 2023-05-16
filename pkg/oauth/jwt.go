package oauth

import (
	"context"
	"errors"
	"fmt"
	"micro/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	JWT oautInterface = &AccessDetails{}
)

type oautInterface interface {
	VerifyToken(ctx context.Context, request string) (*AccessDetails, error)
}

// AccessDetails struct
type AccessDetails struct {
	UserID string
	Dialer string
	JwtID  string
	NonOTP string
	Aud    string
}

func (a *AccessDetails) VerifyToken(ctx context.Context, request string) (*AccessDetails, error) {
	token, err := jwt.Parse(request, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(config.C().JWT.HMACSecret), nil
		}
		if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			return jwt.ParseRSAPublicKeyFromPEM([]byte(config.C().JWT.RSASecret))
		}
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); ok {
			return jwt.ParseECPublicKeyFromPEM([]byte(config.C().JWT.ECDSASecret))
		}

		zap.L().Error(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		return nil, errors.New("new errors")
	})
	if err != nil {
		return nil, errors.New("new errors")
	}
	if token.Claims == nil || !token.Valid {
		return nil, errors.New("new errors")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID, ok := claims["sub"].(string)
		if !ok {
			zap.L().Error(fmt.Sprintf("error in get user sub from token: %v", claims))
			return nil, errors.New("new errors")
		}

		if _, err := uuid.Parse(userID); err != nil {
			zap.L().Error(fmt.Sprintf("error in sub typo: %v", err))
			return nil, errors.New("new errors")
		}

		dialer, ok := claims["mobile"].(string)
		if !ok {
			zap.L().Error(fmt.Sprintf("error in get user mobile from token: %v", claims))
			return nil, errors.New("new errors")
		}
		nonOtp, ok := claims["non_otp"].(string)
		if !ok {
			zap.L().Error(fmt.Sprintf("error in get user nonOtp from token: %v", claims))
			return nil, errors.New("new errors")
		}
		jwtID, ok := claims["jti"].(string)
		if !ok {
			zap.L().Error(fmt.Sprintf("error in get JwtID from token: %v", claims))
			return nil, errors.New("new errors")
		}
		aud, ok := claims["aud"].(string)
		if !ok {
			zap.L().Error(fmt.Sprintf("error in get user Aud from token: %v", claims))
			return nil, errors.New("new errors")
		}

		return &AccessDetails{
			UserID: userID,
			Dialer: dialer,
			JwtID:  jwtID,
			NonOTP: nonOtp,
			Aud:    aud,
		}, nil
	}
	return nil, err
}

func ContextWithUserID(c context.Context, id string) context.Context {
	return context.WithValue(c, USER_ID, id) // the user uuid
}

func ContextWithUserMobile(c context.Context, mobile string) context.Context {
	return context.WithValue(c, DIALER, mobile) // the user mobile
}

func UserIDFromContext(c context.Context) (string, error) {
	userID, ok := c.Value(USER_ID).(string)
	if !ok {
		return "", errors.New("unable to get user id from context")
	}
	return userID, nil
}

func UserMobileFromContext(c context.Context) (string, error) {
	dialer, ok := c.Value(DIALER).(string)
	if !ok {
		return "", errors.New("unable to get user mobile from context")
	}
	return dialer, nil
}

type jwtPayload string

const (
	USER_ID jwtPayload = "user_id"
	DIALER  jwtPayload = "dialer"
)
