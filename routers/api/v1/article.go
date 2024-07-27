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

// 获取文章列表
// @Summary      Get Articles
// @Produce      json
// @Param        state	query	int	false  "State"
// @Param        tag_id	query	int	false  "TagId"
// @Success 	 200 {string} json "{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /api/v1/articles [get]
func GetArticles(c *gin.Context) {

	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("文章状态state 只能为0或1")
		maps["state"] = state
	}
	var tag_id int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tag_id = com.StrTo(arg).MustInt()
		valid.Min(tag_id, 1, "tag_id").Message("文章tag_id 最小为1")
		maps["tag_id"] = tag_id

	}
	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		data["total"] = models.GetArticlesTotal(maps)
		data["list"] = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
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

// 获取指定文章
// @Summary      Get Single Article By Id
// @Produce      json
// @Param        id						query	int		true	"ID"
// @Success 	 200 {string} 			json 	"{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /api/v1/articles/{id} 	[get]
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("文章Id 最小为1")

	var data interface{}
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			code = e.SUCCESS
			maps["id"] = id
			data = models.GetArticle(maps)

		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
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

// 新建文章
// @Summary      Add Article
// @Produce      json
// @Param        title			query	string	true  	"Title"
// @Param        desc			query	string	true  	"Desc"
// @Param        content		query	string	true  	"Content"
// @Param        created_by		query	string	true  	"CreatedBy"
// @Param        tag_id			query	int		false  	"TagId"
// @Success 	 200 {string} 	json 	"{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	valid := validation.Validation{}

	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	created_by := c.Query("created_by")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()

	var tag_id int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tag_id = com.StrTo(arg).MustInt()

		valid.Min(tag_id, 1, "tag_id").Message("文章tag_id 最小为1")
	}

	valid.Required(title, "title").Message("文章标题不为空")
	valid.MaxSize(title, 100, "title").Message("文章标题最大长度为100字符")
	valid.Required(desc, "desc").Message("文章简介不为空")
	valid.MaxSize(desc, 255, "desc").Message("文章简介最大长度为255字符")
	valid.Required(content, "content").Message("文章内容不为空")
	valid.MaxSize(content, 65535, "content").Message("文章内容最大长度为65535字符")
	valid.Required(created_by, "created_by").Message("文章创建人不为空")
	valid.Range(state, 0, 1, "state").Message("文章状态state 只能为0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {

		if models.ExistTagByID(tag_id) {
			code = e.SUCCESS
			maps := make(map[string]interface{})

			maps["tag_id"] = tag_id
			maps["title"] = title
			maps["desc"] = desc
			maps["content"] = content
			maps["created_by"] = created_by
			maps["state"] = state

			models.AddArticle(maps)

		} else {
			code = e.ERROR_NOT_EXIST_TAG
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
		"data": make(map[string]interface{}),
	})
}

// 更新指定文章
// @Summary      Update Article By Id
// @Produce      json
// @Param        id				query		int		true  	"ID"
// @Param        title			query		string	true 	"Title"
// @Param        desc			query		string	true  	"Desc"
// @Param        content		query		string	true  	"Content"
// @Param        modified_by	query		string	true  	"ModifiedBy"
// @Param        tag_id			query		int		false  	"TagId"
// @Param        state			query		int		false  	"State"
// @Success 	 200 			{string}	json 	"{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modified_by := c.Query("modified_by")

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("文章状态state 只能为0或1")
	}
	var tag_id int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tag_id = com.StrTo(arg).MustInt()

		valid.Min(tag_id, 1, "tag_id").Message("文章tag_id 最小为1")
	}
	valid.Min(id, 1, "id").Message("文章id 最小为1")
	valid.Required(title, "title").Message("文章标题不为空")
	valid.MaxSize(title, 100, "title").Message("文章标题最大长度为100字符")
	valid.Required(desc, "desc").Message("文章简介不为空")
	valid.MaxSize(desc, 255, "desc").Message("文章简介最大长度为255字符")
	valid.Required(content, "content").Message("文章内容不为空")
	valid.MaxSize(content, 65535, "content").Message("文章内容最大长度为65535字符")
	valid.Required(modified_by, "modified_by").Message("文章更改人不为空")
	valid.MaxSize(modified_by, 100, "modified_by").Message("文章更改人最大长度为100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(tag_id) {
			if models.ExistArticleByID(id) {
				code = e.SUCCESS
				maps := make(map[string]interface{})

				maps["title"] = title
				maps["desc"] = desc
				maps["content"] = content
				maps["modified_by"] = modified_by
				maps["state"] = state
				maps["tag_id"] = tag_id
				maps["deleted_on"] = 0

				models.EditArticle(id, maps)

			} else {
				code = e.ERROR_NOT_EXIST_ARTICLE
			}
		} else {
			code = e.ERROR_NOT_EXIST_TAG
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
		"data": make(map[string]interface{}),
	})

}

// 删除指定文章
// @Summary      Delete Article By Id
// @Produce      json
// @Param        id		query		int		true  "ID"
// @Success 	 200 	{string} 	json 	"{"code": 200, "data":{}, "msg": "ok"}"
// @Router       /api/v1/articles/{id} 	[delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("文章Id 最小为1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			code = e.SUCCESS
			models.DeleteArticle(id)
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
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
		"data": make(map[string]interface{}),
	})
}
