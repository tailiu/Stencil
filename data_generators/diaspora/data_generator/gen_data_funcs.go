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
	genConfig.UserPopularityScores = ParetoScores(ALPHA, XM, userNum)
	genConfig.UserCommentScores = shuffleSlices(ParetoScores(ALPHA, XM, userNum))
	genConfig.UserLikeScores = shuffleSlices(ParetoScores(ALPHA, XM, userNum))

	return genConfig
}

func getSumOfFloatSlice(s []float64) float64 {
	var sum float64
	for _, num := range s {
		sum += num
	}
	return sum
}

func assignRemaingData(scores []float64, totalScore float64, totalDataNum int, tempResults []float64) []float64 {
	remainingDataNum := float64(totalDataNum) - getSumOfFloatSlice(tempResults)
	for i := 0; i < len(tempResults); i++ {
		tempResults[i] += math.Floor(scores[i] / totalScore * remainingDataNum)
	}
	return tempResults
}

func transformFloat64ToInt(data []float64) []int {
	var data1 []int
	for _, val := range data {
		data1 = append(data1, int(val))
	}
	return data1
}

func AssignDataToUsersByUserPopScores(genConfig *GenConfig, userNum, dataNum int) []int {
	var results []float64

	totalScore := getSumOfFloatSlice(genConfig.UserPopularityScores)
	for i := 0; i < userNum; i++ {
		results = append(results, math.Floor(genConfig.UserPopularityScores[i] / totalScore * float64(dataNum)))
	}

	assignRemaingDataTimes := 200
	for i := 0; i < assignRemaingDataTimes; i++ {
		assignRemaingData(genConfig.UserPopularityScores, totalScore, dataNum, results)
	}

	return transformFloat64ToInt(results)
}

func RandomNonnegativeIntWithUpperBound(upperBound int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(upperBound)
}

func GetSumOfIntSlice(s []int) int {
	var sum int
	for _, num := range s {
		sum += num
	}
	return sum
}

