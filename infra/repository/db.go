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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", user, password, host, port, name, "10s")
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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

func (s *Repositories) AutoMigrate() {
	s.db.AutoMigrate(&entity.User{}, &entity.PublicUser{})
}
