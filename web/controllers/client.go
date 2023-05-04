package controllers

import (
	"ehang.io/nps/lib/common"
	"ehang.io/nps/lib/file"
	"ehang.io/nps/lib/rate"
	"ehang.io/nps/server"
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"strings"
)

type ClientController struct {
	BaseController
}

func (s *ClientController) List() {
	if s.Ctx.Request.Method == "GET" {
		s.Data["menu"] = "client"
		s.SetInfo("client")
		s.display("client/list")
		return
	}
	start, length := s.GetAjaxParams()
	clientIdSession := s.GetSession("clientId")
	var clientId int
	if clientIdSession == nil {
		clientId = 0
	} else {
		clientId = clientIdSession.(int)
	}
	list, cnt := server.GetClientList(start, length, s.getEscapeString("search"), s.getEscapeString("sort"), s.getEscapeString("order"), clientId)
	cmd := make(map[string]interface{})
	ip := s.Ctx.Request.Host
	cmd["ip"] = common.GetIpByAddr(ip)
	cmd["bridgeType"] = beego.AppConfig.String("bridge_type")
	cmd["bridgePort"] = server.Bridge.TunnelPort
	s.AjaxTable(list, cnt, cnt, cmd)
}

// 添加客户端
func (s *ClientController) Add() {
	if s.Ctx.Request.Method == "GET" {
		s.Data["menu"] = "client"
		s.SetInfo("add client")
		s.display()
	} else {
		t := &file.Client{
			VerifyKey: s.getEscapeString("vkey"),
			Id:        int(file.GetDb().JsonDb.GetClientId()),
			Status:    true,
			Remark:    s.getEscapeString("remark"),
			Cnf: &file.Config{
				U:        s.getEscapeString("u"),
				P:        s.getEscapeString("p"),
				Compress: common.GetBoolByStr(s.getEscapeString("compress")),
				Crypt:    s.GetBoolNoErr("crypt"),
			},
			ConfigConnAllow: s.GetBoolNoErr("config_conn_allow"),
			RateLimit:       s.GetIntNoErr("rate_limit"),
			MaxConn:         s.GetIntNoErr("max_conn"),
			WebUserName:     s.getEscapeString("web_username"),
			WebPassword:     s.getEscapeString("web_password"),
			MaxTunnelNum:    s.GetIntNoErr("max_tunnel"),
			Flow: &file.Flow{
				ExportFlow: 0,
				InletFlow:  0,
				FlowLimit:  int64(s.GetIntNoErr("flow_limit")),
			},
		}
		if err := file.GetDb().NewClient(t); err != nil {
			s.AjaxErr(err.Error())
		}
		s.AjaxOk("add success")
	}
}

// 批量添加客户端
func (s *ClientController) Batch() {
	if s.Ctx.Request.Method == "GET" {
		s.Data["menu"] = "client"
		s.SetInfo("batch client")
		s.display()
	} else {
            var count = 0
            var startnumber = 0
	        if s.GetIntNoErr("count") == 0 {
            	 count = 1
            }else{
                count = s.GetIntNoErr("count")
            }

	        if s.GetIntNoErr("startnumber") == 0 {
            	 startnumber = 1
            }else{
                startnumber = s.GetIntNoErr("startnumber")
            }

            //s.AjaxErr(strconv.Itoa(startnumber))
        var port = s.GetIntNoErr("port")
        for i := 0;i < count; i++ {
            var VerifyKey = s.getEscapeString("vkey");
             if VerifyKey != "" {
                 VerifyKey = VerifyKey+strconv.Itoa(startnumber+i)
             }
            var clientId = int(file.GetDb().JsonDb.GetClientId())

            t := &file.Client{
                VerifyKey: VerifyKey,
                Id:        clientId,
                Status:    true,
                Remark:    s.getEscapeString("remark"),
                Cnf: &file.Config{
                    U:        s.getEscapeString("u"),
                    P:        s.getEscapeString("p"),
                    Compress: common.GetBoolByStr(s.getEscapeString("compress")),
                    Crypt:    s.GetBoolNoErr("crypt"),
                },
                ConfigConnAllow: s.GetBoolNoErr("config_conn_allow"),
                RateLimit:       s.GetIntNoErr("rate_limit"),
                MaxConn:         s.GetIntNoErr("max_conn"),
                WebUserName:     s.getEscapeString("web_username"),
                WebPassword:     s.getEscapeString("web_password"),
                MaxTunnelNum:    s.GetIntNoErr("max_tunnel"),
                Flow: &file.Flow{
                    ExportFlow: 0,
                    InletFlow:  0,
                    FlowLimit:  int64(s.GetIntNoErr("flow_limit")),
                },
            }
            if err := file.GetDb().NewClient(t); err != nil {
                s.AjaxErr(err.Error())
            }
            //填写了端口 新增端口
            if(port!=0){
                     t := &file.Tunnel{
            			Port:         port,
            			ServerIp:     s.getEscapeString("server_ip"),
            			Mode:         "socks5",
            			Target:       &file.Target{TargetStr: s.getEscapeString("target"), LocalProxy: s.GetBoolNoErr("local_proxy")},
            			Id:           int(file.GetDb().JsonDb.GetTaskId()),
            			Status:       true,
            			Remark:       s.getEscapeString("remark"),
            			Password:     s.getEscapeString("password"),
            			LocalPath:    s.getEscapeString("local_path"),
            			StripPre:     s.getEscapeString("strip_pre"),
            			Flow:         &file.Flow{},
            			S5User:       s.getEscapeString("S5User"),
            			CreateTime:   time.Now().Format(common.DEFAULT_TIME),
            			ExpireTime:   s.getEscapeString("expire_time"),
            			MultiAccount: &file.MultiAccount{AccountMap: authStrToMap(s.getEscapeString("S5User"))},
            		}
            		if t.Mode == "socks5" && t.S5User == "" {
            			s.AjaxErr("The account number cannot be empty")
            			return
            		}
//             		if !tool.TestServerPort(t.Port, t.Mode) {
//             			s.AjaxErr("The port cannot be opened because it may has been occupied or is no longer allowed.")
//             		}
            		var err error
            		if t.Client, err = file.GetDb().GetClient(clientId); err != nil {
            			s.AjaxErr(err.Error())
            		}
            		if t.Client.MaxTunnelNum != 0 && t.Client.GetTunnelNum() >= t.Client.MaxTunnelNum {
            			s.AjaxErr("The number of tunnels exceeds the limit")
            		}

            		if err := file.GetDb().NewTask(t); err != nil {
            			s.AjaxErr(err.Error())
            		}
            		if err := server.AddTask(t); err != nil {
            			s.AjaxErr(err.Error())
            		}
            		port = port+1;
             }
        }
		s.AjaxOk("batch add success")
	}
}

func (s *ClientController) GetClient() {
	if s.Ctx.Request.Method == "POST" {
		id := s.GetIntNoErr("id")
		data := make(map[string]interface{})
		if c, err := file.GetDb().GetClient(id); err != nil {
			data["code"] = 0
		} else {
			data["code"] = 1
			data["data"] = c
		}
		s.Data["json"] = data
		s.ServeJSON()
	}
}

// 修改客户端
func (s *ClientController) Edit() {
	id := s.GetIntNoErr("id")
	if s.Ctx.Request.Method == "GET" {
		s.Data["menu"] = "client"
		if c, err := file.GetDb().GetClient(id); err != nil {
			s.error()
		} else {
			s.Data["c"] = c
		}
		s.SetInfo("edit client")
		s.display()
	} else {
		if c, err := file.GetDb().GetClient(id); err != nil {
			s.error()
			s.AjaxErr("client ID not found")
			return
		} else {
			if s.getEscapeString("web_username") != "" {
				if s.getEscapeString("web_username") == beego.AppConfig.String("web_username") || !file.GetDb().VerifyUserName(s.getEscapeString("web_username"), c.Id) {
					s.AjaxErr("web login username duplicate, please reset")
					return
				}
			}
			if s.GetSession("isAdmin").(bool) {
				if !file.GetDb().VerifyVkey(s.getEscapeString("vkey"), c.Id) {
					s.AjaxErr("Vkey duplicate, please reset")
					return
				}
				c.VerifyKey = s.getEscapeString("vkey")
				c.Flow.FlowLimit = int64(s.GetIntNoErr("flow_limit"))
				c.RateLimit = s.GetIntNoErr("rate_limit")
				c.MaxConn = s.GetIntNoErr("max_conn")
				c.MaxTunnelNum = s.GetIntNoErr("max_tunnel")
			}
			c.Remark = s.getEscapeString("remark")
			c.Cnf.U = s.getEscapeString("u")
			c.Cnf.P = s.getEscapeString("p")
			c.Cnf.Compress = common.GetBoolByStr(s.getEscapeString("compress"))
			c.Cnf.Crypt = s.GetBoolNoErr("crypt")
			b, err := beego.AppConfig.Bool("allow_user_change_username")
			if s.GetSession("isAdmin").(bool) || (err == nil && b) {
				c.WebUserName = s.getEscapeString("web_username")
			}
			c.WebPassword = s.getEscapeString("web_password")
			c.ConfigConnAllow = s.GetBoolNoErr("config_conn_allow")
			if c.Rate != nil {
				c.Rate.Stop()
			}
			if c.RateLimit > 0 {
				c.Rate = rate.NewRate(int64(c.RateLimit * 1024))
				c.Rate.Start()
			} else {
				c.Rate = rate.NewRate(int64(2 << 23))
				c.Rate.Start()
			}
			file.GetDb().JsonDb.StoreClientsToJsonFile()
		}
		s.AjaxOk("save success")
	}
}

// 更改状态
func (s *ClientController) ChangeStatus() {
	id := s.GetIntNoErr("id")
	if client, err := file.GetDb().GetClient(id); err == nil {
		client.Status = s.GetBoolNoErr("status")
		if client.Status == false {
			server.DelClientConnect(client.Id)
		}
		s.AjaxOk("modified success")
	}
	s.AjaxErr("modified fail")
}

// 删除客户端
func (s *ClientController) Del() {
	id := s.GetIntNoErr("id")
	if err := file.GetDb().DelClient(id); err != nil {
		s.AjaxErr("delete error")
	}
	server.DelTunnelAndHostByClientId(id, false)
	server.DelClientConnect(id)
	s.AjaxOk("delete success")
}

// 批量删除客户端
func (s *ClientController) Batchdel() {
    ids := strings.Split(s.getEscapeString("ids"), ",")
    for i := 0; i < len(ids); i++ {
            id,errs := strconv.Atoi(ids[i])
            if errs != nil {
               s.AjaxErr("delete error")
            }
        	if err := file.GetDb().DelClient(id); err != nil {
        		s.AjaxErr("delete error")
        	}
        	server.DelTunnelAndHostByClientId(id, false)
        	server.DelClientConnect(id)
    }
	s.AjaxOk("delete success")
}
