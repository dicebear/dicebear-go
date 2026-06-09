package style

// ComponentVariant is a read-only view over an entry in a component's variants
// block.
type ComponentVariant struct {
	Elements []Element `json:"elements"`
	Weight   *float64  `json:"weight"`
}

// WeightOr1 returns the weighted-pick weight, defaulting to 1.
func (cv ComponentVariant) WeightOr1() float64 {
	if cv.Weight != nil {
		return *cv.Weight
	}
	return 1
}
