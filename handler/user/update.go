package user

import (
	"github.com/gin-gonic/gin"
	"apiserver/model"
	"apiserver/handler"
	"apiserver/pkg/errno"
	"strconv"
	"github.com/sirupsen/logrus"
	"github.com/lexkong/log/lager"
	"apiserver/util"
)

func Update(c *gin.Context) {
	logrus.Info("Update function called.", lager.Data{"X-Request-Id": util.GetReqId(c)})
	UserId, _ := strconv.Atoi(c.Param(":id"))

	var u model.UserModel
	if err := c.Bind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	u.Id = uint64(UserId)

	if err := u.Validate(); err != nil {
		handler.SendResponse(c, errno.ErrValidation, nil)
		return
	}
	if err := u.Encrypt(); err != nil {
		handler.SendResponse(c, errno.ErrEncrypt, nil)
		return
	}
	if err := u.Update(); err != nil {
		handler.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	handler.SendResponse(c, nil, nil)
}
