package set

import (
	"fmt"
)

type StringerSet[T fmt.Stringer] Set[string]

func (vs StringerSet[T]) Add(p T) {
	Set[string](vs).Add(p.String())
}

func (vs StringerSet[T]) Contains(p T) bool {
	return Set[string](vs).Contains(p.String())
}
