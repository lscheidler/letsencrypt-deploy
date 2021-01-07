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

package ssl_ssh_profile

import (
	"github.com/lscheidler/letsencrypt-deploy/provider/fortios/api/v2"
)

type SslSshProfile struct {
	Results []*result `json:"results"`
}

type result struct {
	Name       *string        `json:"name"`
	ServerCert *v2.ServerCert `json:"server-cert"`
}

/*
{
  "http_method": "<string>",
  "revision": "<string>",
  "results": [
    {
      "name": "<string>",
      "q_origin_key": "<string>",
      "css-class": "<string>",
      "comment": "",
      "ssl": {
        "inspect-all": "<string>",
        "client-cert-request": "<string>",
        "unsupported-ssl": "<string>",
        "invalid-server-cert": "<string>",
        "untrusted-server-cert": "<string>",
        "sni-server-cert-check": "<string>"
      },
      "https": {
        "ports": "",
        "status": "<string>",
        "client-cert-request": "<string>",
        "unsupported-ssl": "<string>",
        "invalid-server-cert": "<string>",
        "untrusted-server-cert": "<string>",
        "sni-server-cert-check": "<string>"
      },
      "ftps": {
        "ports": "",
        "status": "<string>",
        "client-cert-request": "<string>",
        "unsupported-ssl": "<string>",
        "invalid-server-cert": "<string>",
        "untrusted-server-cert": "<string>",
        "sni-server-cert-check": "<string>"
      },
      "imaps": {
        "ports": "",
        "status": "<string>",
        "client-cert-request": "<string>",
        "unsupported-ssl": "<string>",
        "invalid-server-cert": "<string>",
        "untrusted-server-cert": "<string>",
        "sni-server-cert-check": "<string>"
      },
      "pop3s": {
        "ports": "",
        "status": "<string>",
        "client-cert-request": "<string>",
        "unsupported-ssl": "<string>",
        "invalid-server-cert": "<string>",
        "untrusted-server-cert": "<string>",
        "sni-server-cert-check": "<string>"
      },
      "smtps": {
        "ports": "",
        "status": "<string>",
        "client-cert-request": "<string>",
        "unsupported-ssl": "<string>",
        "invalid-server-cert": "<string>",
        "untrusted-server-cert": "<string>",
        "sni-server-cert-check": "<string>"
      },
      "ssh": {
        "ports": "",
        "status": "<string>",
        "inspect-all": "<string>",
        "unsupported-version": "<string>",
        "ssh-tun-policy-check": "<string>",
        "ssh-algorithm": "<string>"
      },
      "whitelist": "<string>",
      "block-blacklisted-certificates": "<string>",
      "ssl-exempt": [ ],
      "server-cert-mode": "<string>",
      "use-ssl-server": "<string>",
      "caname": "<string>",
      "untrusted-caname": "",
      "server-cert": {
        "q_origin_key": "<string>",
        "name": "<string>",
        "datasource": "<string>"
      },
      "ssl-server": [ ],
      "ssl-anomalies-log": "<string>",
      "ssl-exemptions-log": "<string>",
      "rpc-over-https": "<string>",
      "mapi-over-https": "<string>"
    }
  ],
  "vdom": "<string>",
  "path": "<string>",
  "name": "<string>",
  "mkey": "<string>",
  "status": "<string>",
  "http_status": <integer>,
  "serial": "<string>",
  "version": "<string>",
  "build": <integer>
}
*/
