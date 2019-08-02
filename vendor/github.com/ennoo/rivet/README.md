# rivet [![GoDoc](https://godoc.org/github.com/ennoo/rivet?status.svg)](https://godoc.org/github.com/ennoo/rivet) [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fennoo%2Frivet.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fennoo%2Frivet?ref=badge_shield) ![GitHub](https://img.shields.io/github/license/ennoo/rivet.svg)

[![Go Report Card](https://goreportcard.com/badge/github.com/ennoo/rivet)](https://goreportcard.com/report/github.com/ennoo/rivet)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/474e759e4a7b48c3b4aaefda5079f1d3)](https://www.codacy.com/app/aberic/rivet?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=ennoo/rivet&amp;utm_campaign=Badge_Grade)
[![Travis (.org)](https://img.shields.io/travis/ennoo/rivet.svg?label=travis-ci%20build)](https://www.travis-ci.org/ennoo/rivet)
[![CircleCI (all branches)](https://img.shields.io/circleci/project/github/ennoo/rivet.svg?label=circle-ci%20build)](https://circleci.com/gh/ennoo/rivet)
[![Coveralls github](https://img.shields.io/coveralls/github/ennoo/rivet.svg)](https://coveralls.io/github/ennoo/rivet?branch=master)

rivet提供一套用于go开发的微服务解决方案，包括网关、负载均衡、熔断降级、请求转发等功能，目前支持consul作为第三方服务发现组件。除此外也基于第三方开源框架做了简单的封装，如Http/Https、MySQL数据库以及一些常用的工具方法。
在examples中有上述相关组件和实现的Demo。
<br><br>

组件信息
------------

| Bow            | 网关服务                                        |
| :------------- |:-----------------------------------------------------------------|
| docker         | [![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/ennoo/bow.svg)](https://hub.docker.com/r/ennoo/bow/dockerfile) [![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/ennoo/bow.svg)](https://hub.docker.com/r/ennoo/bow/builds) [![](https://images.microbadger.com/badges/image/ennoo/bow.svg)](https://microbadger.com/images/ennoo/bow "Get your own image badge on microbadger.com") [![](https://images.microbadger.com/badges/version/ennoo/bow.svg)](https://microbadger.com/images/ennoo/bow "Get your own version badge on microbadger.com") [![Docker Pulls](https://img.shields.io/docker/pulls/ennoo/bow.svg?label=pulls)](https://hub.docker.com/r/ennoo/bow)|


| Shunt          | 负载均衡                                   |
| :------------- |:-----------------------------------------------------------------|
| docker         | [![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/ennoo/shunt.svg)](https://hub.docker.com/r/ennoo/shunt/dockerfile) [![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/ennoo/shunt.svg)](https://hub.docker.com/r/ennoo/bow/builds) [![](https://images.microbadger.com/badges/image/ennoo/shunt.svg)](https://microbadger.com/images/ennoo/shunt "Get your own image badge on microbadger.com") [![](https://images.microbadger.com/badges/version/ennoo/shunt.svg)](https://microbadger.com/images/ennoo/shunt "Get your own version badge on microbadger.com") [![Docker Pulls](https://img.shields.io/docker/pulls/ennoo/shunt.svg?label=pulls)](https://hub.docker.com/r/ennoo/shunt)|

