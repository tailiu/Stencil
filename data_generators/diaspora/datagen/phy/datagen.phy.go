package phy

import (
	"database/sql"
	"diaspora/datagen"
	"diaspora/db"
	"diaspora/helper"
	"errors"
	"fmt"
	"log"
	"os"
	"stencil/qr"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User datagen.User
type Post datagen.Post

var flog, err = os.Create("./datagen.log")

func NewUser(QR *qr.QR, dbConn *sql.DB) (int32, int32, []int32) {

	var QIs []*qr.QI

	username := helper.RandomString(helper.RandomNumber(4, 10))
	serialized_private_key := "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----\n"
	language := "en"
	email := fmt.Sprintf("%s@%s.com", username, helper.RandomString(helper.RandomNumber(2, 8)))
	encrypted_password := "$2a$10$408zooOxx9.C.sNm9Csg0.uY83YZ.1f6qX1m4tn3D8tD03jbPPs62"
	color_theme := "original"
	guid := uuid.New()
	diaspora_handle := fmt.Sprintf("%s@127.0.0.1", username)
	serialized_public_key := "-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAxebDG8MvymvjtN7RgcO6\nn0KtvNQKrlPUQMt9GsSrujW19ZK5vf4PLUEB8CrRJmkYnEdh7qQtF/6deyOEHxLN\nzSe1vO1H/1bFnk5daZ+I2ttfdS/KfKKG2V16K8T4/Xf4MYjRpHHgbVP8+D1ecaHX\nFddUc5EbSfuZ1mFLpFtM56TsbQBtV0GZTcv68tQklWO3SrnUhEoaIwWWRBYv4BJt\nqNpZHhebfpX4nzJNs15q3VJhBCE6p9Kmu8s217Jb6NqWfdMpJaUk6rQ6iOjzNHxh\nxJnoGCwUEbITyEu4AOd+EWHXpRpMKY+Ytmjlym2+7iJPkJRZLVLwTwFoGOYsODcA\n8bHUI2kBtE2SOUYW9ax5FgvFJ7nltgDkw6EI/bb3Mejl2fmhqZ5ELmjUhswrW/2G\njHqBaVYhwXYyMA1CAkQ9x4Ffj0ojHtshHvGWKblZhzgYuVryLc1Q5aS3DwAic2AJ\nF0eycGZSYsPDPLXbaIk+PKB07G8QhLLA1+MSZ3CAkFVo+5iLI+fADXcSo0uCvqL4\nmlOu4bLsjgtkjpGtyY8frJH+9WUO36+tDCkY1Hvq1/sz4s4iRcID4rwTwI7iZRZA\nGtatWAX+gQ6HpT/Opao1a2EsG9xVh9zQDT9TpUs6bHRV3fXU+LvfHTIjLg8foCHT\ngT0i/y1piQH/7pN/IbTKV+ECAwEAAQ==\n-----END PUBLIC KEY-----\n"
	full_name := fmt.Sprintf("%s %s", helper.RandomString(5), helper.RandomString(5))
	sign_in_count := 1
	current_sign_in_ip := "127.0.0.1"
	last_sign_in_ip := "127.0.0.1"
	var aspect_ids []int32

	// sql := "INSERT INTO users (username, serialized_private_key, language, email, encrypted_password, created_at, updated_at, color_theme, last_seen, sign_in_count, current_sign_in_ip, last_sign_in_ip, current_sign_in_at, last_sign_in_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id "
	// user_id, _ := db.RunTxWQnArgsReturningId(tx, sql, username, serialized_private_key, language, email, encrypted_password, time.Now(), time.Now(), color_theme, time.Now(), sign_in_count, current_sign_in_ip, last_sign_in_ip, time.Now(), time.Now())
	user_id := QR.NewRowId()
	cols := []string{"id", "username", "serialized_private_key", "language", "email", "encrypted_password", "created_at", "updated_at", "color_theme", "last_seen", "sign_in_count", "current_sign_in_ip", "last_sign_in_ip", "current_sign_in_at", "last_sign_in_at"}
	vals := []interface{}{user_id, username, serialized_private_key, language, email, encrypted_password, time.Now(), time.Now(), color_theme, time.Now(), sign_in_count, current_sign_in_ip, last_sign_in_ip, time.Now(), time.Now()}
	qi := qr.CreateQI("users", cols, vals, qr.QTInsert)

	userQIs := QR.ResolveInsert(qi, user_id)
	QIs = append(QIs, userQIs...)

	// sql = "INSERT INTO people (guid,diaspora_handle,serialized_public_key,owner_id,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	// person_id, _ := db.RunTxWQnArgsReturningId(tx, sql, guid, diaspora_handle, serialized_public_key, user_id, time.Now(), time.Now())

	person_id := QR.NewRowId()
	cols = []string{"id", "guid", "diaspora_handle", "serialized_public_key", "owner_id", "created_at", "updated_at"}
	vals = []interface{}{person_id, guid, diaspora_handle, serialized_public_key, user_id, time.Now(), time.Now()}
	qi = qr.CreateQI("people", cols, vals, qr.QTInsert)

	peopleQIs := QR.ResolveInsert(qi, person_id)
	QIs = append(QIs, peopleQIs...)

	// sql = "INSERT INTO profiles (person_id, created_at, updated_at, full_name) VALUES ($1, $2, $3, $4)"
	// db.RunTxWQnArgs(tx, sql, person_id, time.Now(), time.Now(), full_name)

	profile_id := QR.NewRowId()
	cols = []string{"id", "person_id", "created_at", "updated_at", "full_name"}
	vals = []interface{}{profile_id, person_id, time.Now(), time.Now(), full_name}
	qi = qr.CreateQI("profiles", cols, vals, qr.QTInsert)

	profileQIs := QR.ResolveInsert(qi, profile_id)
	QIs = append(QIs, profileQIs...)

	// sql = "INSERT INTO aspects (name, user_id, created_at, updated_at, order_id) VALUES ($1, $2, $3, $4, $5)  RETURNING id"

	cols = []string{"id", "name", "user_id", "created_at", "updated_at", "order_id"}
	for idx, aspect_name := range []string{"Family", "Friends", "Work", "Acquaintances"} {
		aspect_id := QR.NewRowId()
		vals = []interface{}{aspect_id, aspect_name, user_id, time.Now(), time.Now(), idx + 1}
		qi = qr.CreateQI("aspects", cols, vals, qr.QTInsert)
		aspectQI := QR.ResolveInsert(qi, aspect_id)
		QIs = append(QIs, aspectQI...)
		// aspect_id, _ := db.RunTxWQnArgsReturningId(tx, sql, aspect_name, user_id, time.Now(), time.Now(), idx+1)
		aspect_ids = append(aspect_ids, aspect_id)
	}

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("create user transaction can't even begin")
	} else {
		success := true
		for _, qi := range QIs {
			query, args := qi.GenSQL()
			// fmt.Println(query)
			if err := db.RunTxWQnArgs(tx, query, args...); err != nil {
				success = false
				fmt.Println("Some error:", err)
				break
			}
		}
		if success {
			// fmt.Println("SUCCESS!")
			tx.Commit()
		}
	}
	// tx.Rollback()
	// log.Fatal("stop")

	return user_id, person_id, aspect_ids
}

func NewPost(QR *qr.QR, dbConn *sql.DB, user_id, person_id int, aspect_ids []int) int32 {

	var QIs []*qr.QI

	// Params
	guid := uuid.New()
	post_type := "StatusMessage"
	text := helper.RandomText(helper.RandomNumber(20, 200))

	// SQLs

	// sql := "INSERT INTO posts (author_id, guid, type, text, created_at, updated_at, interacted_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	// post_id, _ := db.RunTxWQnArgsReturningId(tx, sql, person_id, guid, post_type, text, time.Now(), time.Now(), time.Now())

	post_id := QR.NewRowId()
	cols := []string{"id", "author_id", "guid", "type", "text", "created_at", "updated_at", "interacted_at"}
	vals := []interface{}{post_id, person_id, guid, post_type, text, time.Now(), time.Now(), time.Now()}
	qi := qr.CreateQI("posts", cols, vals, qr.QTInsert)

	postQIs := QR.ResolveInsert(qi, post_id)
	QIs = append(QIs, postQIs...)

	// sql = "INSERT INTO aspect_visibilities (shareable_id, aspect_id) VALUES ($1, $2)"

	for _, aid := range aspect_ids {
		if helper.RandomNumber(1, 50)%2 == 0 {
			av_id := QR.NewRowId()
			cols = []string{"id", "shareable_id", "aspect_id"}
			vals = []interface{}{av_id, post_id, aid}
			qi = qr.CreateQI("aspect_visibilities", cols, vals, qr.QTInsert)

			avQIs := QR.ResolveInsert(qi, av_id)
			QIs = append(QIs, avQIs...)
			// db.RunTxWQnArgs(tx, sql, post_id, aid)
		}
	}

	// sql = "INSERT INTO share_visibilities (shareable_id, user_id) VALUES ($1, $2)"
	// db.RunTxWQnArgs(tx, sql, post_id, user_id)

	sv_id := QR.NewRowId()
	cols = []string{"id", "shareable_id", "user_id"}
	vals = []interface{}{sv_id, post_id, user_id}
	qi = qr.CreateQI("share_visibilities", cols, vals, qr.QTInsert)

	avQIs := QR.ResolveInsert(qi, sv_id)
	QIs = append(QIs, avQIs...)

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("create post transaction can't even begin")
	} else {
		success := true
		for _, qi := range QIs {
			query, args := qi.GenSQL()
			// fmt.Println(query)
			if err := db.RunTxWQnArgs(tx, query, args...); err != nil {
				success = false
				fmt.Println("Some error:", err)
				break
			}
		}
		if success {
			// fmt.Println("~ success")
			tx.Commit()
		}
	}

	// tx.Commit()

	return post_id
}

func NewComment(QR *qr.QR, dbConn *sql.DB, post_id, person_id, post_owner_id int) {

	var QIs []*qr.QI

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("create comment transaction can't even begin")
	}

	// Params

	text := helper.RandomText(helper.RandomNumber(10, 100))
	guid := uuid.New()
	target_type := "Post"
	notif_type := "Notifications::CommentOnPost"
	notif_id := QR.NewRowId()
	// SQLs

	{
		id := QR.NewRowId()
		cols := []string{"id", "text", "commentable_id", "author_id", "guid", "created_at", "updated_at"}
		vals := []interface{}{id, text, post_id, person_id, guid, time.Now(), time.Now()}
		qi := qr.CreateQI("comments", cols, vals, qr.QTInsert)
		qis := QR.ResolveInsert(qi, id)
		QIs = append(QIs, qis...)
	}

	// {
	// 	qu := qr.CreateQU(QR)
	// 	qu.SetTable("posts")
	// 	qu.SetUpdate_("posts.comments_count", "comments_count::int + 1")
	// 	qu.SetWhere("posts.type", "=", "StatusMessage")
	// 	qu.SetWhere("posts.id", "=", fmt.Sprint(post_id))
	// 	for _, sql := range qu.GenSQL() {
	// 		if err := db.RunTxWQnArgs(tx, sql); err != nil {
	// 			return
	// 		}
	// 	}
	// }

	{
		cols := []string{"id", "target_type", "target_id", "recipient_id", "created_at", "updated_at", "type"}
		vals := []interface{}{notif_id, target_type, post_id, post_owner_id, time.Now(), time.Now(), notif_type}
		qi := qr.CreateQI("notifications", cols, vals, qr.QTInsert)
		qis := QR.ResolveInsert(qi, notif_id)
		QIs = append(QIs, qis...)
	}

	{
		id := QR.NewRowId()
		cols := []string{"id", "notification_id", "person_id", "created_at", "updated_at"}
		vals := []interface{}{id, notif_id, person_id, time.Now(), time.Now()}
		qi := qr.CreateQI("notification_actors", cols, vals, qr.QTInsert)
		qis := QR.ResolveInsert(qi, id)
		QIs = append(QIs, qis...)
	}

	success := true
	for _, qi := range QIs {
		query, args := qi.GenSQL()
		// log.Fatal(query)
		if err := db.RunTxWQnArgs(tx, query, args...); err != nil {
			success = false
			fmt.Println("Some error:", err)
			break
		}
	}
	if success {
		// fmt.Println("~ success")
		tx.Commit()
	}
}

func GetPostsForUser(QR *qr.QR, dbConn *sql.DB, user_id int) []*Post {

	var posts []*Post

	// sql := `SELECT id, guid, author_id, text
	// 		FROM posts
	// 		WHERE author_id = $1 AND id NOT IN (
	// 			SELECT distinct(target_id) FROM likes
	// 			UNION
	// 			SELECT distinct(commentable_id) FROM comments
	// 		)
	// 		order by random()`

	interactedAlready := make(map[string]bool)

	qComments := qr.CreateQS(QR)
	qComments.FromSimple("comments")
	qComments.ColFunction("distinct(%s)", "comments.commentable_id", "post_id")
	for _, row := range db.DataCall(dbConn, qComments.GenSQL()) {
		interactedAlready[row["post_id"]] = true
	}

	qLikes := qr.CreateQS(QR)
	qLikes.FromSimple("likes")
	qLikes.ColFunction("distinct(%s)", "likes.target_id", "post_id")
	for _, row := range db.DataCall(dbConn, qLikes.GenSQL()) {
		interactedAlready[row["post_id"]] = true
	}

	qPosts := qr.CreateQS(QR)
	qPosts.FromSimple("posts")
	qPosts.ColSimple("posts.id")
	qPosts.ColSimple("posts.guid")
	qPosts.ColSimple("posts.author_id")
	qPosts.ColSimple("posts.text")
	qPosts.WhereSimpleVal("posts.author_id", "=", fmt.Sprint(user_id))
	qPosts.OrderBy("random()")

	for _, row := range db.DataCall(dbConn, qPosts.GenSQL()) {
		if _, ok := interactedAlready[row["id"]]; !ok {
			if pid, err := strconv.Atoi(row["id"]); err == nil {
				if uid, err := strconv.Atoi(row["author_id"]); err == nil {
					post := new(Post)
					post.Author = uid
					post.ID = pid
					post.GUID = row["guid"]
					post.Text = row["text"]
					posts = append(posts, post)
					// fmt.Println(post)
				}
			}
		} else {
			// fmt.Println("skipped")
		}
	}

	return posts
}

func NewLike(QR *qr.QR, dbConn *sql.DB, post_id, person_id, post_owner_id int) {

	var QIs []*qr.QI

	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal("create like transaction can't even begin")
	}

	// Params

	guid := uuid.New()
	target_type := "Post"
	notif_type := "Notifications::Liked"
	notif_id := QR.NewRowId()

	{
		id := QR.NewRowId()
		cols := []string{"id", "target_id", "author_id", "guid", "created_at", "updated_at", "target_type"}
		vals := []interface{}{id, post_id, person_id, guid, time.Now(), time.Now(), target_type}
		qi := qr.CreateQI("likes", cols, vals, qr.QTInsert)
		qis := QR.ResolveInsert(qi, id)
		QIs = append(QIs, qis...)
	}

	// {
	// 	qu := qr.CreateQU(QR)
	// 	qu.SetTable("posts")
	// 	qu.SetUpdate_("posts.likes_count", "likes_count::int + 1")
	// 	qu.SetWhere("posts.type", "=", "StatusMessage")
	// 	qu.SetWhere("posts.id", "=", fmt.Sprint(post_id))
	// 	for _, sql := range qu.GenSQL() {
	// 		if err := db.RunTxWQnArgs(tx, sql); err != nil {
	// 			return
	// 		}
	// 	}
	// }

	{
		cols := []string{"id", "target_type", "target_id", "recipient_id", "created_at", "updated_at", "type"}
		vals := []interface{}{notif_id, target_type, post_id, post_owner_id, time.Now(), time.Now(), notif_type}
		qi := qr.CreateQI("notifications", cols, vals, qr.QTInsert)
		qis := QR.ResolveInsert(qi, notif_id)
		QIs = append(QIs, qis...)
	}

	{
		id := QR.NewRowId()
		cols := []string{"id", "notification_id", "person_id", "created_at", "updated_at"}
		vals := []interface{}{id, notif_id, person_id, time.Now(), time.Now()}
		qi := qr.CreateQI("notification_actors", cols, vals, qr.QTInsert)
		qis := QR.ResolveInsert(qi, id)
		QIs = append(QIs, qis...)
	}

	success := true
	for _, qi := range QIs {
		query, args := qi.GenSQL()
		// log.Fatal(query)
		if err := db.RunTxWQnArgs(tx, query, args...); err != nil {
			success = false
			fmt.Println("Some error:", err)
			break
		}
	}
	if success {
		// fmt.Println("~ success")
		tx.Commit()
	}
}

func FollowUser(QR *qr.QR, dbConn *sql.DB, person_id_1, person_id_2, aspect_id int) {

	var QIs []*qr.QI
	// log.Println("Creating new follow!")

	ok1, contact_id1 := ContactExists(QR, person_id_2, person_id_1)
	ok2, contact_id2 := ContactExists(QR, person_id_1, person_id_2)

	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal("create follow transaction can't even begin")
	}

	// Params

	success := true

	// SQLs

	if ok1 {
		// sql := "UPDATE contacts SET sharing = $1, updated_at = $2 WHERE contacts.id = $3"
		qu := qr.CreateQU(QR)
		qu.SetTable("contacts")
		qu.SetUpdate("contacts.sharing", "t")
		qu.SetUpdate_("contacts.updated_at", "current_timestamp")
		qu.SetWhere("contacts.id", "=", contact_id1)
		for _, sql := range qu.GenSQL() {
			// fmt.Println(sql)
			if err := db.RunTxWQnArgs(tx, sql); err != nil {
				return
			}
		}

		if ok2 {
			// sql = "UPDATE contacts SET receiving = $1, updated_at = $2 WHERE contacts.id = $3"
			qu := qr.CreateQU(QR)
			qu.SetTable("contacts")
			qu.SetUpdate("contacts.receiving", "t")
			qu.SetUpdate_("contacts.updated_at", "current_timestamp")
			qu.SetWhere("contacts.id", "=", contact_id2)
			for _, sql := range qu.GenSQL() {
				// log.Fatal(sql)
				if err := db.RunTxWQnArgs(tx, sql); err != nil {

				}
			}

			// sql = "INSERT INTO aspect_memberships (aspect_id,aspect_id,created_at,updated_at) VALUES ($1, $2, $3, $4)"
			am_id := QR.NewRowId()
			cols := []string{"id", "aspect_id", "contact_id", "created_at", "updated_at"}
			vals := []interface{}{am_id, aspect_id, contact_id2, time.Now(), time.Now()}
			qi := qr.CreateQI("aspect_memberships", cols, vals, qr.QTInsert)

			amQIs := QR.ResolveInsert(qi, am_id)
			QIs = append(QIs, amQIs...)
			// db.RunTxWQnArgs(tx, sql, aspect_id, contact_id2, time.Now(), time.Now())
		}
	} else {

		// sql := "INSERT INTO contacts (user_id,person_id,created_at,updated_at,receiving) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		//contact_id, _ := db.RunTxWQnArgsReturningId(tx, sql, person_id_1, person_id_2, time.Now(), time.Now(), "t")

		contact1_id := QR.NewRowId()
		cols1 := []string{"id", "user_id", "person_id", "created_at", "updated_at", "receiving"}
		vals1 := []interface{}{contact1_id, person_id_1, person_id_2, time.Now(), time.Now(), "t"}
		q_contact1 := qr.CreateQI("contacts", cols1, vals1, qr.QTInsert)

		contact1QIs := QR.ResolveInsert(q_contact1, contact1_id)
		QIs = append(QIs, contact1QIs...)

		// sql = "INSERT INTO contacts (user_id,person_id,created_at,updated_at,sharing) VALUES ($1, $2, $3, $4, $5)"
		// db.RunTxWQnArgs(tx, sql, person_id_2, person_id_1, time.Now(), time.Now(), "t")

		contact2_id := QR.NewRowId()
		cols2 := []string{"id", "user_id", "person_id", "created_at", "updated_at", "receiving"}
		vals2 := []interface{}{contact2_id, person_id_2, person_id_1, time.Now(), time.Now(), "t"}
		q_contact2 := qr.CreateQI("contacts", cols2, vals2, qr.QTInsert)

		contact2QIs := QR.ResolveInsert(q_contact2, contact2_id)
		QIs = append(QIs, contact2QIs...)

		// sql = "INSERT INTO aspect_memberships (aspect_id,contact_id,created_at,updated_at) VALUES ($1, $2, $3, $4)"
		// db.RunTxWQnArgs(tx, sql, aspect_id, contact_id, time.Now(), time.Now())

		am_id := QR.NewRowId()
		cols := []string{"id", "aspect_id", "contact_id", "created_at", "updated_at"}
		vals := []interface{}{am_id, aspect_id, contact1_id, time.Now(), time.Now()}
		qi := qr.CreateQI("aspect_memberships", cols, vals, qr.QTInsert)

		amQIs := QR.ResolveInsert(qi, am_id)
		QIs = append(QIs, amQIs...)
	}

	target_type := "Person"
	notif_type := "Notifications::StartedSharing"

	// sql := "INSERT INTO notifications (target_type,target_id,recipient_id,created_at,updated_at,type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	// notif_id, _ := db.RunTxWQnArgsReturningId(tx, sql, target_type, person_id_1, person_id_2, time.Now(), time.Now(), notif_type)

	notif_id := QR.NewRowId()
	cols := []string{"id", "target_type", "target_id", "recipient_id", "created_at", "updated_at", "type"}
	vals := []interface{}{notif_id, target_type, person_id_1, person_id_2, time.Now(), time.Now(), notif_type}
	notif_qi := qr.CreateQI("notifications", cols, vals, qr.QTInsert)

	notifQIs := QR.ResolveInsert(notif_qi, notif_id)
	QIs = append(QIs, notifQIs...)

	// sql = "INSERT INTO notification_actors (notification_id,person_id,created_at,updated_at) VALUES ($1, $2, $3, $4) RETURNING id"
	// db.RunTxWQnArgs(tx, sql, notif_id, person_id_1, time.Now(), time.Now())

	notif_actor_id := QR.NewRowId()
	cols = []string{"id", "notification_id", "person_id", "created_at", "updated_at"}
	vals = []interface{}{notif_actor_id, notif_id, person_id_1, time.Now(), time.Now()}
	notif_actor_qi := qr.CreateQI("notification_actors", cols, vals, qr.QTInsert)

	notifactorQIs := QR.ResolveInsert(notif_actor_qi, notif_actor_id)
	QIs = append(QIs, notifactorQIs...)

	for _, qi := range QIs {
		query, args := qi.GenSQL()
		// fmt.Println(query, args)
		if err := db.RunTxWQnArgs(tx, query, args...); err != nil {
			success = false
			log.Println("Some error:", err)
			break
		}
	}
	if success {
		// fmt.Println("~ success")
		tx.Commit()
	}
	// log.Fatal("person_id_1: ", person_id_1, " | person_id_2: ", person_id_2)
}

func ContactExists(QR *qr.QR, person_id_1, person_id_2 int) (bool, string) {

	// sql := "SELECT id FROM contacts WHERE user_id = $1 AND person_id = $2"\

	q := qr.CreateQS(QR)
	q.FromSimple("contacts")
	q.ColSimple("contacts.id")
	q.WhereSimpleVal("contacts.user_id", "=", fmt.Sprint(person_id_1))
	q.WhereOperatorVal("AND", "contacts.person_id", "=", fmt.Sprint(person_id_2))
	sql := q.GenSQL()
	// fmt.Println(sql)
	// log.Println("checking", sql, person_id_1, person_id_2)
	res := db.DataCall1(QR.StencilDB, sql)
	// log.Println("result of contact exists", res)
	if len(res) > 0 {
		return true, res[0]["contacts.id"]
	}
	return false, ""
}

func AspectMembershipExists(QR *qr.QR, contact_id, aspect_id int) bool {
	sql := "SELECT id FROM aspect_memberships WHERE aspect_id = $1 AND contact_id = $2 LIMIT 1"
	res := db.DataCall(QR.StencilDB, sql, aspect_id, contact_id)
	if len(res) > 0 {
		return true
	}
	return false
}

func NewReshare(QR *qr.QR, dbConn *sql.DB, post Post, person_id int) {

	var QIs []*qr.QI

	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal("create comment transaction can't even begin")
	}

	// Params

	target_type := "Post"
	notif_type := "Notifications::Reshared"
	notif_id := QR.NewRowId()

	{
		id := QR.NewRowId()
		cols := []string{"id", "author_id", "public", "guid", "type", "text", "created_at", "updated_at", "root_guid", "interacted_at"}
		vals := []interface{}{id, person_id, "t", uuid.New(), "Reshare", post.Text, time.Now(), time.Now(), post.GUID, time.Now()}
		qi := qr.CreateQI("posts", cols, vals, qr.QTInsert)
		qis := QR.ResolveInsert(qi, id)
		QIs = append(QIs, qis...)
	}

	// {
	// 	qu := qr.CreateQU(QR)
	// 	qu.SetTable("posts")
	// 	qu.SetUpdate_("posts.reshares_count", "reshares_count::int +1")
	// 	qu.SetWhere("posts.type", "=", "StatusMessage")
	// 	qu.SetWhere("posts.id", "=", fmt.Sprint(post.ID))
	// 	for _, sql := range qu.GenSQL() {
	// 		if err := db.RunTxWQnArgs(tx, sql); err != nil {
	// 			return
	// 		}
	// 	}
	// }

	{
		id := QR.NewRowId()
		cols := []string{"id", "guid", "target_id", "target_type", "author_id", "created_at", "updated_at"}
		vals := []interface{}{id, uuid.New(), post.ID, target_type, post.Author, time.Now(), time.Now()}
		qi := qr.CreateQI("participations", cols, vals, qr.QTInsert)
		qis := QR.ResolveInsert(qi, id)
		QIs = append(QIs, qis...)
	}

	{
		id := QR.NewRowId()
		cols := []string{"id", "guid", "target_id", "target_type", "author_id", "created_at", "updated_at"}
		vals := []interface{}{id, uuid.New(), post.ID, target_type, person_id, time.Now(), time.Now()}
		qi := qr.CreateQI("participations", cols, vals, qr.QTInsert)
		qis := QR.ResolveInsert(qi, id)
		QIs = append(QIs, qis...)
	}

	{
		cols := []string{"id", "target_type", "target_id", "recipient_id", "created_at", "updated_at", "type"}
		vals := []interface{}{notif_id, target_type, post.ID, post.Author, time.Now(), time.Now(), notif_type}
		qi := qr.CreateQI("notifications", cols, vals, qr.QTInsert)
		qis := QR.ResolveInsert(qi, notif_id)
		QIs = append(QIs, qis...)
	}

	{
		id := QR.NewRowId()
		cols := []string{"id", "notification_id", "person_id", "created_at", "updated_at"}
		vals := []interface{}{id, notif_id, person_id, time.Now(), time.Now()}
		qi := qr.CreateQI("notification_actors", cols, vals, qr.QTInsert)
		qis := QR.ResolveInsert(qi, id)
		QIs = append(QIs, qis...)
	}

	success := true
	for _, qi := range QIs {
		query, args := qi.GenSQL()
		// log.Fatal(query)
		if err := db.RunTxWQnArgs(tx, query, args...); err != nil {
			success = false
			fmt.Println("Some error:", err)
			break
		}
	}
	if success {
		// fmt.Println("~ success")
		tx.Commit()
	}
}

func NewConversation(QR *qr.QR, dbConn *sql.DB, person_id_1, person_id_2 string) (int32, error) {

	// log.Println("person_id_1, person_id_2: ", person_id_1, person_id_2)

	var QIs []*qr.QI

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println("create conversation transaction can't even begin")
		return -1, errors.New("create conversation transaction can't begin")
	}

	// Params

	subject := helper.RandomText(helper.RandomNumber(5, 15))
	guid := uuid.New()
	conversation_id := QR.NewRowId()

	{
		cols := []string{"id", "subject", "guid", "author_id", "created_at", "updated_at"}
		vals := []interface{}{conversation_id, subject, guid, person_id_1, time.Now(), time.Now()}
		qi := qr.CreateQI("conversations", cols, vals, qr.QTInsert)
		qis := QR.ResolveInsert(qi, conversation_id)
		QIs = append(QIs, qis...)
	}

	{
		conversation_visibilities_id := QR.NewRowId()
		cols := []string{"id", "conversation_id", "person_id", "created_at", "updated_at"}
		vals := []interface{}{conversation_visibilities_id, conversation_id, person_id_1, time.Now(), time.Now()}
		qi := qr.CreateQI("conversation_visibilities", cols, vals, qr.QTInsert)
		qis := QR.ResolveInsert(qi, conversation_visibilities_id)
		QIs = append(QIs, qis...)
	}

	{
		conversation_visibilities_id := QR.NewRowId()
		cols := []string{"id", "conversation_id", "person_id", "created_at", "updated_at"}
		vals := []interface{}{conversation_visibilities_id, conversation_id, person_id_2, time.Now(), time.Now()}
		qi := qr.CreateQI("conversation_visibilities", cols, vals, qr.QTInsert)
		qis := QR.ResolveInsert(qi, conversation_visibilities_id)
		QIs = append(QIs, qis...)
	}

	msgQIs := GenNewMessage(QR, person_id_1, fmt.Sprint(conversation_id))
	QIs = append(QIs, msgQIs...)

	success := true
	for _, qi := range QIs {
		query, args := qi.GenSQL()
		// fmt.Println(query, args)
		if err := db.RunTxWQnArgs(tx, query, args...); err != nil {
			success = false
			fmt.Println("Some error:", err)
			break
		}
	}
	if success {
		// fmt.Println("~ successconvo")
		tx.Commit()
		return conversation_id, err
	}

	return -1, err
}

func GenNewMessage(QR *qr.QR, person_id, conversation_id string) []*qr.QI {

	msgtext := helper.RandomText(helper.RandomNumber(20, 100))
	message_id := QR.NewRowId()
	cols := []string{"id", "conversation_id", "author_id", "guid", "text", "created_at", "updated_at"}
	vals := []interface{}{message_id, conversation_id, person_id, uuid.New(), msgtext, time.Now(), time.Now()}
	qi := qr.CreateQI("messages", cols, vals, qr.QTInsert)
	return QR.ResolveInsert(qi, message_id)
}

func NewMessage(QR *qr.QR, dbConn *sql.DB, person_id, conversation_id string) error {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println("create message transaction can't even begin")
		return errors.New("Message Transaction can't begin")
	}

	success := true
	for _, qi := range GenNewMessage(QR, person_id, conversation_id) {
		query, args := qi.GenSQL()
		if err := db.RunTxWQnArgs(tx, query, args...); err != nil {
			success = false
			fmt.Println("Some error:", err)
			break
		}
	}
	if success {
		// fmt.Println("~ successmsg")
		tx.Commit()
	}

	return err
}

func UpdateConversation(QR *qr.QR, dbConn *sql.DB, conversation_id int) {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println("UpdateConversation transaction can't even begin")
		return
	}

	// sql := "UPDATE conversations SET updated_at = $1 WHERE conversations.id = $2 "
	// db.RunTxWQnArgs(tx, sql, time.Now(), conversation_id)

	qu := qr.CreateQU(QR)
	qu.SetTable("conversations")
	qu.SetUpdate_("conversations.updated_at", "current_timestamp")
	qu.SetWhere("conversations.id", "=", fmt.Sprint(conversation_id))
	for _, sql := range qu.GenSQL() {
		// fmt.Println(sql)
		if err := db.RunTxWQnArgs(tx, sql); err != nil {
			return
		}
	}

	// sql = "UPDATE conversation_visibilities SET unread = $1, updated_at = $2 WHERE conversation_visibilities.id = $3"
	// db.RunTxWQnArgs(tx, sql, 1, time.Now(), conversation_visibilities_id)

	tx.Commit()
}

func GetAllUsersWithAspectsNew(QR *qr.QR) []*User {

	var users []*User

	log.Println("fetching users")
	qu := qr.CreateQS(QR)
	qu.FromSimple("users")
	qu.ColAlias("users.id", "user_id")
	// qu.LimitResult("10")
	fmt.Println(qu.GenSQL())
	res_users := db.DataCall(QR.StencilDB, qu.GenSQL())

	log.Println("fetching people")
	qp := qr.CreateQS(QR)
	qp.FromSimple("people")
	qp.ColAlias("people.owner_id", "user_id")
	qp.ColAlias("people.id", "person_id")
	// qp.LimitResult("10")
	fmt.Println(qp.GenSQL())
	res_people := db.DataCall(QR.StencilDB, qp.GenSQL())

	log.Println("fetching aspects")
	qa := qr.CreateQS(QR)
	qa.FromSimple("aspects")
	qa.ColSimple("aspects.user_id")
	qa.ColAlias("aspects.id", "aspect_id")
	// qa.LimitResult("100")
	fmt.Println(qa.GenSQL())
	res_aspects := db.DataCall(QR.StencilDB, qa.GenSQL())

	log.Println("looping")
	for _, user := range res_users {
		newuser := new(User)
		newuser.User_ID, _ = strconv.Atoi(user["user_id"])

		for _, person := range res_people {
			if strings.EqualFold(user["user_id"], person["user_id"]) {
				newuser.Person_ID, _ = strconv.Atoi(person["person_id"])
				break
			}
		}

		var aspect_ids []int
		for _, aspect := range res_aspects {
			if strings.EqualFold(user["user_id"], aspect["user_id"]) {
				aspect_id, _ := strconv.Atoi(aspect["aspect_id"])
				aspect_ids = append(aspect_ids, aspect_id)
				if len(aspect_ids) >= 4 {
					break
				}
			}
		}

		newuser.Aspects = aspect_ids
		users = append(users, newuser)
		fmt.Println(newuser)
	}

	return users
}

func GetAllUsersWithAspects(QR *qr.QR) []*User {

	var users []*User

	// sql := `select * from diaspora_users_with_aspects`

	fq := qr.CreateQS(QR)
	fq.FromSimple("users")
	fq.FromJoin("people", "users.id=people.owner_id")
	fq.FromJoin("aspects", "users.id=aspects.user_id")
	fq.ColAlias("users.id", "user_id")
	fq.ColAlias("people.id", "person_id")
	fq.ColAlias("aspects.id", "aspect_id")

	q := qr.CreateQS(QR)
	q.ColSimple("user_id,person_id")
	q.ColFunction("string_agg(%s::text, ',')", "aspect_id", "aspects")
	q.FromQuery(fq)
	q.GroupByString("user_id,person_id")
	q.OrderBy("random()")

	sql := q.GenSQL()

	res := db.DataCall(QR.StencilDB, sql)

	for _, row := range res {
		user := new(User)
		user.User_ID, _ = strconv.Atoi(row["user_id"])
		user.Person_ID, _ = strconv.Atoi(row["person_id"])
		var aspect_ids []int
		for _, aspect_id := range strings.Split(row["aspects"], ",") {
			aspect_id, _ := strconv.Atoi(aspect_id)
			aspect_ids = append(aspect_ids, aspect_id)
		}
		user.Aspects = aspect_ids
		users = append(users, user)
	}
	return users
}

func GetAllUsersWithAspectsExcept(QR *qr.QR, column, table string) []*User {

	var users []*User

	// sql := fmt.Sprintf(`SELECT user_id, person_id, string_agg(aspect_id::text, ',') as aspects
	// 		FROM (
	// 			SELECT users.id as user_id, people.id as person_id, aspects.id as aspect_id
	// 			FROM users JOIN people ON users.id = people.owner_id JOIN aspects ON aspects.user_id = users.id
	// 		) tab
	// 		WHERE user_id NOT IN (
	// 			SELECT DISTINCT %s FROM %s
	// 		)
	// 		GROUP BY user_id, person_id
	// 		ORDER BY random()`, column, table)

	fq := qr.CreateQS(QR)
	wq := qr.CreateQS(QR)
	q := qr.CreateQS(QR)
	q.ColSimple("col_name")
	q.ColAlias("col_name", "alias")
	q.ColFunction("string_agg(col_name::text, ',')", "aspect_id", "aspects")
	q.FromSimple("tab_name")
	q.FromJoin("tab_name", "tab_name.col1 = tab_name2.col2")
	q.FromQuery(fq)
	// q.WhereSimple("col != val")
	// q.WhereOperator("and", "col != val")
	q.WhereQuery("not in", wq)
	q.GroupBy("cols")
	q.OrderBy("cols")

	sql := q.GenSQL()

	res := db.DataCall(QR.StencilDB, sql)

	for _, row := range res {
		user := new(User)
		user.User_ID, _ = strconv.Atoi(row["user_id"])
		user.Person_ID, _ = strconv.Atoi(row["person_id"])
		var aspect_ids []int
		for _, aspect_id := range strings.Split(row["aspects"], ",") {
			aspect_id, _ := strconv.Atoi(aspect_id)
			aspect_ids = append(aspect_ids, aspect_id)
		}
		user.Aspects = aspect_ids
		users = append(users, user)
	}

	return users
}

func GetFriendsOfUser(QR *qr.QR, person_id int) []*User {

	var users []*User

	// sql := `SELECT users.id as user_id, people.id as person_id, string_agg(aspects.id::text, ',') as aspects, contacts.id as contact_id, am.aspect_id as contact_aspect
	// 		FROM contacts
	// 		JOIN aspect_memberships am on contacts.id = am.contact_id
	// 		JOIN people on people.id = contacts.user_id
	// 		JOIN users on users.id = people.owner_id
	// 		JOIN aspects on aspects.user_id = users.id
	// 		WHERE contacts.user_id = $1 AND contacts.sharing = true
	// 		GROUP BY users.id, people.id, contacts.id, am.aspect_id
	// 	`

	fq := qr.CreateQS(QR)
	fq.FromSimple("contacts")
	fq.FromJoin("aspect_memberships", "contacts.id=aspect_memberships.contact_id")
	fq.FromJoin("people", "contacts.user_id=people.id")
	fq.FromJoin("users", "people.owner_id=users.id")
	fq.FromJoin("aspects", "users.id=aspects.user_id")
	fq.ColAlias("users.id", "user_id")
	fq.ColAlias("people.id", "person_id")
	fq.ColAlias("contacts.id", "contact_id")
	fq.ColFunction("string_agg(%s::text, ',')", "aspects.id", "aspects")
	fq.ColAlias("aspect_memberships.aspect_id", "contact_aspect")
	fq.WhereSimpleVal("contacts.user_id", "=", fmt.Sprint(person_id))
	fq.WhereOperatorBool("AND", "contacts.sharing", "=", "true")
	fq.GroupBy("users.id")
	fq.GroupBy("people.id")
	fq.GroupBy("contacts.id")
	fq.GroupBy("aspect_memberships.aspect_id")

	log.Fatal(fq.GenSQL())

	sql := fq.GenSQL()

	// sql := `
	// 	SELECT supplementary_676.id as user_id,
	// 	base_people_1.id as person_id,
	// 	base_contacts_1.id as contact_id,
	// 	string_agg(supplementary_630.id::text, ',') as aspects,
	// 	supplementary_628.aspect_id as contact_aspect
	// 	FROM  base_contacts_1
	// 	JOIN supplementary_638 ON base_contacts_1.pk = supplementary_638.pk
	// 	JOIN supplementary_628 ON base_contacts_1.id::text = supplementary_628.contact_id::text
	// 	JOIN base_people_1 ON base_contacts_1.user_id::text = base_people_1.id::text
	// 	JOIN base_users_1 ON base_people_1.pk = base_users_1.pk
	// 	JOIN supplementary_654 ON base_users_1.pk = supplementary_654.pk
	// 	JOIN supplementary_676 ON supplementary_654.owner_id::text = supplementary_676.id::text
	// 	JOIN base_users_2 ON supplementary_676.pk = base_users_2.pk
	// 	JOIN base_users_1 buser1 ON base_users_2.pk = buser1.pk
	// 	JOIN supplementary_630 ON supplementary_676.id::text = supplementary_630.user_id::text
	// 	WHERE  base_contacts_1.user_id = $1  AND supplementary_638.sharing = true
	// 	GROUP BY supplementary_676.id , base_people_1.id , base_contacts_1.id , supplementary_628.aspect_id
	// `

	res := db.DataCall(QR.StencilDB, sql, person_id)

	for _, row := range res {
		user := new(User)
		user.User_ID, _ = strconv.Atoi(row["user_id"])
		user.Person_ID, _ = strconv.Atoi(row["person_id"])
		var aspect_ids []int
		for _, aspect_id := range strings.Split(row["aspects"], ",") {
			aspect_id, _ := strconv.Atoi(aspect_id)
			aspect_ids = append(aspect_ids, aspect_id)
		}
		user.ContactAspect, _ = strconv.Atoi(row["contact_aspect"])
		user.ContactID, _ = strconv.Atoi(row["contact_id"])
		users = append(users, user)
	}

	return users
}

func GetRandomUser(QR *qr.QR, except_id string) *User {

	if except_id == "" {
		except_id = "9999999"
	}

	sql := "SELECT users.id as user_id, people.id as person_id FROM users JOIN people ON users.id = people.owner_id WHERE users.id NOT IN ($1) ORDER BY random() LIMIT 1"

	res := db.DataCall(QR.StencilDB, sql, except_id)
	row := res[0]

	sql = "SELECT id FROM aspects WHERE user_id = $1"

	aspects := db.DataCall(QR.StencilDB, sql, row["user_id"])

	var aspect_ids []int

	for _, aspect := range aspects {
		aspect_id, _ := strconv.Atoi(aspect["id"])
		aspect_ids = append(aspect_ids, aspect_id)
	}

	user := new(User)
	user.User_ID, _ = strconv.Atoi(row["user_id"])
	user.Person_ID, _ = strconv.Atoi(row["person_id"])
	user.Aspects = aspect_ids

	return user
}
