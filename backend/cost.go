package backend

import "fmt"

// Input frames a haggle: you sell at Price with ProfitPerUnit, and the customer
// wants DiscountPct off. How many more units must you sell to keep total profit?
type Input struct {
	Price         float64 `json:"price"`
	ProfitPerUnit float64 `json:"profitPerUnit"`
	DiscountPct   float64 `json:"discountPct"`
}

// Result answers the haggle with real numbers.
type Result struct {
	DiscountPerUnit   float64 `json:"discountPerUnit"`
	NewProfitPerUnit  float64 `json:"newProfitPerUnit"`
	ExtraUnitsPct     float64 `json:"extraUnitsPct"` // % more units needed to hold total profit
	Feasible          bool    `json:"feasible"`      // false if discount wipes out the margin
}

// Headline is the extra-units percentage required.
func (r Result) Headline() float64 { return r.ExtraUnitsPct }

// Label flags whether holding profit is even possible.
func (r Result) Label() string {
	if r.Feasible {
		return "feasible"
	}
	return "infeasible"
}

// Validate reports whether the Input is well formed.
func (in Input) Validate() error {
	if in.Price <= 0 {
		return fmt.Errorf("price must be positive")
	}
	if in.ProfitPerUnit < 0 {
		return fmt.Errorf("profit per unit cannot be negative")
	}
	if in.DiscountPct < 0 || in.DiscountPct >= 100 {
		return fmt.Errorf("discount %% must be between 0 and 100")
	}
	return nil
}

// Evaluate computes the extra volume needed to keep total profit unchanged after
// granting the discount. If the discount meets or exceeds the margin, no volume
// can recover it (infeasible).
func Evaluate(in Input) Result {
	discountPerUnit := in.Price * in.DiscountPct / 100
	newProfit := in.ProfitPerUnit - discountPerUnit
	if newProfit <= 0 {
		return Result{DiscountPerUnit: discountPerUnit, NewProfitPerUnit: newProfit, Feasible: false}
	}
	extraPct := (in.ProfitPerUnit/newProfit - 1) * 100
	return Result{
		DiscountPerUnit:  discountPerUnit,
		NewProfitPerUnit: newProfit,
		ExtraUnitsPct:    extraPct,
		Feasible:         true,
	}
}
