package datastore

import "context"

type DB interface {
	InsertOne(context.Context, InsertParams) (int64, error)
	SelectOne(context.Context, SelectParams) error
	UpdateOne(context.Context, UpdateParams) error
	Close()
}

type InsertParams struct {
	Query string
	Vars  []interface{}
}

type SelectParams struct {
	Query   string
	Filters []interface{}
	Result  []interface{}
}

type UpdateParams struct {
	Query string
	Vars  []interface{}
}
