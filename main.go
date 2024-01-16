package main

import (
	"context"
	"fmt"
	"os"
	"path"
)

type Deb struct{}

// example usage: "call build --binaries $PWD/example_app/bin/example_app --architecture amd64 --package 'example' --standards-version '0.0.1' --maintainer 'Nic Jackson' --description 'build'"
func (m *Deb) Build(
	ctx context.Context,
	binaries *File,
	Architecture string,
	Package string,
	StandardsVersion string,
	Maintainer string,
	Description string,
	Section Optional[string],
	Priority Optional[string],
	BuildDepends Optional[string],
	Depends Optional[string],
	Homepage Optional[string],
) (string, error) {

	controlContents := fmt.Sprintf(controlTemplate,
		Section.GetOr("unknown"),
		Priority.GetOr(""),
		Maintainer,
		BuildDepends.GetOr(""),
		StandardsVersion,
		Homepage.GetOr(""),
		Package,
		Architecture,
		Depends.GetOr(""),
		Description)

	tmpDir, err := os.CreateTemp("", "dagger")
	if err != nil {
		return "", err
	}

	defer os.RemoveAll(tmpDir.Name())

	controlFile := path.Join(tmpDir.Name(), "control")
	os.WriteFile(controlFile, []byte(controlContents), 0644)

	return dag.Container().
		From("alpine:latest").
		WithMountedFile("/etc/bin", binaries).
		WithMountedFile("/etc/DEBIAN/control", dag.Host().File(controlFile)).
		WithWorkdir("/etc").
		Stdout(ctx)
}

var controlTemplate = `
Section: %s
Priority: %s
Maintainer: %s
Build-Depends: %s
Standards-Version: %s
Homepage: %s
Package: %s
Architecture: %s
Depends: %s
Description: %s
`
