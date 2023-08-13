package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shicli/gin-first/common"
	"github.com/shicli/gin-first/response"
)

type Article struct {
	ID            uint32 `gorm:"primary_key" json:"id"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	CreatedOn     uint32 `json:"created_on"`
	ModifiedOn    uint32 `json:"modified_on"`
	DeletedOn     uint32 `json:"deleted_on"`
	IsDel         uint8  `json:"is_del"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
}

func Register(ctx *gin.Context) {
	//DB := common.GetDB()
	//var requestUser = model.User{}
	//ctx.Bind(&requestUser)
	////获取参数
	//name := requestUser.Name
	//telephone := requestUser.Telephone
	//password := requestUser.Password

	//newUser := model.User{
	//	Name:      name,
	//	Telephone: telephone,
	//	Password:  string(hasePassword),
	//}
	//DB.Create(&newUser)

	token, _ := common.GenerateToken(1, "张三")
	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

// @Summary	获取单个文章
// @Produce	json
// @Param		id	path		int		true	"文章ID"
// @Success	200	{object}	Article	"成功"
// @Failure	400	{object}	string	"请求错误"
// @Failure	500	{object}	string	"内部错误"
// @Router		/api/v1/articles/{id} [get]
func (a Article) Get(c *gin.Context) {
	return
}

// @Summary	获取多个文章
// @Produce	json
// @Param		name		query		string	false	"文章名称"
// @Param		tag_id		query		int		false	"标签ID"
// @Param		state		query		int		false	"状态"
// @Param		page		query		int		false	"页码"
// @Param		page_size	query		int		false	"每页数量"
// @Success	200			{object}	Article	"成功"
// @Failure	400			{object}	string	"请求错误"
// @Failure	500			{object}	string	"内部错误"
// @Router		/api/v1/articles [get]
func (a Article) List(c *gin.Context) {
	return
}

// @Summary	创建文章
// @Produce	json
// @Param		tag_id			body		string	true	"标签ID"
// @Param		title			body		string	true	"文章标题"
// @Param		desc			body		string	false	"文章简述"
// @Param		cover_image_url	body		string	true	"封面图片地址"
// @Param		content			body		string	true	"文章内容"
// @Param		created_by		body		int		true	"创建者"
// @Param		state			body		int		false	"状态"
// @Success	200				{object}	Article	"成功"
// @Failure	400				{object}	string	"请求错误"
// @Failure	500				{object}	string	"内部错误"
// @Router		/api/v1/articles [post]
func (a Article) Create(c *gin.Context) {

}

// @Summary	更新文章
// @Produce	json
// @Param		tag_id			body		string	false	"标签ID"
// @Param		title			body		string	false	"文章标题"
// @Param		desc			body		string	false	"文章简述"
// @Param		cover_image_url	body		string	false	"封面图片地址"
// @Param		content			body		string	false	"文章内容"
// @Param		modified_by		body		string	true	"修改者"
// @Success	200				{object}	Article	"成功"
// @Failure	400				{object}	string	"请求错误"
// @Failure	500				{object}	string	"内部错误"
// @Router		/api/v1/articles/{id} [put]
func (a Article) Update(c *gin.Context) {
	return
}

// @Summary	删除文章
// @Produce	json
// @Param		id	path		int		true	"文章ID"
// @Success	200	{string}	string	"成功"
// @Failure	400	{object}	string	"请求错误"
// @Failure	500	{object}	string	"内部错误"
// @Router		/api/v1/articles/{id} [delete]
func (a Article) Delete(c *gin.Context) {
	return
}
