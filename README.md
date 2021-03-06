# goawaybreaks

Small utility to remove blank (whitespace only) lines from import groups in go code

## Usage
```
goawaybreaks -w file.go
```
Removes breaks and writes the result to the source file instead of stdout

There is also a streaming stdin/stdout mode of operation just like `goimports`

## Motivation

>Make imports more consistent

goimports respects "groups" of imports, defined by there being a gap of at least one line between two import statements. Those lines could be empty or have comments. This means that two standard imports separated by an empty line will form their own groups, if you mix remote (or appengine or local) imports with these then each group will then in turn be split into more groups, for example:
```
import (
 "remote.com/other"
 "local"

 "remote.com/package"
 "local2"
)
```
will become
```
import (
        "local"

        "remote.com/other"

        "local2"

        "remote.com/pack"
)
```
Personally I don't intentionally split my imports into my own groups and I don't think you should as it means that there is less consistency when `goimports` is run.
If you run `goawaybreaks` first then `goimports` the imports would become
```
import (
        "local"
        "local2"

        "remote.com/other"
        "remote.com/pack"
)
```
