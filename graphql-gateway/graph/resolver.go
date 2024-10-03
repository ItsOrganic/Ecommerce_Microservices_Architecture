package graph

import "gpql-gateway/db"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ProductDB *db.MongoInstance
	UserDB    *db.MongoInstance
	OrderDB   *db.MongoInstance
}
