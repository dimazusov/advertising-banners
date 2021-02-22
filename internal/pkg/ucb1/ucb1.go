package ucb1

import (
	"math"
)

type banner struct {
	ID         uint
	Count      uint
	TotalCount uint
}

type ucb1 struct {
	banners []banner
}

type Ucb1 interface {
	Add(ID, count, total uint)
	GetBest() (ID uint)
}

func New() Ucb1 {
	return &ucb1{}
}

func (m *ucb1) Add(ID, count, total uint) {
	m.banners = append(m.banners, banner{
		ID:         ID,
		Count:      count,
		TotalCount: total,
	})
}

func (m *ucb1) GetBest() (ID uint) {
	if len(m.banners) == 0 {
		return 0
	}

	maxID := m.banners[0].ID
	maxValue := float64(0)

	for _, b := range m.banners {
		value := m.Calculate(b.Count, b.TotalCount)

		if value > maxValue {
			maxID = b.ID
			maxValue = value
		}
	}
	return maxID
}

func (m *ucb1) Calculate(count, totalCount uint) float64 {
	const coeff = 2 / 10

	return coeff + math.Sqrt(2*math.Log(float64(count))/float64(totalCount))
}
