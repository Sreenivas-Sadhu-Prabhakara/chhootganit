package backend

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEvaluate_ExtraUnits(t *testing.T) {
	// Price 100, profit 20/unit, 5% discount => ₹5 off => new profit 15.
	// extra units to hold profit = 20/15 - 1 = 33.33%.
	r := Evaluate(Input{Price: 100, ProfitPerUnit: 20, DiscountPct: 5})
	if !r.Feasible {
		t.Fatal("should be feasible")
	}
	if math.Abs(r.NewProfitPerUnit-15) > 1e-9 {
		t.Fatalf("newProfit=%v want 15", r.NewProfitPerUnit)
	}
	if math.Abs(r.ExtraUnitsPct-100.0/3.0) > 1e-6 {
		t.Fatalf("extraPct=%v want 33.33", r.ExtraUnitsPct)
	}
}

func TestEvaluate_InfeasibleWhenDiscountEatsMargin(t *testing.T) {
	// 25% off ₹100 = ₹25 > ₹20 profit => impossible.
	r := Evaluate(Input{Price: 100, ProfitPerUnit: 20, DiscountPct: 25})
	if r.Feasible {
		t.Fatal("should be infeasible")
	}
}

func TestValidate(t *testing.T) {
	if err := (Input{Price: 100, ProfitPerUnit: 10, DiscountPct: 5}).Validate(); err != nil {
		t.Fatalf("valid rejected: %v", err)
	}
	for i, bad := range []Input{{Price: 0}, {Price: 100, DiscountPct: 100}, {Price: 100, ProfitPerUnit: -1}} {
		if err := bad.Validate(); err == nil {
			t.Fatalf("bad %d accepted", i)
		}
	}
}

func TestEvaluateEndpoint(t *testing.T) {
	srv := NewServer(nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/evaluate",
		strings.NewReader(`{"price":100,"profitPerUnit":20,"discountPct":5}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	var r Result
	json.Unmarshal(rec.Body.Bytes(), &r)
	if !r.Feasible {
		t.Fatal("expected feasible")
	}
}
