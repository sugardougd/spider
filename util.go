package spider

import (
	"fmt"
)

const TAB = "\t"
const TAB2 = TAB + TAB
const BLANK = " "
const BLANK2 = BLANK + BLANK
const K = 1024
const M = K * K
const G = K * K * K

func BytesTo(memory uint64) string {
	if memory > G {
		return fmt.Sprintf("%.1fG", BytesToGb(memory))
	} else if memory > M {
		return fmt.Sprintf("%.1fM", BytesToMb(memory))
	} else if memory > K {
		return fmt.Sprintf("%.1fK", BytesToKb(memory))
	}
	return fmt.Sprintf("%dB", memory)
}

func BytesToKb(b uint64) float64 {
	return Percentage(b, K)
}

func BytesToMb(b uint64) float64 {
	return Percentage(b, M)
}

func BytesToGb(b uint64) float64 {
	return Percentage(b, G)
}

func Percentage[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](value1, value2 T) float64 {
	if value2 == 0 {
		return 0
	}
	return float64(value1) / float64(value2)
}
