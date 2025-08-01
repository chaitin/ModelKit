package version

import "fmt"

var (
	Version   = "v0.0.0"
	BuildTime = ""
	GitCommit = ""
)

type VersionInfo struct{}

func NewVersionInfo() *VersionInfo {
	return &VersionInfo{}
}

func (v *VersionInfo) Print() {
	fmt.Printf("🚀 Starting ModelKit Server\n")
	fmt.Printf("📦 Version:    %s\n", Version)
	fmt.Printf("⏰ BuildTime:  %s\n", BuildTime)
	fmt.Printf("📝 GitCommit:  %s\n", GitCommit)
}

func (v *VersionInfo) Version() string {
	return Version
}

func (v *VersionInfo) BuildTime() string {
	return BuildTime
}

func (v *VersionInfo) GitCommit() string {
	return GitCommit
}
