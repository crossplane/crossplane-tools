/*
Copyright 2019 The Crossplane Authors.

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

// Package match identifies Go types as common Crossplane resources.
package match

import (
	"go/types"
	"slices"

	"github.com/crossplane/crossplane-tools/internal/comments"
	"github.com/crossplane/crossplane-tools/internal/fields"
)

// An Object matcher is a function that returns true if the supplied object
// matches.
type Object func(o types.Object) bool

// ManagedLegacy returns an Object matcher that returns true if the supplied
// Object is a legacy (cluster-scoped) Crossplane managed resource using
// crossplane-runtime common/v1 types.
func ManagedLegacy() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsSpec().And(fields.HasFieldThat(
				fields.IsResourceSpec().And(fields.IsEmbedded()),
			)),
			fields.IsStatus().And(fields.HasFieldThat(
				fields.IsResourceStatus().And(fields.IsEmbedded()),
			)),
		)
	}
}

// ManagedModern returns an Object matcher that returns true if the supplied
// Object is a modern (namespaced) Crossplane managed resource using
// crossplane-runtime common/v2 types.
func ManagedModern() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsSpec().And(fields.HasFieldThat(
				fields.IsResourceV2Spec().And(fields.IsEmbedded()),
			)),
			fields.IsStatus().And(fields.HasFieldThat(
				fields.IsResourceStatus().And(fields.IsEmbedded()),
			)),
		)
	}
}

// ManagedModernCore returns an Object matcher that returns true if the supplied
// Object is a namespaced (modern) Crossplane managed resource using the core API
// v2 types.
func ManagedModernCore() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsSpec().And(fields.HasFieldThat(
				fields.IsManagedResourceSpecCore().And(fields.IsEmbedded()),
			)),
			fields.IsStatus().And(fields.HasFieldThat(
				fields.IsManagedResourceStatusCore().And(fields.IsEmbedded()),
			)),
		)
	}
}

// ManagedListLegacy returns an Object matcher that returns true if the supplied
// Object is a list of legacy (cluster-scoped) Crossplane managed resources.
func ManagedListLegacy() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsItems().And(fields.IsSlice()).And(fields.HasFieldThat(
				fields.IsTypeMeta().And(fields.IsEmbedded()),
				fields.IsObjectMeta().And(fields.IsEmbedded()),
				fields.IsSpec().And(fields.HasFieldThat(
					fields.IsResourceSpec().And(fields.IsEmbedded()),
				)),
				fields.IsStatus().And(fields.HasFieldThat(
					fields.IsResourceStatus().And(fields.IsEmbedded()),
				)),
			)),
		)
	}
}

// ManagedListModern returns an Object matcher that returns true if the supplied
// Object is a list of modern (namespaced) Crossplane managed resources using
// crossplane-runtime common/v2 types.
func ManagedListModern() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsItems().And(fields.IsSlice()).And(fields.HasFieldThat(
				fields.IsTypeMeta().And(fields.IsEmbedded()),
				fields.IsObjectMeta().And(fields.IsEmbedded()),
				fields.IsSpec().And(fields.HasFieldThat(
					fields.IsResourceV2Spec().And(fields.IsEmbedded()),
				)),
				fields.IsStatus().And(fields.HasFieldThat(
					fields.IsResourceStatus().And(fields.IsEmbedded()),
				)),
			)),
		)
	}
}

// ManagedListModernCore returns an Object matcher that returns true if the
// supplied Object is a list of namespaced (modern) Crossplane managed resources
// using the core API v2 types.
func ManagedListModernCore() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsItems().And(fields.IsSlice()).And(fields.HasFieldThat(
				fields.IsTypeMeta().And(fields.IsEmbedded()),
				fields.IsObjectMeta().And(fields.IsEmbedded()),
				fields.IsSpec().And(fields.HasFieldThat(
					fields.IsManagedResourceSpecCore().And(fields.IsEmbedded()),
				)),
				fields.IsStatus().And(fields.HasFieldThat(
					fields.IsManagedResourceStatusCore().And(fields.IsEmbedded()),
				)),
			)),
		)
	}
}

// ManagedLegacyCore returns an Object matcher that returns true if the supplied
// Object is a cluster-scoped (legacy) Crossplane managed resource using the core
// API v2 ClusterManagedResourceSpec.
func ManagedLegacyCore() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsSpec().And(fields.HasFieldThat(
				fields.IsClusterManagedResourceSpecCore().And(fields.IsEmbedded()),
			)),
			fields.IsStatus().And(fields.HasFieldThat(
				fields.IsResourceStatus().Or(fields.IsManagedResourceStatusCore()).And(fields.IsEmbedded()),
			)),
		)
	}
}

// ManagedListLegacyCore returns an Object matcher that returns true if the
// supplied Object is a list of cluster-scoped (legacy) Crossplane managed
// resources using the core API v2 types.
func ManagedListLegacyCore() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsItems().And(fields.IsSlice()).And(fields.HasFieldThat(
				fields.IsTypeMeta().And(fields.IsEmbedded()),
				fields.IsObjectMeta().And(fields.IsEmbedded()),
				fields.IsSpec().And(fields.HasFieldThat(
					fields.IsClusterManagedResourceSpecCore().And(fields.IsEmbedded()),
				)),
				fields.IsStatus().And(fields.HasFieldThat(
					fields.IsResourceStatus().Or(fields.IsManagedResourceStatusCore()).And(fields.IsEmbedded()),
				)),
			)),
		)
	}
}

// ProviderConfig returns an Object matcher that returns true if the supplied
// Object is a Crossplane ProviderConfig.
func ProviderConfig() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsSpec(),
			fields.IsStatus().And(fields.HasFieldThat(
				fields.IsProviderConfigStatus().And(fields.IsEmbedded()),
			)),
		)
	}
}

// ProviderConfigCore returns an Object matcher that returns true if the supplied
// Object is a Crossplane ProviderConfig using the core API v2 types.
func ProviderConfigCore() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsSpec(),
			fields.IsStatus().And(fields.HasFieldThat(
				fields.IsProviderConfigStatusCore().And(fields.IsEmbedded()),
			)),
		)
	}
}

// ProviderConfigUsageLegacy returns an Object matcher that returns true if the
// supplied Object is a legacy (non-typed) Crossplane ProviderConfigUsage.
func ProviderConfigUsageLegacy() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsProviderConfigUsage().And(fields.IsEmbedded()),
		)
	}
}

// ProviderConfigUsageModern returns an Object matcher that returns true if the
// supplied Object is a modern (typed) Crossplane ProviderConfigUsage.
func ProviderConfigUsageModern() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsTypedProviderConfigUsage().And(fields.IsEmbedded()),
		)
	}
}

// ProviderConfigUsageLegacyCore returns an Object matcher that returns true if
// the supplied Object is a Crossplane ProviderConfigUsage embedding the core API
// v2 non-typed ProviderConfigUsage.
func ProviderConfigUsageLegacyCore() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsProviderConfigUsageCore().And(fields.IsEmbedded()),
		)
	}
}

// ProviderConfigUsageModernCore returns an Object matcher that returns true if
// the supplied Object is a Crossplane ProviderConfigUsage embedding the core API
// v2 TypedProviderConfigUsage.
func ProviderConfigUsageModernCore() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsObjectMeta().And(fields.IsEmbedded()),
			fields.IsTypedProviderConfigUsageCore().And(fields.IsEmbedded()),
		)
	}
}

// ProviderConfigUsageListLegacy returns an Object matcher that returns true if
// the supplied Object is a list of legacy (non-typed) Crossplane provider config
// usages.
func ProviderConfigUsageListLegacy() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsItems().And(fields.IsSlice()).And(fields.HasFieldThat(
				fields.IsTypeMeta().And(fields.IsEmbedded()),
				fields.IsObjectMeta().And(fields.IsEmbedded()),
				fields.IsProviderConfigUsage().And(fields.IsEmbedded()),
			)),
		)
	}
}

// ProviderConfigUsageListModern returns an Object matcher that returns true if
// the supplied Object is a list of modern (typed) Crossplane provider config
// usages.
func ProviderConfigUsageListModern() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsItems().And(fields.IsSlice()).And(fields.HasFieldThat(
				fields.IsTypeMeta().And(fields.IsEmbedded()),
				fields.IsObjectMeta().And(fields.IsEmbedded()),
				fields.IsTypedProviderConfigUsage().And(fields.IsEmbedded()),
			)),
		)
	}
}

// ProviderConfigUsageListLegacyCore returns an Object matcher that returns true
// if the supplied Object is a list of Crossplane provider config usages
// embedding the core API v2 non-typed ProviderConfigUsage.
func ProviderConfigUsageListLegacyCore() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsItems().And(fields.IsSlice()).And(fields.HasFieldThat(
				fields.IsTypeMeta().And(fields.IsEmbedded()),
				fields.IsObjectMeta().And(fields.IsEmbedded()),
				fields.IsProviderConfigUsageCore().And(fields.IsEmbedded()),
			)),
		)
	}
}

// ProviderConfigUsageListModernCore returns an Object matcher that returns true
// if the supplied Object is a list of Crossplane provider config usages
// embedding the core API v2 TypedProviderConfigUsage.
func ProviderConfigUsageListModernCore() Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta().And(fields.IsEmbedded()),
			fields.IsItems().And(fields.IsSlice()).And(fields.HasFieldThat(
				fields.IsTypeMeta().And(fields.IsEmbedded()),
				fields.IsObjectMeta().And(fields.IsEmbedded()),
				fields.IsTypedProviderConfigUsageCore().And(fields.IsEmbedded()),
			)),
		)
	}
}

// HasMarker returns an Object matcher that returns true if the supplied Object
// has a comment marker k with the value v. Comment markers are read from the
// supplied Comments.
func HasMarker(c comments.Comments, k, v string) Object {
	return func(o types.Object) bool {
		if slices.Contains(comments.ParseMarkers(c.For(o))[k], v) {
			return true
		}

		if slices.Contains(comments.ParseMarkers(c.Before(o))[k], v) {
			return true
		}

		return false
	}
}

// DoesNotHaveMarker returns and Object matcher that returns true if the
// supplied Object does not have a comment marker k with the value v. Comment
// marker are read from the supplied Comments.
func DoesNotHaveMarker(c comments.Comments, k, v string) Object {
	return func(o types.Object) bool {
		return !HasMarker(c, k, v)(o)
	}
}

// AllOf returns an Object matcher that returns true if all of the supplied
// Object matchers return true.
func AllOf(match ...Object) Object {
	return func(o types.Object) bool {
		for _, fn := range match {
			if !fn(o) {
				return false
			}
		}
		return true
	}
}

// AnyOf returns an Object matcher that returns true if any of the supplied
// Object matchers return true.
func AnyOf(match ...Object) Object {
	return func(o types.Object) bool {
		for _, fn := range match {
			if fn(o) {
				return true
			}
		}
		return false
	}
}
