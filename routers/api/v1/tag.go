package v1

import (
	"go-gin-example/models"
	"go-gin-example/pkg/e"
	"go-gin-example/pkg/logging"
	"go-gin-example/pkg/setting"
	"go-gin-example/pkg/util"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// 获取多个文章标签
// @Summary      Get Tags
// @Produce      json
// @Param        name			query		string		true  	"Name"
// @Param        state			query		int			false  	"State"
// @Success 	 200 			{string}	json 		"{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /api/v1/tags/ [get]
func GetTags(c *gin.Context) {
	valid := validation.Validation{}

	name := c.Query("name")
	valid.Required(name, "name").Message("名称不为空")
	valid.MaxSize(name, 100, "name").Message("名称最大长度为100字符")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	maps["deleted_on"] = 0
	if name != "" {
		maps["name"] = name
	}
	var state int = -1
	arg := c.Query("state")
	if arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		data["list"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
		data["total"] = models.GetTagTotal(maps)

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

// 新增文章标签
// @Summary      Add Tags
// @Produce      json
// @Param        name			query		string		true  	"Name"
// @Param        state			query		int			false  	"State"
// @Success 	 200 			{string}	json 		"{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /api/v1/tags/ 	[post]
func AddTag(c *gin.Context) {
	name := c.Query("name")
	created_by := c.Query("created_by")

	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt() // 强制类型转换

	// 表单验证
	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不为空")
	valid.MaxSize(name, 100, "name").Message("名称最大长度为100字符")
	valid.Required(created_by, "created_by").Message("创建人不为空")
	valid.MaxSize(created_by, 100, "created_by").Message("创建人最大长度为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, created_by, state)
		} else {
			code = e.ERROR_EXIST_TAG
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
		"data": make(map[string]string),
	})

}

// 修改文章标签
// @Summary      Edit Tags By ID
// @Produce      json
// @Param        id					query		int			true  	"ID"
// @Param        name				query		string		true  	"Name"
// @Param        modified_by		query		string		true  	"ModifiedBy"
// @Param        state				query		int			false  	"State"
// @Success 	 200 				{string}	json 		"{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /api/v1/tags/{id} 	[put]
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modified_by := c.Query("modified_by")

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	valid := validation.Validation{}
	valid.Required(id, "id").Message("标签id 不为空")
	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.Required(name, "name").Message("名称不为空")
	valid.MaxSize(name, 100, "name").Message("名称长度不超过100字符")
	valid.Required(modified_by, "modified_by").Message("更改人不为空")
	valid.MaxSize(modified_by, 100, "modified_by").Message("更改人长度不超过100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistTagByID(id) {
			code = e.ERROR_NOT_EXIST_TAG
		} else {
			code = e.SUCCESS

			data := make(map[string]interface{})
			data["name"] = name
			data["modified_by"] = modified_by
			data["state"] = state
			data["deleted_on"] = 0
			models.EditTag(id, data)
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
		"data": map[string]string{},
	})
}

// 删除文章标签
// @Summary      Delete Tags By ID
// @Produce      json
// @Param        id					query		int			true  	"ID"
// @Success 	 200 				{string}	json 		"{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /api/v1/tags/{id} 	[delete]
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Required(id, "id").Message("标签id不为空")
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistTagByID(id) {
			code = e.ERROR_NOT_EXIST_TAG
		} else {
			code = e.SUCCESS
			models.DeleteTag(id)
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
		"data": make(map[string]string),
	})

}
