package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"encoding/json"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
	"kuaidian-app/library/common"
	"kuaidian-app/library/ldap"
	"kuaidian-app/models"
	"strconv"
	"strings"
	"time"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Post() {
	//哈希校验成功后 更新 auth_key
	beego.Info(string(c.Ctx.Input.RequestBody))
	postData := map[string]string{"user_password": "", "user_name": ""}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &postData)
	if err != nil {
		c.SetJson(1, nil, "数据格式错误")
		return
	}
	password := postData["user_password"]
	userName := postData["user_name"]
	if userName == "" || password == "" {
		c.SetJson(1, nil, "用户名或密码不存在")
		return
	}
	var user models.User
	o := orm.NewOrm()
	err = o.Raw("SELECT * FROM `user` WHERE username= ?", userName).QueryRow(&user)
	beego.Info(user)

	enableLdap, _ := beego.AppConfig.Bool("enableLdap")
	if enableLdap == true {
		msg, user, isOk := c.ldapLogin(userName, password, user)
		if !isOk {
			c.SetJson(1, nil, msg)
			return
		} else {
			userAuth := common.Md5String(user.Username + common.GetString(time.Now().Unix()))
			user.AuthKey = userAuth
			models.UpdateUserById(&user)
			c.SetJson(0, user, "")
			return

		}
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			c.SetJson(1, nil, "用户名或密码错误")
			return
		} else {
			if user.AuthKey == "" {
				userAuth := common.Md5String(user.Username + common.GetString(time.Now().Unix()))
				user.AuthKey = userAuth
				models.UpdateUserById(&user)
			}
			user.PasswordHash = ""
			c.SetJson(0, user, "")
			return
		}
	}
}

func (c *LoginController) ldapLogin(userName string, password string, coinlab_user models.User) (msg string, user models.User, isOk bool) {
	ldap := ldap.Ldap{}
	e := ldap.Connect()
	if e != nil {
		return "ldap连接失败", coinlab_user, false
	}
	//验证用户身份
	ldap_user, e := ldap.AuthByUidAndPassword(userName, password)
	if e != nil {
		return "ldap身份认证失败", coinlab_user, false
	}
	//验证是否在coinlab用户组
	ldapGroupFilter := beego.AppConfig.String("ldapGroupFilter")
	ldapGroupFilter = strings.Replace(ldapGroupFilter, "{UidNumber}", ldap_user.UidNumber, -1)
	ldapGroupFilter = strings.Replace(ldapGroupFilter, "{uid}", ldap_user.Uid, -1)
	ldapGroupFilter = strings.Replace(ldapGroupFilter, "{cn}", ldap_user.Cn, -1)
	ldapGroupFilter = strings.Replace(ldapGroupFilter, "{sn}", ldap_user.Sn, -1)

	groupCn, e := ldap.SearchGroupCn(ldapGroupFilter)
	if e != nil {
		beego.Info("ldap组身份验证失败")
		return "ldap组身份验证失败", coinlab_user, false
	} else {
		o := orm.NewOrm()

		role_id64, _ := strconv.ParseInt(beego.AppConfig.String("ldapGroupName2roleid_"+groupCn), 10, 64)
		role_id := int16(role_id64)
		// 用户不存在，自动同步进coinlab数据库
		if coinlab_user.Username == "" {
			c.AddUserFromLdap2Coinlab(ldap_user, role_id)
			_ = o.Raw("SELECT * FROM `user` WHERE username= ?", user).QueryRow(&coinlab_user)
			logs.Info(coinlab_user)
		} else {
			//role变更
			if role_id != coinlab_user.Role {
				coinlab_user.Role = role_id
				models.UpdateUserById(&coinlab_user)
			}
		}

		return "", coinlab_user, true
	}

}

func (c *LoginController) AddUserFromLdap2Coinlab(user ldap.Ldap_user, role_id int16) {
	uidNumber, _ := strconv.Atoi(user.UidNumber)

	userModel := models.User{}
	userModel.Id = uidNumber
	userModel.Username = user.Uid
	userModel.Email = user.Email
	userModel.Realname = user.Cn
	userModel.CreatedAt = time.Now()
	userModel.UpdatedAt = time.Now()
	userModel.Avatar = "default.jpg"
	userModel.Role = role_id
	userModel.FromLdap = 1
	uid, _ := models.AddUser(&userModel)
	beego.Info(uid)
}
