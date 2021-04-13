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

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/lscheidler/letsencrypt-lambda/account"
	"github.com/lscheidler/letsencrypt-lambda/dynamodb"

	"github.com/lscheidler/letsencrypt-deploy/config"
	"github.com/lscheidler/letsencrypt-deploy/provider/fortios"
	"github.com/lscheidler/letsencrypt-deploy/provider/local"
)

func executeHooks() {
	for _, hook := range hooks {
		log.Println((*hook))
		(*hook).Run()
	}
}

func main() {
	args := parseArgs()

	conf := config.NewWithDefaults()
	conf.Load(configFile)
	conf.Merge(args)
	if conf.Fortios && conf.FortiosBaseURL == nil {
		log.Fatal("Option -fortiosBaseURL or config value fortiosBaseURL must be set, when -fortios is set.")
	}

	domains := conf.Domains.Values()

	mustSetString(conf.Email, "Option -email or config option must be set.")
	mustSetStringSlice(conf.Domains, "Option -domains or config option must be set.")

	account := account.New(conf.Email, domains, nil)
	account.ClientPassphrase = conf.Passphrase

	dynamodb := dynamodb.New(conf.DynamodbTableName)
	dynamodb.LoadAccount(account)

	if cert := account.Certificates[fmt.Sprintf("%v", domains)]; cert != nil {
		if hours := time.Now().Sub(cert.CreatedAt).Hours(); conf.Delay == nil || hours > float64(*conf.Delay*24) {
			if conf.Local {
				localProvider := local.New(domains, *conf.OutputLocation, *conf.Prefix)
				if (*localProvider).Deploy(cert) {
					executeHooks()
				}
			}

			if conf.Fortios {
				fortiosProvider := fortios.New(domains, conf)
				if (*fortiosProvider).Deploy(cert) {
					executeHooks()
				}
			}
		} else {
			log.Printf("New certificate is %d days old. Skipping deployment because of delay=%d.\n", int(hours/24), *conf.Delay)
		}
	} else {
		log.Printf("No certificate for %v found.", domains)
	}
}

func mustSetString(val *string, message string) {
	if val == nil || len(*val) == 0 {
		log.Fatal(message)
	}
}

func mustSetStringSlice(val []*string, message string) {
	if len(val) == 0 {
		log.Fatal(message)
	}
}
