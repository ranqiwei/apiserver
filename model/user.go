package model

import (
	"apiserver/pkg/constvar"
	"fmt"
	"apiserver/pkg/auth"
	"gopkg.in/go-playground/validator.v9"
)

//数据库表相应的操作

type UserModel struct {
	BaseModel
	Username string `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Password string `json:"password" gorm:"column:password;not null" binding:"required" validate:"min=5,max=128"`
}

func (u *UserModel) TableName() string {
	return "tb_users"
}

func (u *UserModel) Create() error {
	return DB.Self.Create(&u).Error
}

func DeleteUser(id uint64) error {
	user := UserModel{}    //生成空model
	user.BaseModel.Id = id //指定id
	return DB.Self.Delete(&user).Error
}

func (u *UserModel) Update() error {
	return DB.Self.Save(&u).Error
}

func GetUser(username string) (*UserModel, error) {
	u := &UserModel{}
	d := DB.Self.Where("username = ?", username).First(&u)
	return u, d.Error
}

func ListUser(username string, offset, limit int) ([]*UserModel, uint64, error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}
	users := make([]*UserModel, 0)
	var count uint64

	//查询操作
	where := fmt.Sprintf("username like '%%%s%%'", username)
	if err := DB.Self.Model(&UserModel{}).Where(where).Count(&count).Error; err != nil {
		return users, count, err
	}

	//根据条件在查询结果中筛选结果
	if err := DB.Self.Where(where).Offset(offset).Limit(limit).Order("id desc").Find(&users).Error; err != nil {
		return users, count, err
	}

	return users, count, nil
}

func (u *UserModel) Compare(pwd string) (err error) {
	err = auth.Compare(u.Password, pwd)
	return
}

//encrypt the user password
func (u *UserModel) Encrypt() (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	return err
}

//validate the fields
func (u *UserModel) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
