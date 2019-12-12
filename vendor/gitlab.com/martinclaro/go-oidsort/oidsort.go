package oidsort

import (
	"fmt"
	"strconv"
	"strings"
)

// ByOidString is a type that implements the sort.Interface interface
// so that OIDs can be sorted.
type ByOidString []string

func (o ByOidString) Len() int {
	return len(o)
}

func (o ByOidString) Less(i, j int) bool {
	return (CompareOIDs(o[i], o[j]) < 0)
}

func (o ByOidString) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

// CompareOIDs is a function to compare two numerically-formatted OID
// and to obtein whenever the first argument is lesser (-1), equal (0)
// or greater (1) than the segond argument.
func CompareOIDs(i string, j string) int {
	// Fast equality check
	if i == j {
		return 0
	}

	partsI := strings.Split(i, ".")
	partsJ := strings.Split(j, ".")
	max := 0

	if len(partsI) <= len(partsJ) {
		max = len(partsI)
	} else {
		max = len(partsJ)
	}

	for x := 0; x < max; x++ {
		nodeI, err := strconv.ParseUint(partsI[x], 10, 32)
		nodeJ, err := strconv.ParseUint(partsJ[x], 10, 32)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}
		if uint64(nodeI) < uint64(nodeJ) {
			return -1
		} else if uint64(nodeI) > uint64(nodeJ) {
			return 1
		}
	}

	if len(partsI) < len(partsJ) {
		return -1
	}
	return 1
}
