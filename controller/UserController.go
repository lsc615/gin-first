package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shicli/gin-first/common"
	"github.com/shicli/gin-first/dto"
	"github.com/shicli/gin-first/model"
	"github.com/shicli/gin-first/response"
	"github.com/shicli/gin-first/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
)

type Article struct {
	ID        uint32 `gorm:"primary_key" json:"id"`
	Name      string `json:"created_by"`
	Telemetry string `json:"modified_by"`
}

// @Summary 注册信息
// @Param Name controller. string false "名字"
// @Param Password query string false "密码"
// @Param Telemetry query string false "电话"
// @Success 200 {object} Article "成功"
// @Failure 400 {object} string "请求错误"
// @Failure 500 {object} string "内部错误"
// @Router /api/auth/register [POST]
func Register(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = model.User{}
	ctx.Bind(&requestUser)

	//获取参数
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password

	//数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	//如果没有传入名称，给一个10位的随机数
	if len(name) == 0 {
		name = util.RandString(10)
	}

	//判断手机号码是否存在
	if isTelephoneExists(DB, name) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		log.Printf("用户已经存在")
		return
	}

	//创建用户，加密
	hasePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密失败")
		return
	}

	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasePassword),
	}
	DB.Create(&newUser)

	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}

	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

func Login(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = model.User{}
	ctx.ShouldBind(&requestUser)
	//获取参数
	//name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password
	//数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	//判断手机号码是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号不存在")
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}

	//返回结果
	response.Success(ctx, gin.H{"data": token}, "登陆成功")

	client := &http.Client{}

	// 构造请求到 /api/v1/users
	req, err := http.NewRequest("GET", "http://127.0.0.1:8081/version", nil)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	// 将响应返回给原始请求的客户端
	ctx.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)

}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		//"data": gin.H{"user": user},
		//重新定义返回的信息,user包含敏感信息
		"data": gin.H{"user": dto.ToUserDto(user.(model.User))},
	})
}

func isTelephoneExists(db *gorm.DB, name string) bool {
	var user model.User
	db.Where("name = ?", name).First(&user)
	return user.ID != 0
}
