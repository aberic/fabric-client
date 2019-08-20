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

package chains

import (
	"errors"
	"github.com/ennoo/fabric-client/config"
	sdk "github.com/ennoo/fabric-client/core"
	pb "github.com/ennoo/fabric-client/grpc/proto/chain"
	"github.com/ennoo/fabric-client/service"
	str "github.com/ennoo/rivet/utils/string"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"golang.org/x/net/context"
)

type CAServer struct{}

func (ca *CAServer) Enroll(ctx context.Context, req *pb.ReqEnroll) (*pb.Result, error) {
	var (
		conf *config.Config
		err  error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.Result{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if str.IsEmpty(req.Secret) {
		return &pb.Result{Code: pb.Code_Fail, ErrMsg: "secret is nil"}, errors.New("secret is nil")
	}
	if err = sdk.Enroll(req.OrgName, req.EnrollmentID, service.GetBytes(req.ConfigID),
		optionEnroll(req.Secret, req.Type, req.Profile, req.Label)); nil != err {
		return &pb.Result{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &pb.Result{Code: pb.Code_Success}, nil
}

func (ca *CAServer) Reenroll(ctx context.Context, req *pb.ReqReenroll) (*pb.Result, error) {
	var (
		conf *config.Config
		err  error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.Result{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if str.IsEmpty(req.Secret) {
		return &pb.Result{Code: pb.Code_Fail, ErrMsg: "secret is nil"}, errors.New("secret is nil")
	}
	if err = sdk.Reenroll(req.OrgName, req.EnrollmentID, service.GetBytes(req.ConfigID),
		optionEnroll(req.Secret, req.Type, req.Profile, req.Label)); nil != err {
		return &pb.Result{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &pb.Result{Code: pb.Code_Success}, nil
}

func (ca *CAServer) Register(ctx context.Context, req *pb.ReqRegister) (*pb.Result, error) {
	var (
		conf   *config.Config
		result string
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.Result{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	attrs := make([]msp.Attribute, len(req.RegistrationRequest.Attributes))
	for index, attr := range req.RegistrationRequest.Attributes {
		attrs[index] = msp.Attribute{
			Name:  attr.Name,
			Value: attr.Value,
			ECert: attr.ECert,
		}
	}
	if result, err = sdk.Register(req.OrgName, &msp.RegistrationRequest{
		Name:           req.RegistrationRequest.Name,
		Type:           req.RegistrationRequest.Type,                // (e.g. "client, orderer, peer, app, user")
		MaxEnrollments: int(req.RegistrationRequest.MaxEnrollments), // if omitted, this defaults to max_enrollments configured on the server
		Affiliation:    req.RegistrationRequest.Affiliation,         // The identity's affiliation e.g. org1.department1
		Attributes:     attrs,
		CAName:         req.RegistrationRequest.CaName,
		Secret:         req.RegistrationRequest.Secret,
	}, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.Result{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return &pb.Result{Code: pb.Code_Success, Data: result}, nil
}

func (ca *CAServer) AddAffiliation(ctx context.Context, req *pb.ReqAddAffiliation) (*pb.ResultAffiliation, error) {
	var (
		conf   *config.Config
		result *msp.AffiliationResponse
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultAffiliation{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if result, err = sdk.AddAffiliation(req.OrgName, &msp.AffiliationRequest{
		Name:   req.AffiliationRequest.Name,
		CAName: req.AffiliationRequest.CaName,
		Force:  req.AffiliationRequest.Force,
	}, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultAffiliation{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return affiliation(result), nil
}

func (ca *CAServer) RemoveAffiliation(ctx context.Context, req *pb.ReqRemoveAffiliation) (*pb.ResultAffiliation, error) {
	var (
		conf   *config.Config
		result *msp.AffiliationResponse
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultAffiliation{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if result, err = sdk.RemoveAffiliation(req.OrgName, &msp.AffiliationRequest{
		Name:   req.AffiliationRequest.Name,
		CAName: req.AffiliationRequest.CaName,
		Force:  req.AffiliationRequest.Force,
	}, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultAffiliation{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return affiliation(result), nil
}

func (ca *CAServer) ModifyAffiliation(ctx context.Context, req *pb.ReqModifyAffiliation) (*pb.ResultAffiliation, error) {
	var (
		conf   *config.Config
		result *msp.AffiliationResponse
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultAffiliation{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if result, err = sdk.ModifyAffiliation(req.OrgName, &msp.ModifyAffiliationRequest{
		NewName: req.ModifyAffiliationRequest.NewName,
		AffiliationRequest: msp.AffiliationRequest{
			Name:   req.ModifyAffiliationRequest.AffiliationRequest.Name,
			CAName: req.ModifyAffiliationRequest.AffiliationRequest.CaName,
			Force:  req.ModifyAffiliationRequest.AffiliationRequest.Force,
		},
	}, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultAffiliation{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return affiliation(result), nil
}

func (ca *CAServer) GetAllAffiliations(ctx context.Context, req *pb.ReqGetAllAffiliations) (*pb.ResultAffiliation, error) {
	var (
		conf   *config.Config
		result *msp.AffiliationResponse
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultAffiliation{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if result, err = sdk.GetAllAffiliations(req.OrgName, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultAffiliation{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return affiliation(result), nil
}

func (ca *CAServer) GetAllAffiliationsByCaName(ctx context.Context, req *pb.ReqGetAllAffiliationsByCaName) (*pb.ResultAffiliation, error) {
	var (
		conf   *config.Config
		result *msp.AffiliationResponse
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultAffiliation{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if result, err = sdk.GetAllAffiliationsByCaName(req.OrgName, req.CaName, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultAffiliation{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return affiliation(result), nil
}

func (ca *CAServer) GetAffiliation(ctx context.Context, req *pb.ReqGetAffiliation) (*pb.ResultAffiliation, error) {
	var (
		conf   *config.Config
		result *msp.AffiliationResponse
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultAffiliation{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if result, err = sdk.GetAffiliation(req.Affiliation, req.OrgName, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultAffiliation{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return affiliation(result), nil
}

func (ca *CAServer) GetAffiliationByCaName(ctx context.Context, req *pb.ReqGetAffiliationByCaName) (*pb.ResultAffiliation, error) {
	var (
		conf   *config.Config
		result *msp.AffiliationResponse
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultAffiliation{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if result, err = sdk.GetAffiliationByCaName(req.Affiliation, req.OrgName, req.CaName, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultAffiliation{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return affiliation(result), nil
}

func (ca *CAServer) GetAllIdentities(ctx context.Context, req *pb.ReqGetAllIdentities) (*pb.ResultIdentityResponses, error) {
	var (
		conf   *config.Config
		result []*msp.IdentityResponse
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultIdentityResponses{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if result, err = sdk.GetAllIdentities(req.OrgName, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultIdentityResponses{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return allIdentities(result), nil
}

func (ca *CAServer) GetAllIdentitiesByCaName(ctx context.Context, req *pb.ReqGetAllIdentitiesByCaName) (*pb.ResultIdentityResponses, error) {
	var (
		conf   *config.Config
		result []*msp.IdentityResponse
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultIdentityResponses{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if result, err = sdk.GetAllIdentitiesByCaName(req.OrgName, req.CaName, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultIdentityResponses{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return allIdentities(result), nil
}

func (ca *CAServer) CreateIdentity(ctx context.Context, req *pb.ReqCreateIdentity) (*pb.ResultIdentityResponse, error) {
	var (
		conf   *config.Config
		result *msp.IdentityResponse
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultIdentityResponse{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	attrs := make([]msp.Attribute, len(req.IdentityRequest.Attributes))
	for _, attr := range req.IdentityRequest.Attributes {
		attrs = append(attrs, msp.Attribute{
			Name:  attr.Name,
			Value: attr.Value,
			ECert: attr.ECert,
		})
	}
	if result, err = sdk.CreateIdentity(req.OrgName, &msp.IdentityRequest{
		ID:             req.IdentityRequest.Id,          // The enrollment ID which uniquely identifies an identity (required)
		Affiliation:    req.IdentityRequest.Affiliation, // The identity's affiliation e.g. org1.department1
		Attributes:     attrs,
		Type:           req.IdentityRequest.Type,                // (e.g. "client, orderer, peer, app, user")
		MaxEnrollments: int(req.IdentityRequest.MaxEnrollments), // if omitted, this defaults to max_enrollments configured on the server
		Secret:         req.IdentityRequest.Secret,
		CAName:         req.IdentityRequest.CaName,
	}, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultIdentityResponse{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return identity(result), nil
}

func (ca *CAServer) ModifyIdentity(ctx context.Context, req *pb.ReqModifyIdentity) (*pb.ResultIdentityResponse, error) {
	var (
		conf   *config.Config
		result *msp.IdentityResponse
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultIdentityResponse{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	attrs := make([]msp.Attribute, len(req.IdentityRequest.Attributes))
	for _, attr := range req.IdentityRequest.Attributes {
		attrs = append(attrs, msp.Attribute{
			Name:  attr.Name,
			Value: attr.Value,
			ECert: attr.ECert,
		})
	}
	if result, err = sdk.ModifyIdentity(req.OrgName, &msp.IdentityRequest{
		ID:             req.IdentityRequest.Id,          // The enrollment ID which uniquely identifies an identity (required)
		Affiliation:    req.IdentityRequest.Affiliation, // The identity's affiliation e.g. org1.department1
		Attributes:     attrs,
		Type:           req.IdentityRequest.Type,                // (e.g. "client, orderer, peer, app, user")
		MaxEnrollments: int(req.IdentityRequest.MaxEnrollments), // if omitted, this defaults to max_enrollments configured on the server
		Secret:         req.IdentityRequest.Secret,
		CAName:         req.IdentityRequest.CaName,
	}, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultIdentityResponse{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return identity(result), nil
}

func (ca *CAServer) GetIdentity(ctx context.Context, req *pb.ReqGetIdentity) (*pb.ResultIdentityResponse, error) {
	var (
		conf   *config.Config
		result *msp.IdentityResponse
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultIdentityResponse{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if result, err = sdk.GetIdentity(req.Id, req.OrgName, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultIdentityResponse{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return identity(result), nil
}

func (ca *CAServer) GetIdentityByCaName(ctx context.Context, req *pb.ReqGetIdentityByCaName) (*pb.ResultIdentityResponse, error) {
	var (
		conf   *config.Config
		result *msp.IdentityResponse
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultIdentityResponse{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if result, err = sdk.GetIdentityByCaName(req.Id, req.CaName, req.OrgName, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultIdentityResponse{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return identity(result), nil
}

func (ca *CAServer) RemoveIdentity(ctx context.Context, req *pb.ReqRemoveIdentity) (*pb.ResultIdentityResponse, error) {
	var (
		conf   *config.Config
		result *msp.IdentityResponse
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultIdentityResponse{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if result, err = sdk.RemoveIdentity(req.OrgName, &msp.RemoveIdentityRequest{
		ID:     req.RemoveIdentityRequest.Id, // The enrollment ID which uniquely identifies an identity (required)
		Force:  req.RemoveIdentityRequest.Force,
		CAName: req.RemoveIdentityRequest.CaName,
	}, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultIdentityResponse{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return identity(result), nil
}

func (ca *CAServer) CreateSigningIdentity(ctx context.Context, req *pb.ReqCreateSigningIdentity) (*pb.ResultSigningIdentityResponse, error) {
	var (
		conf   *config.Config
		result mspctx.SigningIdentity
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultSigningIdentityResponse{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	var opts []mspctx.SigningIdentityOption
	if nil != req.PrivateKey {
		opts = append(opts, mspctx.WithPrivateKey(req.PrivateKey))
	}
	if nil != req.Cert {
		opts = append(opts, mspctx.WithCert(req.Cert))
	}
	if result, err = sdk.CreateSigningIdentity(req.OrgName, service.GetBytes(req.ConfigID), opts); nil != err {
		return &pb.ResultSigningIdentityResponse{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return signingIdentity(result), nil
}

func (ca *CAServer) GetSigningIdentity(ctx context.Context, req *pb.ReqGetSigningIdentity) (*pb.ResultSigningIdentityResponse, error) {
	var (
		conf   *config.Config
		result mspctx.SigningIdentity
		err    error
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultSigningIdentityResponse{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if result, err = sdk.GetSigningIdentity(req.Id, req.OrgName, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultSigningIdentityResponse{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	return signingIdentity(result), nil
}

func (ca *CAServer) Revoke(ctx context.Context, req *pb.ReqRevoke) (*pb.ResultRevocationResponse, error) {
	var (
		conf         *config.Config
		result       *msp.RevocationResponse
		err          error
		revokedCerts []*pb.RevokedCert
	)
	if conf = service.Configs[req.ConfigID]; nil == conf {
		return &pb.ResultRevocationResponse{Code: pb.Code_Fail, ErrMsg: "config client is not exist"}, errors.New("config client is not exist")
	}
	if result, err = sdk.Revoke(req.OrgName, &msp.RevocationRequest{
		Name:   req.RevocationRequest.Name,
		Serial: req.RevocationRequest.Serial,
		AKI:    req.RevocationRequest.Aki,
		Reason: req.RevocationRequest.Reason,
		CAName: req.RevocationRequest.CaName,
	}, service.GetBytes(req.ConfigID)); nil != err {
		return &pb.ResultRevocationResponse{Code: pb.Code_Fail, ErrMsg: err.Error()}, err
	}
	for _, cert := range result.RevokedCerts {
		revokedCerts = append(revokedCerts, &pb.RevokedCert{
			Serial: cert.Serial,
			Aki:    cert.AKI,
		})
	}
	return &pb.ResultRevocationResponse{Code: pb.Code_Success, Resp: &pb.RevocationResponse{
		RevokedCerts: revokedCerts,
		Crl:          result.CRL,
	}}, nil
}

func optionEnroll(secret, reqType, profile, label string) []msp.EnrollmentOption {
	var eos []msp.EnrollmentOption
	eos = append(eos, msp.WithSecret(secret))
	if str.IsNotEmpty(reqType) {
		eos = append(eos, msp.WithType(reqType))
	}
	if str.IsNotEmpty(profile) {
		eos = append(eos, msp.WithProfile(profile))
	} else {
		eos = append(eos, msp.WithType("x509"))
	}
	if str.IsNotEmpty(label) {
		eos = append(eos, msp.WithLabel(label))
	}
	return eos
}

func affiliation(result *msp.AffiliationResponse) *pb.ResultAffiliation {
	idInfos := make([]*pb.IdentityInfo, len(result.Identities))
	for _, info := range result.Identities {
		attrs := make([]*pb.Attribute, len(info.Attributes))
		for _, attr := range info.Attributes {
			attrs = append(attrs, &pb.Attribute{
				Name:  attr.Name,
				Value: attr.Value,
				ECert: attr.ECert,
			})
		}
		idInfos = append(idInfos, &pb.IdentityInfo{
			Id:             info.ID,
			Type:           info.Type,
			Affiliation:    info.Affiliation,
			MaxEnrollments: int32(info.MaxEnrollments),
			Attributes:     attrs,
		})
	}
	affiInfos := make([]*pb.AffiliationInfo, len(result.Affiliations))
	for _, info := range result.Affiliations {
		affiInfos = append(affiInfos, &pb.AffiliationInfo{
			Name:             info.Name,
			AffiliationInfos: affiInfos,
			IdentityInfos:    idInfos,
		})
	}
	return &pb.ResultAffiliation{Code: pb.Code_Success, Resp: &pb.AffiliationResponse{
		CaName: result.CAName,
		AffiliationInfo: &pb.AffiliationInfo{
			Name:             result.AffiliationInfo.Name,
			AffiliationInfos: affiInfos,
			IdentityInfos:    idInfos,
		},
	}}
}

func allIdentities(result []*msp.IdentityResponse) *pb.ResultIdentityResponses {
	ids := make([]*pb.IdentityResponse, len(result))
	for _, id := range result {
		attrs := make([]*pb.Attribute, len(id.Attributes))
		for _, attr := range id.Attributes {
			attrs = append(attrs, &pb.Attribute{
				Name:  attr.Name,
				Value: attr.Value,
				ECert: attr.ECert,
			})
		}
		ids = append(ids, &pb.IdentityResponse{
			Id:             id.ID,
			Affiliation:    id.Affiliation,
			Type:           id.Type,
			MaxEnrollments: int32(id.MaxEnrollments),
			Secret:         id.Secret,
			CaName:         id.CAName,
			Attributes:     attrs,
		})
	}
	return &pb.ResultIdentityResponses{Code: pb.Code_Success, Resp: ids}
}

func identity(result *msp.IdentityResponse) *pb.ResultIdentityResponse {
	attrs := make([]*pb.Attribute, len(result.Attributes))
	for _, attr := range result.Attributes {
		attrs = append(attrs, &pb.Attribute{
			Name:  attr.Name,
			Value: attr.Value,
			ECert: attr.ECert,
		})
	}
	return &pb.ResultIdentityResponse{Code: pb.Code_Success, Resp: &pb.IdentityResponse{
		Id:             result.ID,
		Affiliation:    result.Affiliation,
		Type:           result.Type,
		MaxEnrollments: int32(result.MaxEnrollments),
		Secret:         result.Secret,
		CaName:         result.CAName,
		Attributes:     attrs,
	}}
}

func signingIdentity(result mspctx.SigningIdentity) *pb.ResultSigningIdentityResponse {
	return &pb.ResultSigningIdentityResponse{Code: pb.Code_Success, Resp: &pb.SigningIdentityResponse{
		Identifier: &pb.Identifier{
			Id:    result.Identifier().ID,
			MspID: result.Identifier().MSPID,
		},
		EnrollmentCertificate: result.EnrollmentCertificate(),
		PrivateKey: &pb.PrivateKey{
			Ski:       result.PrivateKey().SKI(),
			Symmetric: result.PrivateKey().Symmetric(),
			Private:   result.PrivateKey().Private(),
		},
		PublicVersion: &pb.PublicVersion{
			Identifier: &pb.Identifier{
				Id:    result.PublicVersion().Identifier().ID,
				MspID: result.PublicVersion().Identifier().MSPID,
			},
			EnrollmentCertificate: result.PublicVersion().EnrollmentCertificate(),
		},
	}}
}
