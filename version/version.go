package version

import (
	"github.com/storageos/discovery/types"
)

// api
const (
	ProductName string = "StorageOS Discovery"
	APIVersion         = "1"
)

// Revision that was compiled. This will be filled in by the compiler.
var Revision string

// BuildDate is when the binary was compiled.  This will be filled in by the
// compiler.
var BuildDate string

// Version number that is being run at the moment.  Version should use semver.
var Version string

// Experimental is intended to be used to enable alpha features.
var Experimental string

// GetVersion returns version info.
func GetVersion() types.VersionInfo {
	v := types.VersionInfo{
		Name:       ProductName,
		Version:    Version,
		APIVersion: APIVersion,
		BuildDate:  BuildDate,
	}
	return v
}
