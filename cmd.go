// Copyright 2012 Andrey Mirtchovski. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"
)

import (
	//	"launchpad.net/goamz/aws"
	"launchpad.net/~mirtchovski/goamz/ec2/ec2"
)

func create(e *ec2.EC2, args ...string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "create: need image id\n")
		os.Exit(1)
	}

	options := ec2.RunInstances{
		ImageId:      args[0],
		InstanceType: "t1.micro",
	}

	if len(args) == 2 {
		options.InstanceType = args[1]
	}

	resp, err := e.RunInstances(&options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create: %s\n", err)
		os.Exit(1)
	}

	for _, instance := range resp.Instances {
		fmt.Println(instance.InstanceId)
	}
}

func stop(e *ec2.EC2, args ...string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "stop: need instance id\n")
		os.Exit(1)
	}

	for _, v := range args {
		resp, err := e.StopInstances(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "stop: %s\n", err)
			os.Exit(1)
		}

		for _, r := range resp.StateChanges {
			fmt.Printf("state change: %s → %s\n", r.PreviousState.Name, r.CurrentState.Name)
		}
	}
}

func resume(e *ec2.EC2, args ...string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "resume: need instance id\n")
		os.Exit(1)
	}

	resp, err := e.StartInstances(args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "stop: %s\n", err)
		os.Exit(1)
	}

	for _, r := range resp.StateChanges {
		fmt.Printf("state change: %s → %s\n", r.PreviousState.Name, r.CurrentState.Name)
	}

}

func destroy(e *ec2.EC2, args ...string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "destroy: need instance id\n")
		os.Exit(1)
	}

	nresp, err := e.TerminateInstances(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "destroy: %s\n", err)
		os.Exit(1)
	}

	for _, r := range nresp.StateChanges {
		fmt.Printf("state change: %s → %s\n", r.PreviousState.Name, r.CurrentState.Name)
	}

}

func cloneImage(e *ec2.EC2, args ...string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "clone: need a instance id and options\n")
		os.Exit(1)
	}

	options := &ec2.CreateImage{}

	options.InstanceId = args[0]
	for _, v := range args[1:] {
		s := strings.SplitN(v, "=", 2)
		if len(s) != 2 {
			fmt.Fprintf(os.Stderr, "images: bad key=value pair \"%s\", skipping\n", v)
			continue
		}
		switch s[0] {
		case "name", "Name":
			options.Name = s[1]
		case "description", "Description":
			options.Description = s[1]
		case "NoReboot", "Noreboot":
			switch s[1] {
			case "true", "True", "TRUE", "t", "T":
				options.NoReboot = true
			case "false", "False", "FALSE", "f", "F":
			default:
				fmt.Fprintf(os.Stderr, "createimage: expected option for NoReboot true/false, found %s", s[1])
				return
			}
		default:
			fmt.Fprintf(os.Stderr, "createimage: unknown key/value pair: %s (expected name=n, description=d, noreboot=t|f. skipping...", s[1])

		}
	}

	resp, err := e.CreateImage(options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "createimage error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", resp.Id)
}
func instances(e *ec2.EC2, args ...string) {
	filter := ec2.NewFilter()
	for _, v := range args {
		sl := strings.SplitN(v, "=", 2)
		if len(sl) != 2 {
			fmt.Fprintf(os.Stderr, "instances: bad key=value pair \"%s\", skipping\n", v)
			continue
		}
		filter.Add(sl[0], sl[1])
	}

	resp, err := e.Instances(nil, filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "instances: %s\n", err)
		os.Exit(1)
	}

	for _, r := range resp.Reservations {
		fmt.Println("reservation:", r.ReservationId)
		for _, i := range r.Instances {
			fmt.Printf("%s\t%s\t%s\t%s\n", i.InstanceId, i.State.Name, i.DNSName, i.ImageId)
		}
	}
}

func images(e *ec2.EC2, args ...string) {
	isAll := false
	if len(args) > 0 {
		if args[0] == "-all" {
			isAll = true
			args = args[1:]
		}
	}

	filter := ec2.NewFilter()
	if !isAll {
		filter.Add("is-public", "false")
	}
	for _, v := range args {
		s := strings.SplitN(v, "=", 2)
		if len(s) != 2 {
			fmt.Fprintf(os.Stderr, "images: bad key=value pair \"%s\", skipping\n", v)
			continue
		}
		filter.Add(s[0], s[1])
	}

	resp, err := e.Images(nil, filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "images: %s\n", err)
		os.Exit(1)
	}

	for _, i := range resp.Images {
		fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", i.Id, i.Name, i.State, i.Architecture, i.Type, i.OwnerAlias, i.Description, i.Platform)
	}
}

func snapshots(e *ec2.EC2, args ...string) {
	isAll := false
	if len(args) > 0 {
		if args[0] == "-all" {
			isAll = true
			args = args[1:]
		}
	}

	filter := ec2.NewFilter()
	if !isAll {
		filter.Add("owner-id", "self")
	}
	for _, v := range args {
		sl := strings.SplitN(v, "=", 2)
		if len(sl) != 2 {
			fmt.Fprintf(os.Stderr, "snapshots: bad key=value pair \"%s\", skipping\n", v)
			continue
		}
		filter.Add(sl[0], sl[1])
	}

	resp, err := e.Snapshots(nil, filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "snapshots: %s\n", err)
		os.Exit(1)
	}

	for _, i := range resp.Snapshots {
		fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			i.Id,
			i.VolumeId,
			i.Status,
			i.StartTime,
			i.Progress,
			i.OwnerId,
			i.VolumeSize,
			i.Description,
			i.OwnerAlias,
		)
	}
}

func snapshot(e *ec2.EC2, args ...string) {
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "snapshot: need image id and description\n")
		os.Exit(1)
	}

	resp, err := e.CreateSnapshot(args[0], args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "snapshot: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		resp.Id,
		resp.VolumeId,
		resp.Status,
		resp.StartTime,
		resp.Progress,
		resp.OwnerId,
		resp.VolumeSize,
		resp.Description,
		resp.OwnerAlias,
	)

}

func snapdel(e *ec2.EC2, args ...string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "snapdel: need snapshot id\n")
		os.Exit(1)
	}

	_, err := e.DeleteSnapshots(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "snapdel: %s\n", err)
		os.Exit(1)
	}
}

func deregister(e *ec2.EC2, args ...string) {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "deregister: need image id\n")
		os.Exit(1)
	}

	_, err := e.DeregisterImage(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "deregister: %s\n", err)
		os.Exit(1)
	}
}
