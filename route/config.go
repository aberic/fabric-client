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

package route

import (
	conf "github.com/ennoo/fabric-client/config"
	pb "github.com/ennoo/fabric-client/grpc/proto/raft"
	"github.com/ennoo/fabric-client/raft"
	"github.com/ennoo/fabric-client/service"
	"github.com/ennoo/rivet"
	"github.com/ennoo/rivet/trans/response"
	str "github.com/ennoo/rivet/utils/string"
	"net/http"
	"strings"
)

func Config(router *response.Router) {
	// 仓库相关路由设置
	router.Group = router.Engine.Group("/config")
	router.GET("/get/:configID", get)
	router.POST("/client", initClient)
	router.POST("/client/custom", initClientCustom)
	router.POST("/channel/peer", addOrSetPeerForChannel)
	router.POST("/channel/policy/query", addOrSetQueryChannelPolicyForChannel)
	router.POST("/channel/policy/discovery", addOrSetDiscoveryPolicyForChannel)
	router.POST("/channel/policy/event", addOrSetEventServicePolicyForChannel)
	router.POST("/organizations/order", addOrSetOrdererForOrganizations)
	router.POST("/organizations/org", addOrSetOrgForOrganizations)
	router.POST("/order", addOrSetOrderer)
	router.POST("/peer", addOrSetPeer)
	router.POST("/ca", addOrSetCertificateAuthority)
	router.POST("/self/client", initClientSelf)
	router.POST("/self/channel/peer", addOrSetPeerForChannel)
	router.POST("/self/channel/policy/query", addOrSetQueryChannelPolicyForChannel)
	router.POST("/self/channel/policy/discovery", addOrSetDiscoveryPolicyForChannel)
	router.POST("/self/channel/policy/event", addOrSetEventServicePolicyForChannel)
	router.POST("/self/organizations/order", addOrSetOrdererForOrganizationsSelf)
	router.POST("/self/organizations/org", addOrSetOrgForOrganizationsSelf)
	router.POST("/self/order", addOrSetOrdererSelf)
	router.POST("/self/peer", addOrSetPeerSelf)
	router.POST("/self/ca", addOrSetCertificateAuthoritySelf)
	router.POST("/sync", sync)
	router.POST("/sync/ask", askSync)
}

func askSync(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var configs = new(map[string]*conf.Config)
			if err := router.Context.ShouldBindJSON(configs); err != nil {
				result.SayFail(router.Context, err.Error())
				go raft.SyncConfig()
				return
			}
			for keyNet, configNet := range *configs {
				haveOne := false
				for keyLocal := range service.Configs {
					if keyNet == keyLocal {
						haveOne = true
						break
					}
				}
				if !haveOne {
					service.Configs[keyNet] = configNet
				}
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/sync/ask", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func sync(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var configs = new(map[string]*conf.Config)
		if err := router.Context.ShouldBindJSON(configs); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		service.Configs = *configs
		result.Say(router.Context)
	})
}

func get(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		configID := router.Context.Param("configID")
		config := service.Get(configID)
		if nil == config {
			result.SayFail(router.Context, strings.Join([]string{"config ", configID, " is not exist"}, ""))
		} else {
			result.SaySuccess(router.Context, service.Get(configID))
		}
	})
}

func initClient(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var client = new(service.Client)
			if err := router.Context.ShouldBindJSON(client); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.InitClient(client); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/client", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func initClientSelf(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var client = new(service.ClientSelf)
			if err := router.Context.ShouldBindJSON(client); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.InitClientSelf(client); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/self/client", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func initClientCustom(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var clientCustom = new(service.ClientCustom)
			if err := router.Context.ShouldBindJSON(clientCustom); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.InitClientCustom(clientCustom); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/client/custom", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func addOrSetPeerForChannel(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var channelPeer = new(service.ChannelPeer)
			if err := router.Context.ShouldBindJSON(channelPeer); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.AddOrSetPeerForChannel(channelPeer); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/channel/peer", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func addOrSetQueryChannelPolicyForChannel(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var query = new(service.ChannelPolicyQuery)
			if err := router.Context.ShouldBindJSON(query); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.AddOrSetQueryChannelPolicyForChannel(query); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/channel/policy/query", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func addOrSetDiscoveryPolicyForChannel(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var discovery = new(service.ChannelPolicyDiscovery)
			if err := router.Context.ShouldBindJSON(discovery); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.AddOrSetDiscoveryPolicyForChannel(discovery); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/channel/policy/discovery", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func addOrSetEventServicePolicyForChannel(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var event = new(service.ChannelPolicyEvent)
			if err := router.Context.ShouldBindJSON(event); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.AddOrSetEventServicePolicyForChannel(event); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/channel/policy/event", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func addOrSetOrdererForOrganizations(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var order = new(service.OrganizationsOrder)
			if err := router.Context.ShouldBindJSON(order); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.AddOrSetOrdererForOrganizations(order); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/organizations/order", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func addOrSetOrdererForOrganizationsSelf(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var order = new(service.OrganizationsOrderSelf)
			if err := router.Context.ShouldBindJSON(order); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.AddOrSetOrdererForOrganizationsSelf(order); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/self/organizations/order", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func addOrSetOrgForOrganizations(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var org = new(service.OrganizationsOrg)
			if err := router.Context.ShouldBindJSON(org); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.AddOrSetOrgForOrganizations(org); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/organizations/org", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func addOrSetOrgForOrganizationsSelf(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var org = new(service.OrganizationsOrgSelf)
			if err := router.Context.ShouldBindJSON(org); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.AddOrSetOrgForOrganizationsSelf(org); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/self/organizations/org", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func addOrSetOrderer(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var orderer = new(service.Order)
			if err := router.Context.ShouldBindJSON(orderer); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.AddOrSetOrderer(orderer); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/order", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func addOrSetOrdererSelf(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var orderer = new(service.OrderSelf)
			if err := router.Context.ShouldBindJSON(orderer); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.AddOrSetOrdererSelf(orderer); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/self/order", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func addOrSetPeer(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var peer = new(service.Peer)
			if err := router.Context.ShouldBindJSON(peer); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.AddOrSetPeer(peer); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/peer", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func addOrSetPeerSelf(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var peer = new(service.PeerSelf)
			if err := router.Context.ShouldBindJSON(peer); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.AddOrSetPeerSelf(peer); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/self/peer", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func addOrSetCertificateAuthority(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var certificateAuthority = new(service.CertificateAuthority)
			if err := router.Context.ShouldBindJSON(certificateAuthority); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.AddOrSetCertificateAuthority(certificateAuthority); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/ca", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}

func addOrSetCertificateAuthoritySelf(router *response.Router) {
	if (raft.Nodes[raft.ID].Status == pb.Status_LEADER && raft.Leader.BrokerID == raft.ID) || str.IsEmpty(raft.ID) { // 如果相等，则说明自身即为 Leader 节点
		rivet.Response().Do(router.Context, func(result *response.Result) {
			var certificateAuthority = new(service.CertificateAuthoritySelf)
			if err := router.Context.ShouldBindJSON(certificateAuthority); err != nil {
				result.SayFail(router.Context, err.Error())
				return
			}
			if err := service.AddOrSetCertificateAuthoritySelf(certificateAuthority); nil != err {
				result.SayFail(router.Context, err.Error())
				return
			}
			go raft.SyncConfig()
			result.SaySuccess(router.Context, "success")
		})
	} else { // 将该请求转发给Leader节点处理
		leader := raft.Nodes[raft.Leader.BrokerID]
		rivet.Request().Callback(router.Context, http.MethodPost,
			strings.Join([]string{"http://", leader.Addr, ":", leader.Http}, ""), "config/self/ca", func() *response.Result {
				return &response.Result{ResultCode: response.Fail, Msg: "请求失败，检查集群状态"}
			})
	}
}
