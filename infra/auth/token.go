/**
  @author: honor
  @since: 2021/7/9
  @desc: //TODO
**/
package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Token struct {
}

func NewToken() *Token {
	return &Token{}
}

var _ TokenInterface = &Token{}

type TokenInterface interface {
	CreateToken(userid uint64) (*TokenDetails, error)
	ExtractTokenMetadata(*http.Request) (*AccessDetails, error)
}

func (t *Token) CreateToken(userid uint64) (*TokenDetails, error) {
	pAtExpires := time.Now().Add(time.Minute * 15).Unix()
	pUuid := uuid.NewV4().String()
	pRtExpires := time.Now().Add(time.Hour * 24 * 7).Unix()
	pRefreshUuid := pUuid + "++" + strconv.Itoa(int(userid))
	token := &TokenDetails{AtExpires: pAtExpires, TokenUuid: pUuid, RtExpires: pRtExpires, RefreshToken: pRefreshUuid}
	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = token.TokenUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = token.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = token.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = token.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	token.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return token, nil
}
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

//get the token from the request body
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (t *Token) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	fmt.Println("WE ENTERED METADATA")
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			TokenUuid: accessUuid,
			UserId:    userId,
		}, nil
	}
	return nil, err
}
