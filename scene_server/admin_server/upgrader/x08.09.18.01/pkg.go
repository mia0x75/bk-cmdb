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

package x08_09_18_01

import (
	"context"

	"github.com/tencent/bk-cmdb/common/blog"
	"github.com/tencent/bk-cmdb/scene_server/admin_server/upgrader"
	"github.com/tencent/bk-cmdb/storage/dal"
)

func init() {
	upgrader.RegistUpgrader("x08.09.18.01", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	err = fixedHostPlatAssocateRelation(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade x08.09.18.01] fixedHostPlatAssocateRelation error  %s", err.Error())
		return err
	}

	return nil
}