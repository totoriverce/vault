/*
* Copyright 2018 - Present Okta, Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

// AUTO-GENERATED!  DO NOT EDIT FILE DIRECTLY

package query

import (
	"net/url"
	"strconv"
)

type Params struct {
	Q                    string `json:"q,omitempty"`
	After                string `json:"after,omitempty"`
	Limit                int64  `json:"limit,omitempty"`
	Filter               string `json:"filter,omitempty"`
	Expand               string `json:"expand,omitempty"`
	IncludeNonDeleted    *bool  `json:"includeNonDeleted,omitempty"`
	Activate             *bool  `json:"activate,omitempty"`
	ValidityYears        int64  `json:"validityYears,omitempty"`
	TargetAid            string `json:"targetAid,omitempty"`
	QueryScope           string `json:"query_scope,omitempty"`
	SendEmail            *bool  `json:"sendEmail,omitempty"`
	Cursor               string `json:"cursor,omitempty"`
	Mode                 string `json:"mode,omitempty"`
	Search               string `json:"search,omitempty"`
	DisableNotifications string `json:"disableNotifications,omitempty"`
	Type                 string `json:"type,omitempty"`
	TargetIdpId          string `json:"targetIdpId,omitempty"`
	Until                string `json:"until,omitempty"`
	Since                string `json:"since,omitempty"`
	SortOrder            string `json:"sortOrder,omitempty"`
	Status               string `json:"status,omitempty"`
	TemplateType         string `json:"templateType,omitempty"`
	SortBy               string `json:"sortBy,omitempty"`
	Provider             *bool  `json:"provider,omitempty"`
	NextLogin            string `json:"nextLogin,omitempty"`
	Strict               *bool  `json:"strict,omitempty"`
	UpdatePhone          *bool  `json:"updatePhone,omitempty"`
	TemplateId           string `json:"templateId,omitempty"`
	TokenLifetimeSeconds int64  `json:"tokenLifetimeSeconds,omitempty"`
	ScopeId              string `json:"scopeId,omitempty"`
	OauthTokens          *bool  `json:"oauthTokens,omitempty"`
}

func NewQueryParams(paramOpt ...ParamOptions) *Params {
	p := &Params{}

	for _, par := range paramOpt {
		par(p)
	}

	return p
}

type ParamOptions func(*Params)

func WithQ(queryQ string) ParamOptions {
	return func(p *Params) {
		p.Q = queryQ
	}
}
func WithAfter(queryAfter string) ParamOptions {
	return func(p *Params) {
		p.After = queryAfter
	}
}
func WithLimit(queryLimit int64) ParamOptions {
	return func(p *Params) {
		p.Limit = queryLimit
	}
}
func WithFilter(queryFilter string) ParamOptions {
	return func(p *Params) {
		p.Filter = queryFilter
	}
}
func WithExpand(queryExpand string) ParamOptions {
	return func(p *Params) {
		p.Expand = queryExpand
	}
}
func WithIncludeNonDeleted(queryIncludeNonDeleted bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = queryIncludeNonDeleted
		p.IncludeNonDeleted = b
	}
}
func WithActivate(queryActivate bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = queryActivate
		p.Activate = b
	}
}
func WithValidityYears(queryValidityYears int64) ParamOptions {
	return func(p *Params) {
		p.ValidityYears = queryValidityYears
	}
}
func WithTargetAid(queryTargetAid string) ParamOptions {
	return func(p *Params) {
		p.TargetAid = queryTargetAid
	}
}
func WithQueryScope(queryQueryScope string) ParamOptions {
	return func(p *Params) {
		p.QueryScope = queryQueryScope
	}
}
func WithSendEmail(querySendEmail bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = querySendEmail
		p.SendEmail = b
	}
}
func WithCursor(queryCursor string) ParamOptions {
	return func(p *Params) {
		p.Cursor = queryCursor
	}
}
func WithMode(queryMode string) ParamOptions {
	return func(p *Params) {
		p.Mode = queryMode
	}
}
func WithSearch(querySearch string) ParamOptions {
	return func(p *Params) {
		p.Search = querySearch
	}
}
func WithDisableNotifications(queryDisableNotifications string) ParamOptions {
	return func(p *Params) {
		p.DisableNotifications = queryDisableNotifications
	}
}
func WithType(queryType string) ParamOptions {
	return func(p *Params) {
		p.Type = queryType
	}
}
func WithTargetIdpId(queryTargetIdpId string) ParamOptions {
	return func(p *Params) {
		p.TargetIdpId = queryTargetIdpId
	}
}
func WithUntil(queryUntil string) ParamOptions {
	return func(p *Params) {
		p.Until = queryUntil
	}
}
func WithSince(querySince string) ParamOptions {
	return func(p *Params) {
		p.Since = querySince
	}
}
func WithSortOrder(querySortOrder string) ParamOptions {
	return func(p *Params) {
		p.SortOrder = querySortOrder
	}
}
func WithStatus(queryStatus string) ParamOptions {
	return func(p *Params) {
		p.Status = queryStatus
	}
}
func WithTemplateType(queryTemplateType string) ParamOptions {
	return func(p *Params) {
		p.TemplateType = queryTemplateType
	}
}
func WithSortBy(querySortBy string) ParamOptions {
	return func(p *Params) {
		p.SortBy = querySortBy
	}
}
func WithProvider(queryProvider bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = queryProvider
		p.Provider = b
	}
}
func WithNextLogin(queryNextLogin string) ParamOptions {
	return func(p *Params) {
		p.NextLogin = queryNextLogin
	}
}
func WithStrict(queryStrict bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = queryStrict
		p.Strict = b
	}
}
func WithUpdatePhone(queryUpdatePhone bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = queryUpdatePhone
		p.UpdatePhone = b
	}
}
func WithTemplateId(queryTemplateId string) ParamOptions {
	return func(p *Params) {
		p.TemplateId = queryTemplateId
	}
}
func WithTokenLifetimeSeconds(queryTokenLifetimeSeconds int64) ParamOptions {
	return func(p *Params) {
		p.TokenLifetimeSeconds = queryTokenLifetimeSeconds
	}
}
func WithScopeId(queryScopeId string) ParamOptions {
	return func(p *Params) {
		p.ScopeId = queryScopeId
	}
}
func WithOauthTokens(queryOauthTokens bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = queryOauthTokens
		p.OauthTokens = b
	}
}

func (p *Params) String() string {
	qs := url.Values{}

	if p.Q != "" {
		qs.Add(`q`, p.Q)
	}
	if p.After != "" {
		qs.Add(`after`, p.After)
	}
	if p.Limit != 0 {
		qs.Add(`limit`, strconv.FormatInt(p.Limit, 10))
	}
	if p.Filter != "" {
		qs.Add(`filter`, p.Filter)
	}
	if p.Expand != "" {
		qs.Add(`expand`, p.Expand)
	}
	if p.IncludeNonDeleted != nil {
		qs.Add(`includeNonDeleted`, strconv.FormatBool(*p.IncludeNonDeleted))
	}
	if p.Activate != nil {
		qs.Add(`activate`, strconv.FormatBool(*p.Activate))
	}
	if p.ValidityYears != 0 {
		qs.Add(`validityYears`, strconv.FormatInt(p.ValidityYears, 10))
	}
	if p.TargetAid != "" {
		qs.Add(`targetAid`, p.TargetAid)
	}
	if p.QueryScope != "" {
		qs.Add(`query_scope`, p.QueryScope)
	}
	if p.SendEmail != nil {
		qs.Add(`sendEmail`, strconv.FormatBool(*p.SendEmail))
	}
	if p.Cursor != "" {
		qs.Add(`cursor`, p.Cursor)
	}
	if p.Mode != "" {
		qs.Add(`mode`, p.Mode)
	}
	if p.Search != "" {
		qs.Add(`search`, p.Search)
	}
	if p.DisableNotifications != "" {
		qs.Add(`disableNotifications`, p.DisableNotifications)
	}
	if p.Type != "" {
		qs.Add(`type`, p.Type)
	}
	if p.TargetIdpId != "" {
		qs.Add(`targetIdpId`, p.TargetIdpId)
	}
	if p.Until != "" {
		qs.Add(`until`, p.Until)
	}
	if p.Since != "" {
		qs.Add(`since`, p.Since)
	}
	if p.SortOrder != "" {
		qs.Add(`sortOrder`, p.SortOrder)
	}
	if p.Status != "" {
		qs.Add(`status`, p.Status)
	}
	if p.TemplateType != "" {
		qs.Add(`templateType`, p.TemplateType)
	}
	if p.SortBy != "" {
		qs.Add(`sortBy`, p.SortBy)
	}
	if p.Provider != nil {
		qs.Add(`provider`, strconv.FormatBool(*p.Provider))
	}
	if p.NextLogin != "" {
		qs.Add(`nextLogin`, p.NextLogin)
	}
	if p.Strict != nil {
		qs.Add(`strict`, strconv.FormatBool(*p.Strict))
	}
	if p.UpdatePhone != nil {
		qs.Add(`updatePhone`, strconv.FormatBool(*p.UpdatePhone))
	}
	if p.TemplateId != "" {
		qs.Add(`templateId`, p.TemplateId)
	}
	if p.TokenLifetimeSeconds != 0 {
		qs.Add(`tokenLifetimeSeconds`, strconv.FormatInt(p.TokenLifetimeSeconds, 10))
	}
	if p.ScopeId != "" {
		qs.Add(`scopeId`, p.ScopeId)
	}
	if p.OauthTokens != nil {
		qs.Add(`oauthTokens`, strconv.FormatBool(*p.OauthTokens))
	}

	if len(qs) != 0 {
		return "?" + qs.Encode()
	}
	return ""
}
