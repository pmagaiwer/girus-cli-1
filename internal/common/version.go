package common

import (
	"fmt"
	"runtime"
)

var (
	Version   string
	BuildUser string
	BuildDate string
	CommitID  string
	GoVersion string
	GoOS      string
	GoArch    string
)

func getDefaultIfEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func GetVersion() string {
	version := getDefaultIfEmpty(Version, "dev")
	buildUser := getDefaultIfEmpty(BuildUser, "unknown")
	buildDate := getDefaultIfEmpty(BuildDate, "unknown")
	commitID := getDefaultIfEmpty(CommitID, "unknown")
	goVersion := getDefaultIfEmpty(GoVersion, runtime.Version())
	goOS := getDefaultIfEmpty(GoOS, runtime.GOOS)
	goArch := getDefaultIfEmpty(GoArch, runtime.GOARCH)
	if Lang() == "es" {
		return fmt.Sprintf(
			"versión de girus-cli: %s\n"+
				"ID de commit: %s\n"+
				"construido por: %s\n"+
				"fecha de construcción: %s\n"+
				"versión de Go: %s\n"+
				"versión de GOOS: %s\n"+
				"versión de GOARCH: %s\n",
			version, commitID, buildUser, buildDate, goVersion, goOS, goArch,
		)
	}

	return fmt.Sprintf(
		"versão do girus-cli: %s\n"+
			"commit ID: %s\n"+
			"build por: %s\n"+
			"data da versão: %s\n"+
			"versão do Go: %s\n"+
			"versão do GOOS: %s\n"+
			"versão do GOARCH: %s\n",
		version, commitID, buildUser, buildDate, goVersion, goOS, goArch,
	)
}
