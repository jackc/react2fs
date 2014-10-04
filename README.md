# react2fs

react2fs is a simple utility to react to file system changes by running a
command.

## Usage

    usage:  react2fs [options] command
      -dir=".": directories to watch (separate multiple directories with commas)
      -exclude="": don't watch files matching this regexp
      -include="": only watch files matching this regexp
      -version=false: print version and exit

## Development

react2fs uses [godep](https://github.com/tools/godep) to manage its
dependencies. If you don't already have it installed "go get" it.

    go get github.com/tools/godep

Use go get to get react2fs, but don't automatically install.

    go get -d github.com/jackc/react2fs

Then prefix go commands with godep:

    godep go test
    godep go install

## Copyright

Copyright 2014 Jack Christensen

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
