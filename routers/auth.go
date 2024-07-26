package routers

import (
	"go-gin-example/models"
	"go-gin-example/pkg/e"
	"go-gin-example/pkg/logging"
	"go-gin-example/pkg/util"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// @Summary      Get Auth
// @Produce      json
// @Param        username	query	string	true  "username"
// @Param        password	query	string	true  "password"
// @Success 	 200 {string} json "{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /auth [get]
func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{
		Username: username,
		Password: password,
	}
	ok, _ := valid.Valid(&a)
	code := e.INVALID_PARAMS
	data := make(map[string]interface{})

	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {

			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				code = e.SUCCESS
				data["token"] = token
			}
		} else {
			code = e.ERROR_AUTH
		}
	} else {
		for _, err := range valid.Errors {
			// log.Println(err.Key, err.Message)
			logging.Info(err.Key, err.Message)
		}

	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
