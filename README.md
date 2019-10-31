<p align="center"><img src="https://github.com/spring-media/func/blob/master/logo.png" width="360"></p>

# func [![Build Status](https://travis-ci.com/spring-media/func.svg?token=ErJ9PSqPoBz3w7BYQzzq&branch=master)](https://travis-ci.com/spring-media/func) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Func is a CLI app to simplify development and deployment of AWS Lambda functions using Go and Terraform. It'll scaffold an optionated project structure generating code for

- function code
- build automation using make
- infrastructure and deployment automation using [terraform-aws-lambda](https://github.com/spring-media/terraform-aws-lambda)
- continuous integration/deployment providers like Travis

Func is in an early alpha stage so expect bugs and breaking changes but give it a try!

## installation

Before installing `func` please make sure your system meets the following requirements:

- a working Go environment ([Go 1.11+](https://golang.org/))
- a working terraform environment ([Terraform 0.11+](https://www.terraform.io/downloads.html))
- configured [AWS credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) with sufficient IAM permissions for creating/deleting ressources from [terraform-aws-lambda](https://github.com/spring-media/terraform-aws-lambda) module

### Homebrew (macOS)

```
brew install spring-media/tap/func
```

### from release archive - 64 bits

MacOS

```
$ curl -OL https://github.com/spring-media/func/releases/download/v0.0.10/func_0.0.10_darwin_amd64.tar.gz
$ tar -xvzf func_0.0.10_darwin_amd64.tar.gz
$ sudo mv func /usr/local/bin/func
```

GNU/Linux

```
$ wget https://github.com/spring-media/func/releases/download/v0.0.10/func_0.0.10_linux_amd64.tar.gz
$ tar -xvzf func_0.0.10_linux_amd64.tar.gz
$ sudo mv func /usr/local/bin/
```

### verify installation

```
$ func
Func is a CLI app to simplify development and deployment
	of serverless functions using Go, Terraform and AWS.

Usage:
  func [command]

Available Commands:
  help        Help about any command
  new         Creates a new Lambda project
  version     Print version information of func

Flags:
  -h, --help   help for func

Use "func [command] --help" for more information about a command.
```

## generate new project

### quickstart

(outside of `$GOPATH`)

```
$ func new github.com/you/foo
$ cd foo/
$ make init package deploy
```

### all options

```
$ func help new
Creates Terraform, CI and Go ressources for a new AWS Lambda project
in a new directory.

Usage:
  func new [module name] [flags]

Aliases:
  new, initialize, initialise, create, init

Examples:
func new github.com/you/app

Flags:
      --ci string      ci provider config file to generate [none, travis] (default "none")
  -d, --dry-run        dry run
  -e, --event string   event type triggering the Lambda function [cloudwatch-event, dynamodb, s3, sns] (default "cloudwatch-event")
  -h, --help           help for new
```

## shoulders of giants

func would not be possible if not for all of the great projects it depends on. Please see [SHOULDERS.md](SHOULDERS.md) to see a list of them.
