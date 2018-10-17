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

package app

import (
	"context"
	"fmt"
	"os"
	"time"

	restful "github.com/emicklei/go-restful"

	"github.com/tencent/bk-cmdb/apimachinery"
	"github.com/tencent/bk-cmdb/apimachinery/util"
	"github.com/tencent/bk-cmdb/common"
	"github.com/tencent/bk-cmdb/common/backbone"
	cc "github.com/tencent/bk-cmdb/common/backbone/config"
	"github.com/tencent/bk-cmdb/common/blog"
	"github.com/tencent/bk-cmdb/common/types"
	"github.com/tencent/bk-cmdb/common/version"
	"github.com/tencent/bk-cmdb/source_controller/objectcontroller/app/options"
	"github.com/tencent/bk-cmdb/source_controller/objectcontroller/service"
	"github.com/tencent/bk-cmdb/storage/dal/mongo"
	dalredis "github.com/tencent/bk-cmdb/storage/dal/redis"
)

//Run ccapi server
func Run(ctx context.Context, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	c := &util.APIMachineryConfig{
		ZkAddr:    op.ServConf.RegDiscover,
		QPS:       1000,
		Burst:     2000,
		TLSConfig: nil,
	}

	machinery, err := apimachinery.NewApiMachinery(c)
	if err != nil {
		return fmt.Errorf("new api machinery failed, err: %v", err)
	}

	coreService := new(service.Service)
	server := backbone.Server{
		ListenAddr: svrInfo.IP,
		ListenPort: svrInfo.Port,
		Handler:    restful.NewContainer().Add(coreService.WebService()),
		TLS:        backbone.TLSConfig{},
	}

	regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, types.CC_MODULE_OBJECTCONTROLLER, svrInfo.IP)
	bonC := &backbone.Config{
		RegisterPath: regPath,
		RegisterInfo: *svrInfo,
		CoreAPI:      machinery,
		Server:       server,
	}

	objCtr := new(ObjectController)
	objCtr.Service = coreService
	objCtr.Core, err = backbone.NewBackbone(ctx, op.ServConf.RegDiscover,
		types.CC_MODULE_OBJECTCONTROLLER,
		op.ServConf.ExConfig,
		objCtr.onObjectConfigUpdate,
		bonC)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if nil == objCtr.Config {
			time.Sleep(time.Second)
		} else {
			configReady = true
			break
		}
	}
	if false == configReady {
		return fmt.Errorf("Configuration item not found")
	}

	select {
	case <-ctx.Done():
	}
	return nil
}

type ObjectController struct {
	*service.Service
	Config *options.Config
}

func (h *ObjectController) onObjectConfigUpdate(previous, current cc.ProcessConfig) {
	h.Config = &options.Config{
		Mongo: mongo.ParseConfigFromKV("mongodb", current.ConfigMap),
		Redis: dalredis.ParseConfigFromKV("redis", current.ConfigMap),
	}

	instance, err := mongo.NewMgo(h.Config.Mongo.BuildURI())
	if err != nil {
		blog.Errorf("new mongo client failed, err: %v", err)
		return
	}
	h.Service.Instance = instance

	cache, err := dalredis.NewFromConfig(h.Config.Redis)
	if err != nil {
		blog.Errorf("new redis client failed, err: %v", err)
		return
	}
	h.Cache = cache
}

func newServerInfo(op *options.ServerOption) (*types.ServerInfo, error) {
	ip, err := op.ServConf.GetAddress()
	if err != nil {
		return nil, err
	}

	port, err := op.ServConf.GetPort()
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	info := &types.ServerInfo{
		IP:       ip,
		Port:     port,
		HostName: hostname,
		Scheme:   "http",
		Version:  version.GetVersion(),
		Pid:      os.Getpid(),
	}
	return info, nil
}