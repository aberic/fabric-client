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

package response

import (
	"fmt"
)

// RespError 自定义 error 对象
type RespError struct {
	ErrorCode string
	ErrorMsg  string

	// HTTPStatusCode http的错误编码
	HTTPStatusCode int
}

// Error 实现error接口，自定义error
func (resErr *RespError) Error() string {
	return fmt.Sprintf("%s,error code is:%s", resErr.ErrorMsg, resErr.ErrorCode)
}

// FormatMsg 添加参数
func (resErr *RespError) FormatMsg(args ...interface{}) *RespError {
	newMsg := fmt.Sprintf(resErr.ErrorMsg, args...)
	resErr.ErrorMsg = newMsg
	return resErr
}
