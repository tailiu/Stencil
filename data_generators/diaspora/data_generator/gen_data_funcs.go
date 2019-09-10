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

func AssignDataToUsersByUserScores(scores []float64, dataNum int) []int {
	var results []float64

	totalScore := getSumOfFloatSlice(scores)
	for i := 0; i < len(scores); i++ {
		results = append(results, math.Floor(scores[i] / totalScore * float64(dataNum)))
	}

	assignRemaingDataTimes := 200
	for i := 0; i < assignRemaingDataTimes; i++ {
		assignRemaingData(scores, totalScore, dataNum, results)
	}

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

func AssignScoresToPosts(posts []int) map[int]float64 {
	postPopularityScores := shuffleSlices(ParetoScores(ALPHA, XM, len(posts)))
	postScores := make(map[int]float64)
	for i, postID:= range posts {
		postScores[postID] = postPopularityScores[i]
	}
	return postScores
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
	var cumScores []float64
	
	for i := 0; i < len(scores); i++ {
		var cumSum float64
		for j := 0; j <= i; j++ {
			cumSum += scores[j]
		}
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
		// log.Println(rNum)
		for i := 0; i < len(cumScores); i++ {
			if i == 0 {
				if rNum < cumScores[i] {
					if _, ok := results[i]; ok {
						results[i] += 1
					} else {
						results[i] = 1
					}
				}
			} else {
				if rNum >= cumScores[i-1] && rNum < cumScores[i] {
					if _, ok := results[i]; ok {
						results[i] += 1
					} else {
						results[i] = 1
					}
				}
			}
		}
	}

	// log.Println(results)
	return results
}

func TransformCommentsAssignmentToScores(assignment map[int]int, scores []float64) map[float64]int {
	result := make(map[float64]int)
	for k, v := range assignment {
		k1 := scores[k]
		result[k1] = v
	}
	return result
}