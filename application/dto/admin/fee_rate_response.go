package responses

import (
	"strconv"

	"shop/domain/entities"
)

// FeeRate is the admin view of one tariff row.
type FeeRate struct {
	ID               uint
	Kind             string
	KindLabel        string
	Title            string
	Amount           uint
	FreeThreshold    uint
	HasThreshold     bool
	ThresholdText    string
	Status           bool
	StatusText       string
	StartsAt         string
	EndsAt           string
	Schedule         string
	CurrentlyApplies bool
}

type FeeRates struct {
	Data []FeeRate
}

func ToFeeRate(r *entities.FeeRate, nowIsActive bool) *FeeRate {
	statusText := "غیرفعال"
	if r.Status {
		statusText = "فعال"
	}

	startsAt, endsAt := "", ""
	if r.StartsAt != nil {
		startsAt = r.StartsAt.Format("2006-01-02")
	}
	if r.EndsAt != nil {
		endsAt = r.EndsAt.Format("2006-01-02")
	}
	schedule := "بدون محدودیت"
	switch {
	case startsAt != "" && endsAt != "":
		schedule = startsAt + " تا " + endsAt
	case startsAt != "":
		schedule = "از " + startsAt
	case endsAt != "":
		schedule = "تا " + endsAt
	}

	threshold := uint(0)
	thresholdText := "—"
	if r.FreeThreshold != nil && *r.FreeThreshold > 0 {
		threshold = *r.FreeThreshold
		thresholdText = "بالای " + strconv.FormatUint(uint64(threshold), 10) + " تومان رایگان"
	}

	return &FeeRate{
		ID:               r.ID,
		Kind:             r.Kind,
		KindLabel:        entities.FeeKindLabel(r.Kind),
		Title:            r.Title,
		Amount:           r.Amount,
		FreeThreshold:    threshold,
		HasThreshold:     threshold > 0,
		ThresholdText:    thresholdText,
		Status:           r.Status,
		StatusText:       statusText,
		StartsAt:         startsAt,
		EndsAt:           endsAt,
		Schedule:         schedule,
		CurrentlyApplies: nowIsActive,
	}
}

func ToFeeRates(rates []*entities.FeeRate, activeID uint) *FeeRates {
	out := &FeeRates{Data: make([]FeeRate, 0, len(rates))}
	for _, r := range rates {
		out.Data = append(out.Data, *ToFeeRate(r, r.ID == activeID))
	}
	return out
}
