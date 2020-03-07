/*

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

*/
package main
