/*
Copyright 2022 The Crossplane Authors.

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

// Package main lints CRDs for breaking changes.
package main

import (
	"fmt"
	"log"
	"os"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func main() {
	oldfile, err := os.ReadFile("old.yaml")
	if err != nil {
		log.Fatal(err)
	}

	newfile, err := os.ReadFile("new.yaml")
	if err != nil {
		log.Fatal(err)
	}

	old := &v1.CustomResourceDefinition{}
	err = yaml.Unmarshal(oldfile, old)
	if err != nil {
		log.Fatal(err)
	}

	crd := &v1.CustomResourceDefinition{}
	err = yaml.Unmarshal(newfile, crd)
	if err != nil {
		log.Fatal(err)
	}

	list := PrintFields(old.Spec.Versions[0].Schema.OpenAPIV3Schema, "", crd.Spec.Versions[0].Schema.OpenAPIV3Schema)
	for i := range list {
		fmt.Println(list[i]) //nolint:forbidigo // CLI tools are allowed to print their output.
	}
}

// PrintFields function recursively traverses through the keys.
func PrintFields(sch *v1.JSONSchemaProps, prefix string, newSchema *v1.JSONSchemaProps) []string {
	var a []string

	if len(sch.Properties) == 0 {
		return nil
	}

	for key := range sch.Properties {
		val := sch.Properties[key]
		var temp string

		if prefix == "" {
			temp = key
		} else {
			temp = prefix + "." + key
		}

		prop, ok := newSchema.Properties[key]

		if !ok {
			a = append(a, temp)
			continue
		}
		a = append(a, PrintFields(&val, temp, &prop)...)
	}
	return a
}
