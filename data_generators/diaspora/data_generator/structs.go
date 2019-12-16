package data_generator

import (
	"database/sql"
)

// const ALPHA = 2.0 (old)
// const XM = 0.2 (old)
const ALPHA = 0.05
const XM = 3.0

type User struct {
	User_ID       		int
	Person_ID     		int
	Aspects       		[]int
}

type GenConfig struct {
	DBConn					*sql.DB
	UserPopularityScores	[]float64
	UserCommentScores		[]float64
	UserLikeScores			[]float64
	UserMessageScores		[]float64
}

type Post struct {
	ID 						int
	Author					int
	Score					float64
}