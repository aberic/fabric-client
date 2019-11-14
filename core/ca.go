/*
 * Copyright (c) 2019. Aberic - All Rights Reserved.
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
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// GetCAInfo returns generic CA information
func caInfo(orgName string, sdk *fabsdk.FabricSDK) (*msp.GetCAInfoResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.GetCAInfo()
}

// Enroll enrolls a registered user in order to receive a signed X509 certificate.
// A new key pair is generated for the user. The private key and the
// enrollment certificate issued by the CA are stored in SDK stores.
// They can be retrieved by calling IdentityManager.GetSigningIdentity().
//  Parameters:
//  enrollmentID enrollment ID of a registered user
//  opts are optional enrollment options
//
//  Returns:
//  an error if enrollment fails
func enroll(orgName, enrollmentID string, sdk *fabsdk.FabricSDK, opts ...msp.EnrollmentOption) error {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return err
	}
	return mspClient.Enroll(enrollmentID, opts...)
}

// Reenroll reenrolls an enrolled user in order to obtain a new signed X509 certificate
//  Parameters:
//  enrollmentID enrollment ID of a registered user
//
//  Returns:
//  an error if re-enrollment fails
func reenroll(orgName, enrollmentID string, sdk *fabsdk.FabricSDK, opts ...msp.EnrollmentOption) error {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return err
	}
	return mspClient.Reenroll(enrollmentID, opts...)
}

// Register registers a User with the Fabric CA
//  Parameters:
//  request is registration request
//
//  Returns:
//  enrolment secret
func register(orgName string, registerReq *msp.RegistrationRequest, sdk *fabsdk.FabricSDK) (string, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return "", err
	}
	return mspClient.Register(registerReq)
}

// AffiliationRequest represents the request to add/remove affiliation to the fabric-ca-server
func addAffiliation(orgName string, affReq *msp.AffiliationRequest, sdk *fabsdk.FabricSDK) (*msp.AffiliationResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.AddAffiliation(affReq)
}

// RemoveAffiliation removes an existing affiliation from the server
func removeAffiliation(orgName string, affReq *msp.AffiliationRequest, sdk *fabsdk.FabricSDK) (*msp.AffiliationResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.RemoveAffiliation(affReq)
}

// ModifyAffiliation renames an existing affiliation on the server
func modifyAffiliation(orgName string, affReq *msp.ModifyAffiliationRequest, sdk *fabsdk.FabricSDK) (*msp.AffiliationResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.ModifyAffiliation(affReq)
}

// GetAffiliation returns information about the requested affiliation
func getAffiliation(affiliation, orgName string, sdk *fabsdk.FabricSDK) (*msp.AffiliationResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.GetAffiliation(affiliation)
}

// GetAffiliationByCaName returns information about the requested affiliation
func getAffiliationByCaName(affiliation, orgName, caName string, sdk *fabsdk.FabricSDK) (*msp.AffiliationResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.GetAffiliation(affiliation, msp.WithCA(caName))
}

// GetAllAffiliations returns all affiliations that the caller is authorized to see
func getAllAffiliations(orgName string, sdk *fabsdk.FabricSDK) (*msp.AffiliationResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.GetAllAffiliations()
}

// GetAllAffiliationsByCaName returns all affiliations that the caller is authorized to see
func getAllAffiliationsByCaName(orgName, caName string, sdk *fabsdk.FabricSDK) (*msp.AffiliationResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.GetAllAffiliations(msp.WithCA(caName))
}

// GetAllIdentities returns all identities that the caller is authorized to see
//  Parameters:
//  options holds optional request options
//  Returns:
//  Response containing identities
func getAllIdentities(orgName string, sdk *fabsdk.FabricSDK) ([]*msp.IdentityResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.GetAllIdentities()
}

// GetAllIdentitiesByCaName returns all identities that the caller is authorized to see
//  Parameters:
//  options holds optional request options
//  Returns:
//  Response containing identities
func getAllIdentitiesByCaName(orgName, caName string, sdk *fabsdk.FabricSDK) ([]*msp.IdentityResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.GetAllIdentities(msp.WithCA(caName))
}

// CreateIdentity creates a new identity with the Fabric CA server. An enrollment secret is returned which can then be used,
// along with the enrollment ID, to enroll a new identity.
//  Parameters:
//  request holds info about identity
//
//  Returns:
//  Return identity info including the secret
func createIdentity(orgName string, req *msp.IdentityRequest, sdk *fabsdk.FabricSDK) (*msp.IdentityResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.CreateIdentity(req)
}

// ModifyIdentity modifies identity with the Fabric CA server.
//  Parameters:
//  request holds info about identity
//
//  Returns:
//  Return updated identity info
func modifyIdentity(orgName string, req *msp.IdentityRequest, sdk *fabsdk.FabricSDK) (*msp.IdentityResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.ModifyIdentity(req)
}

// GetIdentity retrieves identity information.
//  Parameters:
//  ID is required identity ID
//  options holds optional request options
//
//  Returns:
//  Response containing identity information
func getIdentity(id, orgName string, sdk *fabsdk.FabricSDK) (*msp.IdentityResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.GetIdentity(id)
}

// GetIdentityByCaName retrieves identity information.
//  Parameters:
//  ID is required identity ID
//  options holds optional request options
//
//  Returns:
//  Response containing identity information
func getIdentityByCaName(id, caName, orgName string, sdk *fabsdk.FabricSDK) (*msp.IdentityResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.GetIdentity(id, msp.WithCA(caName))
}

// RemoveIdentity removes identity with the Fabric CA server.
//  Parameters:
//  request holds info about identity to be removed
//
//  Returns:
//  Return removed identity info
func removeIdentity(orgName string, req *msp.RemoveIdentityRequest, sdk *fabsdk.FabricSDK) (*msp.IdentityResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.RemoveIdentity(req)
}

// CreateSigningIdentity creates a signing identity with the given options
func createSigningIdentity(orgName string, sdk *fabsdk.FabricSDK, opts ...mspctx.SigningIdentityOption) (mspctx.SigningIdentity, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.CreateSigningIdentity(opts...)
}

// GetSigningIdentity returns signing identity for id
//  Parameters:
//  id is user id
//
//  Returns:
//  signing identity
func getSigningIdentity(id, orgName string, sdk *fabsdk.FabricSDK) (mspctx.SigningIdentity, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.GetSigningIdentity(id)
}

// Revoke revokes a User with the Fabric CA
//  Parameters:
//  request is revocation request
//
//  Returns:
//  revocation response
func revoke(orgName string, req *msp.RevocationRequest, sdk *fabsdk.FabricSDK) (*msp.RevocationResponse, error) {
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, err
	}
	return mspClient.Revoke(req)
}
