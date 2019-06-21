package sdk

import (
	"encoding/hex"
	pb "github.com/ennoo/fabric-go-client/grpc/proto"
	str "github.com/ennoo/rivet/utils/string"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/core/common/ccprovider"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/hyperledger/fabric/core/ledger/util"
	com "github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"golang.org/x/protobuf/proto"
	"strconv"
	"strings"
)

func parseBlock(commonBlock *common.Block) (*pb.Block, error) {
	var (
		envelopeCount int32
		txAllCount    int32
		err           error
	)
	envelopeCount = int32(len(commonBlock.Data.Data))
	envelopes := make([]*pb.Envelope, envelopeCount)

	metadata := make([]string, len(commonBlock.Metadata.Metadata))
	for index, md := range commonBlock.Metadata.Metadata {
		metadata[index] = string(md)
	}
	flags := util.TxValidationFlags(metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])

	for index, d := range commonBlock.Data.Data {
		var (
			envelope                 *com.Envelope
			payload                  *com.Payload
			channelHeader            *com.ChannelHeader
			signatureHeader          *com.SignatureHeader
			transaction              *peer.Transaction
			chainCodeHeaderExtension *peer.ChaincodeHeaderExtension
			envelopeInfo             *pb.Envelope
			txCount                  int32
			rwCount                  int32
		)
		if envelope, err = utils.GetEnvelopeFromBlock(d); nil != err {
			goto ERR
		}
		if payload, err = utils.GetPayload(envelope); nil != err {
			goto ERR
		}
		if channelHeader, err = utils.UnmarshalChannelHeader(payload.Header.ChannelHeader); nil != err {
			goto ERR
		}
		if signatureHeader, err = utils.GetSignatureHeader(payload.Header.SignatureHeader); nil != err {
			goto ERR
		}
		if transaction, err = utils.GetTransaction(payload.Data); nil != err {
			goto ERR
		}
		if chainCodeHeaderExtension, err = utils.GetChaincodeHeaderExtension(payload.Header); nil != err {
			goto ERR
		}
		//signedData, _ := protoutil.EnvelopeAsSignedData(envelopes)

		isValid := !flags.IsInvalid(index)
		if envelopeInfo, txCount, rwCount, err = parseEnvelope(channelHeader, signatureHeader, transaction); nil != err {
			goto ERR
		}
		txAllCount += txCount
		envelopeInfo.Signature = hex.EncodeToString(envelope.Signature)
		if nil != chainCodeHeaderExtension.ChaincodeId && nil != chainCodeHeaderExtension.PayloadVisibility {
			envelopeInfo.ChainCode = parseChainCode(chainCodeHeaderExtension)
		}
		envelopeInfo.IsValid = isValid

		envelopes[index] = envelopeInfo

		block := &pb.Block{
			Header: &pb.BlockHeader{
				BlockNumber:   commonBlock.Header.Number,
				DataHash:      hex.EncodeToString(commonBlock.Header.DataHash),
				PreviousHash:  hex.EncodeToString(commonBlock.Header.PreviousHash),
				EnvelopeCount: envelopeCount,
				TxCount:       txAllCount,
				RwCount:       rwCount,
			},
			Envelopes: envelopes,
			//Metadata: &pb.BlockMetadata{
			//	Metadata: metadata,
			//},
		}
		return block, nil
	}
ERR:
	return nil, err
}

func parseEnvelope(channelHeader *com.ChannelHeader, signatureHeader *com.SignatureHeader,
	transaction *peer.Transaction) (envelope *pb.Envelope, actionCount int32, rwCount int32, err error) {
	var (
		transactionActionInfoArray []*pb.Action
		identity                   *msp.SerializedIdentity
	)
	if transactionActionInfoArray, actionCount, rwCount, err = parseTransactionActionInfoArray(transaction); nil != err {
		return
	}
	if identity, err = serializedIdentity(signatureHeader.Creator); nil != err {
		return
	}
	envelope = &pb.Envelope{
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
		CreateID:      string(identity.IdBytes),
		MspID:         identity.Mspid,
		Nonce:         hex.EncodeToString(signatureHeader.Nonce),
		TransactionEnvelopeInfo: &pb.Transaction{
			TxCount: actionCount,
		},
	}
	if nil != transactionActionInfoArray {
		envelope.TransactionEnvelopeInfo.TransactionActionInfoArray = transactionActionInfoArray
	}
	return
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

func parseTransactionActionInfoArray(transaction *peer.Transaction) (transactionActionInfoArray []*pb.Action, actionCount int32, rwCount int32, err error) {
	actionCount = int32(len(transaction.Actions))
	transactionActionInfoArray = make([]*pb.Action, actionCount)
	for index, a := range transaction.Actions {
		var (
			chainCodeActionPayload   *peer.ChaincodeActionPayload
			proposalResponsePayload  *peer.ProposalResponsePayload
			pbChainCodeActionPayload *pb.ChainCodeActionPayload
			rwc                      int32
		)
		if chainCodeActionPayload, err = utils.GetChaincodeActionPayload(a.Payload); nil != err {
			goto ERR
		}
		if proposalResponsePayload, err =
			utils.GetProposalResponsePayload(chainCodeActionPayload.Action.ProposalResponsePayload); nil == err {
			if pbChainCodeActionPayload, rwc, err = parseChainCodeActionPayload(chainCodeActionPayload, proposalResponsePayload); nil != err {
				goto ERR
			}
			rwCount += rwc
			action := &pb.Action{
				ChainCodeActionPayload: pbChainCodeActionPayload,
			}
			transactionActionInfoArray[index] = action
		} else {
			if actionCount == 1 {
				transactionActionInfoArray = nil
				break
			}
		}
		return
	}
ERR:
	return
}

func parseChainCodeActionPayload(payload *peer.ChaincodeActionPayload,
	prPayload *peer.ProposalResponsePayload) (*pb.ChainCodeActionPayload, int32, error) {
	es := payload.Action.Endorsements
	endorsements := make([]*pb.Endorsement, len(es))
	var (
		identity                 *msp.SerializedIdentity
		chainCodeProposalPayload *pb.ChainCodeInvocationSpec
		chainCodeAction          *peer.ChaincodeAction
		pbChainCodeAction        *pb.ChainCodeAction
		rwCount                  int32
		err                      error
	)
	for index, e := range es {
		if identity, err = serializedIdentity(e.Endorser); nil != err {
			return nil, 0, err
		}
		endorsements[index] = &pb.Endorsement{
			Signature: hex.EncodeToString(e.Signature),
			CreateID:  string(identity.IdBytes),
			MspID:     identity.Mspid,
		}
	}
	if chainCodeProposalPayload, err = chainCodeInvocationSpec(payload); nil != err {
		return nil, 0, err
	}
	if chainCodeAction, err = utils.GetChaincodeAction(prPayload.Extension); nil != err {
		return nil, 0, err
	}
	if pbChainCodeAction, rwCount, err = parseChainCodeAction(chainCodeAction); nil != err {
		return nil, 0, err
	}
	return &pb.ChainCodeActionPayload{
		ChainCodeProposalPayload: chainCodeProposalPayload,
		ChainCodeEndorsedAction: &pb.ChainCodeEndorsedAction{
			ProposalResponsePayload: &pb.ProposalResponsePayload{
				ProposalHash:    hex.EncodeToString(prPayload.ProposalHash),
				ChainCodeAction: pbChainCodeAction,
			},
			Endorsements: endorsements,
		},
	}, rwCount, nil
}

func parseChainCodeAction(chainCode *peer.ChaincodeAction) (*pb.ChainCodeAction, int32, error) {
	var (
		event   *peer.ChaincodeEvent
		rwCount int
		err     error
	)
	if event, err = utils.GetChaincodeEvents(chainCode.Events); err != nil {
		return nil, 0, err
	}
	txRWSet := &rwsetutil.TxRwSet{}
	if err := txRWSet.FromProtoBytes(chainCode.Results); err != nil {
		return nil, 0, err
	}
	nsRwSets := make([]*pb.NsRwSets, len(txRWSet.NsRwSets))
	for index, ns := range txRWSet.NsRwSets {
		readCount := len(ns.KvRwSet.Reads)
		writeCount := len(ns.KvRwSet.Writes)
		rwCount += readCount
		rwCount += writeCount

		reads := make([]*pb.KVRead, readCount)
		for i, r := range ns.KvRwSet.Reads {
			var version *pb.Version
			if nil != r.Version {
				version = &pb.Version{
					BlockNum: r.Version.BlockNum,
					TxNum:    r.Version.TxNum,
				}
			}
			reads[i] = &pb.KVRead{
				Key:     r.Key,
				Version: version,
			}
		}

		writes := make([]*pb.KVWrite, writeCount)
		for i, w := range ns.KvRwSet.Writes {
			writes[i] = &pb.KVWrite{
				Key:      w.Key,
				IsDelete: w.IsDelete,
			}
			if nil != w.Value {
				if ccv, err := payloadAnalysis(w.Value); nil == err {
					writes[i].Data = &pb.KVWrite_CCValue{
						CCValue: ccv,
					}
				} else {
					writes[i].Data = &pb.KVWrite_Value{Value: string(w.Value)}
				}
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
	pbTxRwSet := &pb.TxRwSet{
		RwCount:  int32(rwCount),
		NsRwSets: nsRwSets,
	}

	chainCodeAction := &pb.ChainCodeAction{
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
		},
	}

	if ccv, err := payloadAnalysis(chainCode.Response.Payload); nil == err {
		chainCodeAction.Response.Data = &pb.Response_CCValue{
			CCValue: ccv,
		}
	} else {
		chainCodeAction.Response.Data = &pb.Response_Value{Value: string(chainCode.Response.Payload)}
	}
	return chainCodeAction, int32(rwCount), nil
}

func payloadAnalysis(buf []byte) (ccv *pb.ChainCodeValue, err error) {
	instantiationPolicy := strings.Split(str.SingleSpace(string(buf)), "\n")
	for index := range instantiationPolicy {
		instantiationPolicy[index] = strconv.QuoteToASCII(instantiationPolicy[index])
	}
	cdRWSet := &ccprovider.ChaincodeData{}
	if err = proto.Unmarshal(buf, cdRWSet); nil == err {
		otherData := &ccprovider.CDSData{}
		err := proto.Unmarshal(cdRWSet.Data, otherData)
		var (
			codeHash     = ""
			metaDataHash = ""
		)
		if err == nil {
			codeHash = hex.EncodeToString(otherData.CodeHash)
			metaDataHash = hex.EncodeToString(otherData.MetaDataHash)
		}
		ccv = &pb.ChainCodeValue{
			Name:    cdRWSet.Name,
			Version: cdRWSet.Version,
			Escc:    cdRWSet.Escc,
			Vscc:    cdRWSet.Vscc,
			Policy:  string(cdRWSet.Policy),
			Data: &pb.ChainCodeValueDate{
				CodeHash:     codeHash,
				MetaDataHash: metaDataHash,
			},
			Id:                  hex.EncodeToString(cdRWSet.Id),
			InstantiationPolicy: instantiationPolicy,
		}
	}
	return
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
