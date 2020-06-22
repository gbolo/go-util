package backend

import (
	"fmt"
	"runtime"
)

var (
	// !! These variables defined by the Makefile and passed in with ldflags !!
	// !! DO NOT CHANGE THESE DEFAULT VALUES !!

	// Version of application
	Version = "devel"
	// CommitSHA is the short SHA hash of the git commit
	CommitSHA = "unknown"
	// BuildDate is the date this application was compiled
	BuildDate = "unknown"
)

type appInfo struct {
	appVersion string
	commitRef  string
	goVersion  string
	goPlatform string
}

// PrintVersion prints the current version information to stdout
func PrintVersion() {
	fmt.Printf(`vme-portal:
  version     : %s
  build date  : %s
  git hash    : %s
  go version  : %s
  go compiler : %s
  platform    : %s/%s
`, Version, BuildDate, CommitSHA, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
}

func getVersionResponse() versionInfo {
	return versionInfo{
		Version:   Version,
		CommitSHA: CommitSHA,
		BuildDate: BuildDate,
	}
}

func getAppInfo() appInfo {
	info := appInfo{
		appVersion: Version,
		commitRef:  CommitSHA,
		goVersion:  runtime.Version(),
		goPlatform: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
	if info.commitRef == "" {
		info.commitRef = "unknown"
	}
	if info.appVersion == "" {
		info.appVersion = "devel"
	}
	return info
}

func (i appInfo) summary() string {
	return fmt.Sprintf("version %s (ref %s) %s [%s]", i.appVersion, i.commitRef, i.goVersion, i.goPlatform)
}
