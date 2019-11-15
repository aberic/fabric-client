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

package geneses

import (
	"fmt"
	"github.com/aberic/fabric-client/grpc/proto/generate"
	"time"
)

type generateCertificateRequest struct {
	generate.EnrollRequest
	NotAfter  time.Time
	NotBefore time.Time
	CR        string `json:"certificate_request"`
}

type enrollmentResponse struct {
	Result   enrollmentResponseResult `json:"result"`
	Success  bool                     `json:"success"`
	Errors   []responseErr            `json:"errors"`
	Messages []string                 `json:"messages"`
}

type responseErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type enrollmentResponseResult struct {
	Cert       string
	ServerInfo enrollmentResponseServerInfo
}

type enrollmentResponseServerInfo struct {
	CAName                    string
	CAChain                   string
	IssuerPublicKey           string
	IssuerRevocationPublicKey string
	Version                   string
}

func (ge enrollmentResponse) error() error {
	errs := ""
	for _, err := range ge.Errors {
		errs += err.Message + ": "
	}
	return fmt.Errorf(errs)
}
