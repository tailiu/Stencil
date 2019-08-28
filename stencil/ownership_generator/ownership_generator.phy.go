package ownership_generator

import (
	"fmt"
	"log"
	"stencil/db"
	"stencil/migrate"
	"stencil/qr"
	"strings"
)

func GetUsersForApp(app_id string) []string {
	dbConn := db.GetDBConn(db.STENCIL_DB)
	var users []string
	query := "SELECT user_id FROM user_table WHERE app_id = $1 AND user_id NOT IN (SELECT user_id FROM user_table WHERE app_id != $2)" //+ "AND user_id NOT IN (SELECT user_id FROM owned_data)"
	if result, err := db.DataCall(dbConn, query, app_id, app_id); err == nil {
		for _, row := range result {
			users = append(users, fmt.Sprint(row["user_id"]))
		}
	} else {
		log.Fatal(err)
	}
	// fmt.Println(users)
	return users
}

func GenOwnership(users []string, srcApp, srcAppID, dstApp, dstAppID string, thead_id int) {
	for _, uid := range users {
		log.Println("started ownership gen uid:", uid, "in thread", thead_id)
		mWorker := migrate.CreateMigrationWorker(uid, srcApp, srcAppID, dstApp, dstAppID, nil, "")
		if err := TraverseDAG(&mWorker, mWorker.GetRoot()); err != nil {
			log.Println("ERROR in GENOWN for uid:", uid, "in thread", thead_id)
		}
		log.Println("finished ownership gen uid:", uid, "in thread", thead_id)
	}
	log.Println("finished ownership gen for all users in thread", thead_id)
}

func TraverseDAG(mWorker *migrate.MigrationWorker, node *migrate.DependencyNode) error {
	for _, dep := range mWorker.SrcAppConfig.GetSubDependencies(node.Tag.Name) {
		if child, err := mWorker.SrcAppConfig.GetTag(dep.Tag); err == nil {
			qs := mWorker.SrcAppConfig.GetTagQS(child)
			mWorker.ResolveDependencyConditions(node, dep, child, qs)
			qs.WhereMFlag(qr.EXISTS, "0", mWorker.SrcAppConfig.AppID)
			sql := qs.GenSQL()
			if result, err := db.DataCall(mWorker.DBConn, sql); err == nil {
				for _, data := range result {
					for col, val := range data {
						if strings.Contains(col, "pk.") && val != nil {
							if success := db.AddOwnedData(mWorker.UserID(), fmt.Sprint(val), mWorker.DBConn); success {
								// fmt.Println(mWorker.UserID(), fmt.Sprint(val))
							}
						}
					}
				}
			} else {
				return err
			}
		}
	}
	return nil
}
