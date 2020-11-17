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
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/lscheidler/letsencrypt-lambda/account/certificate"
	"github.com/lscheidler/letsencrypt-lambda/crypto"
)

func readLocalCertificate() []byte {
	pemFilename := getFilename(nil, ".pem")
	if dat, err := ioutil.ReadFile(pemFilename); err == nil {
		return dat
	}
	return nil
}

func writeCertificate(cert *certificate.Certificate) {
	createdAt := cert.CreatedAt.Format("2006-01-02")
	pemFilename := getFilename(&createdAt, ".pem")
	keyFilename := getFilename(&createdAt, ".key")

	pemdata := cert.Pem
	block, _ := pem.Decode(pemdata)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		log.Fatal("failed to decode PEM block containing private key")
	}

	priv, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	var privpem bytes.Buffer
	if err := crypto.EncodeECDSAKey(&privpem, priv); err != nil {
		log.Fatal(err)
	}

	log.Println("Writing file", pemFilename)
	if err := ioutil.WriteFile(pemFilename, pemdata, 0600); err != nil {
		log.Fatal(err)
	}

	log.Println("Writing file", keyFilename)
	if err := ioutil.WriteFile(keyFilename, privpem.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
}

func rewriteLinks(cert *certificate.Certificate) {
	createdAt := cert.CreatedAt.Format("2006-01-02")

	pemFilename := filepath.Base(getFilename(&createdAt, ".pem"))
	pemSymlink := getFilename(nil, ".pem")
	rewriteLink(pemFilename, pemSymlink)

	keyFilename := filepath.Base(getFilename(&createdAt, ".key"))
	keySymlink := getFilename(nil, ".key")
	rewriteLink(keyFilename, keySymlink)
}

func rewriteLink(source string, destination string) {
	if stat, err := os.Stat(destination); stat != nil || err == nil {
		if err := os.Remove(destination); err != nil {
			log.Println(err)
		}
	}

	log.Printf("Creating symlink: %s %s\n", source, destination)
	if err := os.Symlink(source, destination); err != nil {
		log.Println(err)
	}
}

func getFilename(date *string, ext string) string {
	filename := prefix
	if date != nil {
		filename = filename + *date + "."
	}
	filename = filename + domains[0] + ext

	return filepath.Join(outputLocation, filename)
}
