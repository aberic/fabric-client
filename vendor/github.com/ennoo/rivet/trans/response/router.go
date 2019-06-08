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
 */

package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Router 路由入口对象
type Router struct {
	Engine  *gin.Engine
	Context *gin.Context
	Group   *gin.RouterGroup
}

// Use 路由 Use 请求方法
func (router *Router) Use(relativePath string, f func(router *Router)) {
	router.Group.Use(func(context *gin.Context) {
		router.Context = context
		f(router)
	})
}

// Handle registers a new request handle and middleware with the given path and method.
// The last handler should be the real handler, the other ones should be middleware that can and should be shared among different routes.
// See the example code in GitHub.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (router *Router) Handle(relativePath string, f func(router *Router)) {
	router.Group.Use(func(context *gin.Context) {
		router.Context = context
		f(router)
	})
}

// Any 路由 Any 请求方法
func (router *Router) Any(relativePath string, f func(router *Router)) {
	router.Group.GET(relativePath, func(context *gin.Context) {
		router.Context = context
		f(router)
	})
}

// GET 路由 GET 请求方法
func (router *Router) GET(relativePath string, f func(router *Router)) {
	router.Group.GET(relativePath, func(context *gin.Context) {
		router.Context = context
		f(router)
	})
}

// POST 路由 POST 请求方法
func (router *Router) POST(relativePath string, f func(router *Router)) {
	router.Group.POST(relativePath, func(context *gin.Context) {
		router.Context = context
		f(router)
	})
}

// DELETE 路由 DELETE 请求方法
func (router *Router) DELETE(relativePath string, f func(router *Router)) {
	router.Group.DELETE(relativePath, func(context *gin.Context) {
		router.Context = context
		f(router)
	})
}

// PATCH 路由 PATCH 请求方法
func (router *Router) PATCH(relativePath string, f func(router *Router)) {
	router.Group.PATCH(relativePath, func(context *gin.Context) {
		router.Context = context
		f(router)
	})
}

// PUT 路由 PUT 请求方法
func (router *Router) PUT(relativePath string, f func(router *Router)) {
	router.Group.PUT(relativePath, func(context *gin.Context) {
		router.Context = context
		f(router)
	})
}

// OPTIONS 路由 OPTIONS 请求方法
func (router *Router) OPTIONS(relativePath string, f func(router *Router)) {
	router.Group.OPTIONS(relativePath, func(context *gin.Context) {
		router.Context = context
		f(router)
	})
}

// HEAD 路由 HEAD 请求方法
func (router *Router) HEAD(relativePath string, f func(router *Router)) {
	router.Group.HEAD(relativePath, func(context *gin.Context) {
		router.Context = context
		f(router)
	})
}

// StaticFile 路由 StaticFile 请求方法
// StaticFile registers a single route in order to serve a single file of the local filesystem.
// router.StaticFile("favicon.ico", "./resources/favicon.ico")
func (router *Router) StaticFile(relativePath, filepath string) {
	router.Group.StaticFile(relativePath, filepath)
}

// Static 路由 Static 请求方法
// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//     router.Static("/static", "/var/www")
func (router *Router) Static(relativePath, root string) {
	router.Group.Static(relativePath, root)
}

// StaticFS 路由 StaticFS 请求方法
// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
// Gin by default user: gin.Dir()
func (router *Router) StaticFS(relativePath string, fs http.FileSystem) {
	router.Group.StaticFS(relativePath, fs)
}
