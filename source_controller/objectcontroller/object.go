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

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/pflag"

	"github.com/tencent/bk-cmdb/common"
	"github.com/tencent/bk-cmdb/common/blog"
	"github.com/tencent/bk-cmdb/common/types"
	"github.com/tencent/bk-cmdb/common/util"
	"github.com/tencent/bk-cmdb/source_controller/objectcontroller/app"
	"github.com/tencent/bk-cmdb/source_controller/objectcontroller/app/options"
)

func main() {
	common.SetIdentification(types.CC_MODULE_OBJECTCONTROLLER)
	runtime.GOMAXPROCS(runtime.NumCPU())

	blog.InitLogs()
	defer blog.CloseLogs()

	op := options.NewServerOption()
	op.AddFlags(pflag.CommandLine)

	util.InitFlags()

	if err := app.Run(context.Background(), op); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		blog.Fatal(err)
	}
}