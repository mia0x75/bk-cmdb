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

package plugins

import (
	"github.com/tencent/bk-cmdb/common"
	"github.com/tencent/bk-cmdb/common/metadata"
	"github.com/tencent/bk-cmdb/web_server/middleware/user/plugins/manager"
	_ "github.com/tencent/bk-cmdb/web_server/middleware/user/plugins/register"

	"github.com/gin-gonic/gin"
)

func CurrentPlugin(c *gin.Context, version string) metadata.LoginUserPluginInerface {
	if "" == version {
		version = common.BKDefaultLoginUserPluginVersion
	}

	var selfPlugin *metadata.LoginPluginInfo
	for _, plugin := range manager.LoginPluginInfo {
		if plugin.Version == version {
			return plugin.HandleFunc
		}
		if common.BKDefaultLoginUserPluginVersion == plugin.Version {
			selfPlugin = plugin
		}
	}
	if nil != selfPlugin {
		return selfPlugin.HandleFunc
	}

	return nil
}