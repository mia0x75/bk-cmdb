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

package inst

import (
	"github.com/tencent/bk-cmdb/framework/common"
	"github.com/tencent/bk-cmdb/framework/core/output/module/model"
	"github.com/tencent/bk-cmdb/framework/core/types"
)

func createPlat(target model.Model) (CommonInstInterface, error) {
	return &inst{target: target, datas: types.MapStr{}}, nil
}

// findPlatsLikeName find all insts by inst name
func findPlatsLikeName(target model.Model, platName string) (Iterator, error) {
	cond := common.CreateCondition().Field(PlatName).Like(platName)
	return NewIteratorInst(target, cond)
}

// findPlatsByCondition find all insts by condition
func findPlatsByCondition(target model.Model, cond common.Condition) (Iterator, error) {
	return NewIteratorInst(target, cond)
}
