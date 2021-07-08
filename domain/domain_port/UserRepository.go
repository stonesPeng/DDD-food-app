/**
  @author: honor
  @since: 2021/7/8
  @desc: //TODO
**/
package domain_port

import "DDD-food-app/domain/entity"

type UserRepository interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUser(uint64) (*entity.User, error)
	GetUsers() ([]entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
}
