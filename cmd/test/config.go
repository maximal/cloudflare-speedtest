package test

type TestStep struct {
	Upload bool
	Bytes  uint
	Count  uint
}

var testSteps = []TestStep{
	// Empty download × 20
	{Upload: false, Bytes: 100_000, Count: 10},

	// 100 kB upload × 10
	{Upload: true, Bytes: 100_000, Count: 10},

	// 1 MB download × 8
	{Upload: false, Bytes: 1_000_000, Count: 8},

	// 1 MB upload × 8
	{Upload: true, Bytes: 1_000_000, Count: 8},

	// 10 MB download × 6
	{Upload: false, Bytes: 10_000_000, Count: 6},

	// 10 MB upload × 5
	{Upload: true, Bytes: 10_000_000, Count: 5},

	// 25 MB download × 5
	{Upload: false, Bytes: 25_000_000, Count: 5},

	// 25 MB upload × 4
	{Upload: true, Bytes: 25_000_000, Count: 4},

	// 100 MB download × 4
	{Upload: false, Bytes: 100_000_000, Count: 4},
}
