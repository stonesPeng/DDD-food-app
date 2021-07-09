/**
  @author: honor
  @since: 2021/7/9
  @desc: //TODO
**/
package application

import (
	"DDD-food-app/domain/domain_port"
	"DDD-food-app/domain/entity"
)

type userApp struct {
	userRepo domain_port.UserRepository
}

var _ UserAppInterface = &userApp{}

type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUsers() ([]entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
}

func (u *userApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return u.userRepo.SaveUser(user)
}

func (u *userApp) GetUsers() ([]entity.User, error) {
	return u.userRepo.GetUsers()
}

func (u *userApp) GetUser(id uint64) (*entity.User, error) {
	return u.userRepo.GetUser(id)
}

func (u *userApp) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	return u.userRepo.GetUserByEmailAndPassword(user)
}
