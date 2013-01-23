// Copyright 2012 Andrey Mirtchovski. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

import (
	"launchpad.net/goamz/aws"
	"launchpad.net/~mirtchovski/goamz/ec2/ec2"
)

type cmd struct {
	cmd   func(e *ec2.EC2, arg ...string)
	Name  string
	Usage string
	Help  string
}

var helpTemplate = `Aws is a tool for managing Amazon Web Services images and instances.

Usage:

    aws command [arguments]

The commands are:
{{range .}}
    {{.Name | printf "%-11s"}} {{.Help}}{{end}}

`

var usageTemplate = `Usage:

{{range .}}
    aws {{.Name | printf "%-11s"}} {{.Usage}}{{end}}

`

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

var commands = []*cmd{
	&cmd{create, "create", "image_id [instance_type]", "create an instance from image_id and instance_type (default=t1.micro)"},
	&cmd{cloneImage, "clone", "id [name=... description=... noreboot=t/f]", "create an ami from an EBS-backed running instance"},
	&cmd{stop, "stop", "id [...]", "stop one or more instances"},
	&cmd{resume, "resume", "id [...]", "resume one or more instances"},
	&cmd{destroy, "destroy", "id [...]", "terminate one or more instances"},
	&cmd{instances, "instances", "[key=value]", "list instances matching key=value, available filters: http://goo.gl/4No7c"},
	&cmd{images, "images", "[-all] [key=value]", "list our images matching key=value, -all for all available filters: http://goo.gl/SRBhW"},
	&cmd{deregister, "deregister", "image_id", "deregister an image. nobody will be able to boot from it afterwards"},
	&cmd{snapshots, "snapshots", "[key=value]", "list our snapshots matchng key=value"},
	&cmd{snapshot, "snapshot", "image_id", "create snapshot from image_id"},
	&cmd{snapdel, "delete", "snapshot_id ...", "delete one or more snapshots"},
	&cmd{regions, "regions", "", "list all available regions"},
	&cmd{nil, "help", "", "print this message"},
	&cmd{nil, "usage", "[command]", "print usage information"},
}

func regions(*ec2.EC2, ...string) {
	for _, v := range aws.Regions {
		fmt.Printf("%s\tec2: %s\ts3: %s\n", v.Name, v.EC2Endpoint, v.S3Endpoint)
	}
}

var debug = flag.Bool("d", false, "print debug info")
var region = flag.String("zone", "us-east-1", "service (eu-west-1, us-east-1, etc)")

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "usage: %s cmd [args]; 'help' will print information about each command\n", os.Args[0])
		os.Exit(1)
	}

	args := flag.Args()
	var ucmd *cmd
	for _, v := range commands {
		if args[0] == v.Name {
			ucmd = v
		}
	}
	if ucmd == nil {
		fmt.Fprintf(os.Stderr, "unknown command: %s; 'help' will give more information on available commands\n", args[0])
		os.Exit(1)
	}
	args = args[1:]

	if ucmd.Name == "help" {
		tmpl(os.Stderr, helpTemplate, commands)
		os.Exit(1)
	}
	if ucmd.Name == "usage" {
		cmd := make([]*cmd, 0, 5)
		if len(args) > 0 {
			for _, v1 := range args {
				for _, v2 := range commands {
					if v1 == v2.Name {
						cmd = append(cmd, v2)
					}
				}
			}
		} else {
			cmd = commands
		}
		tmpl(os.Stderr, usageTemplate, cmd)
		os.Exit(1)
	}

	r, ok := aws.Regions[*region]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown region: %s (aws regions to list all available)\n", *region)
		os.Exit(1)
	}

	env := os.Getenv("AWS_ACCESS_KEY_ID")
	if env == "" {
		fmt.Fprintf(os.Stderr, "set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY\n")
		os.Exit(1)
	}

	auth, err := aws.EnvAuth()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not find auth info: %s\n", err)
		os.Exit(1)
	}

	e := ec2.New(auth, r)

	ucmd.cmd(e, args...)

	return
}
