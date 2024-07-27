package models

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {

	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
	return
}

func AddTag(name string, created_by string, state int) bool {
	var tag Tag = Tag{
		Name:      name,
		CreatedBy: created_by,
		State:     state,
	}

	db.Create(&tag)
	return true
}

func EditTag(id int, data map[string]interface{}) bool {
	db.Model(&Tag{}).Where("id=?", id).Updates(data)
	return true
}

func DeleteTag(id int) bool {

	db.Delete(&Tag{}, id)
	return true
}

// // 回调函数
// func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
// 	scope.SetColumn("CreatedOn", time.Now().Unix())
// 	return nil
// }

// // 回调函数
// func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
// 	scope.SetColumn("ModifiedOn", time.Now().Unix())
// 	return nil
// }

func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)
	return
}

func ExistTagByID(id int) bool {
	var tag Tag
	db.Select("id").Where("id=? AND deleted_on=?", id, 0).First(&tag)
	return tag.ID > 0
}

func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name=? AND deleted_on=?", name, 0).First(&tag)
	return tag.ID > 0
}

// 定时清理软删除的数据
func CleanAllTag() bool {
	db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{})

	return true
}
