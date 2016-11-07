package http

import (
	"strings"

	"github.com/asyou-me/lib.v1/utils"
	"github.com/labstack/echo"

	"github.com/asyou-me/http/errors"
	"github.com/asyou-me/http/handlers"
	"github.com/asyou-me/http/types"
)

// POST 创建
func POST(c echo.Context) error {
	reqType := getReqType(c)
	rest := handlers.RestFulls[reqType]
	req, err := recvREQ(c, rest.POST)
	if err != nil {
		er := errors.New(
			errors.NULL, err.Error(),
		)
		er.LogMessage(err.Error())
		return sendErr(c, er)
	}

	if reqType != req.Type {
		return sendErr(c, errors.New(errors.NULL, "报文 type 与 url type 不一致"))
	}
	resp := &types.RESP{}
	err = handlers.POST(req, resp)
	if err != nil {
		return sendErr(c, err)
	}
	return sendData(c, 200, resp)
}

// LIST 获取列表
func LIST(c echo.Context) error {
	req := &types.REQ{
		Query: &types.Query{},
	}
	req.Type = getReqType(c)
	req.Query.Limit = utils.StringToInt(c.Param("limit"), 10)
	resp := &types.RESP{}
	err := handlers.LIST(req, resp)
	if err != nil {
		return sendErr(c, err)
	}
	return sendData(c, 200, resp)
}

// GET 获取
func GET(c echo.Context) error {
	req := &types.REQ{
		Query: &types.Query{},
	}
	req.Query.ID = c.Param("id")
	if req.Query.ID == "" {
		return sendErr(c, errors.New(errors.PARAME_REQUIRE, "请传入正确的 id"))
	}
	include := c.QueryParam("include")
	if include != "" {
		includes := strings.Split(include, ",")
		req.Query.Include = &includes
	}
	req.Type = getReqType(c)
	resp := &types.RESP{}
	err := handlers.GET(req, resp)
	if err != nil {
		return sendErr(c, err)
	}
	return sendData(c, 200, resp)
}

// PATCH 更新
func PATCH(c echo.Context) error {
	reqType := getReqType(c)
	rest := handlers.RestFulls[reqType]
	req, err := recvREQ(c, rest.PATCH)
	if err != nil {
		er := errors.New(
			errors.NULL, err.Error(),
		)
		er.LogMessage(err.Error())
		return sendErr(c, er)
	}

	if reqType != req.Type {
		return sendErr(c, errors.New(errors.NULL, "报文 type 与 url type 不一致"))
	}
	req.Query = &types.Query{}
	req.Query.ID = c.Param("id")
	if req.Query.ID == "" {
		return sendErr(c, errors.New(errors.PARAME_REQUIRE, "请传入正确的 id"))
	}
	req.Type = getReqType(c)

	resp := &types.RESP{}
	err = handlers.PACTH(req, resp)
	if err != nil {
		return sendErr(c, err)
	}
	return sendData(c, 200, resp)
}

// DELETE 删除
func DELETE(c echo.Context) error {
	req := &types.REQ{
		Query: &types.Query{},
	}
	req.Query.ID = c.Param("id")
	if req.Query.ID == "" {
		return sendErr(c, errors.New(errors.PARAME_REQUIRE, "请传入正确的 id"))
	}
	req.Type = getReqType(c)
	resp := &types.RESP{}
	err := handlers.DELETE(req, resp)
	if err != nil {
		return sendErr(c, err)
	}
	return sendData(c, 200, resp)
}
