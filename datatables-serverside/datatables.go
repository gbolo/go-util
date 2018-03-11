// Package datatablessrv handles the server side processing of an AJAX request for DataTables
// For details on the parameters and the results, read the datatables documentation at
// https://datatables.net/manual/server-side
package datatablessrv

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Copyright (c) 2017 Escape Velocity, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

// ErrNotDataTablesReq indicates that this is not being requested by Datatables
var ErrNotDataTablesReq = errors.New("Not a DataTables request")

// SortDir is the direction of the sort (ascending/descending)
type SortDir int

const (
	// Asc for ascending sorting
	Asc SortDir = iota
	// Desc for descending sorting
	Desc
)

// OrderInfo tracks the list of columns to sort by and in which direction to sort them.
type OrderInfo struct {
	// ColNum indicates Which column to apply sorting to (zero based index to the Columns data)
	ColNum int
	// Direction tells us which way to sort
	Direction SortDir
}

// ColData tracks all of the columns requested by DataTables
type ColData struct {
	// columns[i][name] Column's name, as defined by columns.name.
	Name string
	// columns[i][data] Column's data source, as defined by columns.data.
	// It is poss
	Data string
	// columns[i][searchable]	boolean	Flag to indicate if this column is searchable (true) or not (false).
	// This is controlled by columns.searchable.
	Searchable bool
	// columns[i][orderable] Flag to indicate if this column is orderable (true) or not (false).
	// This is controlled by columns.orderable.
	Orderable bool
	// columns[i][search][value] Search value to apply to this specific column.
	Searchval string
	// columns[i][search][regex]
	// Flag to indicate if the search term for this column should be treated as regular expression (true) or not (false).
	// As with global search, normally server-side processing scripts will not perform regular expression searching
	// for performance reasons on large data sets, but it is technically possible and at the discretion of your script.
	UseRegex bool
}

// DataTablesInfo represents all of the information that was requested by DataTables
type DataTablesInfo struct {
	// HasFilter Indicates there is a filter on the data to apply.  It is used to optimize generating
	// the query filters
	HasFilter bool
	// Draw counter. This is used by DataTables to ensure that the Ajax returns
	// from server-side processing requests are drawn in sequence by DataTables
	// (Ajax requests are asynchronous and thus can return out of sequence).
	// This is used as part of the draw return parameter (see below).
	Draw int
	// Start is the paging first record indicator.
	// This is the start point in the current data set (0 index based - i.e. 0 is the first record).
	Start int
	// Length is the number of records that the table can display in the current draw.
	// It is expected that the number of records returned will be equal to this number, unless the server has fewer records to return.
	//  Note that this can be -1 to indicate that all records should be returned (although that negates any benefits of server-side processing!)
	Length int
	// Searchval holds the global search value. To be applied to all columns which have searchable as true.
	Searchval string
	// UseRegex is true if the global filter should be treated as a regular expression for advanced searching.
	//  Note that normally server-side processing scripts will not perform regular expression
	//  searching for performance reasons on large data sets, but it is technically possible and at the discretion of your script.
	UseRegex bool
	// Order provides information about what columns are to be ordered in the results and which direction
	Order []OrderInfo
	// Columns provides a mapping of what fields are to be searched
	Columns []ColData
}

// MySQLFilter generates the filter for a mySQL query based on the request and a map of the strings
// to the database entries.  Note if the field is searchable and we don't have a map, this generates
// an error
// It is assumed that there is a fulltext index on fields in very large tables in order to optimize
// performance of this generated query
// We have several things that we can generate here.
// For the global string, (in di.Searchval) with no Regex we would generate something like
//  MATCH (field1,field2,field3) AGAINST('searchval')
// If an individual entry has a searchval with no Regex we generate
//  MATCH (field1) AGAINST('searchval')
// If we have a global string with a regex, we generate
//   field1 REGEX 'searchval' OR field2 REGEX 'searchval' OR field3 REGEX 'searchval'
// Likewise for individual entries with a searchval we generate
//   field1 REGEX 'searchval'
// Note: If the searchval doesn't actually contain wildcard values (^$.*+|(){}[]?) then the search value bit is actually
// cleared by ParseDatatablesRequest so that we never actually see it
// NOTE: It is the responsibility of the caller to put the " WHERE " in front of the string when it
// is non-null.  This allows the filter to be used in other situations or where it may need to be part
// of a more complex logical operation
// NOTE: We assume that the Searchval strings have all been escaped and quoted so that we can put in the string
// with no potential SQL injection
func (di *DataTablesInfo) MySQLFilter(SQLFieldMap map[string]string) (res string, err error) {
	// In the case where there is no filter at all, we can just return quickly
	if !di.HasFilter {
		return
	}
	extra := ""
	for _, colData := range di.Columns {
		if colData.Searchable {
			// Map the external name to the actual field in the database
			sqlName, isFound := SQLFieldMap[colData.Data]
			if !isFound {
				err = fmt.Errorf("Column Data Name %v not found in SQL FieldMap", colData.Data)
				return
			}
			// If we have a global search val, generate a match against the global value for this field
			if di.Searchval != "" {
				// For wildcards we have to generate a REGEXP request
				if di.UseRegex {
					res += extra + sqlName + " REGEX " + di.Searchval
					extra = " OR "
				} else {
					// In the special case where we have a top level non wild card search value we want
					// to gang all the fields together into a single match string
					res += extra + "MATCH(" + sqlName + ") AGAINST(" + di.Searchval + ")"
					extra = " OR "
				}
			}
			// See if we have a search value specific for this individual element
			if colData.Searchval != "" {
				if colData.UseRegex {
					res += extra + sqlName + " REGEX " + colData.Searchval
				} else {
					res += extra + "MATCH(" + sqlName + ") AGAINST(" + colData.Searchval + ")"
				}
				extra = " OR "
			}
		}
	}
	return
}

// MySQLOrderby generates the order by clause for a mySQL query based on the request and a map of the strings
// to the database entries.  Note if the field is orderable and we don't have a map, this generates
// an error.  The string IS prefixed by a space so that you can just append it.
func (di *DataTablesInfo) MySQLOrderby(SQLFieldMap map[string]string) (res string, err error) {
	extra := " ORDER BY "
	// Go through the list of requested items to order
	for _, orderItem := range di.Order {
		// Make sure that the column is in range
		if orderItem.ColNum >= len(di.Columns) {
			err = fmt.Errorf("Datatables Request order column %v out of range %v of columns", orderItem.ColNum, len(di.Columns))
			return
		}
		// Get the data for that column and figure out if the name is one of the fields that we
		// allow in the table
		colData := di.Columns[orderItem.ColNum]
		// Map the external name to the actual field in the database
		sqlName, isFound := SQLFieldMap[colData.Data]
		if !isFound {
			err = fmt.Errorf("Invalid datatables request column name %v", colData.Data)
			return
		}
		// Make sure we can actually order on the column (in theory this will never happen)
		if !colData.Orderable {
			err = fmt.Errorf("Datatables requested ordering on non-orderable column %v", colData.Data)
			return
		}
		// We have the column in the database, add it to the order by query that we are generating
		// The first time we have " ORDER BY " in the extra string, subsequent times we get a simple ","
		// which allows us to build up the string without backtracking to remove characters
		res += extra + sqlName
		if orderItem.Direction == Desc {
			res += " DESC"
		}
		extra = ","
	}
	// If for some reason we got to the end with no columns, then we give them the order by the first item
	if res == "" {
		res = extra + "1"
	}
	return
}

// parseParts takes the split out parts of the field string, verifies that they are
// syntactically valid and then parses them out
//  for example columns[i][search][regex] would come in as
//       field:  'columns[2][search][regex]'
//       nameparts[0]  'columns'
//       nameparts[1]  '2]'
//       nameparts[2]  'search]'
//       nameparts[3]  'regex]'
func parseParts(field string, nameparts []string) (index int, elem1 string, elem2 string, err error) {
	defaultErr := fmt.Errorf("Invalid order[] element %v", field)
	numRegex, err := regexp.Compile("^[0-9]+]$")
	if err != nil {
		return
	}
	elemRegex, err := regexp.Compile("^[a-z]+]$")
	if err != nil {
		return
	}
	if len(nameparts) != 3 && len(nameparts) != 4 {
		err = defaultErr
		return
	}
	// Make sure it is a number followed by the closing ]
	if !numRegex.MatchString(nameparts[1]) {
		err = defaultErr
		return
	}
	// And parse it as a number to make sure
	numstr := strings.TrimSuffix(nameparts[1], "]")
	index, err = strconv.Atoi(numstr)
	if err != nil {
		return
	}
	// Check that the next index is a name token followed by a ]
	if !elemRegex.MatchString(nameparts[2]) {
		err = defaultErr
		return
	}
	// Strip off the trailing ]
	elem1 = strings.TrimSuffix(nameparts[2], "]")
	// If we had a third element, check to make sure it is also close by a ]
	if len(nameparts) == 4 {
		if !elemRegex.MatchString(nameparts[3]) {
			err = defaultErr
			return
		}
		// And trim off the ]
		elem2 = strings.TrimSuffix(nameparts[3], "]")
	}
	// Let's sanity check and make sure they aren't returning an index that is way out of range.
	// We shall assume that no more than 200 columns are being returned
	if index > 200 || index < 0 {
		err = defaultErr
	}
	return
}

// ParseDatatablesRequest checks the HTTP request to see if it corresponds
// to a datatables AJAX data request and parses the request data into
// the DataTablesInfo structure.
//
// This structure can be used by MySQLFilter and MySQLOrderby to generate a
// MySQL query to run against a database.
//
// For example assuming you are going to fill in a response structure to DataTables
// such as:
//
//   type QueryResponse struct {
//       DateAdded   time.Time
//       Status      string
//       Email       struct {
//           Name      string
//           Email     string
//       }
//   }
//   var emailQueueFields = map[string]string{
//       "DateAdded":          "t1.dateadded",
//       "Status":             "t1.status",
//       "Email.Name":         "t2.Name",
//       "Email.Email":        "t2.Email",
//   }
//
//   const baseQuery = `
//       SELECT t1.dateadded
//             ,t1.status
//             ,t2.Name
//             ,t2.Email
//       FROM infotable t1
//       LEFT JOIN usertable t2
//         ON t1.key = t2.key`
//
//       // See if we have a where clause to add to the base query
//       query := baseQuery
//       sqlPart, err := di.MySQLFilter(sqlFields)
//       // If we did have a where filter, append it.  Note that it doesn't put the " WHERE "
//       // in front because we might be doing a boolean operation.
//       if sqlPart != "" {
//           query += " WHERE " + sqlPart
//       }
//       sqlPart, err = di.MySQLOrderby(sqlFields)
//       query += sqlPart
//
// At that point you have a query that you can send straight to mySQL
//
func ParseDatatablesRequest(r *http.Request) (res *DataTablesInfo, err error) {
	var index int
	var elem string
	var elem2 string
	foundDraw := false
	res = &DataTablesInfo{}
	// Let the request parse the post values into the r.Form structure
	err = r.ParseForm()
	if err != nil {
		return
	}
	for field, value := range r.Form {
		// Remember that HTML sends us an array of values, but for datatables we only have one entry so we
		// we can shortcut and take the first element (which will be the only element) of the field.
		val0 := value[0]
		// Split out on the [ into pieces so we can see what the name is.  Note that we will have another
		// routine split out remainder of the string.
		nameparts := strings.Split(field, "[")
		switch nameparts[0] {
		case "draw":
			foundDraw = true
			res.Draw, err = strconv.Atoi(val0)
		case "start":
			res.Start, err = strconv.Atoi(val0)
		case "length":
			res.Length, err = strconv.Atoi(val0)
		case "search":
			if len(nameparts) != 2 {
				err = fmt.Errorf("Invalid search[] element %v", field)
			} else if nameparts[1] == "value]" {
				res.Searchval = val0
			} else if nameparts[1] == "regex]" {
				res.UseRegex = (val0 == "true")
			} else {
				err = fmt.Errorf("Invalid search[] element %v", field)
			}
		case "order":
			index, elem, _, err = parseParts(field, nameparts)
			if err == nil {
				// Make sure there is a spot to store this one.  Note that we may see
				// order[3][column] before we see order[0][dir]
				for len(res.Order) <= index {
					res.Order = append(res.Order, OrderInfo{})
				}
				switch elem {
				case "column":
					res.Order[index].ColNum, err = strconv.Atoi(val0)
				case "dir":
					res.Order[index].Direction = Asc
					if val0 == "desc" {
						res.Order[index].Direction = Desc
					}
				}
			}
		case "columns":
			index, elem, elem2, err = parseParts(field, nameparts)
			// First make sure we have a valid column number to work against
			if err == nil {
				// Fill up the slice to get to the spot where it is going
				// because the columns may come out of order.. I.e. we may see
				// columns[4][search][value] before we see columns[0][data]
				for len(res.Columns) <= index {
					res.Columns = append(res.Columns, ColData{})
				}
			}
			// Now fill in the field in the column slice
			switch elem {
			case "data":
				res.Columns[index].Data = val0
			case "name":
				res.Columns[index].Name = val0
			case "searchable":
				res.Columns[index].Searchable = (val0 != "false")
			case "orderable":
				res.Columns[index].Orderable = (val0 != "false")
			case "search":
				switch elem2 {
				case "value":
					res.Columns[index].Searchval = val0
				case "regex":
					res.Columns[index].UseRegex = (val0 != "false")
				}
			}
		}
		// Any errors along the way and we get out.
		if err != nil {
			return
		}
	}
	// If no Draw was specified in the request, then this isn't a datatables request and we can safely ignore it
	if !foundDraw {
		res = nil
		err = errors.New("Not a DataTables request")
	} else {
		// We have a valid datatables request.  See if we actually have any filtering
		res.HasFilter = false
		// Check the global search value to see if it has anything on it
		if res.Searchval != "" {
			// We do have a filter so note that for later
			res.HasFilter = true
			// If they ask for a regex but don't use any regular expressions, then turn off regex for efficiency
			if res.UseRegex && !strings.ContainsAny(res.Searchval, "^$.*+|[]?") {
				res.UseRegex = false
			}
			// Escape the single quotes and any escape characters and then quote the string
			res.Searchval = strings.Replace(res.Searchval, "\\", "\\\\", -1)
			res.Searchval = "'" + strings.Replace(res.Searchval, "'", "\\'", -1) + "'"
		}
		// Now we check all of the columns to see if they have search expressions
		for _, colData := range res.Columns {
			if colData.Searchval != "" {
				// We have a search expression so we remember we have a filter
				res.HasFilter = true
				// CHeck for any regular expression characters and turn off regex if not
				if colData.UseRegex && !strings.ContainsAny(colData.Searchval, "[]^$.*?+") {
					colData.UseRegex = false
				}
				// Escape the single quotes and any escape characters and then quote the string
				colData.Searchval = strings.Replace(colData.Searchval, "\\", "\\\\", -1)
				colData.Searchval = "'" + strings.Replace(colData.Searchval, "'", "\\'", -1) + "'"
			}
		}
	}
	return
}
