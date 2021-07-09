/**
  @author: honor
  @since: 2021/7/9
  @desc: //TODO
**/
package auth

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

type AuthInterface interface {
	CreateAuth(uint64, *TokenDetails) error
	FetchAuth(string) (uint64, error)
	DeleteRefresh(string) error
	DeleteTokens(*AccessDetails) error
}

var _ AuthInterface = &ClientData{}

type ClientData struct {
	client *redis.Client
}

func NewAuth(client *redis.Client) *ClientData {
	return &ClientData{client: client}
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	TokenUuid string
	UserId    uint64
}

func (c *ClientData) CreateAuth(userid uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	atCreated, err := c.client.Set(td.TokenUuid, strconv.Itoa(int(userid)), at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := c.client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}
	return nil
}
func (c *ClientData) FetchAuth(tokenUuid string) (uint64, error) {
	userid, err := c.client.Get(tokenUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}
func (c *ClientData) DeleteRefresh(refreshUuid string) error {
	//delete refresh token
	deleted, err := c.client.Del(refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}
	return nil
}
func (c *ClientData) DeleteTokens(authD *AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", authD.TokenUuid, authD.UserId)
	//delete access token
	deletedAt, err := c.client.Del(authD.TokenUuid).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := c.client.Del(refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}
