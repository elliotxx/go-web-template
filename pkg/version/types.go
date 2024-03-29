//go:generate go run gen.go
//go:generate go fmt

package version

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	git "github.com/elliotxx/go-web-template/pkg/util/gitutil"
	goversion "github.com/hashicorp/go-version"
	"gopkg.in/yaml.v3"
)

var info = NewMainOrDefaultVersionInfo()

const EnvSpecifiedVersion = "SPECIFIED_VERSION"

func NewMainOrDefaultVersionInfo() *Info {
	v := NewDefaultVersionInfo()

	if i, ok := debug.ReadBuildInfo(); ok {
		mod := &i.Main
		if mod.Replace != nil {
			mod = mod.Replace
		}

		if mod.Version != "(devel)" {
			v.ReleaseVersion = mod.Version
		}
	}

	return v
}

func NewDefaultVersionInfo() *Info {
	return &Info{
		ReleaseVersion: "default-version",
		GitInfo: &GitInfo{
			LatestTag: "",
			Commit:    "",
			TreeState: "",
		},
		BuildInfo: &BuildInfo{
			GoVersion: runtime.Version(),
			GOOS:      runtime.GOOS,
			GOARCH:    runtime.GOARCH,
			NumCPU:    runtime.NumCPU(),
			Compiler:  runtime.Compiler,
			BuildTime: time.Now().Format("2006-01-02 15:04:05"),
		},
	}
}

// Info contains versioning information.
// following attributes:
//
//	   ReleaseVersion - "vX.Y.Z-00000000" used to indicate the last release version,
//			  containing GitVersion and GitCommitShort.
type Info struct {
	ReleaseVersion string     `json:"releaseVersion" yaml:"releaseVersion"` // Such as "v1.2.3-3836f877"
	GitInfo        *GitInfo   `json:"gitInfo,omitempty" yaml:"gitInfo,omitempty"`
	BuildInfo      *BuildInfo `json:"buildInfo,omitempty" yaml:"buildInfo,omitempty"`
}

// GitInfo contains git information.
// following attributes:
//
//	LatestTag - "vX.Y.Z" used to indicate the last git tag.
//	Commit - The git commit id corresponding to this source code.
//	TreeState - "clean" indicates no changes since the git commit id
//	    "dirty" indicates source code changes after the git commit id
type GitInfo struct {
	LatestTag string `json:"latestTag,omitempty" yaml:"latestTag,omitempty"` // Such as "v1.2.3"
	Commit    string `json:"commit,omitempty" yaml:"commit,omitempty"`       // Such as "3836f8770ab8f488356b2129f42f2ae5c1134bb0"
	TreeState string `json:"treeState,omitempty" yaml:"treeState,omitempty"` // Such as "clean", "dirty"
}

type BuildInfo struct {
	GoVersion string `json:"goVersion,omitempty" yaml:"goVersion,omitempty"`
	GOOS      string `json:"GOOS,omitempty" yaml:"GOOS,omitempty"`
	GOARCH    string `json:"GOARCH,omitempty" yaml:"GOARCH,omitempty"`
	NumCPU    int    `json:"numCPU,omitempty" yaml:"numCPU,omitempty"`
	Compiler  string `json:"compiler,omitempty" yaml:"compiler,omitempty"`
	BuildTime string `json:"buildTime,omitempty" yaml:"buildTime,omitempty"` // Such as "2021-10-20 18:24:03"
}

func NewInfo() (*Info, error) {
	var (
		isHeadAtTag    bool
		headHash       string
		headHashShort  string
		latestTag      string
		gitVersion     *goversion.Version
		releaseVersion string
		isDirty        bool
		gitTreeState   string
		err            error
	)

	// Get git info
	if headHash, err = git.GetHeadHash(); err != nil {
		return nil, err
	}

	if headHashShort, err = git.GetHeadHashShort(); err != nil {
		return nil, err
	}

	if latestTag, err = git.GetLatestTag(); err != nil {
		return nil, err
	}

	if gitVersion, err = goversion.NewVersion(latestTag); err != nil {
		return nil, err
	}

	if isHeadAtTag, err = git.IsHeadAtTag(latestTag); err != nil {
		return nil, err
	}

	if isDirty, err = git.IsDirty(); err != nil {
		return nil, err
	}

	// Get git tree state
	if isDirty {
		gitTreeState = "dirty"
	} else {
		gitTreeState = "clean"
	}

	// Get release version
	if specifiedVersion := os.Getenv(EnvSpecifiedVersion); strings.TrimSpace(specifiedVersion) == "" {
		if isHeadAtTag {
			releaseVersion = gitVersion.Original()
		} else {
			releaseVersion = fmt.Sprintf("%s-%s", gitVersion.Original(), headHashShort)
		}
	} else {
		releaseVersion = specifiedVersion
	}

	return &Info{
		ReleaseVersion: releaseVersion,
		GitInfo: &GitInfo{
			LatestTag: gitVersion.Original(),
			Commit:    headHash,
			TreeState: gitTreeState,
		},
		BuildInfo: &BuildInfo{
			GoVersion: runtime.Version(),
			GOOS:      runtime.GOOS,
			GOARCH:    runtime.GOARCH,
			NumCPU:    runtime.NumCPU(),
			Compiler:  runtime.Compiler,
			BuildTime: time.Now().Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (v *Info) String() string {
	return v.YAML()
}

func (v *Info) ShortString() string {
	return fmt.Sprintf("%s; git: %s; build time: %s", v.ReleaseVersion, v.GitInfo.Commit, v.BuildInfo.BuildTime)
}

func (v *Info) JSON() string {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return ""
	}

	return string(data)
}

func (v *Info) YAML() string {
	data, err := yaml.Marshal(v)
	if err != nil {
		return ""
	}

	return string(data)
}
