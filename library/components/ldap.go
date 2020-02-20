package components

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	ldap "gopkg.in/ldap.v3"
)

type Ldap struct {
	link *ldap.Conn
}

func new_ldap() (l Ldap) {
	l.connect()
	return l
}

func (l *Ldap) connect() (err bool) {
	ldapHost := beego.AppConfig.String("ldapHost")
	ldapPort, _ := beego.AppConfig.Int("ldapPort")
	link, e := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapHost, ldapPort))

	if e != nil {
		logs.Info("ldap connect error")
		return false
	}

	e = link.Bind(beego.AppConfig.String("ldapManagerDn"), beego.AppConfig.String("ldapManagerPassword"))
	if e != nil {
		logs.Info("ldap login error")
		return false
	}
	l.link = link
	return true
}
