/**
  @author: honor
  @since: 2021/7/9
  @desc: //TODO
**/
package intefaces

import (
	"DDD-food-app/application"
	"DDD-food-app/domain/entity"
	"DDD-food-app/infra/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserInterface struct {
	userApp      application.UserAppInterface
	authService  auth.AuthInterface
	tokenService auth.TokenInterface
}

func NewUsers(userApp application.UserAppInterface, authService auth.AuthInterface, tokenService auth.TokenInterface) *UserInterface {
	return &UserInterface{
		userApp:      userApp,
		authService:  authService,
		tokenService: tokenService,
	}
}

func (userInterface *UserInterface) SaveUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid json",
		})
		return

	}
	validateErr := user.Validate("")
	if len(validateErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, validateErr)
		return
	}
	if newUser, err := userInterface.userApp.SaveUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	} else {
		c.JSON(http.StatusCreated, newUser.PublishUser())
	}
}

func (userInterface *UserInterface) GetUsers(c *gin.Context) {
	users := entity.Users{} //customize user
	var err error
	//us, err = application.UserApp.GetUsers()
	users, err = userInterface.userApp.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users.PublishUsers())
}

func (userInterface *UserInterface) GetUser(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := userInterface.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user.PublishUser())
}
