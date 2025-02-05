// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/middleware"
	"github.com/zhenorzz/goploy/model"
	"github.com/zhenorzz/goploy/response"
	"net/http"
)

type Agent Controller

func (a Agent) Routes() []core.Route {
	return []core.Route{
		core.NewWhiteRoute("/agent/report", http.MethodPost, a.Report).Middleware(middleware.CheckSign),
		core.NewWhiteRoute("/agent/getServerID", http.MethodPost, a.GetServerID).Middleware(middleware.CheckSign),
		core.NewWhiteRoute("/agent/getCronList", http.MethodPost, a.GetCronList).Middleware(middleware.CheckSign),
		core.NewWhiteRoute("/agent/getCronLogs", http.MethodPost, a.GetCronLogs).Middleware(middleware.CheckSign),
		core.NewWhiteRoute("/agent/cronReport", http.MethodPost, a.CronReport).Middleware(middleware.CheckSign),
	}
}

func (Agent) GetServerID(gp *core.Goploy) core.Response {
	type ReqData struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	server, err := model.Server{
		Name: reqData.Name,
		IP:   reqData.IP,
	}.GetData()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			ID int64 `json:"id"`
		}{ID: server.ID},
	}
}

func (Agent) GetCronList(gp *core.Goploy) core.Response {
	type ReqData struct {
		ServerID int64 `json:"serverId" validate:"gt=0"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	crons, err := model.Cron{ServerID: reqData.ServerID}.GetList()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model.Crons `json:"list"`
		}{List: crons},
	}
}

func (Agent) GetCronLogs(gp *core.Goploy) core.Response {
	type ReqData struct {
		ServerID int64  `json:"serverId" validate:"gt=0"`
		CronID   int64  `json:"cronId" validate:"gt=0"`
		Page     uint64 `json:"page" validate:"gt=0"`
		Rows     uint64 `json:"rows" validate:"gt=0"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	crons, err := model.CronLog{ServerID: reqData.ServerID, CronID: reqData.CronID}.GetList(reqData.Page, reqData.Rows)

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{
		Data: struct {
			List model.CronLogs `json:"list"`
		}{List: crons},
	}
}

func (Agent) CronReport(gp *core.Goploy) core.Response {
	type ReqData struct {
		ServerId   int64  `json:"serverId" validate:"gt=0"`
		CronId     int64  `json:"cronId" validate:"gt=0"`
		ExecCode   int    `json:"execCode"`
		Message    string `json:"message" validate:"required"`
		ReportTime string `json:"reportTime" validate:"required"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err := model.CronLog{
		ServerID:   reqData.ServerId,
		CronID:     reqData.CronId,
		ExecCode:   reqData.ExecCode,
		Message:    reqData.Message,
		ReportTime: reqData.ReportTime,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}

func (Agent) Report(gp *core.Goploy) core.Response {
	type ReqData struct {
		ServerId   int64  `json:"serverId" validate:"gt=0"`
		Type       int    `json:"type" validate:"gt=0"`
		Item       string `json:"item" validate:"required"`
		Value      string `json:"value" validate:"required"`
		ReportTime string `json:"reportTime" validate:"required"`
	}

	var reqData ReqData
	if err := decodeJson(gp.Body, &reqData); err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}

	err := model.ServerAgentLog{
		ServerID:   reqData.ServerId,
		Type:       reqData.Type,
		Item:       reqData.Item,
		Value:      reqData.Value,
		ReportTime: reqData.ReportTime,
	}.AddRow()

	if err != nil {
		return response.JSON{Code: response.Error, Message: err.Error()}
	}
	return response.JSON{}
}
