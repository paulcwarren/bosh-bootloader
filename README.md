# bosh-bootloader
---

This is a command line utility for standing up a CloudFoundry or Concourse installation 
on an IAAS of your choice. This CLI is currently under heavy development, and the
initial goal is to support bootstrapping a CloudFoundry installation on AWS.

* [CI](https://mega.ci.cf-app.com/pipelines/bosh-bootloader)
* [Tracker](https://www.pivotaltracker.com/n/projects/1488988)

## Prerequisites

### Install Dependencies

The following should be installed on your local machine
- Golang >= 1.5 (install with `brew install go`)
- bosh-init ([installation instructions](http://bosh.io/docs/install-bosh-init.html))

### Install bosh-bootloader

bosh-bootloader can be installed with go get:

```
go get github.com/pivotal-cf-experimental/bosh-bootloader/bbl
```

### Configure AWS

The AWS IAM user that is provided to bbl will need the following policy:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:*",
                "cloudformation:*",
                "elasticloadbalancing:*",
                "iam:*"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
```

## Usage

The `bbl` command can be invoked on the command line and will display its usage.

```
$ bbl
Usage:
  bbl [GLOBAL OPTIONS] COMMAND [OPTIONS]

Global Options:
  --help    [-h] "print usage"
  --version [-v] "print version"

  --state-dir             "Directory that stores the state.json"

Commands:
  destroy [--no-confirm]                                                                                               "tears down a BOSH Director environment on AWS"
  director-address                                                                                                     "prints the BOSH director address"
  director-username                                                                                                    "prints the BOSH director username"
  director-password                                                                                                    "prints the BOSH director password"
  bosh-ca-cert                                                                                                         "prints the BOSH director CA certificate"
  env-id                                                                                                               "prints the environment ID"
  help                                                                                                                 "prints usage"
  lbs                                                                                                                  "prints any attached load balancers"
  ssh-key                                                                                                              "prints the ssh private key"
  create-lbs --type=<concourse,cf> --cert=<path> --key=<path> [--chain=<path>] [--skip-if-exists]                      "attaches a load balancer with the supplied certificate, key, and optional chain"
  update-lbs --cert=<path> --key=<path> [--chain=<path>] [--skip-if-missing]                                           "updates a load balancer with the supplied certificate, key, and optional chain"
  delete-lbs [--skip-if-missing]                                                                                       "deletes the attached load balancer"
  up --aws-access-key-id <aws_access_key_id> --aws-secret-access-key <aws_secret_access_key> --aws-region <aws_region> "deploys a BOSH Director on AWS"
  version                                                                                                              "prints version"
```
