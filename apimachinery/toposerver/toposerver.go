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

package toposerver

import (
	"fmt"

	"github.com/tencent/bk-cmdb/apimachinery/rest"
	"github.com/tencent/bk-cmdb/apimachinery/toposerver/inst"
	"github.com/tencent/bk-cmdb/apimachinery/toposerver/object"
	"github.com/tencent/bk-cmdb/apimachinery/toposerver/openapi"
	"github.com/tencent/bk-cmdb/apimachinery/toposerver/privilege"
	"github.com/tencent/bk-cmdb/apimachinery/util"
)

type TopoServerClientInterface interface {
	Instance() inst.InstanceInterface
	Object() object.ObjectInterface
	OpenAPI() openapi.OpenApiInterface
	Privilege() privilege.PrivilegeInterface
}

func NewTopoServerClient(c *util.Capability, version string) TopoServerClientInterface {
	base := fmt.Sprintf("/topo/%s", version)
	return &topoServer{
		restCli: rest.NewRESTClient(c, base),
	}
}

type topoServer struct {
	restCli rest.ClientInterface
}

func (t *topoServer) Instance() inst.InstanceInterface {
	return inst.NewInstanceClient(t.restCli)
}

func (t *topoServer) Object() object.ObjectInterface {
	return object.NewObjectInterface(t.restCli)
}

func (t *topoServer) OpenAPI() openapi.OpenApiInterface {
	return openapi.NewOpenApiInterface(t.restCli)
}

func (t *topoServer) Privilege() privilege.PrivilegeInterface {
	return privilege.NewPrivilegeInterface(t.restCli)
}
