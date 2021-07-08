/**
  @author: honor
  @since: 2021/7/8
  @desc: //TODO
**/
package repository

import (
	"DDD-food-app/domain/domain_port"
	"DDD-food-app/domain/entity"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Repositories struct {
	User domain_port.UserRepository
	db   *gorm.DB
}

func NewRepositories(user, password, port, host, name string) (*Repositories, error) {
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, name, password)
	if db, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, //禁用外键约束
	}); err != nil {
		return nil, err
	} else {
		return &Repositories{
			User: NewUserRepository(db),
			db:   db,
		}, nil
	}
}

func (s *Repositories) AutoMigrate() string {
	return s.db.AutoMigrate(&entity.User{}).Error()
}
