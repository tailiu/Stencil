package data_generator

import (
	"database/sql"
)

// If ALPHA >= 1.0, the bigger ALPHA is, the flat the distribution is
// If 0 < ALPHA < 1.0, the bigger ALPHA is, the flat the distribution is 

// const ALPHA = 2.0 (old)
// const XM = 0.2 (old)
const ALPHA = 2.0
const XM = 0.2

//********* Generator Structs *********//

type DataGen struct {
	DBConn					*sql.DB
	UserPopularityScores	[]float64
	UserCommentScores		[]float64
	UserLikeScores			[]float64
	UserMessageScores		[]float64
}

//********* Diaspora Structs *********//

type DUser struct {
	User_ID       		int
	Person_ID     		int
	Aspects       		[]int
}

type Post struct {
	ID 						int
	Author					int
	Score					float64
}