### Example Parser

```go
	query := `(status LIKE home AND (status IN ("To Do", "In Progress", "Closed") AND artifact='')) OR metric > 0.98 OR metric >= 0.98 OR metric != 0.98`

	parser := NewParser(test)
	parser.ParserToArray() 
	parser.ParserToGroups() 
	query, values := parser.ParserToSQL()

	queryResult = "(status LIKE ? AND (status IN(?,?,?) AND artifact = ?)) OR metric > ? OR metric >= ? OR metric != ? "
	valResult = []string{"home", "To Do", "In Progress", "Closed", "", "0.98", "0.98", "0.98"}
```


## Grammar

Goven has a simple syntax that allows for powerful queries.

Fields can be compared using the following operators:

`=`, `!=`, `>=`, `<=`, `<`, `>`, `LIKE` , `IN`, `NOT IN`

Example for : 
LIKE : `name LIKE hh`
IN : `status IN ("To Do", "In Progress", "Closed")`
NOT IN : `status NOT IN ("To Do", "In Progress", "Closed")`

More advanced queries can be built up using bracketed expressions:
`(status LIKE home AND (status IN ("To Do", "In Progress", "Closed") AND artifact='')) OR metric > 0.98 OR metric >= 0.98 OR metric != 0.98`
