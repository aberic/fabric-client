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

package sdk

import (
	"encoding/pem"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/mocks"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"net/http"
)

// Peer
// "peer1", "http://peer1.com", "Org1MSP", nil,
func Peer(ccName, path, version, name, url, msp string, cert *pem.Block, roles ...string) *mocks.MockPeer {
	return &mocks.MockPeer{
		MockName:  name,
		MockURL:   url,
		MockRoles: roles,
		MockCert:  cert,
		MockMSP:   msp,
		Status:    http.StatusOK,
	}
}

// PeerChannel
// "peer1", "http://peer1.com", "Org1MSP", nil,
func PeerChannel(channelID, name, url, msp string, cert *pem.Block, roles ...string) (*mocks.MockPeer, error) {
	//prepare sample response
	response := new(pb.ChannelQueryResponse)
	channels := make([]*pb.ChannelInfo, 1)
	channels[0] = &pb.ChannelInfo{ChannelId: channelID}
	response.Channels = channels

	responseBytes, err := proto.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sample response")
	}

	return &mocks.MockPeer{
		MockName:  name,
		MockURL:   url,
		MockRoles: roles,
		MockCert:  cert,
		MockMSP:   msp,
		Status:    http.StatusOK,
		Payload:   responseBytes,
	}, nil
}
