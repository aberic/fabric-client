/*
 * Copyright (c) 2019. ENNOO - All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package bow

import (
	"github.com/ennoo/rivet/trans/response"
	"github.com/gin-gonic/gin"
	"strings"
)

// Route 网关服务路由
func Route(engine *gin.Engine, filter func(result *response.Result) bool) {
	// 仓库相关路由设置
	vRepo := engine.Group("/")
	for index := range routeServices {
		bowService := routeServices[index]
		vRepo.Any(strings.Join([]string{bowService.InURI, "/*do"}, ""), func(context *gin.Context) {
			RunBow(context, bowService.Name, filter)
		})
	}
}
