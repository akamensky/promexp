package prommetric

import (
	"fmt"
	"sort"
	"strings"
)

type labelsT []*label

func newLabels(ls map[string]string) labelsT {
	result := make(labelsT, 0)
	for k, v := range ls {
		result = append(result, newLabel(k, v))
	}
	return result
}

func (ls labelsT) String() string {
	if ls == nil || len(ls) == 0 {
		return ""
	}

	lines := make([]string, 0)
	for _, l := range ls {
		lines = append(lines, l.String())
	}

	sort.Strings(lines)

	return fmt.Sprintf("{%s}", strings.Join(lines, ","))
}

type label struct {
	key   string
	value string
}

func newLabel(k, v string) *label {
	return &label{
		key:   k,
		value: v,
	}
}

func (l *label) String() string {
	return fmt.Sprintf(`%s="%s"`, l.key, l.value)
}
