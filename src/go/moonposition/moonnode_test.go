package moonposition_test

import (
	"testing"

	moonposition "webeph/moonposition"
	testutils "webeph/testutils"
)

// Expected values from AstroApp. Seeing that we don't know how Swiss Eph calculates the node, and we have no examples
// using our algorithm, one test with a tiny error over tolerance out of five is quite good.

func TestTrueNode1(t *testing.T) {
	jd := 2448000.5 //
	// Expected value, matching AstroApp
	expected := 313.216667
	got := moonposition.FindAscendingNode(jd)
	// 0˚ 1' 0.0084".
	if !testutils.CheckTolerance(got, expected, testutils.StandardTolerance) {
		t.Errorf("TestTrueNode1: expected %v to be %v", testutils.FloatToDMS(got), testutils.FloatToDMS(expected))
	}
}

func TestTrueNode2(t *testing.T) {
	jd := 2460053.5 // 4/19/2023@00:00 UTC
	expected := 34.016667
	got := moonposition.FindAscendingNode(jd)
	if !testutils.CheckTolerance(got, expected, testutils.StandardTolerance) {
		t.Errorf("TestTrueNode2: expected %v to be %v", testutils.FloatToDMS(got), testutils.FloatToDMS(expected))
	}
}

func TestTrueNode3(t *testing.T) {
	jd := 2463950.5 // 12/19/2033@00:00 UTC
	expected := 189.066667
	got := moonposition.FindAscendingNode(jd)
	if !testutils.CheckTolerance(got, expected, testutils.StandardTolerance) {
		t.Errorf("TestTrueNode3: expected %v to be %v", testutils.FloatToDMS(got), testutils.FloatToDMS(expected))
	}
}

func TestTrueNode4(t *testing.T) {
	jd := 2471072.5 // 6/19/2053@00:00 UTC
	expected := 170.4
	got := moonposition.FindAscendingNode(jd)
	if !testutils.CheckTolerance(got, expected, testutils.StandardTolerance) {
		t.Errorf("TestTrueNode4: expected %v to be %v", testutils.FloatToDMS(got), testutils.FloatToDMS(expected))
	}
}

func TestTrueNode5(t *testing.T) {
	jd := 2478038.5 // 7/15/2072@00:00 UTC
	expected := 160.633334
	got := moonposition.FindAscendingNode(jd)
	if !testutils.CheckTolerance(got, expected, testutils.StandardTolerance) {
		t.Errorf("TestTrueNode5: expected %v to be %v", testutils.FloatToDMS(got), testutils.FloatToDMS(expected))
	}
}
