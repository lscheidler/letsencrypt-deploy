/*
Copyright 2021 Lars Eric Scheidler

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fortios

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/lscheidler/letsencrypt-lambda/account/certificate"
	"github.com/lscheidler/letsencrypt-lambda/crypto"

	"github.com/lscheidler/letsencrypt-deploy/config"
	"github.com/lscheidler/letsencrypt-deploy/provider"
	"github.com/lscheidler/letsencrypt-deploy/provider/fortios/api/v2/cmdb/firewall/ssl_ssh_profile"
	"github.com/lscheidler/letsencrypt-deploy/provider/fortios/api/v2/cmdb/system/global"
	"github.com/lscheidler/letsencrypt-deploy/provider/fortios/api/v2/monitor/system/available_certificates"
)

type FortiOS struct {
	client  *http.Client
	domains []string
	conf    *config.Config
}

func New(domains []string, conf *config.Config) *provider.Provider {
	f := &FortiOS{
		client:  &http.Client{},
		conf:    conf,
		domains: domains,
	}
	prov := provider.Provider(f)
	return &prov
}

func (f *FortiOS) Deploy(cert *certificate.Certificate) bool {
	if ac, err := f.availableCertificates(); err != nil {
		log.Fatal(err)
	} else {
		createdAt := cert.CreatedAt.Format("2006-01-02")
		name := createdAt + "_letsencrypt_" + f.domains[0]

		if c := ac.Has(name); c == nil {
			log.Println(name)
			key, err := getPrivateKey(cert.Pem)
			if err != nil {
				log.Fatal(err)
			}

			f.insertCertificate(name, cert.Pem, key)
		} else {
			log.Println("[fortios] Certificate already uptodate.")
		}

		// Admin-Interface:
		//   /api/v2/cmdb/system/global?access_token=
		//   PUT {"admin-server-cert": cert_name}
		if f.conf.FortiosAdminServerCert {
			if serv, err := f.getAdminServerCert(); err != nil {
				log.Println(err)
				f.setAdminServerCert(name)
			} else {
				if *serv.Results.AdminServerCert.Name != name {
					f.setAdminServerCert(name)
				} else {
					log.Println("[fortios] Admin Certificate already uptodate.")
				}
			}
		}

		// deep inspection:
		//   /api/v2/cmdb/firewall/ssl-ssh-profile/<profile-name>?datasource=1
		//   PUT {"server-cert": cert_name}
		for _, profile := range f.conf.FortiosSslSshProfiles {
			if serv, err := f.getSslSshProfileServerCert(profile); err != nil {
				log.Println(err)
				f.setSslSshProfileServerCert(profile, name)
			} else {
				if *serv.Results[0].ServerCert.Name != name {
					f.setSslSshProfileServerCert(profile, name)
				} else {
					log.Printf("[fortios:%s] SslSshProfile Certificate already uptodate.", *profile)
				}
			}
		}
	}
	return false
}

func (f *FortiOS) availableCertificates() (*available_certificates.AvailableCertificates, error) {
	u, _ := url.Parse(*f.conf.FortiosBaseUrl)
	u.Path = path.Join(u.Path, "/api/v2/monitor/system/available-certificates")

	if body, err := f.newRequest("GET", u.String(), "?scope=vdom&with_remote=1&with_ca=1&with_crl=1&access_token="+*f.conf.FortiosAccessToken, nil); err != nil {
		return nil, err
	} else {
		var apiResp available_certificates.AvailableCertificates
		if err := json.Unmarshal(body, &apiResp); err != nil {
			log.Println("Unmarshal failed", string(body))
			return nil, err
		}
		return &apiResp, nil
	}
}

func (f *FortiOS) insertCertificate(name string, pem []byte, key []byte) error {
	u, _ := url.Parse(*f.conf.FortiosBaseUrl)
	u.Path = path.Join(u.Path, "/api/v2/monitor/vpn-certificate/local/import")

	content := map[string]string{
		"certname":         name,
		"file_content":     base64.StdEncoding.EncodeToString(pem),
		"key_file_content": base64.StdEncoding.EncodeToString(key),
		"scope":            "global",
		"type":             "regular",
	}
	//log.Println(content)

	_, err := f.newRequest("POST", u.String(), "?access_token="+*f.conf.FortiosAccessToken, content)
	return err
}

func (f *FortiOS) getAdminServerCert() (*global.Global, error) {
	// /api/v2/cmdb/system/global?datasource=1&with_meta=1
	u, _ := url.Parse(*f.conf.FortiosBaseUrl)
	u.Path = path.Join(u.Path, "/api/v2/cmdb/system/global")

	if body, err := f.newRequest("GET", u.String(), "?datasource=1&with_meta=1&access_token="+*f.conf.FortiosAccessToken, nil); err != nil {
		return nil, err
	} else {
		var apiResp global.Global
		if err := json.Unmarshal(body, &apiResp); err != nil {
			log.Println("Unmarshal failed", string(body))
			return nil, err
		}
		return &apiResp, nil
	}
}

func (f *FortiOS) setAdminServerCert(name string) error {
	u, _ := url.Parse(*f.conf.FortiosBaseUrl)
	u.Path = path.Join(u.Path, "/api/v2/cmdb/system/global")

	content := map[string]string{
		"admin-server-cert": name,
	}

	_, err := f.newRequest("PUT", u.String(), "?access_token="+*f.conf.FortiosAccessToken, content)
	return err
}

func (f *FortiOS) getSslSshProfileServerCert(profname *string) (*ssl_ssh_profile.SslSshProfile, error) {
	// /api/v2/cmdb/firewall/ssl-ssh-profile/<profname>?datasource=1
	u, _ := url.Parse(*f.conf.FortiosBaseUrl)
	u.Path = path.Join(u.Path, "/api/v2/cmdb/firewall/ssl-ssh-profile", *profname)

	if body, err := f.newRequest("GET", u.String(), "?datasource=1&access_token="+*f.conf.FortiosAccessToken, nil); err != nil {
		return nil, err
	} else {
		var apiResp ssl_ssh_profile.SslSshProfile
		if err := json.Unmarshal(body, &apiResp); err != nil {
			log.Println("Unmarshal failed", string(body))
			return nil, err
		}
		return &apiResp, nil
	}
}

func (f *FortiOS) setSslSshProfileServerCert(profname *string, certname string) error {
	u, _ := url.Parse(*f.conf.FortiosBaseUrl)
	u.Path = path.Join(u.Path, "/api/v2/cmdb/firewall/ssl-ssh-profile/"+*profname)

	content := map[string]string{
		"name":        *profname,
		"server-cert": certname,
	}

	_, err := f.newRequest("PUT", u.String(), "?datasource=1&access_token="+*f.conf.FortiosAccessToken, content)
	return err
}

func getPrivateKey(pemdata []byte) ([]byte, error) {
	block, _ := pem.Decode(pemdata)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	priv, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	var privpem bytes.Buffer
	if err := crypto.EncodeECDSAKey(&privpem, priv); err != nil {
		return nil, err
	}

	return privpem.Bytes(), nil
}

func (f *FortiOS) newRequest(method string, url string, requestParameter string, content map[string]string) ([]byte, error) {
	var body io.Reader

	if content != nil {
		if contentJson, err := json.Marshal(content); err != nil {
			return nil, err
		} else {
			body = bytes.NewBuffer(contentJson)
		}
	}

	req, err := http.NewRequest(method, url+requestParameter, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	if resp, err := f.client.Do(req); err != nil {
		return nil, err
	} else {
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			return nil, err
		} else {
			if resp.StatusCode == 200 {
				//log.Println(string(body))
				return body, nil
			} else {
				return nil, fmt.Errorf("Request %s failed with %d: %s", url, resp.StatusCode, string(body))
			}
		}
	}
}
