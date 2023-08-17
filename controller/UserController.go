package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shicli/gin-first/common"
	"github.com/shicli/gin-first/dto"
	"github.com/shicli/gin-first/model"
	"github.com/shicli/gin-first/response"
	"github.com/shicli/gin-first/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Article struct {
	ID        uint32 `gorm:"primary_key" json:"id"`
	Name      string `json:"created_by"`
	Telemetry string `json:"modified_by"`
}

// @Summary 注册信息
// @Param name formData string false "名字"
// @Param password formData string false "密码"
// @Param telemetry formData string false "电话"
// @Success 200 {object} Article "成功"
// @Failure 400 {object} string "请求错误"
// @Failure 500 {object} string "内部错误"
// @Router /api/auth/register [POST]
func Register(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = model.User{}
	ctx.ShouldBind(&requestUser)

	//获取参数
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password

	//数据验证
	fmt.Println(telephone, "手机号码长度", len(telephone))
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	if len(name) == 0 {
		name = util.RandString(10)
	}

	//判断手机号码是否存在
	if isTelephoneExists(DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}

	//创建用户
	hasePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "加密失败")
		return
	}

	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasePassword),
	}
	DB.Create(&newUser)

	token, _ := common.GenerateToken(newUser)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": dto.ToUserDto(user.(model.User))},
	})
}

func isTelephoneExists(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	return user.ID != 0
}
