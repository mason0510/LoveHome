package controllers

import (
	"github.com/astaxie/beego"
	"LoveHome/models"
	"encoding/json"
	"github.com/astaxie/beego/orm"
)

type SessionController struct {
	beego.Controller
}

func(this*SessionController)RetData(resp map[string]interface{}){
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this*SessionController)GetSessionData(){
	resp := make(map[string]interface{})
	defer this.RetData(resp)
	user := models.User{}
	//user.Name = "wyj"
	//
	resp["errno"] =models.RECODE_DBERR
	resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)

	name := this.GetSession("name")
	if name != nil {
		user.Name = name.(string)
		resp["errno"] =models.RECODE_OK
		resp["errmsg"] = models.RecodeText(models.RECODE_OK)
		resp["data"] = user

	}

}

func (this*SessionController)DeleteSessionData(){
	resp := make(map[string]interface{})
	defer this.RetData(resp)
	this.DelSession("name")

	resp["errno"] =models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)


}

func(this*SessionController) login()  {
	//得到用户信息
	resp:=make(map[string]interface{})
	//获取前端数据
	json.Unmarshal(this.Ctx.Input.RequestBody,&resp)

	//合法
	if resp["mobile"]==nil||resp["password"]==nil {
	resp["errno"]=models.RECODE_DATAERR
	resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
		return
	}


	//取数据库对比

	o:=orm.NewOrm()
	user:=models.User{Name:resp["mobile"].(string)}
	qs := o.QueryTable("user")
	err := qs.Filter("mobile", "7777").One(&user)
	if err!=nil {
		resp["errno"]=models.RECODE_DATAERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
		return
		}
	if user.Password_hash!=resp["password"]{
		resp["err"]=models.RECODE_DATAERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
		beego.Info("name=",resp["mobile"],"password=",resp["password"])
		return
		}

	//添加session
this.SetSession("name",resp["name"])
this.SetSession("mobile",resp["mobile"])
this.SetSession("user_id",user.Id)


	//返回json给前端
resp["errno"]=models.RECODE_OK
resp["errmsg"]=models.RecodeText(models.RECODE_OK)


	}

