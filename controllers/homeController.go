package controllers

import (
	"encoding/json"
	"gin_admin_demo/comment"
	"gin_admin_demo/comment/menu"
	"gin_admin_demo/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HomeController struct {
	BaseController
}

func (con *HomeController) Home() gin.HandlerFunc {
	menuList := menu.GetMenu()
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userInfoJson := session.Get("userInfo")
		userData := make(map[string]interface{})
		err := json.Unmarshal([]byte(userInfoJson.(string)), &userData)
		privs := make(map[string]struct{})
		json.Unmarshal([]byte(userData["privs"].(string)), &privs)
		if err != nil {
			// 取不到就是没有登录
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(200, `<script type="text/javascript">top.location.href="/admin/login"</script>`)
			return
		}
		c.HTML(http.StatusOK, "views/template/home/home.html", gin.H{
			"menuList":  menuList,
			"userInfo":  userData,
			"userPrivs": privs,
		})
	}
}

func (con *HomeController) Welcome() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "views/template/home/welcome.html", gin.H{})
	}
}

func (con *HomeController) EditPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		c.HTML(http.StatusOK, "views/template/home/password_form.html", gin.H{
			"id": id,
		})
	}
}

func (con *HomeController) SavePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.PostForm("id")
		oldPassword := c.PostForm("old_password")
		newPassword := c.PostForm("new_password")
		subPassword := c.PostForm("sub_password")

		if newPassword != subPassword {
			con.Error(c, "请确认新密码")
			return
		}

		adminUser, _ := models.GetAdminUserById(id)
		oldPass := comment.Encryption(oldPassword, adminUser.Salt)
		if oldPass != adminUser.Password {
			con.Error(c, "原密码不正确")
			return
		}

		newPass := comment.Encryption(newPassword, adminUser.Salt)
		err := models.AlterAdminUserPass(id, newPass)
		if err != nil {
			con.Error(c, "修改密码失败")
			return
		} else {
			con.Success(c, "", "修改成功")
			return
		}
	}
}
