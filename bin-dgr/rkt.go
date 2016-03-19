package main

import (
	"github.com/blablacar/dgr/bin-dgr/common"
	"github.com/blablacar/dgr/dist"
	"github.com/n0rad/go-erlog/logs"
	"io/ioutil"
	"os"
)

var ACI_BUILDER = common.NewACFullName("dgrtool.com/aci-builder:1")
var ACI_TESTER = common.NewACFullName("dgrtool.com/aci-tester:1")

//var internalAcis = []*spec.ACFullname{ACI_BATS, ACI_BUILDER}

func ImportInternalBuilderIfNeeded(manifest *AciManifest) {
	if manifest.Builder.Image.String() == "" {
		manifest.Builder.Image = *ACI_BUILDER
		importInternalAci("aci-builder.aci") // TODO
	}
}

func ImportInternalTesterIfNeeded(manifest *AciManifest) {
	if manifest.TestBuilder.String() == "" {
		manifest.TestBuilder = *ACI_TESTER
		importInternalAci("aci-tester.aci") // TODO
	}
}

func importInternalAci(filename string) {
	filepath := "dist/bindata/" + filename
	content, err := dist.Asset(filepath)
	if err != nil {
		logs.WithE(err).WithField("aci", filepath).Fatal("Cannot found internal aci")
	}
	if err := ioutil.WriteFile("/tmp/tmp.aci", content, 0644); err != nil {
		logs.WithE(err).WithField("aci", filepath).Fatal("Failed to write tmp aci to /tmp/tmp.aci")
	}
	if _, err := Home.Rkt.Fetch("/tmp/tmp.aci"); err != nil { // TODO does not support concurrency
		logs.WithE(err).Fatal("Failed to import internal image to rkt")
	}
	os.Remove("/tmp/tmp.aci")
	return
}