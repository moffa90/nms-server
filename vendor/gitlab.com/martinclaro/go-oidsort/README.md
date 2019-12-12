# go-oidsort

[![GoDoc](https://godoc.org/gitlab.com/martinclaro/go-oidsort?status.png)](https://godoc.org/gitlab.com/martinclaro/go-oidsort) [![Go Report Card](https://goreportcard.com/badge/gitlab.com/martinclaro/go-oidsort)](https://goreportcard.com/report/gitlab.com/martinclaro/go-oidsort)

go-oidsort is a sorting interface for lists of SNMP OIDs.

## Example

```go
package main

import (
    "sort"
    "gitlab.com/martinclaro/go-oidsort"
)

func main() {
    // ...
    oids := []string
    // ...
    sort.Sort(oidsort.ByOidString(oids))
    // ...
}
```
