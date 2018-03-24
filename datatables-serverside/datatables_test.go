package datatablessrv

import (
	"testing"
)

func TestMySQLGenerateQueryFromColNames(t *testing.T) {


	testData := DataTablesInfo{
		Draw: 1,
		HasFilter: false,
		Start: 0,
		Length: 10,
		Columns: []ColData{
			ColData{
				Name: "first_name FROM test'; DROP TABLE test; SELECT first_name",
				Orderable: true,
				Searchable: true,
				Searchval: "george",
			},
			ColData{
				Name: "last_name",
				Orderable: true,
				Searchable: true,
				Searchval: "bolo",
			},
			ColData{
				Name: "middle_name",
				Orderable: true,
				Searchable: true,
			},
			ColData{
				Name: "age",
				Orderable: true,
				Searchable: true,
			},
		},
		Order: []OrderInfo{
			OrderInfo{
				ColNum: 0,
				Direction: Asc,
			},
			OrderInfo{
				ColNum: 3,
				Direction: Desc,
			},
		},
	}

	query, err := testData.MySQLGenerateQueryFromColNames("person")
	if err != nil {
		t.Error("Error:",err)
		return
	}


	t.Log("Query:", query)

}
