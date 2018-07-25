package graph

import (
	"sort"
	"testing"
)

func testDegreeDistributionEquality(t *testing.T, expected, actual *DegreeDistribution) {
	testDegreeMapEquality(t, expected.m, actual.m)
	testDegreesEquality(t, expected.degrees, actual.degrees)
}

func testDegreeMapEquality(t *testing.T, expected, actual map[int]int) {
	if len(actual) != len(expected) {
		t.Errorf("expected: %d, actual: %d", len(expected), len(actual))
	}
	for k, numExpected := range expected {
		numActual := actual[k]
		if numActual != numExpected {
			t.Errorf("expected: %d, actual: %d", numExpected, numActual)
		}
	}
}

func testDegreesEquality(t *testing.T, expected, actual []int) {
	if len(actual) != len(expected) {
		t.Errorf("expected: %d, actual: %d", len(expected), len(actual))
	}
	for i, kExpected := range expected {
		kActual := actual[i]
		if kActual != kExpected {
			t.Errorf("expected: %d, actual: %d", kExpected, kActual)
		}
	}
}

func TestNewDegreeDistribution(t *testing.T) {
	dist := NewDegreeDistribution()
	if len(dist.m) > 0 {
		t.Errorf("expected: %d, actual: %d", 0, len(dist.m))
	}
	if len(dist.degrees) > 0 {
		t.Errorf("expected: %d, actual: %d", 0, len(dist.degrees))
	}
}

func TestDegreeDistribution_Sort(t *testing.T) {
	dist := &DegreeDistribution{
		degrees: []int{2, 0, 1},
	}

	sort.Sort(dist)
	testDegreesEquality(t, []int{0, 1, 2}, dist.degrees)
}

func TestDegreeDistribution_GetNum(t *testing.T) {
	type input struct {
		k int
	}
	type output struct {
		num int
	}
	testCases := []struct {
		name string
		dist *DegreeDistribution
		in   input
		out  output
	}{
		{
			"success",
			&DegreeDistribution{
				m: map[int]int{0: 1, 1: 2, 2: 3},
			},
			input{1},
			output{2},
		},
		{
			"success: no key",
			&DegreeDistribution{
				m: map[int]int{0: 1, 1: 2, 2: 3},
			},
			input{3},
			output{0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dist, in, out := tc.dist, tc.in, tc.out

			num := dist.GetNum(in.k)
			if num != out.num {
				t.Errorf("expected: %d, actual: %d", out.num, num)
			}
		})
	}
}

func TestDegreeDistribution_GetDegrees(t *testing.T) {
	expected := []int{0, 1, 2}
	dist := &DegreeDistribution{
		degrees: expected,
	}

	actual := dist.GetDegrees()
	testDegreesEquality(t, expected, actual)
}

func TestDegreeDistribution_Add(t *testing.T) {
	type input struct {
		k int
	}
	testCases := []struct {
		name     string
		actual   *DegreeDistribution
		expected *DegreeDistribution
		in       input
	}{
		{
			"success",
			&DegreeDistribution{
				m:       map[int]int{},
				degrees: []int{},
			},
			&DegreeDistribution{
				m:       map[int]int{0: 1},
				degrees: []int{0},
			},
			input{0},
		},
		{
			"success: duplication",
			&DegreeDistribution{
				m:       map[int]int{0: 1, 1: 2},
				degrees: []int{0, 1},
			},
			&DegreeDistribution{
				m:       map[int]int{0: 2, 1: 2},
				degrees: []int{0, 1},
			},
			input{0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, expected, in := tc.actual, tc.expected, tc.in

			actual.Add(in.k)
			testDegreeDistributionEquality(t, expected, actual)
		})
	}
}

func TestDegreeDistribution_CalcAverageDegree(t *testing.T) {
	expected := 2.0
	dist := &DegreeDistribution{
		m: map[int]int{0: 1, 1: 2, 2: 3, 3: 4},
	}

	actual := dist.CalcAverageDegree()
	if actual != expected {
		t.Errorf("expected: %f, actual: %f", expected, actual)
	}
}
