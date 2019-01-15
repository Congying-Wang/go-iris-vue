package routes

import (
	"go-iris/middleware/jwts"
	"go-iris/utils"
	"go-iris/web/models"
	"go-iris/web/supports"
	"go-iris/web/supports/vo"
	"time"

	"github.com/kataras/iris"
)

func Registe(ctx iris.Context) {
	user := new(models.User)
	ctx.ReadJSON(&user)

	user.CreateTime = time.Now()
	user.Password = utils.AESEncrypt([]byte(user.Password))

	effect, err := models.CreateUser(user)
	//err := u.DoRegiste(user)
	if effect <= 0 || err != nil {
		ctx.Application().Logger().Errorf("用户[%s]注册失败。%s", user.Username, err.Error())
		supports.Error(ctx, iris.StatusInternalServerError, supports.Registe_failur, nil)
	} else {
		supports.Ok_(ctx, supports.Registe_success)
	}
}

func Login(ctx iris.Context) {
	user := new(models.User)
	if err := ctx.ReadJSON(&user); err != nil {
		ctx.Application().Logger().Errorf("用户[%s]登录失败。%s", "", err.Error())
		supports.Error(ctx, iris.StatusInternalServerError, supports.Login_failur, nil)
		return
	}

	mUser := new(models.User)
	mUser.Username = user.Username
	has, err := models.GetUserByUsername(mUser)
	//has, err := u.DoLogin(mUser)
	//golog.Error(mUser)
	if err != nil {
		ctx.Application().Logger().Errorf("用户[%s]登录失败。%s", user.Username, err.Error())
		supports.Error(ctx, iris.StatusInternalServerError, supports.Login_failur, nil)
		return
	}

	if !has { // 用户名不正确
		supports.Unauthorized(ctx, supports.Username_failur, nil)
		return
	}

	ckPassword := utils.CheckPWD(user.Password, mUser.Password)
	if !ckPassword {
		supports.Unauthorized(ctx, supports.Password_failur, nil)
		return
	}

	token, err := jwts.GenerateToken(mUser);
	if err != nil {
		ctx.Application().Logger().Errorf("用户[%s]登录，生成token出错。%s", user.Username, err.Error())
		supports.Error(ctx, iris.StatusInternalServerError, supports.Token_create_failur, nil)
		return
	}

	supports.Ok(ctx, supports.Login_success, &vo.UserVO{
		Username:   mUser.Username,
		Appid:      mUser.Appid,
		Name:       mUser.Name,
		Phone:      mUser.Phone,
		Email:      mUser.Email,
		Userface:   mUser.Userface,
		CreateTime: mUser.CreateTime,
		Token:      token,
	})
	return
}

// 修改角色的权限


