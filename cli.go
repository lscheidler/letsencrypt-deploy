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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/lscheidler/letsencrypt-deploy/hook"
	"github.com/lscheidler/letsencrypt-deploy/hook/aws/sns"
	"github.com/lscheidler/letsencrypt-deploy/hook/exec"
)

type Domains []string

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (i *Domains) String() string {
	return fmt.Sprint(*i)
}

// Set is the method to set the flag value, part of the flag.Value interface
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (i *Domains) Set(value string) error {
	for _, t := range strings.Split(value, ",") {
		*i = append(*i, t)
	}
	return nil
}

type Hooks []*hook.Hook

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (i *Hooks) String() string {
	return fmt.Sprint(*i)
}

// Set is the method to set the flag value, part of the flag.Value interface
// Set's argument is a string to be parsed to set the flag.
func (i *Hooks) Set(value string) error {
	t := strings.Split(value, ";")

	switch t[0] {
	case "exec":
		switch l := len(t); l {
		case 2:
			*i = append(*i, exec.New(t[1]))
		default:
			log.Printf("Hook \"%s\" must have one argument (-hook \"exec;<command>\")", value)
		}
	case "sns":
		switch l := len(t); l {
		case 2, 3, 4:
			*i = append(*i, sns.New(t[1:]))
		default:
			log.Printf("Hook \"%s\" must have 1-2 arguments (-hook \"sns;<sns-topic>[;<sns-subject>[;<sns-message>]]\")", value)
		}
	default:
		log.Printf("No hook with the name \"%s\" available", t[0])
	}
	return nil
}

var (
	delay                 int
	domains               Domains
	dynamodbTableName     *string
	dynamodbTableNameFlag string
	email                 string
	hooks                 Hooks
	outputLocation        string
	passphraseFile        string
	passphrase            string
	prefix                string
)

const (
	delayUsage                  = "deploy certificates with a delay after creation (days)"
	delayDefaultVal             = 0
	domainUsage                 = "domains"
	dynamodbTableNameDefaultVal = ""
	dynamodbTableNameUsage      = "dynamodb table name"
	emailDefaultVal             = ""
	emailUsage                  = "account email"
	outputLocationDefaultVal    = "/etc/ssl/private"
	outputLocationUsage         = "output location for certificates"
	passphraseFileDefaultVal    = ""
	passphraseFileUsage         = "passphrase file"
	prefixDefaultVal            = "letsencrypt."
	prefixUsage                 = "file prefix for letsencrypt certificates"
	hooksUsage                  = "run hook after certificates has updated"
)

func parseArgs() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s (%s):\n", os.Args[0], version)
		flag.PrintDefaults()

		const (
			examples = `
hooks:

  exec;<command>
        execute <command> after updating certificates
  sns;<sns-topic>[;<sns-subject>[;<sns-message>]]
        publish a auto-generated message to <sns-topic> after updating certificates,
        if sns-message is set, use this message for publishing
`
		)

		fmt.Fprintf(flag.CommandLine.Output(), "%s\n", examples)
	}

	flag.IntVar(&delay, "delay", delayDefaultVal, delayUsage)
	flag.Var(&domains, "domain", domainUsage)
	flag.Var(&domains, "d", domainUsage)
	flag.StringVar(&dynamodbTableNameFlag, "dynamodbTableName", dynamodbTableNameDefaultVal, dynamodbTableNameUsage)
	flag.StringVar(&dynamodbTableNameFlag, "t", dynamodbTableNameDefaultVal, dynamodbTableNameUsage)
	flag.StringVar(&email, "email", emailDefaultVal, emailUsage)
	flag.StringVar(&email, "e", emailDefaultVal, emailUsage)
	flag.StringVar(&outputLocation, "outputLocation", outputLocationDefaultVal, outputLocationUsage)
	flag.StringVar(&outputLocation, "o", outputLocationDefaultVal, outputLocationUsage)
	flag.StringVar(&passphraseFile, "passphraseFile", passphraseFileDefaultVal, passphraseFileUsage)
	flag.StringVar(&passphraseFile, "p", passphraseFileDefaultVal, passphraseFileUsage)
	flag.StringVar(&prefix, "prefix", prefixDefaultVal, prefixUsage)

	flag.Var(&hooks, "hook", hooksUsage)
	flag.Var(&hooks, "H", hooksUsage)

	flag.Parse()

	if len(dynamodbTableNameFlag) > 0 {
		dynamodbTableName = &dynamodbTableNameFlag
	}

	if len(email) == 0 {
		log.Fatal("Option -email must be set.")
	}

	if len(domains) == 0 {
		log.Fatal("Option -domain must be set.")
	}

	if len(passphraseFile) == 0 {
		log.Fatal("Option -passphraseFile must be set.")
	}

	if dat, err := ioutil.ReadFile(passphraseFile); err != nil {
		log.Fatal(err)
	} else {
		passphrase = strings.TrimSpace(string(dat))
	}
}
