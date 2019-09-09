package data_generator

import (
	"database/sql"
)

const ALPHA = 2.0
const XM = 0.2

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
	DBConn					*sql.DB
	UserPopularityScores	[]float64
	UserCommentScores		[]float64
	UserLikeScores			[]float64
	PostPopularityScores	[]float64
}