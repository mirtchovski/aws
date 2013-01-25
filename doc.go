// Copyright 2012 Andrey Mirtchovski. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Aws manages Elastic Compute Cloud instances and related Amazon Web Services
components.

The aws command assumes that AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are
set in the environment before running the program.

Install:
	go get github.com/mirtchovski/aws

Usage:

    aws command [arguments]

The commands are:

    create      create an instance from image_id and instance_type (default=t1.micro)
    clone       create an ami from an EBS-backed running instance
    stop        stop one or more instances
    resume      resume one or more instances
    destroy     terminate one or more instances
    instances   list instances matching key=value, available filters: http://goo.gl/4No7c
    images      list our images matching key=value, -all for all available filters: http://goo.gl/SRBhW
    deregister  deregister an image. nobody will be able to boot from it afterwards
    snapshots   list our snapshots matching key=value
    snapshot    create snapshot from image
    delete      delete one or more snapshots
    regions     list all available regions
    help        print this message
    usage       print usage information

Each command has individual usage patterns available through "aws usage command" or "aws usage":

    aws create      image_id [instance_type]
    aws clone       id [name=... description=... noreboot=t/f]
    aws stop        id [...]
    aws resume      id [...]
    aws destroy     id [...]
    aws instances   [key=value]
    aws images      [-all] [key=value]
    aws deregister  image_id
    aws snapshots   [key=value]
    aws snapshot    image_id
    aws delete      snapshot_id ...
    aws regions
    aws help
    aws usage       [command]


Examples

Find the currently available instances:

    $ ./aws instances
    reservation: r-11xxxxxx
    i-aa000000	running	ec2-1.amazonaws.com	ami-xxxxxxxx
    reservation: r-22xxxxxx
    i-bb000000	stopped		ami-yyyyyyyy
    reservation: r-33xxxxxx
    i-cc000000	running	ec2-2.amazonaws.com	ami-zzzzzzzz

Find all images belonging to us:

    $ ./aws images
    ami-xxxxxxxx	GoogleDocs	available	x86_64	machine		Browser for Google Docs
    ami-yyyyyyyy	GoogleEarth	available	x86_64	machine		Google Earth
    ami-zzzzzzzz	OpenOffice	available	x86_64	machine		Open Office

Clone the running instance i-cc000000's image:

    $ ./aws clone i-cc000000 name=test description='this is a test'
    ami-ffffffff

Destroy the newly created image:

    $ ./aws deregister ami-ffffffff

*/
package documentation

// Bug(aam): Incomplete. AWS has way more commands than we need.
