package builder

import (
	"github.com/blablacar/cnt/cnt"
	"github.com/blablacar/cnt/spec"
	"github.com/ghodss/yaml"
	"github.com/juju/errors"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/logs"
	"io/ioutil"
	"path/filepath"
)

const POD_MANIFEST = "/cnt-pod-manifest.yml"

type Pod struct {
	fields   data.Fields
	path     string
	args     BuildArgs
	target   string
	manifest spec.PodManifest
}

func NewPod(path string, args BuildArgs) (*Pod, error) {
	fullPath, err := filepath.Abs(path)
	if err != nil {
		logs.WithE(err).WithField("path", path).Fatal("Cannot get fullpath")
	}

	manifest, err := readPodManifest(fullPath + POD_MANIFEST)
	if err != nil {
		return nil, errors.Annotate(err, "Failed to read pod manifest")
	}
	fields := data.WithField("pod", manifest.Name.String())

	target := path + PATH_TARGET
	if cnt.Home.Config.TargetWorkDir != "" {
		currentAbsDir, err := filepath.Abs(cnt.Home.Config.TargetWorkDir + "/" + manifest.Name.ShortName())
		if err != nil {
			logs.WithEF(err, fields).Panic("invalid target path")
		}
		target = currentAbsDir
	}

	pod := &Pod{
		fields:   fields,
		path:     fullPath,
		args:     args,
		target:   target,
		manifest: *manifest,
	}

	return pod, nil
}

func readPodManifest(manifestPath string) (*spec.PodManifest, error) {
	source, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	manifest := &spec.PodManifest{}
	err = yaml.Unmarshal([]byte(source), manifest)
	if err != nil {
		return nil, err
	}

	for i, app := range manifest.Pod.Apps {
		if app.Name == "" {
			manifest.Pod.Apps[i].Name = app.Dependencies[0].ShortName()
		}
	}
	//TODO check that there is no app name conflict
	return manifest, nil
}

func (p *Pod) toAciManifest(e spec.RuntimeApp) spec.AciManifest {
	fullname := spec.NewACFullName(p.manifest.Name.Name() + "_" + e.Name + ":" + p.manifest.Name.Version())
	return spec.AciManifest{
		Aci: spec.AciDefinition{
			Annotations:   e.Annotations,
			App:           e.App,
			Dependencies:  e.Dependencies,
			PathWhitelist: nil, // TODO
		},
		NameAndVersion: *fullname,
	}
}
