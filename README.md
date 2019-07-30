# func

[![Build Status](https://travis-ci.com/spring-media/func.svg?token=ErJ9PSqPoBz3w7BYQzzq&branch=master)](https://travis-ci.com/spring-media/func) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Func is a CLI app to simplify development and deployment of AWS Lambda functions using Go and Terraform.

## requirements

for installing and developing func

- [Go 1.11+](https://golang.org/)
- `$GOPATH/bin` is [added](https://golang.org/doc/code.html#GOPATH) to `$PATH`

for using generated projects with func

- [Terraform 0.11+](https://www.terraform.io/downloads.html)
- [aws cli](https://docs.aws.amazon.com/cli/latest/userguide/installing.html) with configured [credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) and sufficient IAM permissions for creating/deleting ressources from [terraform-aws-lambda](https://github.com/spring-media/terraform-aws-lambda) module

## installation

(TODO: provide homebrew tap)

clone the repository

```
git clone git@github.com:spring-media/func.git
```

run

```
make install
```

verify your installation

```
$ func

Func is a CLI app to simplify development and deployment
	of serverless functions using Go, Terraform and AWS.

Usage:
  func [command]

Available Commands:
  help        Help about any command
  new         Creates a new Lambda project

Flags:
  -h, --help   help for func

Use "func [command] --help" for more information about a command.
```

## generate new project

(outside of `$GOPATH`)

```
$ func new github.com/you/foo
$ cd foo/
$ make init package deploy
```

## shoulders of giants

func would not be possible if not for all of the great projects it depends on. Please see [SHOULDERS.md](SHOULDERS.md) to see a list of them.
