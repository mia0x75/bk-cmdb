/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package user

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/tencent/bk-cmdb/common"
	"github.com/tencent/bk-cmdb/common/blog"
	"github.com/tencent/bk-cmdb/common/core/cc/api"
	"github.com/tencent/bk-cmdb/common/metadata"
	"github.com/tencent/bk-cmdb/web_server/application/middleware/user/plugins"
	webCommon "github.com/tencent/bk-cmdb/web_server/common"
)

type publicUser struct {
}

// LoginUser  user login
func (m *publicUser) LoginUser(c *gin.Context) bool {

	ccapi := api.NewAPIResource()
	config, _ := ccapi.ParseConfig()
	user := plugins.CurrentPlugin(c)
	isMultiOwner := false
	multipleOwner, ok := config["session.multiple_owner"]
	if ok && common.LoginSystemMultiSupplierTrue == multipleOwner {
		isMultiOwner = true
	}
	userInfo, loginSucc := user.LoginUser(c, config, isMultiOwner)
	if !loginSucc {
		return false
	}

	if true == isMultiOwner || true == userInfo.MultiSupplier {
		err := NewOwnerManager(userInfo.UserName, userInfo.OnwerUin, userInfo.Language).InitOwner()
		if nil != err {
			blog.Error("InitOwner error: %v", err)
			return false
		}
	}
	strOwnerUinlist := []byte("")
	if 0 != len(userInfo.OwnerUinArr) {
		strOwnerUinlist, _ = json.Marshal(userInfo.OwnerUinArr)
	}

	cookielanguage, _ := c.Cookie("blueking_language")
	session := sessions.Default(c)
	session.Set(common.WEBSessionUinKey, userInfo.UserName)
	session.Set(common.WEBSessionChineseNameKey, userInfo.ChName)
	session.Set(common.WEBSessionPhoneKey, userInfo.Phone)
	session.Set(common.WEBSessionEmailKey, userInfo.Email)
	session.Set(common.WEBSessionRoleKey, userInfo.Role)
	session.Set(common.HTTPCookieBKToken, userInfo.BkToken)
	session.Set(common.WEBSessionOwnerUinKey, userInfo.OnwerUin)
	session.Set(common.WEBSessionAvatarUrlKey, userInfo.AvatarUrl)
	session.Set(common.WEBSessionOwnerUinListeKey, string(strOwnerUinlist))
	if userInfo.MultiSupplier {
		session.Set(common.WEBSessionMultiSupplierKey, common.LoginSystemMultiSupplierTrue)
	} else {
		session.Set(common.WEBSessionMultiSupplierKey, common.LoginSystemMultiSupplierFalse)
	}

	session.Set(webCommon.IsSkipLogin, "0")
	if "" != cookielanguage {
		session.Set(common.WEBSessionLanguageKey, cookielanguage)
	} else {
		session.Set(common.WEBSessionLanguageKey, userInfo.Language)
	}
	session.Save()
	return true
}

// GetUserList get user list from paas
func (m *publicUser) GetUserList(c *gin.Context) (int, interface{}) {

	ccapi := api.NewAPIResource()
	config, _ := ccapi.ParseConfig()
	user := plugins.CurrentPlugin(c)
	userList, err := user.GetUserList(c, config)
	rspBody := metadata.LonginSystemUserListResult{}
	if nil != err {
		rspBody.Code = common.CCErrCommHTTPDoRequestFailed
		rspBody.ErrMsg = err.Error()
		rspBody.Result = false
	}
	rspBody.Result = true
	rspBody.Data = userList
	return 200, rspBody
}

func (m *publicUser) GetLoginUrl(c *gin.Context) string {
	ccapi := api.NewAPIResource()
	config, _ := ccapi.ParseConfig()
	siteUrl, ok := config["site.domain_url"]
	if ok {
		siteUrl = strings.Trim(siteUrl, "/")
	}
	params := new(metadata.LogoutRequestParams)
	err := json.NewDecoder(c.Request.Body).Decode(params)
	if nil != err || (common.LogoutHTTPSchemeHTTP != params.HTTPScheme && common.LogoutHTTPSchemeHTTPS != params.HTTPScheme) {
		params.HTTPScheme, err = c.Cookie(common.LogoutHTTPSchemeCookieKey)
		if nil != err || (common.LogoutHTTPSchemeHTTP != params.HTTPScheme && common.LogoutHTTPSchemeHTTPS != params.HTTPScheme) {
			params.HTTPScheme = common.LogoutHTTPSchemeHTTP
		}
	}
	user := plugins.CurrentPlugin(c)
	return user.GetLoginUrl(c, config, params)

}
