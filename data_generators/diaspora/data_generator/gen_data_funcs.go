package data_generator

import (
	"math"
	"math/rand"
	"diaspora/db"
	"time"
	// "log"
	"sort"
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

			// val = 0.5, alhpa = 2 score = 1.414
			// val = 0.5, alpha = 1, score = 2
			// val = 0.5, alpha = 0.5, score = 4
			// val = 0.5, alpha = 0.25, score = 16
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

func Initialize(app string) *GenConfig {

	genConfig := new(GenConfig)
	genConfig.DBConn = db.GetDBConn(app)

	return genConfig

}

func InitializeWithUserNum(genConfig *GenConfig, userNum int) {

	genConfig.UserPopularityScores = ParetoScores(ALPHA, XM, userNum)
	genConfig.UserCommentScores = shuffleSlices(ParetoScores(ALPHA, XM, userNum))
	genConfig.UserLikeScores = shuffleSlices(ParetoScores(ALPHA, XM, userNum))
	genConfig.UserMessageScores = shuffleSlices(ParetoScores(ALPHA, XM, userNum))

}

func getSumOfFloatSlice(s []float64) float64 {

	var sum float64

	for _, num := range s {

		sum += num

	}

	return sum

}

func assignRemaingData(scores []float64, totalScore float64, 
	totalDataNum int, tempResults []float64) []float64 {
	
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

func AssignDataToUsersByUserScores(scores []float64, dataNum int) []int {
	var results []float64

	totalScore := getSumOfFloatSlice(scores)

	for i := 0; i < len(scores); i++ {

		results = append(results, math.Floor(scores[i] / totalScore * float64(dataNum)))

	}

	// assignRemaingDataTimes := 200
	// for i := 0; i < assignRemaingDataTimes; i++ {
	// 	assignRemaingData(scores, totalScore, dataNum, results)
	// }

	return transformFloat64ToInt(results)
}

func RandomNonnegativeIntWithUpperBound(upperBound int) int {

	rand.Seed(time.Now().UnixNano())

	return rand.Intn(upperBound)

}

func RandomNonnegativeFloat64WithUpperBound(upperBound float64) float64 {

	rand.Seed(time.Now().UnixNano())

	return rand.Float64() * upperBound

}

func GetSumOfIntSlice(s []int) int {
	var sum int
	for _, num := range s {
		sum += num
	}
	return sum
}

func AssignParetoDistributionScoresToData(data []int) map[int]float64 {

	scores := shuffleSlices(ParetoScores(ALPHA, XM, len(data)))

	dataScores := make(map[int]float64)

	for i, data1:= range data {

		dataScores[data1] = scores[i]
	}

	return dataScores

}

func AssignParetoDistributionScoresToDataReturnSlice(dataLen int) []float64 {

	return shuffleSlices(ParetoScores(ALPHA, XM, dataLen))

}

func GetSeqsByPersonIDs(users []User, personIDs []int) []int {

	var seq []int

	for _, personID := range personIDs {

		for i, user := range users {

			if personID == user.Person_ID {

				seq = append(seq, i)
			}
		}
	}

	return seq
}

func RandomNumWithProbGenerator(scores []float64, sum int) map[int]int { 
	
	if len(scores) == 0 {
		return nil
	}

	var cumScores []float64

	var cumSum float64
	
	for i := 0; i < len(scores); i++ {

		cumSum += scores[i]

		cumScores = append(cumScores, cumSum)

	}
	// log.Println(scores)
	// log.Println(len(scores))
	// log.Println("**********************************")
	// log.Println(cumScores)

	results := make(map[int]int)
	upperBound := cumScores[len(cumScores) - 1]

	for k := 0; k < sum; k++ {

		rNum := RandomNonnegativeFloat64WithUpperBound(upperBound)

		index := sort.SearchFloat64s(cumScores, rNum)

		if _, ok := results[index]; ok {
			results[index] += 1
		} else {
			results[index] = 1
		}

	}
	
	return results

	// for k := 0; k < sum; k++ {

	// 	rNum := RandomNonnegativeFloat64WithUpperBound(upperBound)
	// 	// log.Println(rNum)

	// 	for i := 0; i < len(cumScores); i++ {

	// 		if i == 0 {
	// 			if rNum < cumScores[i] {
	// 				if _, ok := results[i]; ok {
	// 					results[i] += 1
	// 				} else {
	// 					results[i] = 1
	// 				}
	// 			}
	// 		} else {
	// 			if rNum >= cumScores[i-1] && rNum < cumScores[i] {
	// 				if _, ok := results[i]; ok {
	// 					results[i] += 1
	// 				} else {
	// 					results[i] = 1
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// log.Println(results)

}

func TransformCommentsAssignmentToScores(assignment map[int]int, 
	scores []float64) map[float64]int {
	
	result := make(map[float64]int)
	
	for k, v := range assignment {

		k1 := scores[k]

		result[k1] = v

	}
	
	return result

}

func MakeRange(min, max int) []int {

	a := make([]int, max-min+1)
	
    for i := range a {
        a[i] = min + i
	}
	
    return a
}

func CalculateNextPostSeqStart(lastPostSeqStart, userSeqStart, 
	userSeqEnd int, postAssignment []int) int {

	postSeq := lastPostSeqStart
	
	for i := userSeqStart; i < userSeqEnd; i++ {

		postSeq += postAssignment[userSeqStart]

	}

	return postSeq

}

func MergeTwoMaps(m1, m2 map[int]float64) map[int]float64 {

	res := make(map[int]float64)

	for k1, v1 := range m1 {
		res[k1] = v1
	}

	for k2, v2 := range m2 {
		res[k2] = v2
	}

	return res

}