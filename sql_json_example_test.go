package unusual_generics_test

import (
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/xakep666/unusual_generics"
)

func ExampleSQLJSON_Scan() {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	mock.ExpectQuery("SELECT").WillReturnRows(
		sqlmock.NewRows([]string{"x"}).
			AddRow(`{"ints": [1, 2, 3], "str": "xxx"}`).
			AddRow(`{"ints": [4, 5, 6], "str": "yyy"}`),
	)

	rows, err := db.Query("SELECT x FROM y")
	if err != nil {
		panic(err)
	}

	type JSONDesc struct {
		Ints []int  `json:"ints"`
		Str  string `json:"str"`
	}

	for rows.Next() {
		var row unusual_generics.SQLJSON[JSONDesc]
		if err := rows.Scan(&row); err != nil {
			panic(err)
		}

		fmt.Println(row.X.Ints, row.X.Str)
	}

	// Output:
	// [1 2 3] xxx
	// [4 5 6] yyy
}

func ExampleSQLJSON_Value() {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	mock.ExpectExec(`INSERT INTO x VALUES .*`).
		WithArgs(`{"ints":[1,2,3],"str":"xxx"}`, `{"ints":[4,5,6],"str":"yyy"}`).
		WillReturnResult(sqlmock.NewResult(2, 2))

	type JSONDesc struct {
		Ints []int  `json:"ints"`
		Str  string `json:"str"`
	}

	res, err := db.Exec(`INSERT INTO x VALUES ($1), ($2)`,
		unusual_generics.SQLJSONOf(JSONDesc{Ints: []int{1, 2, 3}, Str: "xxx"}),
		unusual_generics.SQLJSONOf(JSONDesc{Ints: []int{4, 5, 6}, Str: "yyy"}),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.LastInsertId())
	fmt.Println(res.RowsAffected())

	// Output:
	// 2 <nil>
	// 2 <nil>
}
