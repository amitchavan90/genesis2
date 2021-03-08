package store

import (
	"context"
	"database/sql"
	"fmt"
	"genesis/db"

	"github.com/ninja-software/terror"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/gofrs/uuid"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/jmoiron/sqlx"
)

// NewLoyaltyStore handle loyalty methods
func NewLoyaltyStore(conn *sqlx.DB) *Loyalty {
	ts := &Loyalty{conn}
	return ts
}

// Loyalty for persistence
type Loyalty struct {
	Conn *sqlx.DB
}

// GetByProductID get loyalty activity by product id
func (s *Loyalty) GetByProductID(id uuid.UUID) (*db.UserLoyaltyActivity, error) {
	return db.UserLoyaltyActivities(
		db.UserLoyaltyActivityWhere.ProductID.EQ(null.StringFrom(id.String())),
	).One(s.Conn)
}

// Insert a user loyalty activity
func (s *Loyalty) Insert(record *db.UserLoyaltyActivity, txes ...*sql.Tx) (*db.UserLoyaltyActivity, error) {
	err := handleTransactions(s.Conn, func(tx *sql.Tx) error {
		return record.Insert(tx, boil.Infer())
	}, txes...)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return record, nil
}

// UserActivity returns all loyalty activity for the user
func (s *Loyalty) UserActivity(userID uuid.UUID) (db.UserLoyaltyActivitySlice, error) {
	return db.UserLoyaltyActivities(
		db.UserLoyaltyActivityWhere.UserID.EQ(userID.String()),
		qm.OrderBy(db.UserLoyaltyActivityColumns.CreatedAt+" DESC"),
	).All(s.Conn)
}

// TotalPointsResult - custom struct for grabbing total points
type TotalPointsResult struct {
	AmountSum int `boil:"amount_sum"`
}

// TotalPoints calculates and returns the current total of loyalty points for a user
func (s *Loyalty) TotalPoints(ctx context.Context, userID uuid.UUID) (int, error) {
	var data TotalPointsResult

	err := queries.Raw(
		fmt.Sprintf(
			`SELECT sum(%s) as "amount_sum" FROM %s WHERE %s='%s'`,
			db.UserLoyaltyActivityColumns.Amount,
			db.TableNames.UserLoyaltyActivities,
			db.UserLoyaltyActivityColumns.UserID,
			userID.String(),
		),
	).Bind(ctx, s.Conn, &data)
	if err != nil {
		return 0, terror.New(err, "")
	}

	return data.AmountSum, nil
}
