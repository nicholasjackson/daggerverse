package main

import (
	"context"
	"fmt"
)

type Deb struct{}

// example usage: "call build --binaries $PWD/example_app/bin/example_app --architecture amd64 --package-name 'example' --version '0.0.1' --maintainer 'Nic Jackson' --description 'build'"
func (m *Deb) Build(
	ctx context.Context,
	binaries *File,
	architecture string,
	packageName string,
	version string,
	maintainer string,
	description string,
	section Optional[string],
	priority Optional[string],
	depends Optional[string],
	homepage Optional[string],
) (*File, error) {

	controlContents := fmt.Sprintf(controlTemplate,
		section.GetOr("unknown"),
		priority.GetOr(""),
		maintainer,
		version,
		homepage.GetOr(""),
		packageName,
		architecture,
		depends.GetOr(""),
		description,
	)

	packageDirectory := dag.Directory().WithNewDirectory(packageName)
	packageDirectory = packageDirectory.WithNewFile("DEBIAN/control", controlContents)
	packageDirectory = packageDirectory.WithFile("bin/example", binaries)

	return dag.Container().
		From("debian:latest").
		WithMountedDirectory("/working/package", packageDirectory).
		WithWorkdir("/working").
		WithExec([]string{"dpkg-deb", "--root-owner-group", "--build", "package"}).
		File("/working/package.deb"), nil
}

var controlTemplate = `
Section: %s
Priority: %s
Maintainer: %s
Version: %s
Homepage: %s
Package: %s
Architecture: %s
Depends: %s
Description: %s
Package-Type: deb
`
