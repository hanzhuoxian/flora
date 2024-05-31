package scheme

import (
	"fmt"
	"strings"
)

type ObjectKind interface {
	GroupVersionKind() GroupVersionKind
	SetGroupVersionKind(GroupVersionKind)
}

type emptyObjectKind struct{}

func (emptyObjectKind) GroupVersionKind() GroupVersionKind {
	return GroupVersionKind{}
}

func (emptyObjectKind) SetGroupVersionKind(gvk GroupVersionKind) {
	// no-op
}

var EmptyObjectKind = emptyObjectKind{}

// GroupVersion is a type that holds a Group and a Version.  It is intended to be embedded in other types to provide a group and version.
type GroupKind struct {
	Group string
	Kind  string
}

// Empty return true if Group and Kind both 0.
func (gk GroupKind) Empty() bool {
	return len(gk.Group) == 0 && len(gk.Kind) == 0
}

// String returns the string representation of the GroupKind.
func (gk GroupKind) String() string {
	return gk.Kind + "." + gk.Group
}

func (gk GroupKind) WithVersion(version string) GroupVersionKind {
	return GroupVersionKind{
		Group:   gk.Group,
		Version: version,
		Kind:    gk.Kind,
	}
}

// ParseGroupKind parses a GroupKind from a string.
func ParseGroupKind(gk string) GroupKind {
	i := strings.Index(gk, ".")
	if i == -1 {
		return GroupKind{Kind: gk}
	}

	return GroupKind{Kind: gk[:i], Group: gk[i+1:]}
}

func ParseKindArg(arg string) (*GroupVersionKind, GroupKind) {
	var gvk GroupVersionKind
	if strings.Count(arg, ".") >= 2 {
		s := strings.SplitN(arg, ".", 3)
		gvk = GroupVersionKind{
			Group:   s[2],
			Version: s[1],
			Kind:    s[0],
		}
	}

	return &gvk, ParseGroupKind(arg)
}

type GroupVersion struct {
	Group   string
	Version string
}

// Empty return true if Group and Version both 0.
func (gv GroupVersion) Empty() bool {
	return len(gv.Group) == 0 && len(gv.Version) == 0
}

// String puts "group" and "version" into a single "group/version" string. For the legacy v1
// it returns "v1".
func (gv GroupVersion) String() string {
	if len(gv.Group) > 0 {
		return gv.Group + "/" + gv.Version
	}

	return gv.Version
}

// Identifier implements runtime.GroupVersioner interface.
func (gv GroupVersion) Identifier() string {
	return gv.String()
}

// ParseGroupVersion parses a GroupVersion from a string.
func ParseGroupVersion(gv string) (GroupVersion, error) {
	if len(gv) == 0 || gv == "/" {
		return GroupVersion{}, nil
	}

	switch strings.Count(gv, "/") {
	case 0:
		return GroupVersion{Group: "", Version: gv}, nil
	case 1:
		i := strings.Index(gv, "/")

		return GroupVersion{Group: gv[:i], Version: gv[i+1:]}, nil
	default:
		return GroupVersion{}, fmt.Errorf("unexpected GroupVersion string: %v", gv)
	}
}

// KindForGroupVersionKinds identifies the preferred GroupVersionKind out of a list. It returns ok false.
func (gv GroupVersion) KindForGroupVersionKinds(kinds []GroupVersionKind) (target GroupVersionKind, ok bool) {
	for _, kind := range kinds {
		if kind.Group == gv.Group && kind.Version == gv.Version {
			return kind, true
		}
	}
	for _, kind := range kinds {
		if kind.Group == gv.Group {
			return kind, true
		}
	}

	return GroupVersionKind{}, false
}

type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

func (gvk GroupVersionKind) Empty() bool {
	return len(gvk.Group) == 0 && len(gvk.Version) == 0 && len(gvk.Kind) == 0
}

func (gvk GroupVersionKind) GroupVersion() GroupVersion {
	return GroupVersion{Group: gvk.Group, Version: gvk.Version}
}

func (gvk GroupVersionKind) GroupKind() GroupKind {
	return GroupKind{Group: gvk.Group, Kind: gvk.Kind}
}

type GroupResource struct {
	Group    string
	Resource string
}

// Empty return true if Group and Resource both 0.
func (gr GroupResource) Empty() bool {
	return len(gr.Group) == 0 && len(gr.Resource) == 0
}

// String returns the string representation of the GroupResource.
func (gr GroupResource) String() string {
	if len(gr.Group) == 0 {
		return gr.Resource
	}

	return gr.Resource + "." + gr.Group
}

// WithVersion returns a GroupVersionResource with the provided version.
func (gr GroupResource) WithVersion(version string) GroupVersionResource {
	return GroupVersionResource{
		Group:    gr.Group,
		Version:  version,
		Resource: gr.Resource,
	}
}

func ParseGroupResource(gr string) GroupResource {
	if i := strings.Index(gr, "."); i >= 0 {
		return GroupResource{Group: gr[i+1:], Resource: gr[:i]}
	}

	return GroupResource{Resource: gr}
}

type GroupVersionResource struct {
	Group    string
	Version  string
	Resource string
}
