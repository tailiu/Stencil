package data_generator

import (
	"database/sql"
)

type User struct {
	User_ID       		int
	Person_ID     		int
	Aspects       		[]int
	ContactID     		int
	ContactAspect 		int
	PopularityScore 	float64
	CommentScore 		float64
	LikeScore 			float64
}

type GenConfig struct {
	DBConn				*sql.DB
	PopularityScores	[]float64
	CommentScores		[]float64
	LikeScores			[]float64
}

const ALPHA = 2.0
const XM = 0.2