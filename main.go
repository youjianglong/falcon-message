package main

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/youjianglong/falcon-message/config"
	"github.com/youjianglong/falcon-message/sender"
	"github.com/youjianglong/falcon-message/util"
	"github.com/tylerb/graceful"
)

// OK
// P3
// Endpoint:SMARTMATRIX_MONITOR
// Metric:api
// Tags:act=webapi.bindDevice,loc=GZ
// all(#2): 1==0
// Note:act=webapi.bindDevice,loc=GZ
// Max:3, Current:1
// Timestamp:2017-06-02 08:02:00
// http://127.0.0.1:8081/portal/template/view/37

const (
	// IMDingPrefix 钉钉 前缀
	IMDingPrefix = "[ding]:"
	// IMWexinPrefix 微信前缀
	IMWexinPrefix = "[wexin]:"
)

var (
	cfg  config.Config
	ding *sender.DingTalk
	wx   *sender.Weixin
	qs *sender.QCloudSMS
)

func main() {
	cfg = config.Read()

	if cfg.DingTalk.Enable {
		ding = sender.NewDingTalk()
	}

	if cfg.Weixin.Enable {
		wx = sender.NewWeixin(cfg.Weixin.CorpID, cfg.Weixin.Secret)
		go wx.GetAccessToken()
	}
	if cfg.QCloudSms.Enable {
		qs = sender.NewQCloudSMS(cfg.QCloudSms)
	}
	engine := echo.New()
	engine.Server.Addr = cfg.Addr
	server := &graceful.Server{Timeout: time.Second * 10, Server: engine.Server, Logger: graceful.DefaultLogger()}
	engine.Use(middleware.Recover())
	engine.Use(middleware.Logger())
	api := engine.Group("/api/v1")
	api.GET("/wechat/auth", wxAuth)
	api.POST("/message", func(c echo.Context) error {
		log.Println("message comming")
		tos := c.FormValue("tos")
		content := c.FormValue("content")
		log.Println("tos:", tos, " content:", content)
		if content == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "content is requied")
		}

		content = util.HandleContent(content)
		if strings.HasPrefix(tos, IMDingPrefix) { //是钉钉
			if cfg.DingTalk.Enable {
				token := tos[len(IMDingPrefix):]
				ding.Send(token, content)
			}
		} else if strings.HasPrefix(tos, IMWexinPrefix) { //微信
			if cfg.Weixin.Enable {
				tos := tos[len(IMWexinPrefix):]
				if err := wx.Send(tos, content); err != nil {
					return echo.NewHTTPError(500, err.Error())
				}
			}
		}

		return nil
	})

	api.POST("/sms", func(c echo.Context) error {
		tos := c.FormValue("tos")
		content := c.FormValue("content")
		log.Println("tos:", tos, " content:", content)
		if content == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "content is requied")
		}
		content = util.HandleContent(content)
		params:=[]string{
			time.Now().Format("于15:04:05"),
			content,
		}
		if cfg.QCloudSms.Enable {
			if err := qs.SendWithTpl(tos,params); err != nil {
				log.Println("[sms]",err)
				return echo.NewHTTPError(500, err.Error())
			}
		}

		return nil
	})

	log.Println("listening on ", cfg.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// WxAuth 开启回调模式验证
func wxAuth(context echo.Context) error {
	if cfg.Weixin.Enable {
		echostr := context.FormValue("echostr")
		if echostr == "" {
			return errors.New("无法获取请求参数, echostr 为空")
		}
		var buf []byte
		var err error
		if buf, err = wx.Auth(echostr); err != nil {
			return err
		}

		return context.JSONBlob(200, buf)
	}

	return context.String(200, "微信没有启用")
}
