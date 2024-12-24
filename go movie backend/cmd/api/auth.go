package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	Issuer      string
	Audience    string
	Secret      string
	TokenExpiry time.Duration
	RefreshExpiry time.Duration
	CookieDomain string
	CookiePath string
	CookieName string
}
type jwtUser struct{
	ID int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}
type TokenPairs struct{
	Token string `json:"access_token"`
	Refresh string `json:"refresh_token"`

}
type Claims struct{
	jwt.RegisteredClaims
}

func (j *Auth) GenerateTokenPair(user *jwtUser)(TokenPairs,error){
	token:=jwt.New(jwt.SigningMethodHS256)
	claims:=token.Claims.(jwt.MapClaims)
	claims["name"]=fmt.Sprintf("%s %s",user.FirstName,user.LastName)
    claims["sub"]=fmt.Sprint(user.ID)
	claims["aud"]=j.Audience
	claims["iss"]=j.Issuer
	claims["iat"]=time.Now().UTC().Unix()
	claims["typ"]="JWT"

	claims["exp"]=time.Now().UTC().Add(j.TokenExpiry).Unix()

	signedAccessToken,err:=token.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{},err
	}
	fmt.Println("ok signed")
	refreshToken:=jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims:=refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"]=fmt.Sprint(user.ID)
	refreshTokenClaims["iat"]=time.Now().UTC().Unix()
	refreshTokenClaims["exp"]=time.Now().UTC().Add(j.RefreshExpiry).Unix()


	signedRefreshToken,err:=refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{},err
	}
	fmt.Println("token pair generated")

	var tokenPairs =TokenPairs{
		Token: signedAccessToken,
		Refresh: signedRefreshToken,

	}
	return tokenPairs,nil
	}
	func (j *Auth) GetRefershCookie (refreshToken string) *http.Cookie{
		return &http.Cookie{
			Name:j.CookieName,
			Path: j.CookiePath,
			Value: refreshToken,
			Expires:time.Now().Add(j.RefreshExpiry),
			MaxAge: int(j.RefreshExpiry.Seconds()),
			Secure: true,
			Domain: j.CookieDomain,
			SameSite: http.SameSiteNoneMode,
			HttpOnly: true,
		}
	}
	func (j *Auth) GetExpiredRefershCookie (refreshToken string) *http.Cookie{
		return &http.Cookie{
			Name:j.CookieName,
			Path: j.CookiePath,
			Value: "",
			Expires:time.Unix(0,0),
			MaxAge: -1,
			Secure: true,
			Domain: j.CookieDomain,
			SameSite: http.SameSiteNoneMode,
			HttpOnly: true,
		}
	}