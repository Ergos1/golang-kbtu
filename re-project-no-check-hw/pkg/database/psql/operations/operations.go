package operations

import (
	"database/sql"
	"fmt"

	"example.com/pkg/utils"
	"github.com/jmoiron/sqlx"
)

// func Select(db sqlx.DB, table string, arg interface{}) ([]*interface{}, error) {
// 	fileds := utils.DBFields(arg)
// 	csv := utils.FieldsCSV(fileds)
// 	sql := "SELECT " + csv + " FROM " + table
// 	objs := make([]*interface{}, 0)
// 	err := db.Select(&objs, sql)
// 	return objs, err
// }

// func SelectByID(db sqlx.DB, table string, arg interface{}) (*interface{}, error) {
// 	fileds := utils.DBFields(arg)
// 	csv := utils.FieldsCSV(fileds)
// 	id := utils.DBValueByField(arg, "id")
// 	sql := "SELECT " + csv + " FROM " + table + " WHERE id="+id
// 	obj := struct{}
// 	err := db.Get(obj, sql)
// 	return obj, err
// }

func Insert(db *sqlx.DB, table string, arg interface{}) (sql.Result, error) {
	fileds := utils.DBFields(arg)
	csv := utils.FieldsCSVIgnoreId(fileds)
	csvc := utils.FieldsCSVColonsIgnoreId(fileds)
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, csv, csvc)
	return db.NamedExec(sql, arg)
}


func Update(db *sqlx.DB, table string, arg interface{}) (sql.Result, error) {
	fileds := utils.DBFields(arg)
	values := utils.DBValues(arg)
	csv := utils.SetCsvIgnoreId(fileds, values)
	id := utils.DBValueByField(arg, "id")
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE id=%s", table, csv, id)
	return db.Exec(sql)
}

func Delete(db *sqlx.DB, table string, id uint) (sql.Result, error) {
	sql := fmt.Sprintf("DElETE FROM %s WHERE id=%d", table, id)
	return db.Exec(sql)
}
