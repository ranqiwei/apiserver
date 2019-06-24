package user

import (
	"apiserver/handler"
	"apiserver/model"
	"apiserver/pkg/errno"
	"apiserver/util"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

//在控制台和日志打印错误，用户返回的错误由c.JSON写内容
/*func Create(c *gin.Context) {
	//定义一个变量
	var r CreateRequest

	if err := c.Bind(&r); err != nil { //当Bind error的时候就已经添加code了，下面的操作等于重写code
		handler.SendResponse(c, errno.ErrBind, nil)
		//c.JSON(http.StatusOK, gin.H{"error": errno.ErrBind}) //GIN报错Headers were already written. Wanted to override status code 400 with 200
		return
	}

	admin2 := c.Param("username") //URL做参数
	log.Infof("URL username:%s", admin2)
	desc := c.Query("desc") //URL中带的参数
	log.Infof("URL key param desc: %s", desc)

	contentType := c.GetHeader("Content-Type")
	log.Infof("Header Content-Type:%s", contentType)

	log.Debugf("username is: [%s], password is [%s]", r.Username, r.Password)

	if r.Username == "" { //带code
		handler.SendResponse(c, errno.New(errno.ErrUserNotFound, fmt.Errorf("username can not found in db: xx.xx.xx.xx")), nil)
		return
	}

	/*if errno.IsErrUserNotFound(err) {
		log.Debug("err type is ErrUserNotFound")
	}

	if r.Password == "" {
		handler.SendResponse(c, fmt.Errorf("password is empty"), nil)
		return
	}

	rsp := CreateResponse{Username: r.Username,}
	handler.SendResponse(c, nil, rsp)
}
*/

// @Summary Add new user to the database
// @Description Add a new user
// @Tags user
// @Accept json
// @Param user body user.CreateRequest true "Create a new user"
// @Success 200 {object} user.CreateResponse "{"code":0,"message":"OK","data":{"username":"kong"}}"
// @Router /user [post]
func Create(c *gin.Context) {
	log.Info("User Create function called", lager.Data{"X-Request-Id": util.GetReqId(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	u := model.UserModel{Username: r.Username, Password: r.Password}

	if err := u.Validate(); err != nil {
		handler.SendResponse(c, errno.ErrValidation, nil)
		return
	}
	if err := u.Encrypt(); err != nil {
		handler.SendResponse(c, errno.ErrEncrypt, nil)
		return
	}

	if err := u.Create(); err != nil {
		log.Debugf("err is:%s", err.Error())
		handler.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	rsp := CreateResponse{Username: u.Username}
	handler.SendResponse(c, errno.OK, rsp)
}
