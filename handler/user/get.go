package user

import (
	"github.com/gin-gonic/gin"
	"apiserver/model"
	"apiserver/handler"
)

func Get(c *gin.Context) {
	username := c.Param("username")
	user, err := model.GetUser(username)
	if err != nil {
		handler.SendResponse(c, err, nil)
		return
	}
	handler.SendResponse(c, nil, user)

}
