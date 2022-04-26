// +build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/mcandre/slick"
	"github.com/mcandre/mage-extras"
)

// artifactsPath describes where artifacts are produced.
var artifactsPath = "bin"

// Default references the default build task.
var Default = Test

// Audit runs a security audit.
func Audit() error { return mageextras.SnykTest() }

// UnitTests runs the unit test suite.
func UnitTest() error { return mageextras.UnitTest() }

// IntegrationTest executes the integration test suite.
func IntegrationTest() error {
	mg.Deps(Install)

	cmd := exec.Command("slick", "-n", "examples")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err == nil {
		return errors.New("Expected an exit status error")
	}

	return nil
}

// Text runs unit and integration tests.
func Test() error { mg.Deps(UnitTest); mg.Deps(IntegrationTest); return nil }

// CoverHTML denotes the HTML formatted coverage filename.
var CoverHTML = "cover.html"

// CoverProfile denotes the raw coverage data filename.
var CoverProfile = "cover.out"

// CoverageHTML generates HTML formatted coverage data.
func CoverageHTML() error { mg.Deps(CoverageProfile); return mageextras.CoverageHTML(CoverHTML, CoverProfile) }

// CoverageProfile generates raw coverage data.
func CoverageProfile() error { return mageextras.CoverageProfile(CoverProfile) }

// GoVet runs go vet with shadow checks enabled.
func GoVet() error { return mageextras.GoVetShadow() }

// GoLint runs golint.
func GoLint() error { return mageextras.GoLint() }

// Gofmt runs gofmt.
func GoFmt() error { return mageextras.GoFmt("-s", "-w") }

// GoImports runs goimports.
func GoImports() error { return mageextras.GoImports("-w") }

// Errcheck runs errcheck.
func Errcheck() error { return mageextras.Errcheck("-blank") }

// Nakedret runs nakedret.
func Nakedret() error { return mageextras.Nakedret("-l", "0") }

// Lint runs the lint suite.
func Lint() error {
	mg.Deps(GoVet)
	mg.Deps(GoLint)
	mg.Deps(GoFmt)
	mg.Deps(GoImports)
	mg.Deps(Errcheck)
	mg.Deps(Nakedret)
	return nil
}

// portBasename labels the artifact basename.
var portBasename = fmt.Sprintf("slick-%s", slick.Version)

// repoNamespace identifies the Go namespace for this project.
var repoNamespace = "github.com/mcandre/slick"

// Factorio cross-compiles Go binaries for a multitude of platforms.
func Factorio() error { return mageextras.Factorio(portBasename) }

// Port builds and compresses artifacts.
func Port() error { mg.Deps(Factorio); return mageextras.Archive(portBasename, artifactsPath) }

// Install builds and installs Go applications.
func Install() error { return mageextras.Install() }

// Uninstall deletes installed Go applications.
func Uninstall() error { return mageextras.Uninstall("slick") }

// CleanCoverage deletes coverage data.
func CleanCoverage() error {
	if err := os.RemoveAll(CoverHTML); err != nil {
		return err
	}

	return os.RemoveAll(CoverProfile)
}

// Clean deletes artifacts.
func Clean() error { mg.Deps(CleanCoverage); return os.RemoveAll(artifactsPath) }
