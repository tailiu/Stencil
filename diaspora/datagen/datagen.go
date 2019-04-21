package datagen

import (
	"database/sql"
	"diaspora/db"
	"diaspora/helper"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	User_ID       int
	Person_ID     int
	Aspects       []int
	ContactID     int
	ContactAspect int
}

type Post struct {
	ID     int
	GUID   string
	Author int
	Text   string
}

func NewUser(dbConn *sql.DB) (int, int, []int) {

	// log.Println("Creating new user!")

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println(err)
		log.Fatal("create user transaction can't even begin")
	}

	// Params

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
	var aspect_ids []int

	// SQLs

	sql := "INSERT INTO users (username, serialized_private_key, language, email, encrypted_password, created_at, updated_at, color_theme, last_seen, sign_in_count, current_sign_in_ip, last_sign_in_ip, current_sign_in_at, last_sign_in_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id "
	user_id, _ := db.RunTxWQnArgsReturningId(tx, sql, username, serialized_private_key, language, email, encrypted_password, time.Now(), time.Now(), color_theme, time.Now(), sign_in_count, current_sign_in_ip, last_sign_in_ip, time.Now(), time.Now())

	sql = "INSERT INTO people (guid,diaspora_handle,serialized_public_key,owner_id,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	person_id, _ := db.RunTxWQnArgsReturningId(tx, sql, guid, diaspora_handle, serialized_public_key, user_id, time.Now(), time.Now())

	sql = "INSERT INTO profiles (person_id, created_at, updated_at, full_name) VALUES ($1, $2, $3, $4)"
	db.RunTxWQnArgs(tx, sql, person_id, time.Now(), time.Now(), full_name)

	// sql = "UPDATE users SET unconfirmed_email = NULL, confirm_email_token = NULL WHERE users.unconfirmed_email = $1"
	// db.RunTxWQnArgs(tx, sql, email)

	sql = "INSERT INTO aspects (name, user_id, created_at, updated_at, order_id) VALUES ($1, $2, $3, $4, $5)  RETURNING id"

	for idx, aspect_name := range []string{"Family", "Friends", "Work", "Acquaintances"} {
		aspect_id, _ := db.RunTxWQnArgsReturningId(tx, sql, aspect_name, user_id, time.Now(), time.Now(), idx+1)
		aspect_ids = append(aspect_ids, aspect_id)
	}
	// aspect_ids = append(aspect_ids, db.RunTxWQnArgsReturningId(tx, sql, "Family", user_id, time.Now(), time.Now(), 1))
	// aspect_ids = append(aspect_ids, db.RunTxWQnArgsReturningId(tx, sql, "Friends", user_id, time.Now(), time.Now(), 2))
	// aspect_ids = append(aspect_ids, db.RunTxWQnArgsReturningId(tx, sql, "Work", user_id, time.Now(), time.Now(), 3))
	// aspect_ids = append(aspect_ids, db.RunTxWQnArgsReturningId(tx, sql, "Acquaintances", user_id, time.Now(), time.Now(), 4))

	// sql = "UPDATE users SET sign_in_count = $1, current_sign_in_at = $2, last_sign_in_at = $3, current_sign_in_ip = $4, last_sign_in_ip = $5, updated_at = $6 WHERE users.id = $7"
	// db.RunTxWQnArgs(tx, sql, sign_in_count, time.Now(), time.Now(), current_sign_in_ip, last_sign_in_ip, time.Now(), user_id)

	// sql = "UPDATE users SET updated_at = $1, last_seen = $2 WHERE users.id = $3 "
	// db.RunTxWQnArgs(tx, sql, time.Now(), time.Now(), user_id)

	tx.Commit()

	// log.Println("New user created with id", user_id)

	return user_id, person_id, aspect_ids
}

func NewPost(dbConn *sql.DB, user_id, person_id int, aspect_ids []int) int {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println("create post transaction can't even begin")
		return -1
	}

	// Params
	guid := uuid.New()
	post_type := "StatusMessage"
	text := helper.RandomText(helper.RandomNumber(20, 200))

	// SQLs

	sql := "INSERT INTO posts (author_id, guid, type, text, created_at, updated_at, interacted_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	post_id, _ := db.RunTxWQnArgsReturningId(tx, sql, person_id, guid, post_type, text, time.Now(), time.Now(), time.Now())

	sql = "INSERT INTO aspect_visibilities (shareable_id, aspect_id) VALUES ($1, $2)"

	for _, aid := range aspect_ids {
		if helper.RandomNumber(1, 50)%2 == 0 {
			db.RunTxWQnArgs(tx, sql, post_id, aid)
		}
	}

	sql = "INSERT INTO share_visibilities (shareable_id, user_id) VALUES ($1, $2)"
	db.RunTxWQnArgs(tx, sql, post_id, user_id)

	tx.Commit()

	return post_id
}

func NewComment(dbConn *sql.DB, post_id, person_id, post_owner_id int) (int, error) {

	// log.Println("Creating new comment!")

	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal("create comment transaction can't even begin")
	}

	// Params

	text := helper.RandomText(helper.RandomNumber(10, 100))
	guid := uuid.New()
	target_type := "Post"
	notif_type := "Notifications::CommentOnPost"
	// SQLs

	sql := "INSERT INTO comments (text,commentable_id,author_id,guid,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	if comment_id, err := db.RunTxWQnArgsReturningId(tx, sql, text, post_id, person_id, guid, time.Now(), time.Now()); err == nil {
		sql = "UPDATE posts SET updated_at = $1 WHERE posts.id = $2 "
		if err := db.RunTxWQnArgs(tx, sql, time.Now(), post_id); err == nil {
			sql = "UPDATE posts SET comments_count = comments_count+1 WHERE posts.type IN ('StatusMessage') AND posts.id = $1"
			if err := db.RunTxWQnArgs(tx, sql, post_id); err == nil {
				sql = "UPDATE posts SET updated_at = $1, interacted_at = $2 WHERE posts.id = $3 "
				if err := db.RunTxWQnArgs(tx, sql, time.Now(), time.Now(), post_id); err == nil {
					sql = "INSERT INTO notifications (target_type,target_id,recipient_id,created_at,updated_at,type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
					if notif_id, err := db.RunTxWQnArgsReturningId(tx, sql, target_type, post_id, post_owner_id, time.Now(), time.Now(), notif_type); err == nil {
						sql = "INSERT INTO notification_actors (notification_id,person_id,created_at,updated_at) VALUES ($1, $2, $3, $4)"
						if err := db.RunTxWQnArgs(tx, sql, notif_id, person_id, time.Now(), time.Now()); err == nil {
							tx.Commit()
							// log.Println("New comment created with id", comment_id)
							return comment_id, err
						}
						return -1, err
					}
					return -1, err
				}
				return -1, err
			}
			return -1, err
		}
		return -1, err
	}
	return -1, errors.New("No Comment Created")
}

func GetPostsForUser(dbConn *sql.DB, user_id int) []*Post {

	var posts []*Post

	sql := `SELECT id, guid, author_id, text 
			FROM POSTS 
			WHERE author_id = $1 AND id NOT IN (
				SELECT distinct(target_id) FROM likes
				UNION
				SELECT distinct(commentable_id) FROM comments
			)
			order by random()`

	for _, row := range db.DataCall(dbConn, sql, user_id) {
		if pid, err := strconv.Atoi(row["id"]); err == nil {
			if uid, err := strconv.Atoi(row["author_id"]); err == nil {
				post := new(Post)
				post.Author = uid
				post.ID = pid
				post.GUID = row["guid"]
				post.Text = row["text"]
				posts = append(posts, post)
			}
		}
	}

	return posts
}

func NewLike(dbConn *sql.DB, post_id, person_id, post_owner_id int) (int, error) {

	// log.Println("Creating new like!")

	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal("create like transaction can't even begin")
	}

	// Params

	guid := uuid.New()
	target_type := "Post"
	notif_type := "Notifications::Liked"

	// SQLs

	sql := "INSERT INTO likes (target_id, author_id, guid, created_at, updated_at, target_type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	if like_id, err := db.RunTxWQnArgsReturningId(tx, sql, post_id, person_id, guid, time.Now(), time.Now(), target_type); err == nil {
		sql = "UPDATE posts SET likes_count = likes_count+1 WHERE posts.type IN ('StatusMessage') AND posts.id = $1"
		if err := db.RunTxWQnArgs(tx, sql, post_id); err == nil {
			sql = "INSERT INTO notifications (target_type,target_id,recipient_id,created_at, updated_at, type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
			if notif_id, err := db.RunTxWQnArgsReturningId(tx, sql, target_type, post_id, post_owner_id, time.Now(), time.Now(), notif_type); err == nil {
				sql = "INSERT INTO notification_actors (notification_id,person_id,created_at,updated_at) VALUES ($1, $2, $3, $4)"
				if err := db.RunTxWQnArgs(tx, sql, notif_id, person_id, time.Now(), time.Now()); err == nil {
					sql = "INSERT INTO participations (guid,target_id,target_type,author_id,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
					if err := db.RunTxWQnArgs(tx, sql, guid, post_id, target_type, person_id, time.Now(), time.Now()); err == nil {
						tx.Commit()
						// log.Println("New like created with id", like_id)
						return like_id, err
					}
					return -1, err
				}
				return -1, err
			}
			return -1, err
		}
		return -1, err
	}
	return -1, errors.New("No Like Created")
}

func FollowUser(dbConn *sql.DB, person_id_1, person_id_2, aspect_id int) {

	// log.Println("Creating new follow!")

	ok1, contact_id1 := ContactExists(dbConn, person_id_2, person_id_1)
	ok2, contact_id2 := ContactExists(dbConn, person_id_1, person_id_2)

	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal("create follow transaction can't even begin")
	}

	// Params

	target_type := "Person"
	notif_type := "Notifications::StartedSharing"

	// SQLs

	if ok1 {
		sql := "UPDATE contacts SET sharing = $1, updated_at = $2 WHERE contacts.id = $3"
		db.RunTxWQnArgs(tx, sql, "t", time.Now(), contact_id1)
		if ok2 {
			sql = "UPDATE contacts SET receiving = $1, updated_at = $2 WHERE contacts.id = $3"
			db.RunTxWQnArgs(tx, sql, "t", time.Now(), contact_id2)
			sql = "INSERT INTO aspect_memberships (aspect_id,contact_id,created_at,updated_at) VALUES ($1, $2, $3, $4)"
			db.RunTxWQnArgs(tx, sql, aspect_id, contact_id2, time.Now(), time.Now())
		}
	} else {

		sql := "INSERT INTO contacts (user_id,person_id,created_at,updated_at,receiving) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		contact_id, _ := db.RunTxWQnArgsReturningId(tx, sql, person_id_1, person_id_2, time.Now(), time.Now(), "t")

		sql = "INSERT INTO contacts (user_id,person_id,created_at,updated_at,sharing) VALUES ($1, $2, $3, $4, $5)"
		db.RunTxWQnArgs(tx, sql, person_id_2, person_id_1, time.Now(), time.Now(), "t")

		sql = "INSERT INTO aspect_memberships (aspect_id,contact_id,created_at,updated_at) VALUES ($1, $2, $3, $4)"
		db.RunTxWQnArgs(tx, sql, aspect_id, contact_id, time.Now(), time.Now())
	}

	sql := "INSERT INTO notifications (target_type,target_id,recipient_id,created_at,updated_at,type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	notif_id, _ := db.RunTxWQnArgsReturningId(tx, sql, target_type, person_id_1, person_id_2, time.Now(), time.Now(), notif_type)

	sql = "INSERT INTO notification_actors (notification_id,person_id,created_at,updated_at) VALUES ($1, $2, $3, $4) RETURNING id"
	db.RunTxWQnArgs(tx, sql, notif_id, person_id_1, time.Now(), time.Now())

	tx.Commit()

}

func ContactExists(dbConn *sql.DB, person_id_1, person_id_2 int) (bool, string) {

	sql := "SELECT id FROM contacts WHERE user_id = $1 AND person_id = $2"
	// log.Println("checking", sql, person_id_1, person_id_2)
	res := db.DataCall1(dbConn, sql, person_id_1, person_id_2)
	// log.Println("result of contact exists", res)
	if len(res) > 0 {
		return true, res[0]["id"]
	}
	return false, ""
}

func AspectMembershipExists(dbConn *sql.DB, contact_id, aspect_id int) bool {
	sql := "SELECT id FROM aspect_memberships WHERE aspect_id = $1 AND contact_id = $2 LIMIT 1"
	res := db.DataCall(dbConn, sql, aspect_id, contact_id)
	if len(res) > 0 {
		return true
	}
	return false
}

func NewReshare(dbConn *sql.DB, post Post, person_id int) (int, error) {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatal("create comment transaction can't even begin")
	}

	// Params

	target_type := "Post"
	notif_type := "Notifications::Reshared"

	// SQLs

	sql := "INSERT INTO posts (author_id, public, guid, type, text, created_at, updated_at, root_guid, interacted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"
	if reshare_id, err := db.RunTxWQnArgsReturningId(tx, sql, person_id, "t", uuid.New(), "Reshare", post.Text, time.Now(), time.Now(), post.GUID, time.Now()); err == nil {
		sql = "UPDATE posts SET reshares_count = reshares_count+1 WHERE posts.type IN ('StatusMessage') AND posts.id = $1 "
		if err := db.RunTxWQnArgs(tx, sql, post.ID); err == nil {
			sql = "INSERT INTO participations (guid, target_id, target_type, author_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
			if err := db.RunTxWQnArgs(tx, sql, uuid.New(), reshare_id, target_type, post.Author, time.Now(), time.Now()); err == nil {
				sql = "INSERT INTO participations (guid, target_id, target_type, author_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
				if err := db.RunTxWQnArgs(tx, sql, uuid.New(), post.ID, target_type, person_id, time.Now(), time.Now()); err == nil {
					sql = "INSERT INTO notifications (target_type, target_id, recipient_id, created_at, updated_at, type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
					if notif_id, err := db.RunTxWQnArgsReturningId(tx, sql, target_type, post.ID, post.Author, time.Now(), time.Now(), notif_type); err == nil {
						sql := "INSERT INTO notification_actors (notification_id, person_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id"
						if err := db.RunTxWQnArgs(tx, sql, notif_id, person_id, time.Now(), time.Now()); err == nil {
							tx.Commit()
							return reshare_id, err
						}
						return -1, err
					}
					return -1, err
				}
				return -1, err
			}
			return -1, err
		}
		return -1, err
	}
	return -1, errors.New("No Reshare Created")

}

func NewConversation(dbConn *sql.DB, person_id_1, person_id_2 int) (int, error) {

	// log.Println("person_id_1, person_id_2: ", person_id_1, person_id_2)

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println("create conversation transaction can't even begin")
		return -1, errors.New("create conversation transaction can't begin")
	}

	// Params

	subject := helper.RandomText(helper.RandomNumber(5, 15))
	guid := uuid.New()

	// SQLs

	sql := "INSERT INTO conversations (subject,guid,author_id,created_at,updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	conversation_id, err := db.RunTxWQnArgsReturningId(tx, sql, subject, guid, person_id_1, time.Now(), time.Now())

	if err == nil && conversation_id != -1 {

		log.Println("New conversation created with id", conversation_id)

		sql = "INSERT INTO conversation_visibilities (conversation_id,person_id,created_at,updated_at) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)"
		err = db.RunTxWQnArgs(tx, sql, conversation_id, person_id_1, time.Now(), time.Now(), conversation_id, person_id_2, time.Now(), time.Now())

		// sql = "INSERT INTO conversation_visibilities (conversation_id,person_id,created_at,updated_at) VALUES ($1, $2, $3, $4) RETURNING id"
		// db.RunTxWQnArgs(tx, sql, conversation_id, person_id_2, time.Now(), time.Now())

		if err == nil && conversation_id != -1 {
			tx.Commit()
			NewMessage(dbConn, person_id_1, conversation_id)
			return conversation_id, err
		}
	}

	return -1, err
}

func NewMessage(dbConn *sql.DB, person_id, conversation_id int) (int, error) {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println("create message transaction can't even begin")
		return -1, errors.New("Message Transaction can't begin")
	}

	msgtext := helper.RandomText(helper.RandomNumber(20, 100))

	// SQLs

	sql := "INSERT INTO messages (conversation_id,author_id,guid,text,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	msgid, err := db.RunTxWQnArgsReturningId(tx, sql, conversation_id, person_id, uuid.New(), msgtext, time.Now(), time.Now())

	if err == nil {
		tx.Commit()
	}

	return msgid, err
}

func UpdateConversation(dbConn *sql.DB, conversation_id int) {

	tx, err := dbConn.Begin()
	if err != nil {
		log.Println("UpdateConversation transaction can't even begin")
		return
	}

	sql := "UPDATE conversations SET updated_at = $1 WHERE conversations.id = $2 "
	db.RunTxWQnArgs(tx, sql, time.Now(), conversation_id)

	// sql = "UPDATE conversation_visibilities SET unread = $1, updated_at = $2 WHERE conversation_visibilities.id = $3"
	// db.RunTxWQnArgs(tx, sql, 1, time.Now(), conversation_visibilities_id)

	tx.Commit()
}

func GetAllUsersWithAspects(dbConn *sql.DB) []*User {

	var users []*User

	sql := `SELECT user_id, person_id, string_agg(aspect_id::text, ',') as aspects
			FROM (
				SELECT users.id as user_id, people.id as person_id, aspects.id as aspect_id
				FROM users JOIN people ON users.id = people.owner_id JOIN aspects ON aspects.user_id = users.id
			) tab
			GROUP BY user_id, person_id
			ORDER BY random()`

	res := db.DataCall(dbConn, sql)

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

func GetAllUsersWithAspectsExcept(dbConn *sql.DB, column, table string) []*User {

	var users []*User

	sql := fmt.Sprintf(`SELECT user_id, person_id, string_agg(aspect_id::text, ',') as aspects
			FROM (
				SELECT users.id as user_id, people.id as person_id, aspects.id as aspect_id
				FROM users JOIN people ON users.id = people.owner_id JOIN aspects ON aspects.user_id = users.id
			) tab
			WHERE user_id NOT IN (
				SELECT DISTINCT %s FROM %s
			)
			GROUP BY user_id, person_id
			ORDER BY random()`, column, table)
	res := db.DataCall(dbConn, sql)

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

func GetFriendsOfUser(dbConn *sql.DB, user_id int) []*User {

	var users []*User

	sql := `SELECT users.id as user_id, people.id as person_id, string_agg(aspects.id::text, ',') as aspects, contacts.id as contact_id, am.aspect_id as contact_aspect
			FROM contacts 
			JOIN aspect_memberships am on contacts.id = am.contact_id
			JOIN people on people.id = contacts.user_id
			JOIN users on users.id = people.owner_id
			JOIN aspects on aspects.user_id = users.id
			WHERE contacts.user_id = $1 AND contacts.sharing = true
			GROUP BY users.id, people.id, contacts.id, am.aspect_id
		`

	res := db.DataCall(dbConn, sql, user_id)

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

func GetRandomUser(dbConn *sql.DB, except_id string) *User {

	if except_id == "" {
		except_id = "9999999"
	}

	sql := "SELECT users.id as user_id, people.id as person_id FROM users JOIN people ON users.id = people.owner_id WHERE users.id NOT IN ($1) ORDER BY random() LIMIT 1"

	res := db.DataCall(dbConn, sql, except_id)
	row := res[0]

	sql = "SELECT id FROM aspects WHERE user_id = $1"

	aspects := db.DataCall(dbConn, sql, row["user_id"])

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

func BlockUser() {

}

func FetchPublicFeed() {

}
