package main

// Group ...
func Group(challsByDate map[string][]Challenge, scenHighscore float64, scenName string) (map[string]Challenge, map[string]DateAvg, map[string][]float64) {
	ByDateMax := map[string]Challenge{}
	ByDateAvg := map[string]DateAvg{}
	ByDateAll := map[string][]float64{}

	for date, challs := range challsByDate {
		challsAmount := len(challs)
		var maxScore float64
		var maxChall Challenge
		var avgScore float64
		var sum float64
		var allScores []float64

		for i, chall := range challs {
			if i == 0 || chall.Score > maxScore {
				maxScore = chall.Score
				maxChall = chall
			}
			sum += chall.Score
			allScores = append(allScores, chall.Score)
		}

		avgScore = sum / float64(challsAmount)

		ByDateMax[date] = maxChall
		ByDateAvg[date] = DateAvg{
			Score:        float64(int(avgScore*10)) / 10,
			Grouped:      challsAmount,
			PercentagePB: int((avgScore / scenHighscore) * 100),
		}
		ByDateAll[date] = allScores
	}

	return ByDateMax, ByDateAvg, ByDateAll
}

// Two decimals precision snippet
// float64(int(number*10)) / 10
