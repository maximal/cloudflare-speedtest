package format

import "fmt"

func valueWithUnit(number float64, unit string) string {
	if number < 100 {
		// %.3g
		// value < 10: two decimals, trailing zeroes omitted
		// value < 100: one decimal, trailing zeroes omitted
		return fmt.Sprintf("%.3g %s", number, unit)
	}
	// value >= 100: no decimals
	return fmt.Sprintf("%.0f %s", number, unit)
}

func siUnitPrefix(value uint64) string {
	if value < 1_000 {
		return fmt.Sprintf("%d ", value)
	}
	if value < 1_000_000 {
		return valueWithUnit(float64(value)/1_000, "k")
	}
	if value < 1_000_000_000 {
		return valueWithUnit(float64(value)/1_000/1_000, "M")
	}
	if value < 1_000_000_000_000 {
		return valueWithUnit(float64(value)/1_000/1_000/1_000, "G")
	}
	return valueWithUnit(float64(value)/1_000/1_000/1_000/1_000, "T")
}

func BitsPerSecondSi(bits uint64) string {
	return siUnitPrefix(bits) + "bit/s"
}

func BytesSi(bytes uint64) string {
	return siUnitPrefix(bytes) + "B"
}

func BytesIec(bytes uint64) string {
	if bytes < 1_024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1_024*1_024 {
		return valueWithUnit(float64(bytes)/1_024, "KiB")
	}
	if bytes < 1_024*1_024*1_024 {
		return valueWithUnit(float64(bytes)/1_024/1_024, "MiB")
	}
	if bytes < 1_024*1_024*1_024*1_024 {
		return valueWithUnit(float64(bytes)/1_024/1_024/1_024, "GiB")
	}
	return valueWithUnit(float64(bytes)/1_024/1_024/1_024/1_024, "TiB")
}
