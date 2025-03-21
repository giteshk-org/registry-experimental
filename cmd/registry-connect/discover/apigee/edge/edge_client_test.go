// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package edge

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func oauthTestServer(t *testing.T) *httptest.Server {
	m := http.NewServeMux()

	resp := OAuthResponse{
		AccessToken: "token",
	}

	m.HandleFunc("/noauth", (func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusUnauthorized)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	m.HandleFunc("/oauth", (func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				t.Fatalf("want no error %v", err)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	return httptest.NewServer(m)
}

func TestNewEdgeClient(t *testing.T) {
	ts := oauthTestServer(t)
	defer ts.Close()

	opts := &EdgeClientOptions{
		InsecureSkipVerify: true,
		Auth: &EdgeAuth{
			SkipAuth: false,
			Username: "hi",
			Password: "secret",
			MFAToken: "mfa",
		},
		Debug: true,
	}

	var err error

	SetOAuthURL(ts.URL + "/noauth")

	_, err = NewEdgeClient(opts)
	errorContains(t, err, "401")

	SetOAuthURL(ts.URL + "/oauth")

	_, err = NewEdgeClient(opts)
	if err != nil {
		t.Errorf("want no error got %v", err)
	}
}

func TestStreamToString(t *testing.T) {
	in := "test"
	out := StreamToString(strings.NewReader(in))
	if in != out {
		t.Errorf("want %s got %s", in, out)
	}
}

func TestBool(t *testing.T) {
	in := true
	out := *Bool(in)
	if in != out {
		t.Errorf("want %v got %v", in, out)
	}
}

func TestInt(t *testing.T) {
	in := 123
	out := *Int(in)
	if in != out {
		t.Errorf("want %d got %d", in, out)
	}
}

func TestString(t *testing.T) {
	in := "test"
	out := *String(in)
	if in != out {
		t.Errorf("want %s got %s", in, out)
	}
}

func TestOnRequestCompleted(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{}"))
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	c := &EdgeClient{
		client:     http.DefaultClient,
		BaseURL:    u,
		BaseURLEnv: u,
		auth: &EdgeAuth{
			BearerToken: "token",
		},
	}

	count := 0
	c.OnRequestCompleted(func(req *http.Request, res *http.Response) {
		count += 1
	})

	req, err := c.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Do(req, nil)
	if err != nil {
		t.Errorf("want no error got %v", err)
	}

	if count != 1 {
		t.Errorf("want count to be 1, got %d", count)
	}
}

func TestAuthHeader(t *testing.T) {
	var c *EdgeClient

	u, err := url.Parse("dummy.url")
	if err != nil {
		t.Fatal(err)
	}
	c = &EdgeClient{
		client:     http.DefaultClient,
		BaseURL:    u,
		BaseURLEnv: u,
		auth: &EdgeAuth{
			BearerToken: "token",
		},
	}
	req, err := c.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	authHeader := req.Header.Get("Authorization")
	if authHeader != "Bearer token" {
		t.Errorf("want authorization header to be 'Bearer token', got %s", authHeader)
	}

	c = &EdgeClient{
		client:     http.DefaultClient,
		BaseURL:    u,
		BaseURLEnv: u,
		auth: &EdgeAuth{
			Username: "hi",
			Password: "secret",
		},
	}
	req, err = c.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	authHeader = req.Header.Get("Authorization")
	if authHeader != "Basic aGk6c2VjcmV0" {
		t.Errorf("want authorization header to be 'Basic aGk6c2VjcmV0', got %s", authHeader)
	}
}

func TestNetrcRetrieval(t *testing.T) {
	cred := []byte(`machine api.enterprise.apigee.com
	login hi
	password secret`)

	tmpFile, err := os.CreateTemp("", ".netrc")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := tmpFile.Write(cred); err != nil {
		t.Fatalf("%v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = retrieveAuthFromNetrc("not a path", "dummy")
	errorContains(t, err, "open not a path: no such file or directory")

	_, err = retrieveAuthFromNetrc(tmpFile.Name(), "dummy")
	errorContains(t, err, "cannot find machine:dummy")

	auth, err := retrieveAuthFromNetrc(tmpFile.Name(), "api.enterprise.apigee.com")
	if err != nil {
		t.Errorf("want no error got %v", err)
	}
	if auth.Username != "hi" || auth.Password != "secret" {
		t.Errorf("want username to be hi got %s\n want password to be secret got %s",
			auth.Username, auth.Password)
	}
}

func TestMutualTLSWithCerts(t *testing.T) {
	ts := newMutualTLSServer()
	defer ts.Close()

	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(ts.Certificate())

	opts := &EdgeClientOptions{
		MgmtURL:      ts.URL,
		Org:          "org",
		Env:          "env",
		RootCAs:      caCertPool,
		Certificates: ts.TLS.Certificates,
		Auth: &EdgeAuth{
			SkipAuth: true,
		},
	}

	c, err := NewEdgeClient(opts)
	if err != nil {
		t.Fatal(err)
	}

	req, err := c.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Do(req, nil)
	if err != nil {
		t.Errorf("want no error got %v", err)
	}
}

func TestMutualTLSNoCerts(t *testing.T) {
	ts := newMutualTLSServer()
	defer ts.Close()

	opts := &EdgeClientOptions{
		MgmtURL: ts.URL,
		Org:     "org",
		Env:     "env",
		Auth: &EdgeAuth{
			SkipAuth: true,
		},
		InsecureSkipVerify: true,
	}

	c, err := NewEdgeClient(opts)
	if err != nil {
		t.Fatal(err)
	}

	req, err := c.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Do(req, nil)
	errorContains(t, err, "remote error: tls: bad certificate")
}

func newMutualTLSServer() *httptest.Server {
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{}"))
	}))
	// require mTLS
	ts.TLS = &tls.Config{
		RootCAs:    x509.NewCertPool(),
		ClientAuth: tls.RequireAnyClientCert,
	}
	ts.StartTLS()
	ts.TLS.RootCAs.AddCert(ts.Certificate())

	return ts
}

// errorContains checks if the error string contains the wanted pattern
func errorContains(t *testing.T, out error, want string) {
	t.Helper()
	if out == nil {
		if want != "" {
			t.Errorf("got no error want %s", want)
		}
	} else if !strings.Contains(out.Error(), want) {
		t.Errorf("want %s, got %v ", want, out)
	}
}
