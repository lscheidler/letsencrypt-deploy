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
	"fmt"
	"strings"
)

// Domains list
type Domains []*string

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (i *Domains) String() string {
	return fmt.Sprint(*i)
}

// Set is the method to set the flag value, part of the flag.Value interface
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (i *Domains) Set(value string) error {
	fmt.Println(*i)
	for _, t := range strings.Split(value, ",") {
		// introduced indirection, because it doesn't work directly
		t2 := t
		*i = append(*i, &t2)
	}
	return nil
}

// Values returns a slice with the domains as string
func (i *Domains) Values() []string {
	var domains []string
	for _, v := range []*string(*i) {
		domains = append(domains, *v)
	}
	return domains
}
