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
	"bytes"
	"errors"
	"github.com/aberic/fabric-client/geneses"
	"github.com/aberic/gnomon"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/common/tools/configtxlator/update"
	common2 "github.com/hyperledger/fabric/protos/common"
	"golang.org/x/protobuf/proto"
)

func marshalCommonEnvelope(block *common.Block) (*common.Envelope, error) {
	envelope := &common.Envelope{}
	err := proto.Unmarshal(block.Data.Data[0], envelope)
	return envelope, err
}

func marshalCommonPayload(envelope *common.Envelope) (*common.Payload, error) {
	payload := &common.Payload{}
	err := proto.Unmarshal(envelope.Payload, payload)
	return payload, err
}

func marshalCommonConfigEnvelope(payload *common.Payload) (*common2.ConfigEnvelope, error) {
	configEnvelope := &common2.ConfigEnvelope{}
	err := proto.Unmarshal(payload.Data, configEnvelope)
	return configEnvelope, err
}

//func getConfigGroupData(block *common.Block, consortium, groupName string) ([]byte, error) {
//	var (
//		envelope       *common.Envelope
//		payload        *common.Payload
//		configEnvelope *common.ConfigEnvelope
//		err            error
//	)
//	if envelope, err = marshalCommonEnvelope(block); nil != err {
//		return nil, err
//	}
//	if payload, err = marshalCommonPayload(envelope); nil != err {
//		return nil, err
//	}
//	if configEnvelope, err = marshalCommonConfigEnvelope(payload); nil != err {
//		return nil, err
//	}
//	group := configEnvelope.Config.ChannelGroup.Groups["Consortiums"].Groups[consortium].Groups[groupName]
//	return proto.Marshal(group)
//}

func getApplication(block *common.Block, isApp bool) (*common2.Config, *common2.ConfigGroup, error) {
	var (
		envelope       *common.Envelope
		payload        *common.Payload
		configEnvelope *common2.ConfigEnvelope
		configGroup    *common2.ConfigGroup
		err            error
	)
	if envelope, err = marshalCommonEnvelope(block); nil != err {
		return nil, nil, err
	}
	if payload, err = marshalCommonPayload(envelope); nil != err {
		return nil, nil, err
	}
	if configEnvelope, err = marshalCommonConfigEnvelope(payload); nil != err {
		return nil, nil, err
	}
	if isApp {
		configGroup = configEnvelope.Config.ChannelGroup.Groups["Application"]
	} else {
		configGroup = configEnvelope.Config.ChannelGroup.Groups["Consortiums"]
	}
	if nil == configGroup {
		return nil, nil, errors.New("config group is nil")
	}
	return configEnvelope.Config, configGroup, nil
}

func AddGroup(originalBlock, updateBlock *common.Block, channelID, consortium, groupName string) ([]byte, error) {
	var (
		originalConfig                                             *common2.Config
		updateApplication                                          *common2.ConfigGroup
		configUpdateEnvBytes, configUpdateBytes, updateConfigBytes []byte
		err                                                        error
	)
	if originalConfig, _, err = getApplication(originalBlock, true); nil != err {
		return nil, err
	}
	if _, updateApplication, err = getApplication(updateBlock, false); nil != err {
		return nil, err
	}
	updateConfig := &common2.Config{}
	if updateConfigBytes, err = proto.Marshal(originalConfig); nil != err {
		return nil, err
	}
	if err = proto.Unmarshal(updateConfigBytes, updateConfig); nil != err {
		return nil, err
	}
	newGroup := updateApplication.Groups[consortium].Groups[groupName]
	updateConfig.ChannelGroup.Groups["Application"].Groups[groupName] = newGroup

	configUpdate, err := update.Compute(originalConfig, updateConfig)
	if nil != err {
		return nil, err
	}
	configUpdate.ChannelId = channelID
	gnomon.Log().Info("AddGroup", gnomon.Log().Field("configUpdate", configUpdate))

	if configUpdateBytes, err = proto.Marshal(configUpdate); nil != err {
		return nil, err
	}

	configUpdateEnv := &common2.ConfigUpdateEnvelope{}
	configUpdateEnv.ConfigUpdate = configUpdateBytes

	if configUpdateEnvBytes, err = proto.Marshal(configUpdateEnv); nil != err {
		return nil, err
	}

	return configUpdateEnvBytes, nil
}

func CreateConfigUpdateBytes(configUpdateEnvBytes []byte, leagueName, channelID string) error {
	var (
		envelope      *common.Envelope
		envelopeBytes []byte
		err           error
	)
	if envelope, err = createConfigEnvelopeReader(configUpdateEnvBytes, channelID); nil != err {
		return err
	}
	if envelopeBytes, err = proto.Marshal(envelope); nil != err {
		return err
	}
	channelUpdateFilePath := geneses.ChannelUpdateTXFilePath(leagueName, channelID)
	if _, err = gnomon.File().Append(channelUpdateFilePath, envelopeBytes, true); nil != err {
		return err
	}
	return nil
}

func Sign(configBytes, envelopeBytes []byte, leagueName, orgName, orgUser, channelID string, sdkOpts ...fabsdk.Option) error {
	var (
		ctx context.Client
		err error
	)
	clientContext, _ := clientContext(orgName, orgUser, configBytes, sdkOpts...)
	if ctx, err = clientContext(); nil != err {
		return err
	}
	signingMgr := ctx.SigningManager()
	signature, err := signingMgr.Sign(envelopeBytes, ctx.PrivateKey())
	if err != nil {
		return errors.New("sign failed")
	}
	envelope := &common.Envelope{}
	if err = proto.Unmarshal(envelopeBytes, envelope); nil != err {
		return err
	}
	envelope.Signature = signature
	if envelopeBytes, err = proto.Marshal(envelope); nil != err {
		return err
	}
	channelUpdateFilePath := geneses.ChannelUpdateTXFilePath(leagueName, channelID)
	if _, err = gnomon.File().Append(channelUpdateFilePath, envelopeBytes, true); nil != err {
		return err
	}
	return nil
}

func UpdateChannel(resmgmtClient *resmgmt.Client, envelopeBytes []byte, orderURL, channelID string) error {
	var (
		//envelope      *common.Envelope
		//envelopeBytes []byte
		err error
	)
	//if envelope, err = createConfigEnvelopeReader(block.Data.Data[0], configUpdateBytes, channelID);nil != err {
	//	return err
	//}
	//if envelopeBytes, err = proto.Marshal(envelope); nil != err {
	//	return err
	//}
	//channelUpdateFilePath := geneses.ChannelUpdateTXFilePath(leagueName, channelID)
	//if _, err = gnomon.File().Append(channelUpdateFilePath, envelopeBytes, true);nil!=err {
	//	return err
	//}
	reader := bytes.NewReader(envelopeBytes)
	//org1MspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg("org1"))
	//if nil != err {
	//	return err
	//}
	//org1AdminIdentity, err := org1MspClient.GetSigningIdentity("admin")
	//if nil != err {
	//	return err
	//}
	//org2MspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg("org2"))
	//if nil != err {
	//	return err
	//}
	//org2AdminIdentity, err := org2MspClient.GetSigningIdentity("admin")
	//if nil != err {
	//	return err
	//}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
		ChannelConfig: reader,
		//SigningIdentities: []msp.SigningIdentity{org1AdminIdentity, org2AdminIdentity},
	}
	//configOrg1Signature, err := resmgmtClient.CreateConfigSignature(org1AdminIdentity, channelUpdateFilePath)
	//configOrg2Signature, err := resmgmtClient.CreateConfigSignature(org2AdminIdentity, channelUpdateFilePath)
	//clientContext := sdk.Context(fabsdk.WithUser("admin"), fabsdk.WithOrg("org1"))
	//ctx, err := clientContext()
	//if nil != err {
	//	return err
	//}
	//configOrg1Signature, err := resource.SignChannelConfig(ctx, envelopeBytes, org1AdminIdentity)
	//configOrg2Signature, err := resource.SignChannelConfig(ctx, envelopeBytes, org2AdminIdentity)
	//, resmgmt.WithConfigSignatures(configOrg1Signature, configOrg2Signature)
	_, err = resmgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(orderURL)) // "orderer.example.com"
	return err
}

func createConfigEnvelopeReader(configUpdateEnvBytes []byte, channelID string) (*common.Envelope, error) {
	var (
		payloadBytes, payloadChannelHeaderBytes []byte
		err                                     error
	)
	envelope := &common.Envelope{}

	payload := &common.Payload{}
	//if err = proto.Unmarshal(envelope.Payload, payload); nil != err {
	//	return nil, err
	//}

	payload.Data = configUpdateEnvBytes

	payloadChannelHeader := &common.ChannelHeader{}
	payloadChannelHeader.ChannelId = channelID
	payloadChannelHeader.Type = 2
	if payloadChannelHeaderBytes, err = proto.Marshal(payloadChannelHeader); nil != err {
		return nil, err
	}

	header := &common.Header{}
	header.ChannelHeader = payloadChannelHeaderBytes

	payload.Header = header

	if payloadBytes, err = proto.Marshal(payload); nil != err {
		return nil, err
	}

	envelope.Payload = payloadBytes
	return envelope, nil
}
