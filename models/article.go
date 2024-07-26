package models

type Article struct {
	Model
	Tag Tag `json:"tag"`

	TagId      int    `json:"tag_id" gorm:"index"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// 获取文章列表
func GetArticles(pageNum int, pageSize int, maps map[string]interface{}) []Article {
	var articles []Article
	// 查询articles时预加载tags

	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)
	return articles
}

// 获取指定文章
func GetArticle(maps map[string]interface{}) (article Article) {
	db.Where(maps).First(&article)
	db.Model(&article).Related(&article.Tag)
	return
}

// 新建文章
func AddArticle(maps map[string]interface{}) bool {
	db.Create(&Article{

		TagId:     maps["tag_id"].(int),
		Title:     maps["title"].(string),
		Desc:      maps["desc"].(string),
		Content:   maps["content"].(string),
		CreatedBy: maps["created_by"].(string),
		State:     maps["state"].(int),
	})
	return true
}

// 更新指定文章
func EditArticle(id int, data map[string]interface{}) bool {
	db.Model(&Article{}).Where("id=?", id).Updates(data)
	return true
}

// 删除指定文章
func DeleteArticle(id int) bool {
	db.Delete(&Article{}, id)
	return true
}

// // 回调函数
// func (article *Article) BeforeCreate(scope *gorm.Scope) {
// 	scope.SetColumn("created_on", time.Now().Unix())
// }

// // 回调函数
// func (article *Article) BeforeUpdate(scope *gorm.Scope) {
// 	scope.SetColumn("modified_on", time.Now().Unix())
// }

func GetArticlesTotal(maps map[string]interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)
	return
}

func ExistArticleByID(id int) bool {

	var article Article
	db.Select("id").Where("id=? AND deleted_on=?", id, 0).First(&article)
	return article.ID > 0
}

func CleanAllArticle() bool {
	db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Article{})

	return true
}
