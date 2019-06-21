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
		envelopeCount int64
		err           error
	)
	envelopeSize := len(commonBlock.Data.Data)
	strInt := strconv.Itoa(envelopeSize)
	if envelopeCount, err = strconv.ParseInt(strInt, 10, 64); nil != err {
		goto ERR
	} else {
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
			if envelopeInfo, err = parseEnvelope(channelHeader, signatureHeader, transaction); nil != err {
				goto ERR
			}
			envelopeInfo.Signature = hex.EncodeToString(envelope.Signature)
			if nil != chainCodeHeaderExtension.ChaincodeId && nil != chainCodeHeaderExtension.PayloadVisibility {
				envelopeInfo.ChainCode = parseChainCode(chainCodeHeaderExtension)
			}
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
		return block, nil
	}
ERR:
	return nil, err
}

func parseEnvelope(channelHeader *com.ChannelHeader, signatureHeader *com.SignatureHeader,
	transaction *peer.Transaction) (*pb.Envelope, error) {
	var (
		transactionActionInfoArray []*pb.Action
		actionCount                int64
		identity                   *msp.SerializedIdentity
		err                        error
	)
	if transactionActionInfoArray, actionCount, err = parseTransactionActionInfoArray(transaction); nil != err {
		return nil, err
	}
	if identity, err = serializedIdentity(signatureHeader.Creator); nil != err {
		return nil, err
	}
	envelope := &pb.Envelope{
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
	return envelope, nil
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

func parseTransactionActionInfoArray(transaction *peer.Transaction) ([]*pb.Action, int64, error) {
	var (
		actionCount int64
		err         error
	)
	actionSize := len(transaction.Actions)
	strInt := strconv.Itoa(actionSize)
	if actionCount, err = strconv.ParseInt(strInt, 10, 64); nil != err {
		goto ERR
	} else {
		transactionActionInfoArray := make([]*pb.Action, actionCount)
		for index, a := range transaction.Actions {
			var (
				chainCodeActionPayload   *peer.ChaincodeActionPayload
				proposalResponsePayload  *peer.ProposalResponsePayload
				pbChainCodeActionPayload *pb.ChainCodeActionPayload
			)
			if chainCodeActionPayload, err = utils.GetChaincodeActionPayload(a.Payload); nil != err {
				goto ERR
			}
			if proposalResponsePayload, err =
				utils.GetProposalResponsePayload(chainCodeActionPayload.Action.ProposalResponsePayload); nil == err {
				if pbChainCodeActionPayload, err = parseChainCodeActionPayload(chainCodeActionPayload, proposalResponsePayload); nil != err {
					goto ERR
				}
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
		}
		return transactionActionInfoArray, actionCount, nil
	}
ERR:
	return nil, 0, err
}

func parseChainCodeActionPayload(payload *peer.ChaincodeActionPayload,
	prPayload *peer.ProposalResponsePayload) (*pb.ChainCodeActionPayload, error) {
	es := payload.Action.Endorsements
	endorsements := make([]*pb.Endorsement, len(es))
	var (
		identity                 *msp.SerializedIdentity
		chainCodeProposalPayload *pb.ChainCodeInvocationSpec
		chainCodeAction          *peer.ChaincodeAction
		pbChainCodeAction        *pb.ChainCodeAction
		err                      error
	)
	for index, e := range es {
		if identity, err = serializedIdentity(e.Endorser); nil != err {
			return nil, err
		}
		endorsements[index] = &pb.Endorsement{
			Signature: hex.EncodeToString(e.Signature),
			CreateID:  string(identity.IdBytes),
			MspID:     identity.Mspid,
		}
	}
	if chainCodeProposalPayload, err = chainCodeInvocationSpec(payload); nil != err {
		return nil, err
	}
	if chainCodeAction, err = utils.GetChaincodeAction(prPayload.Extension); nil != err {
		return nil, err
	}
	if pbChainCodeAction, err = parseChainCodeAction(chainCodeAction); nil != err {
		return nil, err
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
	}, nil
}

func parseChainCodeAction(chainCode *peer.ChaincodeAction) (*pb.ChainCodeAction, error) {
	var (
		event     *peer.ChaincodeEvent
		rwCount64 int64
		err       error
	)
	if event, err = utils.GetChaincodeEvents(chainCode.Events); err != nil {
		return nil, err
	}
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
	strInt := strconv.Itoa(rwCount)
	if rwCount64, err = strconv.ParseInt(strInt, 10, 64); err != nil {
		return nil, err
	}
	pbTxRwSet := &pb.TxRwSet{
		RwCount:  rwCount64,
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
	return chainCodeAction, nil
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
