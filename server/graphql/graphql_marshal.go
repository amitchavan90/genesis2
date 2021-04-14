package graphql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/null"
)

// MarshalNullString to handle custom type
func MarshalNullString(ns null.String) graphql.Marshaler {
	if !ns.Valid {
		return graphql.Null
	}
	return graphql.MarshalString(ns.String)
}

// UnmarshalNullString to handle custom type
func UnmarshalNullString(v interface{}) (null.String, error) {
	if v == nil {
		return null.String{Valid: false}, nil
	}
	s, err := graphql.UnmarshalString(v)

	if err != nil {
		return null.StringFrom(s), terror.New(err, "")
	}

	return null.StringFrom(s), nil
}

// MarshalNullInt to handle custom type
func MarshalNullInt(ns null.Int) graphql.Marshaler {
	if !ns.Valid {
		return graphql.Null
	}
	return graphql.MarshalInt(ns.Int)
}

// UnmarshalNullInt to handle custom type
func UnmarshalNullInt(v interface{}) (null.Int, error) {
	if v == nil {
		return null.Int{Valid: false}, nil
	}
	s, err := graphql.UnmarshalInt(v)
	if err != nil {
		return null.IntFrom(s), terror.New(err, "")
	}

	return null.IntFrom(s), nil
}

// MarshalNullTime to handle custom type
func MarshalNullTime(ns null.Time) graphql.Marshaler {
	if !ns.Valid {
		return graphql.Null
	}
	return graphql.MarshalTime(ns.Time)
}

// UnmarshalNullTime to handle custom type
func UnmarshalNullTime(v interface{}) (null.Time, error) {
	if v == nil {
		return null.Time{Valid: false}, nil
	}
	s, err := graphql.UnmarshalTime(v)
	if err != nil {
		return null.TimeFrom(s), terror.New(err, "")
	}
	return null.TimeFrom(s), nil
}

// MarshalNullBool to handle custom type
func MarshalNullBool(ns null.Bool) graphql.Marshaler {
	if !ns.Valid {
		return graphql.Null
	}
	return graphql.MarshalBoolean(ns.Bool)
}

// UnmarshalNullBool to handle custom type
func UnmarshalNullBool(v interface{}) (null.Bool, error) {
	if v == nil {
		return null.Bool{Valid: false}, nil
	}
	s, err := graphql.UnmarshalBoolean(v)
	if err != nil {
		return null.BoolFrom(s), terror.New(err, "")
	}
	return null.BoolFrom(s), nil
}
