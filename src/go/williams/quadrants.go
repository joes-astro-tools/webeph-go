package williams

import (
	unit "webeph/unit"
)

var (
	angle30  = unit.AngleFromDeg(30)
	angle60  = unit.AngleFromDeg(60)
	angle90  = unit.AngleFromDeg(90)
	angle120 = unit.AngleFromDeg(120)
	angle150 = unit.AngleFromDeg(150)
	angle180 = unit.AngleFromDeg(180)
	angle270 = unit.AngleFromDeg(270)
	angle360 = unit.AngleFromDeg(360)
	q1C      = angle90.Rad()
	q2C      = angle180.Rad()
	q3C      = angle270.Rad()
	q4C      = 0.
	// Helper functions to find which quadrant houses the medium coeli
	q1 = func(angle unit.Angle) bool { return (angle.Rad() >= q4C) && (angle.Rad() < q1C) }
	q2 = func(angle unit.Angle) bool { return (angle.Rad() >= q1C) && (angle.Rad() < q2C) }
	q3 = func(angle unit.Angle) bool { return (angle.Rad() >= q2C) && (angle.Rad() < q3C) }
)
