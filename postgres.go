package barkup

import (
	"fmt"
	"os/exec"
	"time"
)

var (
	// TarCmd is the path to the `tar` executable
	TarCmd = "tar"
	// PGDumpCmd is the path to the `pg_dump` executable
	PGDumpCmd = "pg_dump"
)

// Postgres is an `Exporter` interface that backs up a Postgres database via the `pg_dump` command
type Postgres struct {
	// DB Host (e.g. 127.0.0.1)
	Host string
	// DB Port (e.g. 5432)
	Port string
	// DB Name
	DB string
	// Connection Username
	Username string
	// Extra pg_dump options
	// e.g []string{"--inserts"}
	Options []string
}

// Export produces a `pg_dump` of the specified database, and creates a gzip compressed tarball archive.
func (x Postgres) Export() *ExportResult {
	result := &ExportResult{MIME: "application/x-tar"}
	dumpPath := fmt.Sprintf(`%v_%v.sql`, x.DB, time.Now().Unix())
	options := append(x.dumpOptions(), "-Fp", fmt.Sprintf(`-r%v`, dumpPath))
	out, err := exec.Command(PGDumpCmd, options...).Output()
	if err != nil {
		result.Error = makeErr(err, string(out))
		return result
	}

	result.Path = dumpPath + ".tar.gz"
	_, err = exec.Command(TarCmd, "-czf", result.Path, dumpPath).Output()
	if err != nil {
		result.Error = makeErr(err, string(out))
		return result
	}

	return result
}

func (x Postgres) dumpOptions() []string {
	options := x.Options

	if x.DB != "" {
		options = append(options, fmt.Sprintf(`-d%v`, x.DB))
	}

	if x.Host != "" {
		options = append(options, fmt.Sprintf(`-h%v`, x.Host))
	}

	if x.Port != "" {
		options = append(options, fmt.Sprintf(`-p%v`, x.Port))
	}

	if x.Username != "" {
		options = append(options, fmt.Sprintf(`-U%v`, x.Username))
	}

	return options
}
