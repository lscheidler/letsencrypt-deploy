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

package sns

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/sns"

	awshelper "github.com/lscheidler/letsencrypt-lambda/helper/aws"

	"github.com/lscheidler/letsencrypt-deploy/hook"
)

// Sns struct
type Sns struct {
	topic   *string
	subject *string
	message *string
}

// New returns new sns hook struct
func New(args []string) *hook.Hook {
	e := Sns{
		topic: &args[0],
	}

	if len(args) > 1 {
		e.subject = &args[1]
	}

	if len(args) > 2 {
		e.message = &args[2]
	}

	h := hook.Hook(&e)
	return &h
}

// Run hook
func (e *Sns) Run() error {
	sess, conf := awshelper.GetAwsSession()
	svc := sns.New(sess, conf)
	input := &sns.PublishInput{
		TopicArn: e.topic,
		Subject:  e.getSubject(),
		Message:  e.getMessage(),
	}

	if _, err := svc.Publish(input); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// String return topic
func (e *Sns) String() string {
	return *e.topic
}

func (e *Sns) getSubject() *string {
	if e.subject != nil {
		return e.subject
	}

	subject := "[" + e.getHostname() + "] Updated letsencrypt certificates"
	return &subject
}

func (e *Sns) getMessage() *string {
	if e.message != nil {
		return e.message
	}

	message := `Host: ` + e.getHostname() + `

Updated letsencrypt certificates`
	return &message
}

func (e *Sns) getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "<unknown>"
	}

	return hostname
}
