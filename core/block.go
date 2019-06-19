package sdk

import (
	"encoding/hex"
	pb "github.com/ennoo/fabric-go-client/grpc/proto"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	com "github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"golang.org/x/protobuf/proto"
	"strconv"
)

func parseBlock(commonBlock *common.Block) *pb.Block {
	envelopeSize := len(commonBlock.Data.Data)
	strInt32 := strconv.Itoa(envelopeSize)
	envelopeCount, _ := strconv.ParseInt(strInt32, 10, 64)
	envelopes := make([]*pb.Envelope, envelopeCount)

	for index, d := range commonBlock.Data.Data {
		envelope, _ := utils.GetEnvelopeFromBlock(d)
		payload, _ := utils.GetPayload(envelope)
		channelHeader, _ := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		signatureHeader, _ := utils.GetSignatureHeader(payload.Header.SignatureHeader)
		transaction, _ := utils.GetTransaction(payload.Data)
		chainCodeHeaderExtension, _ := utils.GetChaincodeHeaderExtension(payload.Header)
		//signedData, _ := protoutil.EnvelopeAsSignedData(envelopes)

		envelopeInfo := parseEnvelope(channelHeader, signatureHeader, transaction)
		envelopeInfo.Signature = hex.EncodeToString(envelope.Signature)
		envelopeInfo.ChainCode = parseChainCode(chainCodeHeaderExtension)

		envelopes[index] = envelopeInfo
	}

	metadata := make([]string, len(commonBlock.Metadata.Metadata))
	for index, md := range commonBlock.Metadata.Metadata {
		metadata[index] = string(md)
	}

	block := &pb.Block{
		Header: &pb.BlockHeader{
			BlockNumber:   commonBlock.Header.Number,
			DataHash:      hex.EncodeToString(commonBlock.Header.DataHash),
			PreviousHash:  hex.EncodeToString(commonBlock.Header.PreviousHash),
			EnvelopeCount: envelopeCount,
		},
		Envelopes: envelopes,
		//Metadata: &pb.BlockMetadata{
		//	Metadata: metadata,
		//},
	}
	return block
}

func parseEnvelope(channelHeader *com.ChannelHeader, signatureHeader *com.SignatureHeader, transaction *peer.Transaction) *pb.Envelope {
	transactionActionInfoArray, actionCount := parseTransactionActionInfoArray(transaction)
	serializedIdentity, _ := serializedIdentity(signatureHeader.Creator)
	return &pb.Envelope{
		ChannelID: channelHeader.ChannelId,
		Type:      channelHeader.Type,
		Version:   channelHeader.Version,
		Timestamp: &pb.Timestamp{
			Seconds: channelHeader.Timestamp.Seconds,
			Nanos:   channelHeader.Timestamp.Nanos,
		},
		TransactionID: channelHeader.TxId,
		Epoch:         channelHeader.Epoch,
		Extension:     string(channelHeader.Extension),
		TlsCertHash:   string(channelHeader.TlsCertHash),
		CreateID:      string(serializedIdentity.IdBytes),
		MspID:         serializedIdentity.Mspid,
		Nonce:         hex.EncodeToString(signatureHeader.Nonce),
		TransactionEnvelopeInfo: &pb.Transaction{
			TxCount:                    actionCount,
			TransactionActionInfoArray: transactionActionInfoArray,
		},
	}
}

func parseChainCode(chainCode *peer.ChaincodeHeaderExtension) *pb.ChainCodeHeaderExtension {
	return &pb.ChainCodeHeaderExtension{
		PayloadVisibility: string(chainCode.PayloadVisibility),
		ChainCodeID: &pb.ChainCodeID{
			Path:    chainCode.ChaincodeId.Path,
			Name:    chainCode.ChaincodeId.Name,
			Version: chainCode.ChaincodeId.Version,
		},
	}
}

func parseTransactionActionInfoArray(transaction *peer.Transaction) ([]*pb.Action, int64) {
	actionSize := len(transaction.Actions)
	strInt32 := strconv.Itoa(actionSize)
	actionCount, _ := strconv.ParseInt(strInt32, 10, 64)

	transactionActionInfoArray := make([]*pb.Action, actionCount)
	for index, a := range transaction.Actions {
		chainCodeActionPayload, _ := utils.GetChaincodeActionPayload(a.Payload)
		proposalResponsePayload, _ := utils.GetProposalResponsePayload(chainCodeActionPayload.Action.ProposalResponsePayload)
		action := &pb.Action{
			ChainCodeActionPayload: parseChainCodeActionPayload(chainCodeActionPayload, proposalResponsePayload),
		}
		transactionActionInfoArray[index] = action
	}
	return transactionActionInfoArray, actionCount
}

func parseChainCodeActionPayload(payload *peer.ChaincodeActionPayload, prPayload *peer.ProposalResponsePayload) *pb.ChainCodeActionPayload {
	es := payload.Action.Endorsements
	endorsements := make([]*pb.Endorsement, len(es))
	for index, e := range es {
		serializedIdentity, _ := serializedIdentity(e.Endorser)
		endorsements[index] = &pb.Endorsement{
			Signature: hex.EncodeToString(e.Signature),
			CreateID:  string(serializedIdentity.IdBytes),
			MspID:     serializedIdentity.Mspid,
		}
	}

	chainCodeProposalPayload, _ := chainCodeInvocationSpec(payload)
	chainCodeAction, _ := utils.GetChaincodeAction(prPayload.Extension)
	pbChainCodeAction, _ := parseChainCodeAction(chainCodeAction)
	return &pb.ChainCodeActionPayload{
		ChainCodeProposalPayload: chainCodeProposalPayload,
		ChainCodeEndorsedAction: &pb.ChainCodeEndorsedAction{
			ProposalResponsePayload: &pb.ProposalResponsePayload{
				ProposalHash:    hex.EncodeToString(prPayload.ProposalHash),
				ChainCodeAction: pbChainCodeAction,
			},
			Endorsements: endorsements,
		},
	}
}

func parseChainCodeAction(chainCode *peer.ChaincodeAction) (*pb.ChainCodeAction, error) {
	event, _ := utils.GetChaincodeEvents(chainCode.Events)
	txRWSet := &rwsetutil.TxRwSet{}
	if err := txRWSet.FromProtoBytes(chainCode.Results); err != nil {
		return nil, err
	}
	nsRwSets := make([]*pb.NsRwSets, len(txRWSet.NsRwSets))
	for index, ns := range txRWSet.NsRwSets {
		nsRwSets[index] = &pb.NsRwSets{
			NameSpace: ns.NameSpace,
			KVRWSet:   ns.KvRwSet,
		}
	}
	pbTxRwSet := &pb.TxRwSet{
		NsRwSets: nsRwSets,
	}
	return &pb.ChainCodeAction{
		ChainCodeID: &pb.ChainCodeID{
			Path:    chainCode.ChaincodeId.Path,
			Name:    chainCode.ChaincodeId.Name,
			Version: chainCode.ChaincodeId.Version,
		},
		TxRwSet: pbTxRwSet,
		Event: &pb.ChainCodeEvent{
			ChainCodeID:   event.ChaincodeId,
			TransactionID: event.TxId,
			EventName:     event.EventName,
			Payload:       string(event.Payload),
		},
		Response: &pb.Response{
			Status:  chainCode.Response.Status,
			Message: chainCode.Response.Message,
			Payload: string(chainCode.Response.Payload),
		},
	}, nil
}

func serializedIdentity(bytes []byte) (*msp.SerializedIdentity, error) {
	serializedIdentity := &msp.SerializedIdentity{}
	if err := proto.Unmarshal(bytes, serializedIdentity); nil != err {
		return nil, err
	}
	return serializedIdentity, nil
}

func chainCodeInvocationSpec(payload *peer.ChaincodeActionPayload) (*peer.ChaincodeInvocationSpec, error) {
	var (
		cpp *peer.ChaincodeProposalPayload
		err error
	)
	if cpp, err = utils.GetChaincodeProposalPayload(payload.ChaincodeProposalPayload); err != nil {
		return nil, err
	}
	cis := &peer.ChaincodeInvocationSpec{}
	if err := proto.Unmarshal(cpp.Input, cis); err != nil {
		return nil, err
	}
	return cis, nil
}
