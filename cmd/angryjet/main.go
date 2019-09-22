package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/negz/angryjet/internal/generators/managed"
)

func main() {
	var (
		app     = kingpin.New(filepath.Base(os.Args[0]), "Generates Crossplane API type methods.").DefaultEnvars()
		pattern = app.Arg("packages", "Package(s) for which to generate Crossplane methods, for example github.com/crossplaneio/crossplane/apis/...").String()
		base    = app.Flag("base-dir", "Generated files are written to their package paths relative to this directory.").Default(filepath.Join(os.Getenv("GOPATH"), "src")).ExistingDir()
		prefix  = app.Flag("prefix", "This string is prepended to the names of all generated files.").Default("zz_generated.").String()
		header  = app.Flag("header-file", "The contents of this file will be added to the top of all generated files.").ExistingFile()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	m := packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax
	pkgs, err := packages.Load(&packages.Config{Mode: m}, *pattern)
	kingpin.FatalIfError(err, "cannot load packages %s", *pattern)

	h, err := ioutil.ReadFile(*header)
	kingpin.FatalIfError(err, "cannot read header file %s", *header)

	for _, p := range pkgs {
		for _, err := range p.Errors {
			kingpin.FatalIfError(err, "error loading packages using pattern %s", *pattern)
		}
		kingpin.FatalIfError(managed.WriteMethods(*base, *prefix, string(h), p), "cannot write managed resource methods for package %s", p.PkgPath)
	}
}
