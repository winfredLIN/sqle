// Package okr code generated by oapi sdk gen
/*
 * MIT License
 *
 * Copyright (c) 2022 Lark Technologies Pte. Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice, shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package larkokr

import (
	"context"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v3/core"
)

func NewService(config *larkcore.Config) *OkrService {
	o := &OkrService{config: config}
	o.Image = &image{service: o}
	o.MetricSource = &metricSource{service: o}
	o.MetricSourceTable = &metricSourceTable{service: o}
	o.MetricSourceTableItem = &metricSourceTableItem{service: o}
	o.Okr = &okr{service: o}
	o.Period = &period{service: o}
	o.ProgressRecord = &progressRecord{service: o}
	o.UserOkr = &userOkr{service: o}
	return o
}

type OkrService struct {
	config                *larkcore.Config
	Image                 *image                 // 图片
	MetricSource          *metricSource          // 指标库
	MetricSourceTable     *metricSourceTable     // 指标表
	MetricSourceTableItem *metricSourceTableItem // 指标项
	Okr                   *okr                   // OKR
	Period                *period                // OKR周期
	ProgressRecord        *progressRecord        // OKR进展记录
	UserOkr               *userOkr               // 用户OKR
}

type image struct {
	service *OkrService
}
type metricSource struct {
	service *OkrService
}
type metricSourceTable struct {
	service *OkrService
}
type metricSourceTableItem struct {
	service *OkrService
}
type okr struct {
	service *OkrService
}
type period struct {
	service *OkrService
}
type progressRecord struct {
	service *OkrService
}
type userOkr struct {
	service *OkrService
}

// 上传图片
//
// - 上传图片
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/okr-v1/image/upload
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/okrv1/upload_image.go
func (i *image) Upload(ctx context.Context, req *UploadImageReq, options ...larkcore.RequestOptionFunc) (*UploadImageResp, error) {
	options = append(options, larkcore.WithFileUpload())
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/okr/v1/images/upload"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser}
	apiResp, err := larkcore.Request(ctx, apiReq, i.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &UploadImageResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, i.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取指标库
//
// - 获取租户下全部 OKR 指标库（仅限 OKR 企业版使用）
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/okr-v1/metric_source/list
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/okrv1/list_metricSource.go
func (m *metricSource) List(ctx context.Context, req *ListMetricSourceReq, options ...larkcore.RequestOptionFunc) (*ListMetricSourceResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/okr/v1/metric_sources"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser}
	apiResp, err := larkcore.Request(ctx, apiReq, m.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListMetricSourceResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, m.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取指标表
//
// - 获取指定指标库下有哪些指标表（仅限 OKR 企业版使用）
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/okr-v1/metric_source-table/list
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/okrv1/list_metricSourceTable.go
func (m *metricSourceTable) List(ctx context.Context, req *ListMetricSourceTableReq, options ...larkcore.RequestOptionFunc) (*ListMetricSourceTableResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/okr/v1/metric_sources/:metric_source_id/tables"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser}
	apiResp, err := larkcore.Request(ctx, apiReq, m.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListMetricSourceTableResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, m.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 批量更新指标项
//
// - - 该接口用于批量更新多项指标，单次调用最多更新 100 条记录。接口仅限 OKR 企业版使用。;;  更新成功后 OKR 系统会给以下人员发送消息通知：;;	- 首次更新目标值的人员 ;;	- 已经将指标添加为 KR、且本次目标值/起始值/支撑的上级有变更的人员，不包含仅更新了进度值的人员
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/okr-v1/metric_source-table-item/batch_update
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/okrv1/batchUpdate_metricSourceTableItem.go
func (m *metricSourceTableItem) BatchUpdate(ctx context.Context, req *BatchUpdateMetricSourceTableItemReq, options ...larkcore.RequestOptionFunc) (*BatchUpdateMetricSourceTableItemResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/okr/v1/metric_sources/:metric_source_id/tables/:metric_table_id/items/batch_update"
	apiReq.HttpMethod = http.MethodPatch
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser}
	apiResp, err := larkcore.Request(ctx, apiReq, m.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &BatchUpdateMetricSourceTableItemResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, m.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取指标项详情
//
// - 获取某项指标的具体内容（仅限 OKR 企业版使用）
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/okr-v1/metric_source-table-item/get
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/okrv1/get_metricSourceTableItem.go
func (m *metricSourceTableItem) Get(ctx context.Context, req *GetMetricSourceTableItemReq, options ...larkcore.RequestOptionFunc) (*GetMetricSourceTableItemResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/okr/v1/metric_sources/:metric_source_id/tables/:metric_table_id/items/:metric_item_id"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser}
	apiResp, err := larkcore.Request(ctx, apiReq, m.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &GetMetricSourceTableItemResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, m.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取指标项
//
// - 获取指定指标表下的所有指标项（仅限 OKR 企业版使用）
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/okr-v1/metric_source-table-item/list
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/okrv1/list_metricSourceTableItem.go
func (m *metricSourceTableItem) List(ctx context.Context, req *ListMetricSourceTableItemReq, options ...larkcore.RequestOptionFunc) (*ListMetricSourceTableItemResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/okr/v1/metric_sources/:metric_source_id/tables/:metric_table_id/items"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser}
	apiResp, err := larkcore.Request(ctx, apiReq, m.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListMetricSourceTableItemResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, m.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 更新指标项
//
// - - 该接口用于更新某项指标，接口仅限 OKR 企业版使用。;;	更新成功后 OKR 系统会给以下人员发送消息通知：;;	- 首次更新目标值的人员 ;;	- 已经将指标添加为 KR、且本次目标值/起始值/支撑的上级有变更的人员，不包含仅更新了进度值的人员
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/okr-v1/metric_source-table-item/patch
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/okrv1/patch_metricSourceTableItem.go
func (m *metricSourceTableItem) Patch(ctx context.Context, req *PatchMetricSourceTableItemReq, options ...larkcore.RequestOptionFunc) (*PatchMetricSourceTableItemResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/okr/v1/metric_sources/:metric_source_id/tables/:metric_table_id/items/:metric_item_id"
	apiReq.HttpMethod = http.MethodPatch
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser}
	apiResp, err := larkcore.Request(ctx, apiReq, m.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &PatchMetricSourceTableItemResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, m.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 批量获取OKR
//
// - 根据OKR id批量获取OKR
//
// - 使用<md-tag mode="inline" type="token-tenant">tenant_access_token</md-tag>需要额外申请权限<md-perm ;href="https://open.feishu.cn/document/ukTMukTMukTM/uQjN3QjL0YzN04CN2cDN">以应用身份访问OKR信息</md-perm>
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/okr-v1/okr/batch_get
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/okrv1/batchGet_okr.go
func (o *okr) BatchGet(ctx context.Context, req *BatchGetOkrReq, options ...larkcore.RequestOptionFunc) (*BatchGetOkrResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/okr/v1/okrs/batch_get"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeUser, larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, o.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &BatchGetOkrResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, o.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取OKR周期列表
//
// - 获取OKR周期列表
//
// - 使用<md-tag mode="inline" type="token-tenant">tenant_access_token</md-tag>需要额外申请权限<md-perm ;href="https://open.feishu.cn/document/ukTMukTMukTM/uQjN3QjL0YzN04CN2cDN">以应用身份访问OKR信息</md-perm>
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/okr-v1/period/list
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/okrv1/list_period.go
func (p *period) List(ctx context.Context, req *ListPeriodReq, options ...larkcore.RequestOptionFunc) (*ListPeriodResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/okr/v1/periods"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, p.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListPeriodResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, p.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 创建OKR进展记录
//
// - 创建OKR进展记录
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/okr-v1/progress_record/create
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/okrv1/create_progressRecord.go
func (p *progressRecord) Create(ctx context.Context, req *CreateProgressRecordReq, options ...larkcore.RequestOptionFunc) (*CreateProgressRecordResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/okr/v1/progress_records"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser}
	apiResp, err := larkcore.Request(ctx, apiReq, p.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &CreateProgressRecordResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, p.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 删除OKR进展记录
//
// - 根据ID删除OKR进展记录
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/okr-v1/progress_record/delete
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/okrv1/delete_progressRecord.go
func (p *progressRecord) Delete(ctx context.Context, req *DeleteProgressRecordReq, options ...larkcore.RequestOptionFunc) (*DeleteProgressRecordResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/okr/v1/progress_records/:progress_id"
	apiReq.HttpMethod = http.MethodDelete
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser}
	apiResp, err := larkcore.Request(ctx, apiReq, p.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &DeleteProgressRecordResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, p.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取OKR进展记录
//
// - 根据ID获取OKR进展记录详情
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/okr-v1/progress_record/get
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/okrv1/get_progressRecord.go
func (p *progressRecord) Get(ctx context.Context, req *GetProgressRecordReq, options ...larkcore.RequestOptionFunc) (*GetProgressRecordResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/okr/v1/progress_records/:progress_id"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser}
	apiResp, err := larkcore.Request(ctx, apiReq, p.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &GetProgressRecordResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, p.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 更新OKR进展记录
//
// - 根据OKR进展记录ID更新进展详情
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/okr-v1/progress_record/update
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/okrv1/update_progressRecord.go
func (p *progressRecord) Update(ctx context.Context, req *UpdateProgressRecordReq, options ...larkcore.RequestOptionFunc) (*UpdateProgressRecordResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/okr/v1/progress_records/:progress_id"
	apiReq.HttpMethod = http.MethodPut
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant, larkcore.AccessTokenTypeUser}
	apiResp, err := larkcore.Request(ctx, apiReq, p.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &UpdateProgressRecordResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, p.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取用户的OKR列表
//
// - 根据用户的id获取OKR列表
//
// - 使用<md-tag mode="inline" type="token-tenant">tenant_access_token</md-tag>需要额外申请权限<md-perm ;href="https://open.feishu.cn/document/ukTMukTMukTM/uQjN3QjL0YzN04CN2cDN">以应用身份访问OKR信息</md-perm>
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/okr-v1/user-okr/list
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/okrv1/list_userOkr.go
func (u *userOkr) List(ctx context.Context, req *ListUserOkrReq, options ...larkcore.RequestOptionFunc) (*ListUserOkrResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/okr/v1/users/:user_id/okrs"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeUser, larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, u.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListUserOkrResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, u.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}