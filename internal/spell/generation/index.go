package generation

import (
	"iter"

	"github.com/alan-b-lima/sun/internal/spell/semantics"
)

type index[I integer] map[semantics.Symbol]I

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func make_index[I integer](symbols iter.Seq[semantics.Symbol], len int) index[I] {
	index := make(index[I], len)

	var i I
	for symbol := range symbols {
		index[symbol] = i
		i++
	}

	return index
}
