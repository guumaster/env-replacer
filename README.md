[![Tests](https://img.shields.io/github/workflow/status/guumaster/env-replacer/Test)](https://github.com/guumaster/env-replacer/actions?query=workflow%3ATest)
[![GitHub Release](https://img.shields.io/github/release/guumaster/env-replacer.svg?logo=github&labelColor=262b30)](https://github.com/guumaster/env-replacer/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/guumaster/env-replacer)](https://goreportcard.com/report/github.com/guumaster/env-replacer)
[![License](https://img.shields.io/github/license/guumaster/env-replacer)](https://github.com/guumaster/env-replacer/LICENSE)
# Env Replacer

A powerful tool to replace environment variables on template files.

## Why?

I don't like to maintain `bash` scripts, I don't like `gettext` and `envsubst`, so this is a replacement.
 

## Installation

Go to [release page](https://github.com/guumaster/env-replacer/releases) and download the binary you need.


## Examples

```bash
$> env-replacer --path /some-path-to-tpl-files
```


## Usage

    NAME:
       env-replacer - replace env variables on files
    
    USAGE:
       env-replacer [global options] [arguments...]
    
    VERSION:
       dev
    
    AUTHOR:
       guumaster <guuweb@gmail.com>
    
    GLOBAL OPTIONS:
       --path value    path to scan for templates
    
       --files value   file pattern to search (default: "*.tpl")
    
       --quiet, -q     hide report (default: false)
    
       --no-recursive  don't scan path recursively (default: false)
    
       --no-truncate   don't truncate dst files (default: false)
    
       --no-error      early quit on error (default: false)
    
       --no-missing    don't write files with missing variables (default: false)
    
       --no-empty      fail also with variables (default: false)
    
       --strict        same as --no-empty and --no-missing combined (default: false)
    
       --version       print the version (default: false)


### Dependencies & Refs
  - [bmatcuk/doublestar](https://github.com/bmatcuk/doublestar)
  - [urfave/cli](https://github.com/urfave/cli)


### LICENSE

 [MIT license](LICENSE)


### Author(s)

* [guumaster](https://github.com/guumaster)

