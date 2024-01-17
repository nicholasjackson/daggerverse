package main

import (
	"context"
	"fmt"
)

type Deb struct{}

// example usage: "call build --directory $PWD/example_app/files --architecture amd64 --package-name 'example' --version '0.0.1' --maintainer 'Nic Jackson' --description 'build'"
func (m *Deb) Build(
	ctx context.Context,
	files *Directory,
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

	// add the control file to the package directory
	packageDirectory := files.WithNewFile("DEBIAN/control", controlContents)

	platform := Platform(fmt.Sprintf("linux/%s", architecture))

	return dag.Container(ContainerOpts{Platform: platform}).
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
