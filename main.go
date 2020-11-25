/*
Copyright 2020 Lars Eric Scheidler

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
)

func executeHooks() {
	for _, hook := range hooks {
		log.Println((*hook))
		(*hook).Run()
	}
}

func main() {
	parseArgs()

	account := account.New(&email, domains, nil)
	account.ClientPassphrase = &passphrase

	dynamodb := dynamodb.New(dynamodbTableName)
	dynamodb.LoadAccount(account)

	if cert := account.Certificates[fmt.Sprintf("%v", domains)]; cert != nil {
		if hours := time.Now().Sub(cert.CreatedAt).Hours(); hours > float64(delay*24) {
			if string(cert.Pem) != string(readLocalCertificate()) {
				writeCertificate(cert)
				rewriteLinks(cert)
				executeHooks()
			} else {
				log.Println("Local certificate already uptodate.")
			}
		} else {
			log.Printf("New certificate is %d days old. Skipping deployment because of delay=%d.\n", int(hours/24), delay)
		}
	} else {
		log.Printf("No certificate for %v found.", domains)
	}
}
