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

package openapi

import (
	"context"
	"net/http"

	"github.com/tencent/bk-cmdb/apimachinery/rest"
	"github.com/tencent/bk-cmdb/common/metadata"
)

type OpenApiInterface interface {
	GetProcessesByModuleName(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.ProcInfoArrResult, err error)
	DeleteSetHost(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
}

func NewOpenApiInterface(client rest.ClientInterface) OpenApiInterface {
	return &openAPI{client: client}
}

type openAPI struct {
	client rest.ClientInterface
}
