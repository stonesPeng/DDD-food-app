/**
  @author: honor
  @since: 2021/7/8
  @desc: //TODO
**/
package entity

import (
	"DDD-food-app/infra/security"
	"github.com/badoux/checkmail"
	"gorm.io/gorm"
	"html"
	"strings"
	"time"
)

type User struct {
	gorm.Model
	FirstName string `gorm:"size:100;not null;" json:"first_name"`
	LastName  string `gorm:"size:100;not null;" json:"last_name"`
	Email     string `gorm:"size:100;not null;unique" json:"email"`
	Password  string `gorm:"size:100;not null;" json:"password"`
}

type PublicUser struct {
	ID        uint64 `gorm:"primary_key;auto_increment" json:"id,omitempty,string"`
	FirstName string `gorm:"size:100;not null;" json:"first_name"`
	LastName  string `gorm:"size:100;not null;" json:"last_name"`
}

/**
 * @Description: 特殊处理密码，不可逆
 * @receiver u
 * @return error
 */
func (u *User) BeforeSave(tx *gorm.DB) error {
	if result, err := security.Hash(u.Password); err != nil {
		return err
	} else {
		u.Password = *result
		return nil
	}
}

type Users []User

func (users Users) PublishUsers() []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.PublishUser()
	}
	return result
}

func (u *User) PublishUser() interface{} {
	return &PublicUser{
		ID:        uint64(u.ID),
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func (u *User) Prepare() {
	html.EscapeString(strings.TrimSpace(u.FirstName))
	html.EscapeString(strings.TrimSpace(u.LastName))
	html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) map[string]string {
	var errMsg = make(map[string]string)
	var err error
	switch strings.ToLower(action) {
	case "update":
		if u.Email == "" {
			errMsg["email_required"] = "email required"
		} else {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errMsg["invalid_email"] = "invalid email"
			}
		}
	case "login":
		if u.Password == "" {
			errMsg["password_required"] = "password is required"
		}
		if u.Email == "" {
			errMsg["email_required"] = "email is required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errMsg["invalid_email"] = "please provide a valid email"
			}
		}
	case "forgotpassword":
		if u.Email == "" {
			errMsg["email_required"] = "email required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errMsg["invalid_email"] = "please provide a valid email"
			}
		}
	default:
		if u.FirstName == "" {
			errMsg["firstname_required"] = "first name is required"
		}
		if u.LastName == "" {
			errMsg["lastname_required"] = "last name is required"
		}
		if u.Password == "" {
			errMsg["password_required"] = "password is required"
		}
		if u.Password != "" && len(u.Password) < 6 {
			errMsg["invalid_password"] = "password should be at least 6 characters"
		}
		if u.Email == "" {
			errMsg["email_required"] = "email is required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errMsg["invalid_email"] = "please provide a valid email"
			}
		}
	}
	return errMsg
}
