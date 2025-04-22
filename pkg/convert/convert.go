/*
Copyright 2025 The Crossplane Authors.

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

// Package convert contains utilities for converting to and from pointers.
package convert

import (
	"strconv"

	"k8s.io/utils/ptr"
)

// NOTE(negz): There are many equivalents of FromPtrValue and ToPtrValue
// throughout the Crossplane codebase. We duplicate them here to reduce the
// number of packages our API types have to import to support references.

// FromPtrValue adapts a string pointer field for use as a CurrentValue.
//
// Deprecated: Use ptr.Deref from k8s.io/utils/ptr.
func FromPtrValue(v *string) string {
	return ptr.Deref(v, "")
}

// FromFloatPtrValue adapts a float pointer field for use as a CurrentValue.
func FromFloatPtrValue(v *float64) string {
	if v == nil {
		return ""
	}
	return strconv.FormatFloat(*v, 'f', 0, 64)
}

// FromIntPtrValue adapts an int pointer field for use as a CurrentValue.
func FromIntPtrValue(v *int64) string {
	if v == nil {
		return ""
	}
	return strconv.FormatInt(*v, 10)
}

// ToPtrValue adapts a ResolvedValue for use as a string pointer field.
//
// Deprecated: Use ptr.To from k8s.io/utils/ptr.
func ToPtrValue(v string) *string {
	return ptr.To(v)
}

// ToFloatPtrValue adapts a ResolvedValue for use as a float64 pointer field.
func ToFloatPtrValue(v string) *float64 {
	if v == "" {
		return nil
	}
	vParsed, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil
	}
	return &vParsed
}

// ToIntPtrValue adapts a ResolvedValue for use as an int pointer field.
func ToIntPtrValue(v string) *int64 {
	if v == "" {
		return nil
	}
	vParsed, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil
	}
	return &vParsed
}

// FromPtrValues adapts a slice of string pointer fields for use as CurrentValues.
// NOTE: Do not use this utility function unless you have to.
// Using pointer slices does not adhere to our current API practices.
// The current use case is where generated code creates reference-able fields in a provider which are
// string pointers and need to be resolved as part of `ResolveMultiple`.
func FromPtrValues(v []*string) []string {
	res := make([]string, len(v))
	for i := range v {
		res[i] = FromPtrValue(v[i])
	}
	return res
}

// FromFloatPtrValues adapts a slice of float64 pointer fields for use as CurrentValues.
func FromFloatPtrValues(v []*float64) []string {
	res := make([]string, len(v))
	for i := range v {
		res[i] = FromFloatPtrValue(v[i])
	}
	return res
}

// FromIntPtrValues adapts a slice of int64 pointer fields for use as CurrentValues.
func FromIntPtrValues(v []*int64) []string {
	res := make([]string, len(v))
	for i := range v {
		res[i] = FromIntPtrValue(v[i])
	}
	return res
}

// ToPtrValues adapts ResolvedValues for use as a slice of string pointer fields.
// NOTE: Do not use this utility function unless you have to.
// Using pointer slices does not adhere to our current API practices.
// The current use case is where generated code creates reference-able fields in a provider which are
// string pointers and need to be resolved as part of `ResolveMultiple`.
func ToPtrValues(v []string) []*string {
	res := make([]*string, len(v))
	for i := range v {
		res[i] = ToPtrValue(v[i])
	}
	return res
}

// ToFloatPtrValues adapts ResolvedValues for use as a slice of float64 pointer fields.
func ToFloatPtrValues(v []string) []*float64 {
	res := make([]*float64, len(v))
	for i := range v {
		res[i] = ToFloatPtrValue(v[i])
	}
	return res
}

// ToIntPtrValues adapts ResolvedValues for use as a slice of int64 pointer fields.
func ToIntPtrValues(v []string) []*int64 {
	res := make([]*int64, len(v))
	for i := range v {
		res[i] = ToIntPtrValue(v[i])
	}
	return res
}
