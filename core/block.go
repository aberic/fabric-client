package sdk

import (
	"encoding/hex"
	pb "github.com/ennoo/fabric-go-client/grpc/proto"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/hyperledger/fabric/core/ledger/util"
	com "github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"golang.org/x/protobuf/proto"
	"strconv"
)

func parseBlock(commonBlock *common.Block) *pb.Block {
	envelopeSize := len(commonBlock.Data.Data)
	strInt := strconv.Itoa(envelopeSize)
	envelopeCount, _ := strconv.ParseInt(strInt, 10, 64)
	envelopes := make([]*pb.Envelope, envelopeCount)

	metadata := make([]string, len(commonBlock.Metadata.Metadata))
	for index, md := range commonBlock.Metadata.Metadata {
		metadata[index] = string(md)
	}
	flags := util.TxValidationFlags(metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])

	for index, d := range commonBlock.Data.Data {
		envelope, _ := utils.GetEnvelopeFromBlock(d)
		payload, _ := utils.GetPayload(envelope)
		channelHeader, _ := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		signatureHeader, _ := utils.GetSignatureHeader(payload.Header.SignatureHeader)
		transaction, _ := utils.GetTransaction(payload.Data)
		chainCodeHeaderExtension, _ := utils.GetChaincodeHeaderExtension(payload.Header)
		//signedData, _ := protoutil.EnvelopeAsSignedData(envelopes)

		isValid := !flags.IsInvalid(index)
		envelopeInfo := parseEnvelope(channelHeader, signatureHeader, transaction)
		envelopeInfo.Signature = hex.EncodeToString(envelope.Signature)
		envelopeInfo.ChainCode = parseChainCode(chainCodeHeaderExtension)
		envelopeInfo.IsValid = isValid

		envelopes[index] = envelopeInfo
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
		Type:      headerType(channelHeader.Type),
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

func headerType(ht int32) string {
	switch ht {
	case 0: // Used for messages which are signed but opaque
		return "MESSAGE"
	case 1:
		return "CONFIG" // Used for messages which express the channel config
	case 2:
		return "CONFIG_UPDATE" // Used for transactions which update the channel config
	case 3:
		return "ENDORSER_TRANSACTION" // Used by the SDK to submit endorser based transactions
	case 4:
		return "ORDERER_TRANSACTION" // Used internally by the orderer for management
	case 5:
		return "DELIVER_SEEK_INFO" // Used as the type for Envelope messages submitted to instruct the Deliver API to seek
	case 6:
		return "CHAINCODE_PACKAGE" // Used for packaging chaincode artifacts for install
	case 8:
		return "PEER_ADMIN_OPERATION" // Used for invoking an administrative operation on a peer
	default:
		return "TOKEN_TRANSACTION" // Used to denote transactions that invoke token management operations
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
	strInt := strconv.Itoa(actionSize)
	actionCount, _ := strconv.ParseInt(strInt, 10, 64)

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
	rwCount := 0
	nsRwSets := make([]*pb.NsRwSets, len(txRWSet.NsRwSets))
	for index, ns := range txRWSet.NsRwSets {
		readCount := len(ns.KvRwSet.Reads)
		writeCount := len(ns.KvRwSet.Writes)
		rwCount += readCount
		rwCount += writeCount

		reads := make([]*pb.KVRead, readCount)
		for i, r := range ns.KvRwSet.Reads {
			reads[i] = &pb.KVRead{
				Key: r.Key,
				Version: &pb.Version{
					BlockNum: r.Version.BlockNum,
					TxNum:    r.Version.TxNum,
				},
			}
		}

		writes := make([]*pb.KVWrite, writeCount)
		for i, w := range ns.KvRwSet.Writes {
			writes[i] = &pb.KVWrite{
				Key:      w.Key,
				IsDelete: w.IsDelete,
				Value:    string(w.Value),
			}
		}

		nsRwSets[index] = &pb.NsRwSets{
			NameSpace: ns.NameSpace,
			KVRWSet: &pb.KVRWSet{
				Reads:            reads,
				Writes:           writes,
				RangeQueriesInfo: ns.KvRwSet.RangeQueriesInfo,
				MetadataWrites:   ns.KvRwSet.MetadataWrites,
			},
		}
	}
	strInt := strconv.Itoa(rwCount)
	rwCount64, _ := strconv.ParseInt(strInt, 10, 64)
	pbTxRwSet := &pb.TxRwSet{
		RwCount:  rwCount64,
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

func chainCodeInvocationSpec(payload *peer.ChaincodeActionPayload) (*pb.ChainCodeInvocationSpec, error) {
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
	args := make([]string, len(cis.ChaincodeSpec.Input.Args))
	for index, arg := range cis.ChaincodeSpec.Input.Args {
		args[index] = string(arg)
	}
	return &pb.ChainCodeInvocationSpec{
		ChainCodeSpec: &pb.ChainCodeSpec{
			Type: ccType(cis.ChaincodeSpec.Type),
			ChainCodeID: &pb.ChainCodeID{
				Name:    cis.ChaincodeSpec.ChaincodeId.Name,
				Path:    cis.ChaincodeSpec.ChaincodeId.Path,
				Version: cis.ChaincodeSpec.ChaincodeId.Version,
			},
			Input: &pb.ChainCodeInput{
				Args: args,
			},
			Timeout: cis.ChaincodeSpec.Timeout,
		},
	}, nil
}

func ccType(tp peer.ChaincodeSpec_Type) string {
	switch tp {
	case peer.ChaincodeSpec_GOLANG:
		return "GOLANG"
	case peer.ChaincodeSpec_NODE:
		return "NODE"
	case peer.ChaincodeSpec_CAR:
		return "CAR"
	case peer.ChaincodeSpec_JAVA:
		return "JAVA"
	default:
		return "UNDEFINED"
	}
}
