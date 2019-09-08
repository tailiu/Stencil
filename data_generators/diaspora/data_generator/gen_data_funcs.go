package data_generator

import (
	"math"
	"math/rand"
	"diaspora/db"
	"time"
	// "log"
)

func ParetoScores(alpha, Xm float64, num int) []float64 {
	var uniform []float64
	var scores []float64

	var curr float64
	interval := float64(1) / float64(num)

	for i := 0; i < num; i++ {
		uniform = append(uniform, curr)
		curr += interval
	}  
	
	for _, val := range uniform {
		if val >= 0 {
			score := Xm / math.Pow((1.0 - val), (1.0 / alpha))
			scores = append(scores, score)
		}
	}

	return scores
}

func shuffleSlices(s []float64) []float64 {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return s
}

func Initialize(app string, userNum int) *GenConfig {
	genConfig := new(GenConfig)
	genConfig.DBConn = db.GetDBConn(app)
	genConfig.PopularityScores = ParetoScores(ALPHA, XM, userNum)
	genConfig.CommentScores = shuffleSlices(ParetoScores(ALPHA, XM, userNum))
	genConfig.LikeScores = shuffleSlices(ParetoScores(ALPHA, XM, userNum))

	return genConfig
}