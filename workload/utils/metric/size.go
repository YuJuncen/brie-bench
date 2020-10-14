package metric

import "fmt"

var iecUnits = []string{"", "K", "M", "G", "T", "P"}

// ToIce convert a size to iec format("xxxM")
func ToIec(num int64) string {
	unitIdx := 0
	floatNum := float64(num)
	for floatNum > 1024 && unitIdx < len(iecUnits)-1 {
		floatNum = floatNum / 1024
		unitIdx += 1
	}
	return fmt.Sprintf("%.1f%s", floatNum, iecUnits[unitIdx])
}

type SizeStringer uint64

func (s SizeStringer) String() string {
	return ToIec(int64(s))
}
