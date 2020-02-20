package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["kuaidian-app/controllers/api:ProjectController"] = append(beego.GlobalControllerRouter["kuaidian-app/controllers/api:ProjectController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["kuaidian-app/controllers/api:ProjectController"] = append(beego.GlobalControllerRouter["kuaidian-app/controllers/api:ProjectController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["kuaidian-app/controllers/api:ProjectController"] = append(beego.GlobalControllerRouter["kuaidian-app/controllers/api:ProjectController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["kuaidian-app/controllers/api:ProjectController"] = append(beego.GlobalControllerRouter["kuaidian-app/controllers/api:ProjectController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["kuaidian-app/controllers/api:ProjectController"] = append(beego.GlobalControllerRouter["kuaidian-app/controllers/api:ProjectController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["kuaidian-app/controllers/api:TaskController"] = append(beego.GlobalControllerRouter["kuaidian-app/controllers/api:TaskController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["kuaidian-app/controllers/api:TaskController"] = append(beego.GlobalControllerRouter["kuaidian-app/controllers/api:TaskController"],
        beego.ControllerComments{
            Method: "GetAllAndProName",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["kuaidian-app/controllers/api:TaskController"] = append(beego.GlobalControllerRouter["kuaidian-app/controllers/api:TaskController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["kuaidian-app/controllers/api:TaskController"] = append(beego.GlobalControllerRouter["kuaidian-app/controllers/api:TaskController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["kuaidian-app/controllers/api:TaskController"] = append(beego.GlobalControllerRouter["kuaidian-app/controllers/api:TaskController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["kuaidian-app/controllers/api:TaskController"] = append(beego.GlobalControllerRouter["kuaidian-app/controllers/api:TaskController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["kuaidian-app/controllers/api:TokenController"] = append(beego.GlobalControllerRouter["kuaidian-app/controllers/api:TokenController"],
        beego.ControllerComments{
            Method: "IssueToken",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
