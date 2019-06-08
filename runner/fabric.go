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

package main

import (
	"github.com/ennoo/fabric-go-client/route"
	"github.com/ennoo/rivet"
)

func main() {
	rivet.Initialize(false, false, false)
	// rivet.UseDiscovery(discovery.ComponentConsul, "127.0.0.1:8500", "test", "127.0.0.1", 8081)
	rivet.ListenAndServe(&rivet.ListenServe{
		Engine:      rivet.SetupRouter(route.ChainCode, route.Config),
		DefaultPort: "19865",
	})
}
