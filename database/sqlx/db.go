package sqlx

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

func InsertMap(params map[string]interface{}) (string, string) {
	var cols, vals []string
	for k, _ := range params {
		cols = append(cols, fmt.Sprintf("%s", k))
		vals = append(vals, fmt.Sprintf(":%s", k))
	}

	return strings.Join(cols, ","), strings.Join(vals, ",")
}

func UpdateMap(params map[string]interface{}, ignore ...string) string {
	var updateFields []string

	for k, _ := range params {
		ok := true
		for _, is := range ignore {
			if is == k {
				ok = false
				break
			}
		}
		if ok {
			updateFields = append(updateFields, fmt.Sprintf("%s=:%s", k, k))
		}
	}
	return strings.Join(updateFields, ",")
}

func Insert(db *sqlx.DB, table string, data map[string]interface{}) (err error) {

	insertKey, insertValue := InsertMap(data)
	sqlStr := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`, table, insertKey, insertValue)
	//logger.WithFields(logger.Fields{"sql": sqlStr}).Debugf("Insert db record")
	_, err = db.NamedExec(sqlStr, data)
	return
}

func InsertTableData(db *sqlx.DB, table string, data map[string]interface{}) (lastInsertID int64, err error) {

	var cols, vals []string
	var x []interface{}
	id := 1
	for k, v := range data {
		cols = append(cols, fmt.Sprintf("%s", k))
		vals = append(vals, fmt.Sprintf("$%d", id))
		x = append(x, v)
		id++
	}
	sqlStr := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s) RETURNING id`, table, strings.Join(cols, ","), strings.Join(vals, ","))
	//fmt.Println(sqlStr)
	err = db.QueryRow(sqlStr, x...).Scan(&lastInsertID)
	return

}

func InsertTableDataContext(ctx context.Context, db *sqlx.DB, table string, data map[string]interface{}) (lastInsertID int64, err error) {

	var cols, vals []string
	var x []interface{}
	id := 1
	for k, v := range data {
		cols = append(cols, fmt.Sprintf("%s", k))
		vals = append(vals, fmt.Sprintf("$%d", id))
		x = append(x, v)
		id++
	}
	sqlStr := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s) RETURNING id`, table, strings.Join(cols, ","), strings.Join(vals, ","))
	//fmt.Println(sqlStr)
	err = db.QueryRowContext(ctx, sqlStr, x...).Scan(&lastInsertID)
	return
}

//func InsertaddData( table string, data map[string]interface{}) (lastInsertID int, err error) {
//
//	var cols, vals []string
//
//	for k, _ := range data {
//		cols = append(cols, fmt.Sprintf("%s", k))
//		vals = append(vals, fmt.Sprintf(":%s", k))
//
//	}
//	sqlStr := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s) RETURNING id`, table, strings.Join(cols, ","), strings.Join(vals, ","))
//	//fmt.Println(sqlStr)
//	err = db.QueryRow(sqlStr).Scan(&lastInsertID)
//	return
//
//}

//func LastInsertId(db *sqlx.DB) (error,int64) {
//	var id int64
//	err = db.QueryRow("SELECT LASTVAL() id").Scan(&id)
//	if err != nil {
//		logrus.Errorf("Unable To Get lastInsertId")
//		return err, 0
//	}
//
//	result.LastInsertId()
//	return err,id
//}

func InsertByTx(tx *sqlx.Tx, table string, data map[string]interface{}) (err error) {
	insertKey, insertValue := InsertMap(data)
	sqlStr := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`, table, insertKey, insertValue)
	//logger.WithFields(logger.Fields{"sql": sqlStr}).Debugf("Insert db record")
	_, err = tx.NamedExec(sqlStr, data)
	return
}

func InsertTableDataByTx(tx *sqlx.Tx, table string, data map[string]interface{}) (lastInsertID int64, err error) {
	var cols, vals []string
	var x []interface{}
	id := 1
	for k, v := range data {
		cols = append(cols, fmt.Sprintf("%s", k))
		vals = append(vals, fmt.Sprintf("$%d", id))
		x = append(x, v)
		id++
	}
	sqlStr := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s) RETURNING id`, table, strings.Join(cols, ","), strings.Join(vals, ","))
	//fmt.Println(sqlStr)
	err = tx.QueryRow(sqlStr, x...).Scan(&lastInsertID)
	return
}

func InsertTableDataTxContext(ctx context.Context, tx *sqlx.Tx, table string, data map[string]interface{}) (lastInsertID int64, err error) {
	var cols, vals []string
	var x []interface{}
	id := 1
	for k, v := range data {
		cols = append(cols, fmt.Sprintf("%s", k))
		vals = append(vals, fmt.Sprintf("$%d", id))
		x = append(x, v)
		id++
	}
	sqlStr := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s) RETURNING id`, table, strings.Join(cols, ","), strings.Join(vals, ","))
	//fmt.Println(sqlStr)
	err = tx.QueryRowContext(ctx, sqlStr, x...).Scan(&lastInsertID)
	return
}

func InsertorUpdateTableDataByTx(tx *sqlx.Tx, table string, data map[string]interface{}) (lastInsertID int64, err error) {
	var cols, vals, set []string
	var x []interface{}
	id := 1
	for k, v := range data {
		cols = append(cols, fmt.Sprintf("%s", k))
		vals = append(vals, fmt.Sprintf("$%d", id))
		x = append(x, v)
		id++
		set = append(set, fmt.Sprintf("%s=excluded.%s", k, k))
	}
	sqlStr := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s) 
	on conflict(id) do update set %s  RETURNING id`,
		table, strings.Join(cols, ","), strings.Join(vals, ","),
		strings.Join(set, " , "))
	err = tx.QueryRow(sqlStr, x...).Scan(&lastInsertID)
	return
}
func InsertorUpdateTableDataTxContext(ctx context.Context, tx *sqlx.Tx, table string, data map[string]interface{}) (lastInsertID int64, err error) {
	var cols, vals, set []string
	var x []interface{}
	id := 1
	for k, v := range data {
		cols = append(cols, fmt.Sprintf("%s", k))
		vals = append(vals, fmt.Sprintf("$%d", id))
		x = append(x, v)
		id++
		set = append(set, fmt.Sprintf("%s=excluded.%s", k, k))
	}
	sqlStr := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s) 
	on conflict(id) do update set %s  RETURNING id`,
		table, strings.Join(cols, ","), strings.Join(vals, ","),
		strings.Join(set, " , "))
	err = tx.QueryRowContext(ctx, sqlStr, x...).Scan(&lastInsertID)
	return
}
func InsertorUpdateTableDataContext(ctx context.Context, db *sqlx.DB, table string, data map[string]interface{}) (lastInsertID int64, err error) {
	var cols, vals, set []string
	var x []interface{}
	id := 1
	for k, v := range data {
		cols = append(cols, fmt.Sprintf("%s", k))
		vals = append(vals, fmt.Sprintf("$%d", id))
		x = append(x, v)
		id++
		set = append(set, fmt.Sprintf("%s=excluded.%s", k, k))
	}
	sqlStr := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s) 
	on conflict(id) do update set %s  RETURNING id`,
		table, strings.Join(cols, ","), strings.Join(vals, ","),
		strings.Join(set, " , "))
	err = db.QueryRowContext(ctx, sqlStr, x...).Scan(&lastInsertID)
	return
}

func UpdateTable(db *sqlx.DB, table string, updateData map[string]interface{}, whereData ...string) (int64, error) {
	set := UpdateMap(updateData, whereData...)
	where := "true"
	for _, v := range whereData {
		where += fmt.Sprintf(" and %s = :%s", v, v)
	}
	sqlStr := fmt.Sprintf(`UPDATE "%s" SET %s WHERE %s`, table, set, where)
	result, err := db.NamedExec(sqlStr, updateData)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func UpdateTableContext(ctx context.Context, db *sqlx.DB, table string, updateData map[string]interface{}, whereData ...string) (int64, error) {
	set := UpdateMap(updateData, whereData...)
	where := "true"
	for _, v := range whereData {
		where += fmt.Sprintf(" and %s = :%s", v, v)
	}
	sqlStr := fmt.Sprintf(`UPDATE "%s" SET %s WHERE %s`, table, set, where)
	result, err := db.NamedExecContext(ctx, sqlStr, updateData)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func UpdateTableByTx(tx *sqlx.Tx, table string, updateData map[string]interface{}, whereData ...string) (int64, error) {
	set := UpdateMap(updateData, whereData...)
	where := "true"
	for _, v := range whereData {
		where += fmt.Sprintf(" and %s = :%s", v, v)
	}
	sqlStr := fmt.Sprintf(`UPDATE "%s" SET %s WHERE %s`, table, set, where)
	result, err := tx.NamedExec(sqlStr, updateData)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func UpdateTableTxContext(ctx context.Context, tx *sqlx.Tx, table string, updateData map[string]interface{}, whereData ...string) (int64, error) {
	set := UpdateMap(updateData, whereData...)
	where := "true"
	for _, v := range whereData {
		where += fmt.Sprintf(" and %s = :%s", v, v)
	}
	sqlStr := fmt.Sprintf(`UPDATE "%s" SET %s WHERE %s`, table, set, where)
	result, err := tx.NamedExecContext(ctx, sqlStr, updateData)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
