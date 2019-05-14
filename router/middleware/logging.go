package middleware

import (
	"apiserver/handler"
	"apiserver/pkg/errno"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/willf/pad"
	"io/ioutil"
	"regexp"
	"time"
)

//used for capture Response
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

//----------------------------------------------
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get start time of middleware run
		start := time.Now().UTC()
		//Get and validate path
		path := c.Request.URL.Path
		reg := regexp.MustCompile("(/v1/user|/login)")
		if !reg.MatchString(path) {
			return
		}
		if path == "/sd/health" || path == "sd/ram" || path == "sd/cpu" || path == "sd/disk" {
			return
		}
		//read the Body content
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		//Get Method
		method := c.Request.Method
		//Get Ip
		ip := c.ClientIP()

		log.Debugf("New request come in,ip:%s,path:%s,Method: %s, body:`%s`", ip, path, method, string(bodyBytes))

		// capture Response
		blw := bodyLogWriter{
			c.Writer,
			bytes.NewBufferString(""),
		}
		c.Writer = blw

		c.Next()

		//calculate the latency
		end := time.Now().UTC()
		latency := end.Sub(start)

		//get code and message
		code, message := -1, ""
		var response handler.Response

		if err := json.Unmarshal(blw.body.Bytes(), &response); err != nil {
			log.Errorf(err, "response body can't unmarsha1 to model.Response struct,body:`%s`", blw.body.Bytes())
			code = errno.ErrBind.Code
			message = err.Error()
		} else {
			code = response.Code
			message = response.Message
		}

		//log info
		log.Infof("%-13s|%-12s|%s %s|{code: %d,message:%s}", latency, ip, pad.Right(method, 5, ""), path, code, message)
	}
}
