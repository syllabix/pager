# SQL Pager

This package leverages the Paged Details and Response structs defined in the top level pager pacakge for use with SQL. 

#### Currently only supports the MSSQL driver

`go get github.com/syllabix/pager/sqlpager`

The sqlpager can be used in 2 ways, either with a "fluent migrator" syntax, or using manually writing queries.

## Example Usage

```go

//Step 1: import the page pkg dependencies  

import (
    "github.com/syllabix/pager"
    "github.com/syllabix/pager/sqlpager"
)

//Step 2: define the model you wish unmarshall the sql response to

type ExampleItem struct {
    Name string `db:"ItemName"`
    Description string `db:"ItemDescription"`
}

//Example repository/data access function with package scoped access to a sqlx.DB connection pool


var db = database.Get()

func GetItems(pageDetails *page.PageDetails) (*pager.PagedResult, error) {
    //Step 3: Use the sqlpager.New factory method to get a useful instance of a PagedQueryBuilder
    responder := sqlpager.New(db, pageDetails)
    //Step 4: Instantiate a slice of type ExampleItem using Go's make func, and set length 0, and capacity the PageSize (recommended to for performance)
    items := make([]ExampleItem, 0, pageDetails.PageSize)
    //Step 5: Use the SQL Builder methods to assemple your query and then pass in a pointer to the instantiated slice from the previous line
    return responder.Select("*").From("Item").OrderBy("ItemName DESC").ToPagedResponse(&items)
}

//You can also use the PagedQueryBuilder Query method to write a custom query - and then use the requred OrderBy Method to specific order:

func GetItems(pageDetails *page.PageDetails) (*pager.PagedResult, error) {
    responder := sqlpager.New(db, pageDetails)
    items := make([]ExampleItem, 0, pageDetails.PageSize)
    return responder.Query("SELECT ItemName FROM Item WHERE ItemName = ?, "Widget").OrderBy("ItemName ASC").ToPagedResponse(&items)
} 

``` 
