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
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lscheidler/letsencrypt-deploy/config"
	"github.com/lscheidler/letsencrypt-deploy/hook"
	"github.com/lscheidler/letsencrypt-deploy/hook/aws/sns"
	"github.com/lscheidler/letsencrypt-deploy/hook/exec"
)

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
	configFile string
	hooks      Hooks
)

const (
	configFileUsage             = "set config file"
	configFileDefaultVal        = "/etc/letsencrypt-deploy/config.json"
	delayUsage                  = "deploy certificates with a delay after creation (days)"
	delayDefaultVal             = 0
	domainUsage                 = "domains"
	dynamodbTableNameDefaultVal = ""
	dynamodbTableNameUsage      = "dynamodb table name"
	emailDefaultVal             = ""
	emailUsage                  = "account email"
	fortiosUsage                = "deploy certificate on fortios firewall"
	fortiosDefaultVal           = false
	fortiosBaseUrlUsage         = "fortios base url"
	fortiosBaseUrlDefaultVal    = ""
	localUsage                  = "deploy certificate on local machine"
	localDefaultVal             = false
	outputLocationDefaultVal    = "/etc/ssl/private"
	outputLocationUsage         = "output location for certificates"
	prefixDefaultVal            = "letsencrypt."
	prefixUsage                 = "file prefix for letsencrypt certificates"
	hooksUsage                  = "run hook after certificates has updated"
)

func parseArgs() *config.Config {
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

	var (
		domains config.Domains
	)

	conf := &config.Config{}
	flag.StringVar(&configFile, "configFile", configFileDefaultVal, configFileUsage)
	flag.StringVar(&configFile, "c", configFileDefaultVal, configFileUsage)
	if delay := flag.Int("delay", delayDefaultVal, delayUsage); *delay != delayDefaultVal {
		conf.Delay = delay
	}
	flag.Var(&domains, "domain", domainUsage)
	flag.Var(&domains, "d", domainUsage)
	if len(domains) > 0 {
		conf.Domains = domains
	}
	if dynamodbTableName := flag.String("dynamodbTableName", dynamodbTableNameDefaultVal, dynamodbTableNameUsage); *dynamodbTableName != dynamodbTableNameDefaultVal {
		conf.DynamodbTableName = dynamodbTableName
	}
	if dynamodbTableName := flag.String("t", dynamodbTableNameDefaultVal, dynamodbTableNameUsage); *dynamodbTableName != dynamodbTableNameDefaultVal {
		conf.DynamodbTableName = dynamodbTableName
	}
	if fortiosDeployment := flag.Bool("fortios", fortiosDefaultVal, fortiosUsage); *fortiosDeployment != fortiosDefaultVal {
		conf.Fortios = *fortiosDeployment
	}
	if fortiosBaseUrl := flag.String("fortios-base-url", fortiosBaseUrlDefaultVal, fortiosBaseUrlUsage); *fortiosBaseUrl != fortiosBaseUrlDefaultVal {
		conf.FortiosBaseUrl = fortiosBaseUrl
	}
	if email := flag.String("email", emailDefaultVal, emailUsage); *email != emailDefaultVal {
		conf.Email = email
	}
	if email := flag.String("e", emailDefaultVal, emailUsage); *email != emailDefaultVal {
		conf.Email = email
	}
	if localDeployment := flag.Bool("local", localDefaultVal, localUsage); *localDeployment != localDefaultVal {
		conf.Local = *localDeployment
	}
	if outputLocation := flag.String("outputLocation", outputLocationDefaultVal, outputLocationUsage); *outputLocation != outputLocationDefaultVal {
		conf.OutputLocation = outputLocation
	}
	if outputLocation := flag.String("o", outputLocationDefaultVal, outputLocationUsage); *outputLocation != outputLocationDefaultVal {
		conf.OutputLocation = outputLocation
	}
	if prefix := flag.String("prefix", prefixDefaultVal, prefixUsage); *prefix != prefixDefaultVal {
		conf.Prefix = prefix
	}

	flag.Var(&hooks, "hook", hooksUsage)
	flag.Var(&hooks, "H", hooksUsage)

	flag.Parse()

	return conf
}
