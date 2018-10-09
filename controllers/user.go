package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"loveHome/models"
	"path"
	"github.com/weilaihui/fdfs_client"
)

//type
type UserController struct {
	beego.Controller
}
//restdata 接受前端传过来的数据
func (this*UserController)RetData(resp map[string] interface{})  {
	this.Data["json"]=resp
	this.ServeJSON(true)
}
func (this*UserController)Reg()  {
	resp:=make(map[string]interface{})
	//handle json
	json.Unmarshal(this.Ctx.Input.RequestBody,&resp)

	this.RetData(resp)
	beego.Info(`resp["mobile"]=`,resp["mobile"])
	beego.Info(`resp["password"]=`,resp["password"])
	beego.Info(`resp["sms_code"]=`,resp["sms_code"])
	//插入数据库
	o:=orm.NewOrm()
	user:=models.User{}
	user.Password_hash=resp["password"].(string)
	user.Name=resp["name"].(string)
	user.Mobile=resp["mobile"].(string)
	id,err:=o.Insert(&user)
	if err!=nil {
		resp["errno"]=4001
		resp["errmsg"]="注册失败"
		return
	}
	beego.Info("success id=",id)
	resp["errno"]=4001
	resp["errmsg"]="注册成功"

}

//上传头像
func (this*UserController)Postavatar () {
	//获取流数据
	resp := make(map[string]interface{})
	defer this.RetData(resp)

	//获取文件
	fileData,hd,err := this.GetFile("avatar")
	if err != nil {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		beego.Info("===========11111")
		return
	}

	//获取文件后缀
	suffix := path.Ext(hd.Filename) //a.jpg.avi
	//存储到fastdfs上
	fdfsClient,err := fdfs_client.NewFdfsClient("conf/client.conf")
	if err != nil{
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		beego.Info("===========22222222")

		return
	}
	fileBuffer := make([]byte, hd.Size)
	_, err = fileData.Read(fileBuffer)
	if err != nil {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		beego.Info("===========33333")

		return
	}

	DataResponse, err := fdfsClient.UploadByBuffer(fileBuffer, suffix[1:])//aa.jpg

	if err != nil {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		beego.Info("===========44444")

		return
	}
	//获取user_id
	user_id := this.GetSession("user_id")
	var user models.User
	//获取数据库信息和url
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	qs.Filter("Id",user_id).One(&user)
	user.Avatar_url = DataResponse.RemoteFileId

	//更新数据库
	_,errUpdate := o.Update(&user)
	if errUpdate != nil{
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}
	urlMap:= make(map[string]string)
	//Avaurl := "192.168.152.138:8899"+DataResponse.RemoteFileId
	urlMap["avatar_url"] = "http://192.168.152.138:8899/"+DataResponse.RemoteFileId
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	resp["data"] = urlMap




}
