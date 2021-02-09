package main

// Group ...
func Group(challsByDate map[string][]Challenge) (map[string]Challenge, map[string][]interface{}) {
	ByDateMax := map[string]Challenge{}
	ByDateAvg := map[string][]interface{}{}

	for date, challs := range challsByDate {
		challsAmount := len(challs)
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

		avgScore = sum / float64(challsAmount)

		ByDateMax[date] = maxChall
		ByDateAvg[date] = []interface{}{float64(int(avgScore*10)) / 10, challsAmount}
	}

	return ByDateMax, ByDateAvg
}

// Two decimals precision snippet
// float64(int(number*10)) / 10
