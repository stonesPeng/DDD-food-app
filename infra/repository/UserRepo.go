/**
  @author: honor
  @since: 2021/7/8
  @desc: 数据库持久化层
**/
package repository

import (
	"DDD-food-app/domain/domain_port"
	"DDD-food-app/domain/entity"
	"DDD-food-app/infra/security"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

type UserRepo struct {
	db *gorm.DB
}

//UserRepo 是 UserRepository 接口的具体实现
var _ domain_port.UserRepository = &UserRepo{}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	errMsg := map[string]string{}
	if err := repo.db.Debug().Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			errMsg["email_taken"] = "database error"
			return nil, errMsg
		}
		errMsg["db_error"] = "database error"
		return nil, errMsg
	} else {
		return user, nil
	}
}

func (repo *UserRepo) GetUser(id uint64) (*entity.User, error) {
	var user entity.User
	if err := repo.db.Debug().Where("id = ?", id).Take(&user).Error; err != nil {
		return nil, err
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
	}
	return &user, nil
}

func (repo *UserRepo) GetUsers() ([]entity.User, error) {
	var users []entity.User
	if err := repo.db.Debug().Find(&users).Error; err != nil {
		return nil, err
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
	}
	return users, nil
}

func (repo *UserRepo) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	var u entity.User
	errMsg := map[string]string{}
	if err := repo.db.Debug().Where("email = ?", user.Email).Take(&u).Error; err != nil {
		errMsg["db_error"] = "db error"
		return nil, errMsg
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errMsg["user_notfound"] = "user not found"
			return nil, errMsg
		}
	}
	if err := security.VerifyPassword(user.Password, u.Password); err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		errMsg["incorrect_password"] = "incorrect_ password"
		return nil, errMsg
	}
	return &u, nil
}
