package global

import "github.com/samber/lo"

type PeriodType string

const (
	MONTHLY PeriodType = "monthly"
	YEARLY  PeriodType = "yearly"
)

var periodTypeDescriptions = map[PeriodType]string{
	MONTHLY: "Monthly",
	YEARLY:  "Yearly",
}

func (s PeriodType) Description() string {
	return periodTypeDescriptions[s]
}

func PeriodTypeValues() []PeriodType {
	return lo.Keys(periodTypeDescriptions)
}

func PeriodTypeDescriptions() []string {
	return lo.Values(periodTypeDescriptions)
}
