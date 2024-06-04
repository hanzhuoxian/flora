package version

import (
	"github.com/hanzhuoxian/flora/internal/floractl/util/templates"
	"github.com/hanzhuoxian/flora/pkg/cli/options"
	"github.com/hanzhuoxian/flora/pkg/version"
)

type Version struct {
	ClientVersion *version.Info `json:"clientVersion,omitempty" yaml:"clientVersion,omitempty"`
	ServerVersion *version.Info `json:"serverVersion,omitempty" yaml:"serverVersion,omitempty"`
}

var versionExample = templates.Examples(`
# Print the client and server versions for the current context
iamctl version`)

type Options struct {
	ClientOnly bool
	Short      bool
	Output     string
	options.IOStreams
}
