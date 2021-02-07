package main

// Group ...
func Group(challsByDate map[string][]Challenge) (map[string]Challenge, map[string]float64) {
	ByDateMax := map[string]Challenge{}
	ByDateAvg := map[string]float64{}

	for date, challs := range challsByDate {
		var maxScore float64
		var maxChall Challenge
		var avgScore float64
		var sum float64

		for i, chall := range challs {
			if i == 0 || chall.Score > maxScore {
				maxScore = chall.Score
				maxChall = chall
			}
			sum += chall.Score
		}

		avgScore = sum / float64(len(challs))

		// ByDateMax[k] = float64(int(max*10)) / 10
		// ByDateAvg[k] = float64(int(avg*10)) / 10
		ByDateMax[date] = maxChall
		ByDateAvg[date] = float64(int(avgScore*10)) / 10
	}

	return ByDateMax, ByDateAvg
}
