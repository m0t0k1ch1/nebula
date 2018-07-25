package graph

type DegreeDistribution struct {
	m       map[int]int
	degrees []int
}

func NewDegreeDistribution() *DegreeDistribution {
	return &DegreeDistribution{
		m:       map[int]int{},
		degrees: []int{},
	}
}

func (dist *DegreeDistribution) Len() int {
	return len(dist.degrees)
}

func (dist *DegreeDistribution) Less(i, j int) bool {
	return dist.degrees[i] < dist.degrees[j]
}

func (dist *DegreeDistribution) Swap(i, j int) {
	dist.degrees[i], dist.degrees[j] = dist.degrees[j], dist.degrees[i]
}

func (dist *DegreeDistribution) GetNum(k int) int {
	if _, ok := dist.m[k]; !ok {
		return 0
	}
	return dist.m[k]
}

func (dist *DegreeDistribution) GetDegrees() []int {
	return dist.degrees
}

func (dist *DegreeDistribution) Add(k int) {
	if _, ok := dist.m[k]; !ok {
		dist.degrees = append(dist.degrees, k)
	}
	dist.m[k]++
}

func (dist *DegreeDistribution) CalcAverageDegree() float64 {
	kTotal, numTotal := 0, 0
	for k, num := range dist.m {
		kTotal += k * num
		numTotal += num
	}
	return float64(kTotal) / float64(numTotal)
}
