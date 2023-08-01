package http

import (
	"github.com/harley9293/nebulus"
	"github.com/harley9293/nebulus/pkg/httpd"
)

type LoginReq struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type LoginRsp struct {
	Result int    `json:"result"`
	Token  string `json:"token"`
}

func HandleLoginReq(req *LoginReq, ctx *httpd.Context) LoginRsp {
	// do something
	return LoginRsp{}
}

type GetInfoReq struct {
}

type GetInfoRsp struct {
	Info string `json:"info"`
}

func HandleGetInfoReq(req *GetInfoReq, ctx *httpd.Context) GetInfoRsp {
	// do something
	return GetInfoRsp{}
}

func Example() {
	s := httpd.NewService(&httpd.Config{})
	if s == nil {
		panic("NewService() failed")
	}

	s.AddGlobalMiddleWare(httpd.LogMW, httpd.CookieMW, httpd.CorsMW)
	s.AddHandler("POST", "/login", HandleLoginReq)
	s.AddHandler("GET", "/info", HandleGetInfoReq, httpd.AuthMW)

	err := nebulus.Register("http", s, "127.0.0.1:36000")
	if err != nil {
		panic(err)
	}

	nebulus.Run()
}
