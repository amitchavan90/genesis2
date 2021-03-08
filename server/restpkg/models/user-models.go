package models

import (
	"time"
)

// User is an object representing the database table.
type User struct {
	ID                 string    `db:"id" boil:"id" json:"id" toml:"id" yaml:"id"`
	OrganisationID     string    `db:"organisation_id" boil:"organisation_id" json:"organisation_id,omitempty" toml:"organisation_id" yaml:"organisation_id,omitempty"`
	RoleID             string    `db:"role_id" boil:"role_id" json:"role_id" toml:"role_id" yaml:"role_id"`
	Email              string    `db:"email" boil:"email" json:"email,omitempty" toml:"email" yaml:"email,omitempty"`
	FirstName          string    `db:"first_name" boil:"first_name" json:"first_name,omitempty" toml:"first_name" yaml:"first_name,omitempty"`
	LastName           string    `db:"last_name" boil:"last_name" json:"last_name,omitempty" toml:"last_name" yaml:"last_name,omitempty"`
	AffiliateOrg       string    `db:"affiliate_org" boil:"affiliate_org" json:"affiliate_org,omitempty" toml:"affiliate_org" yaml:"affiliate_org,omitempty"`
	ReferralCode       string    `db:"referral_code" boil:"referral_code" json:"referral_code,omitempty" toml:"referral_code" yaml:"referral_code,omitempty"`
	MobilePhone        string    `db:"mobile_phone" boil:"mobile_phone" json:"mobile_phone,omitempty" toml:"mobile_phone" yaml:"mobile_phone,omitempty"`
	MobileVerified     bool      `db:"mobile_verified" boil:"mobile_verified" json:"mobile_verified" toml:"mobile_verified" yaml:"mobile_verified"`
	WechatID           string    `db:"wechat_id" boil:"wechat_id" json:"wechat_id,omitempty" toml:"wechat_id" yaml:"wechat_id,omitempty"`
	Verified           bool      `db:"verified" boil:"verified" json:"verified" toml:"verified" yaml:"verified"`
	VerifyToken        string    `db:"verify_token" boil:"verify_token" json:"verify_token" toml:"verify_token" yaml:"verify_token"`
	VerifyTokenExpires time.Time `db:"verify_token_expires" boil:"verify_token_expires" json:"verify_token_expires" toml:"verify_token_expires" yaml:"verify_token_expires"`
	ResetToken         string    `db:"reset_token" boil:"reset_token" json:"reset_token" toml:"reset_token" yaml:"reset_token"`
	ResetTokenExpires  time.Time `db:"reset_token_expires" boil:"reset_token_expires" json:"reset_token_expires" toml:"reset_token_expires" yaml:"reset_token_expires"`
	PasswordHash       string    `db:"password_hash" boil:"password_hash" json:"password_hash" toml:"password_hash" yaml:"password_hash"`
	Archived           bool      `db:"archived" boil:"archived" json:"archived" toml:"archived" yaml:"archived"`
	ArchivedAt         time.Time `db:"archived_at" boil:"archived_at" json:"archived_at,omitempty" toml:"archived_at" yaml:"archived_at,omitempty"`
	UpdatedAt          time.Time `db:"updated_at" boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	CreatedAt          time.Time `db:"created_at" boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
}
