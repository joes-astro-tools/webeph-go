package williams_test

import (
	"testing"
	testutils "webeph/testutils"
	unit "webeph/unit"
	williams "webeph/williams"
)

var (
	angle180 = unit.AngleFromDeg(180)
)

func arrayToString(index int) string {
	switch index {
	case 0:
		return "2"
	case 1:
		return "3"
	case 2:
		return "5"
	case 3:
		return "6"
	case 4:
		return "8"
	case 5:
		return "9"
	case 6:
		return "11"
	default:
		return "12"
	}
}

func TestFindRegioHouses1(t *testing.T) {
	// Example 1, p129
	// Q1, Q2, Q3 northern
	lst := unit.NewTime('+', 3, 42, 25).Angle()
	ε := unit.AngleFromDeg(23.448333)
	φ := unit.AngleFromDeg(39.745)
	h2 := unit.AngleFromDeg(176.343611)
	h3 := unit.AngleFromDeg(203.546111)
	h5 := unit.AngleFromDeg(275.388333)
	h6 := unit.AngleFromDeg(307.064444)
	h8 := angle180.Add(h2).Mod1()
	h9 := angle180.Add(h3).Mod1()
	h11 := angle180.Add(h5).Mod1()
	h12 := angle180.Add(h6).Mod1()
	expect := [8]unit.Angle{h2, h3, h5, h6, h8, h9, h11, h12}
	got := williams.FindHouses(lst, ε, φ)

	for i, val := range got {
		if !testutils.CheckTolerance(val.Deg(), expect[i].Deg(), 0.01) {
			t.Errorf("FindRegioHouses1: house %v. Expected %v to be %v", arrayToString(i), val.Deg(), expect[i].Deg())
		}
	}

}

func TestFindRegioHouses2(t *testing.T) {
	// Q3, Q4 northern
	// Example, p131
	lst := unit.NewTime('+', 12, 22, 38).Angle()
	ε := unit.AngleFromDeg(23.439444)
	φ := unit.AngleFromDeg(21.55)
	h2 := unit.AngleFromDeg(296.195)
	h3 := unit.AngleFromDeg(331.460833)
	h5 := unit.AngleFromDeg(35.295171) // Williams shows 35.81723, but manual calculation and other successful testing show another calculation.
	h6 := unit.AngleFromDeg(60.683333)
	h8 := angle180.Add(h2).Mod1()
	h9 := angle180.Add(h3).Mod1()
	h11 := angle180.Add(h5).Mod1()
	h12 := angle180.Add(h6).Mod1()
	expect := [8]unit.Angle{h2, h3, h5, h6, h8, h9, h11, h12}
	got := williams.FindHouses(lst, ε, φ)

	for i, val := range got {
		if !testutils.CheckTolerance(val.Deg(), expect[i].Deg(), 0.25) {
			t.Errorf("FindRegioHouses2: house %v. Expected %v to be %v", arrayToString(i), val.Deg(), expect[i].Deg())
		}
	}
}

func TestFindRegioHouses3(t *testing.T) {
	// Q3, Q4 southern
	// Example, p130
	lst := unit.NewTime('+', 10, 33, 14).Angle()
	ε := unit.AngleFromDeg(23.448889)
	φ := unit.AngleFromDeg(-52.4)
	h2 := unit.AngleFromDeg(300.249722)
	h3 := unit.AngleFromDeg(316.445833)
	h5 := unit.AngleFromDeg(12.511944)
	h6 := unit.AngleFromDeg(66.046111)
	h8 := angle180.Add(h2).Mod1()
	h9 := angle180.Add(h3).Mod1()
	h11 := angle180.Add(h5).Mod1()
	h12 := angle180.Add(h6).Mod1()
	expect := [8]unit.Angle{h2, h3, h5, h6, h8, h9, h11, h12}
	got := williams.FindHouses(lst, ε, φ)

	for i, val := range got {
		if !testutils.CheckTolerance(val.Deg(), expect[i].Deg(), 0.01) {
			t.Errorf("FindRegioHouses3: house %v. Expected %v to be %v", arrayToString(i), val.Deg(), expect[i].Deg())
		}
	}
}

func TestFindRegioHouses4(t *testing.T) {
	// Q1, Q2 southern
	// Example, p132
	lst := unit.NewTime('+', 22, 47, 7).Angle()
	ε := unit.AngleFromDeg(23.450278)
	φ := unit.AngleFromDeg(-40.283333)
	h2 := unit.AngleFromDeg(84.001875) // Williams shows 83.977106, but manual calculation shows a different answer.
	h3 := unit.AngleFromDeg(120.753611)
	h5 := unit.AngleFromDeg(190.845)
	h6 := unit.AngleFromDeg(214.366667)
	h8 := angle180.Add(h2).Mod1()
	h9 := angle180.Add(h3).Mod1()
	h11 := angle180.Add(h5).Mod1()
	h12 := angle180.Add(h6).Mod1()
	expect := [8]unit.Angle{h2, h3, h5, h6, h8, h9, h11, h12}
	got := williams.FindHouses(lst, ε, φ)

	for i, val := range got {
		if !testutils.CheckTolerance(val.Deg(), expect[i].Deg(), 0.01) {
			t.Errorf("FindRegioHouses4: house %v. Expected %v to be %v", arrayToString(i), val.Deg(), expect[i].Deg())
		}
	}
}

func TestFindRegioHouses5(t *testing.T) {
	// Q2, Q3 arctic
	// Example, p132
	lst := unit.NewTime('+', 16, 13, 6).Angle()
	ε := unit.AngleFromDeg(23.44375)
	φ := unit.AngleFromDeg(72.85)
	h2 := unit.AngleFromDeg(54.080673)
	h3 := unit.AngleFromDeg(59.509274) // Williams shows 59.719577, but manual calculation differs
	h5 := unit.AngleFromDeg(76.99142)
	h6 := unit.AngleFromDeg(163.5144)
	h8 := angle180.Add(h2).Mod1()
	h9 := angle180.Add(h3).Mod1()
	h11 := angle180.Add(h5).Mod1()
	h12 := angle180.Add(h6).Mod1()
	expect := [8]unit.Angle{h2, h3, h5, h6, h8, h9, h11, h12}
	got := williams.FindHouses(lst, ε, φ)

	for i, val := range got {
		if !testutils.CheckTolerance(val.Deg(), expect[i].Deg(), 0.01) {
			t.Errorf("FindRegioHouses4: house %v. Expected %v to be %v", arrayToString(i), val.Deg(), expect[i].Deg())
		}
	}
}
