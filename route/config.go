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
	"github.com/ennoo/fabric-client/grpc/server/chains"
	"github.com/ennoo/fabric-client/rafts"
	"github.com/ennoo/fabric-client/service"
	"github.com/ennoo/rivet"
	"github.com/ennoo/rivet/trans/response"
	"strings"
	"time"
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
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.Client)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.InitClient(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.InitClient(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func initClientSelf(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.ClientSelf)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.InitClientSelf(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.InitClientSelf(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func initClientCustom(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.ClientCustom)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.InitClientCustom(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.InitClientCustom(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func addOrSetPeerForChannel(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.ChannelPeer)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.AddOrSetPeerForChannel(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.AddOrSetPeerForChannel(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func addOrSetQueryChannelPolicyForChannel(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.ChannelPolicyQuery)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.AddOrSetQueryChannelPolicyForChannel(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.AddOrSetQueryChannelPolicyForChannel(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func addOrSetDiscoveryPolicyForChannel(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.ChannelPolicyDiscovery)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.AddOrSetDiscoveryPolicyForChannel(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.AddOrSetDiscoveryPolicyForChannel(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func addOrSetEventServicePolicyForChannel(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.ChannelPolicyEvent)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.AddOrSetEventServicePolicyForChannel(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.AddOrSetEventServicePolicyForChannel(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func addOrSetOrdererForOrganizations(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.OrganizationsOrder)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.AddOrSetOrdererForOrganizations(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.AddOrSetOrdererForOrganizations(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func addOrSetOrdererForOrganizationsSelf(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.OrganizationsOrderSelf)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.AddOrSetOrdererForOrganizationsSelf(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.AddOrSetOrdererForOrganizationsSelf(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func addOrSetOrgForOrganizations(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.OrganizationsOrg)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.AddOrSetOrgForOrganizations(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.AddOrSetOrgForOrganizations(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func addOrSetOrgForOrganizationsSelf(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.OrganizationsOrgSelf)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.AddOrSetOrgForOrganizationsSelf(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.AddOrSetOrgForOrganizationsSelf(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func addOrSetOrderer(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.Order)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.AddOrSetOrderer(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.AddOrSetOrderer(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func addOrSetOrdererSelf(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.OrderSelf)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.AddOrSetOrdererSelf(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.AddOrSetOrdererSelf(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func addOrSetPeer(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.Peer)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.AddOrSetPeer(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.AddOrSetPeer(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func addOrSetPeerSelf(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.PeerSelf)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.AddOrSetPeerSelf(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.AddOrSetPeerSelf(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func addOrSetCertificateAuthority(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.CertificateAuthority)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.AddOrSetCertificateAuthority(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.AddOrSetCertificateAuthority(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func addOrSetCertificateAuthoritySelf(router *response.Router) {
	rivet.Response().Do(router.Context, func(result *response.Result) {
		var in = new(service.CertificateAuthoritySelf)
		if err := router.Context.ShouldBindJSON(in); err != nil {
			result.SayFail(router.Context, err.Error())
			return
		}
		proxy(
			router,
			true,
			func() (err error) {
				if err = service.AddOrSetCertificateAuthoritySelf(in); nil != err {
					result.SayFail(router.Context, err.Error())
					return
				}
				result.SaySuccess(router.Context, "success")
				return
			},
			func(result *response.Result) {
				if _, err := chains.AddOrSetCertificateAuthoritySelf(rafts.LeaderURL(), in.Trans2pb()); nil != err {
					result.SayFail(router.Context, err.Error())
				} else {
					result.SaySuccess(router.Context, "success")
				}
			},
		)
	})
}

func proxy(router *response.Router, sleep bool, exec func() error, trans func(result *response.Result)) {
	result := &response.Result{}
	switch rafts.Character() {
	case rafts.RoleLeader: // 自身即为 Leader 节点
		err := exec()
		if nil == err {
			rafts.VersionAdd()
		}
	case rafts.RoleCandidate: // 等待选举结果，如果超时则返回
		if sleep {
			time.Sleep(1000 * time.Millisecond)
			proxy(router, false, exec, trans)
		} else {
			result.SayFail(router.Context, "leader is nil")
		}
	case rafts.RoleFollower: // 将该请求转发给Leader节点处理
		trans(result)
	default:
		result.SayFail(router.Context, "unknown err")
	}
}
