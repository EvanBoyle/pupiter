package write

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/evanboyle/pupiter/session"
	"github.com/pkg/errors"
)

// TODO write out to a given src directory in ~/.pulumi/pupiter

// fmt with session (maps to project)
const pulumiYaml = `
name: %s
runtime: nodejs
description: A minimal AWS JavaScript Pulumi program
`

// fmt with userinput, varName
const indexjs = `
"use strict";
const pulumi = require("@pulumi/pulumi");
const aws = require("@pulumi/aws");
const awsx = require("@pulumi/awsx");

// ---------------------------------------------

%s // user program

// ---------------------------------------------

for (var propName in %[2]s) {
    if(%[2]s.hasOwnProperty(propName)) {
        if (typeof %[2]s[propName] === "object") {
            %[2]s[propName].isSecret = Promise.resolve(false);
        }
       
    }
}
exports.%[2]s = %[2]s;

`

// fmt with session
const packageJSON = `
{
    "name": %q,
    "main": "index.js",
    "dependencies": {
        "@pulumi/pulumi": "^2.0.0",
        "@pulumi/aws": "^2.0.0",
        "@pulumi/awsx": "^0.20.0"
    }
}
`

// fmt name to be Pulumi.<varName>.yaml
const pulumiStackYaml = `
config:
  aws:region: us-west-2
`

func writePulumiYaml(dir, sessionName string) error {
	fname := filepath.Join(dir, "Pulumi.yaml")
	text := fmt.Sprintf(pulumiYaml, sessionName)
	err := ioutil.WriteFile(fname, []byte(text), 0777)
	if err != nil {
		return err
	}
	return nil
}

func writeIndexJS(dir, input, varName string) error {
	fname := filepath.Join(dir, "index.js")
	text := fmt.Sprintf(indexjs, input, varName)
	err := ioutil.WriteFile(fname, []byte(text), 0777)
	if err != nil {
		return err
	}
	return nil
}

func writePulumiStackYaml(dir, varName string) error {
	suffix := fmt.Sprintf("Pulumi.%s.yaml", varName)
	fname := filepath.Join(dir, suffix)
	err := ioutil.WriteFile(fname, []byte(pulumiStackYaml), 0777)
	if err != nil {
		return err
	}
	return nil
}

func writePackageJSON(dir, sessionName string) error {
	fname := filepath.Join(dir, "package.json")
	text := fmt.Sprintf(packageJSON, sessionName)
	err := ioutil.WriteFile(fname, []byte(text), 0777)
	if err != nil {
		return err
	}
	return nil
}

func symlinkNodeModules(parentDir, targetDir string) error {
	nodeMods := filepath.Join(filepath.Dir(parentDir), "node_modules")
	if _, err := os.Lstat(filepath.Join(targetDir, "node_modules")); err == nil {
		return nil
	}

	err := os.Symlink(nodeMods, filepath.Join(targetDir, "node_modules"))
	if err != nil {
		return err
	}
	return nil
}

func WriteFiles(varName, text string, session session.Session) error {
	targetDir := path.Join(session.Dir(), varName)
	err := os.MkdirAll(targetDir, 0777)
	if err != nil {
		return err
	}
	err = writePulumiYaml(targetDir, session.Name())
	if err != nil {
		return errors.Wrap(err, "unable to write Pulumi.yaml")
	}
	err = writeIndexJS(targetDir, text, varName)
	if err != nil {
		return errors.Wrap(err, "unable to write index.js")
	}
	err = writePulumiStackYaml(targetDir, varName)
	if err != nil {
		return errors.Wrap(err, "unable to write Pulumi.<stack>.yaml")
	}
	err = writePackageJSON(targetDir, session.Name())
	if err != nil {
		return errors.Wrap(err, "unable to write package.json")
	}
	err = symlinkNodeModules(session.Dir(), targetDir)
	if err != nil {
		return errors.Wrap(err, "unable symlink node_modules")
	}
	return nil
}
