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

package available_certificates

import (
	"github.com/lscheidler/letsencrypt-deploy/provider/fortios/api/v2/monitor/system/available_certificates/certificate"
)

type AvailableCertificates struct {
	Results []*certificate.Certificate `json:"results"`
	Vdom    *string                    `json:"vdom"`
	Path    *string                    `json:"path"`
	Name    *string                    `json:"name"`
	Status  *string                    `json:"status"`
	Serial  *string                    `json:"serial"`
	Version *string                    `json:"version"`
	Build   int                        `json:"build"`
}

func (ac *AvailableCertificates) Has(name string) *certificate.Certificate {
	for _, c := range ac.Results {
		if *c.Name == name {
			return c
		}
	}
	return nil
}
