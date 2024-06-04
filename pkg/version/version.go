package version

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/gosuri/uitable"
)

var (
	// GitVersion is semantic version
	GitVersion = "v0.0.0-master+$Format:%h$"
	// BuildDate in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ').
	BuildDate = "1970-01-01T00:00:00Z"
	// GitCommit sha1 from git, output of $(git rev-parse HEAD).
	GitCommit = "$Format:%H$"
	// GitTreeState state of git tree, either "clean" or "dirty".
	GitTreeState = ""
)

// Info contains versioning information.
type Info struct {
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// String returns the version information in human readable form.
func (i Info) String() string {
	if s, err := i.Text(); err == nil {
		return string(s)
	}
	return i.GitVersion
}

// ToJSON returns the version information in JSON format.
func (i Info) ToJSON() string {
	s, _ := json.Marshal(i)
	return string(s)
}

// Text returns the version information in plain text.
func (i Info) Text() ([]byte, error) {
	table := uitable.New()
	table.MaxColWidth = 80
	table.RightAlign(0)
	table.Separator = " "
	table.AddRow("gitVersion:", i.GitVersion)
	table.AddRow("gitCommit:", i.GitCommit)
	table.AddRow("gitTreeState:", i.GitTreeState)
	table.AddRow("buildDate:", i.BuildDate)
	table.AddRow("goVersion:", i.GoVersion)
	table.AddRow("compiler:", i.Compiler)
	table.AddRow("platform:", i.Platform)

	return table.Bytes(), nil
}

func Get() Info {
	return Info{
		GitVersion:   GitVersion,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
