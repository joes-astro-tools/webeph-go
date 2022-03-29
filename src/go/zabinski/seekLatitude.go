package zabinski

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	julian "webeph/julian"
	"webeph/moonposition"
)

type LatitudeInfo struct {
	β  float64
	tm time.Time
}

// Describes the day as a decimal day, the fractional portion being the time out of 24 hours that has passed.
// Receives:
//	tm: time.Time
// Returns:
//	A float64 representing the decimal day.
func fractionalDay(tm time.Time) float64 {
	midnight := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, time.UTC)
	sinceMidnight := tm.Sub(midnight)
	return float64(tm.Day()) + (sinceMidnight.Hours() / 24)
}

// Finds lunar latitude at a time.
// Receives:
//	tm: time.Time
//	c: a channel collecting latitude information
//	wg: the WaitGroup awaiting results
// Returns:
//	nothing
// Notes:
//	Finds the lunar latitude, embeds that information in a struct, and send the struct to the channel.
func getLunarLatitude(tm time.Time, c chan LatitudeInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	jd := julian.CalendarGregorianToJD(tm.Year(), int(tm.Month()), fractionalDay(tm))
	_, geocentricβ, _ := moonposition.Position(jd)
	info := &LatitudeInfo{
		β:  geocentricβ.Deg(),
		tm: tm,
	}
	c <- *info
}

// Transforms a float64 to a string representation in degrees, minutes and seconds.
// Receives:
//	val: a float64
// Returns:
//	A string representing the float as a DMS. For example, 4.5 becomes '4° 30' 0"'
func floatToDMS(val float64) string {
	degs := math.Floor(val)
	mins := (val - math.Floor(val)) * 60
	secs := (mins - math.Floor(mins)) * 60
	secString := fmt.Sprintf("%.6f", secs)
	sep := " "
	return fmt.Sprint(int(degs)) + "°" + sep + fmt.Sprint(int(mins)) + "'" + sep + secString + "\""
}

// Compares lunar latitudes every hour until a given number of iterations are met.
// Receives:
//	iterations: the number of times to check latitudes. For example, send 3 to check latitudes every hour for 3 hours.
// Returns
func SeekLatitudes(iterations int) (maxD, minD, aveD, maxT string) {
	var wg sync.WaitGroup
	times := []time.Time{}
	b := time.Date(2022, time.July, 15, 0, 0, 0, 0, time.UTC)
	e := b.Add(time.Hour * time.Duration(iterations))
	now := b
	for now.Before(e) {
		times = append(times, now)
		now = now.Add(time.Hour)
	}
	c := make(chan LatitudeInfo, len(times))
	for _, val := range times {
		wg.Add(1)
		go getLunarLatitude(val, c, &wg)
	}
	// Wait for all latitudes.
	wg.Wait()
	// Shut down the channel collecting results.
	close(c)
	// Sort the information.
	sorted := []LatitudeInfo{}
	for info := range c {
		sorted = append(sorted, info)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].tm.Before(sorted[j].tm)
	})
	// Measure latitude differences, and describe max difference, min difference and average.
	max := 0.
	min := 6.
	ave := 0.
	sumDiffs := 0.
	var prev LatitudeInfo
	tot := 0
	for i, info := range sorted {
		// Only start checking diff after prev has been set: that is, after the first time.
		if i > 0 {
			diff := math.Abs(info.β - prev.β)
			if diff > max {
				max = diff
			}
			if diff < min {
				min = diff
			}
			sumDiffs += diff
		}
		prev = info
		tot++
	}
	ave = sumDiffs / float64(tot)
	maxD = fmt.Sprintf("Max: %v. Decimal: %.6v", floatToDMS(max), max)
	minD = fmt.Sprintf("Min: %v. Decimal: %.6v", floatToDMS(min), min)
	aveD = fmt.Sprintf("Average: %v. Decimal: %.6v", floatToDMS(ave), ave)
	return
}
