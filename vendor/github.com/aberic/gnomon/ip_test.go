/*
 *  Copyright (c) 2019. aberic - All Rights Reserved.
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

package gnomon

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const checkMark = " OK! "
const ballotX = " ERROR! "

func mockServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ip address", IP().Get(r))
	}
	return httptest.NewServer(http.HandlerFunc(f))
}

func TestIPCommon_Get(t *testing.T) {
	statusCode := http.StatusOK

	server := mockServer()
	defer server.Close()

	t.Log("Given the need to test downloading content.")
	{

		t.Logf("\tWhen checking \"%s\" for status code \"%d\"", server.URL, statusCode)
		{
			resp, err := http.Get(server.URL)
			if err != nil {
				t.Fatal("\tShould be able to make the Get call.", ballotX, err)
			}
			t.Log("\t\tShould be able to make the Get call.", checkMark)
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode == statusCode {
				t.Logf("\t\tShould receive a \"%d\" status, %v", statusCode, checkMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", statusCode, ballotX, resp.StatusCode)
			}
		}
	}
}

func TestIPCommon_Get_XRealIP(t *testing.T) {
	statusCode := http.StatusOK
	server := mockServer()
	defer server.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL, nil)
	cookie1 := &http.Cookie{Name: "X-Xsrftoken", Value: "df41ba54db5011e89861002324e63af81", HttpOnly: true}
	req.AddCookie(cookie1)

	req.Header.Add("X-Real-IP", "192.168.0.1")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("\tShould be able to make the Get call.", ballotX, err)
	}
	t.Log("\t\tShould be able to make the Get call.", checkMark)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == statusCode {
		t.Logf("\t\tShould receive a \"%d\" status, %v", statusCode, checkMark)
	} else {
		t.Errorf("\t\tShould receive a \"%d\" status. %v %v", statusCode, ballotX, resp.StatusCode)
	}
}

func TestIPCommon_Get_XForwardedFor(t *testing.T) {
	statusCode := http.StatusOK
	server := mockServer()
	defer server.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL, nil)
	cookie1 := &http.Cookie{Name: "X-Xsrftoken", Value: "df41ba54db5011e89861002324e63af81", HttpOnly: true}
	req.AddCookie(cookie1)

	req.Header.Add("X-Forwarded-For", "192.168.0.2")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("\tShould be able to make the Get call.", ballotX, err)
	}
	t.Log("\t\tShould be able to make the Get call.", checkMark)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == statusCode {
		t.Logf("\t\tShould receive a \"%d\" status, %v", statusCode, checkMark)
	} else {
		t.Errorf("\t\tShould receive a \"%d\" status. %v %v", statusCode, ballotX, resp.StatusCode)
	}
}

func TestIPCommon_Get_XForwardedForIP6(t *testing.T) {
	statusCode := http.StatusOK
	server := mockServer()
	defer server.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL, nil)
	cookie1 := &http.Cookie{Name: "X-Xsrftoken", Value: "df41ba54db5011e89861002324e63af81", HttpOnly: true}
	req.AddCookie(cookie1)

	req.Header.Add("X-Forwarded-For", "::1")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("\tShould be able to make the Get call.", ballotX, err)
	}
	t.Log("\t\tShould be able to make the Get call.", checkMark)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == statusCode {
		t.Logf("\t\tShould receive a \"%d\" status, %v", statusCode, checkMark)
	} else {
		t.Errorf("\t\tShould receive a \"%d\" status. %v %v", statusCode, ballotX, resp.StatusCode)
	}
}
