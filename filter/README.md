### Example Model

```go
type User struct {
	ID   uint   `gorm:"id" filter:"=;>;>=;#;%"`
	Name string `gorm:"name" filter:"*"`
}
```

You can allow all operator by `*`

### Protecting Fields

fields is unsearchable by default.

## Grammar

Goven has a simple syntax that allows for powerful queries.

Fields can be compared using the following operators:

`=`, `!=`, `>=`, `<=`, `<`, `>`, `%` , `#`

The `%` operator allows you to do partial string matching using LIKE.

The `#` operator is IN operator followed by a combination of values disrupted by a non-spaced comma
ex: `name#"(duckhue01,duckhue02)"`

Multiple queries can be combined using `AND`, `OR`.

Together this means you can build up a query like this:

`model_name=iris AND version>=2.0`

More advanced queries can be built up using bracketed expressions:

`(model_name=iris AND version>=2.0) OR artifact_type=TENSORFLOW`
