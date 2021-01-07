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

package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Delay                  *int      `json:"delay,omitempty"`
	Domains                Domains   `json:"domains,omitempty"`
	DynamodbTableName      *string   `json:"dynamodbTableName,omitempty"`
	Email                  *string   `json:"email,omitempty"`
	Fortios                bool      `json:"fortios"`
	FortiosAccessToken     *string   `json:"fortiosAccessToken,omitempty"`
	FortiosAdminServerCert bool      `json:"fortiosAdminServerCert,omitempty"`
	FortiosBaseUrl         *string   `json:"fortiosBaseUrl,omitempty"`
	FortiosSslSshProfiles  []*string `json:"fortiosSslSshProfiles,omitempty"`
	Local                  bool      `json:"local"`
	OutputLocation         *string   `json:"outputLocation,omitempty"`
	Passphrase             *string   `json:"passphrase,omitempty"`
	Prefix                 *string   `json:"prefix,omitempty"`
}

var (
	delayDefaultVal                  = 0
	fortiosDefaultVal                = false
	fortiosAdminServerCertDefaultVal = true
	localDefaultVal                  = false
	outputLocationDefaultVal         = "/etc/ssl/private"
	prefixDefaultVal                 = "letsencrypt."
)

func New() *Config {
	config := Config{
		FortiosSslSshProfiles: []*string{},
	}
	return &config
}

func NewWithDefaults() *Config {
	config := Config{
		Delay:                  &delayDefaultVal,
		FortiosAdminServerCert: fortiosAdminServerCertDefaultVal,
		FortiosSslSshProfiles:  []*string{},
		Local:                  localDefaultVal,
		OutputLocation:         &outputLocationDefaultVal,
		Prefix:                 &prefixDefaultVal,
	}
	return &config
}

func (config *Config) Load(path string) {
	if dat, err := ioutil.ReadFile(path); err == nil {
		if err := json.Unmarshal(dat, config); err != nil {
			log.Printf("%v", err)
		}
	} else {
		log.Printf("Configfile %s doesn't exist", path)
	}
}

func (config *Config) Merge(conf2 *Config) {
	config.Delay = setInt(config.Delay, conf2.Delay)
	config.Domains = setStringSlice([]*string(config.Domains), []*string(conf2.Domains))
	config.DynamodbTableName = setString(config.DynamodbTableName, conf2.DynamodbTableName)
	config.Email = setString(config.Email, conf2.Email)
	config.Fortios = setBool(config.Fortios, conf2.Fortios)
	config.FortiosAdminServerCert = setBool(config.FortiosAdminServerCert, conf2.FortiosAdminServerCert)
	config.FortiosBaseUrl = setString(config.FortiosBaseUrl, conf2.FortiosBaseUrl)
	config.FortiosSslSshProfiles = setStringSlice(config.FortiosSslSshProfiles, conf2.FortiosSslSshProfiles)
	config.Local = setBool(config.Local, conf2.Local)
	config.OutputLocation = setString(config.OutputLocation, conf2.OutputLocation)
	config.Prefix = setString(config.Prefix, conf2.Prefix)
}

func setBool(val1 bool, val2 bool) bool {
	if val2 {
		return val2
	} else {
		return val1
	}
}

func setInt(val1 *int, val2 *int) *int {
	if val2 != nil && *val2 > 0 {
		return val2
	} else {
		return val1
	}
}

func setString(val1 *string, val2 *string) *string {
	if val2 != nil && len(*val2) > 0 {
		return val2
	} else {
		return val1
	}
}

func setStringSlice(val1 []*string, val2 []*string) []*string {
	if val2 != nil && len(val2) > 0 {
		return val2
	} else {
		return val1
	}
}
