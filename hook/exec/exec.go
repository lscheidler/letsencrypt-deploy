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

package exec

import (
	"log"
	"os/exec"

	"github.com/lscheidler/letsencrypt-deploy/hook"
)

type Exec struct {
	command *string
}

func New(args string) *hook.Hook {
	e := Exec{
		command: &args,
	}
	h := hook.Hook(&e)
	return &h
}

func (e *Exec) Run() error {
	out, err := exec.Command("sh", "-c", *e.command).Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(out))
	return nil
}

func (e *Exec) String() string {
	return *e.command
}
