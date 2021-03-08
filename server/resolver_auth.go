package genesis

import (
	"context"
	"fmt"
	"genesis/crypto"
	"genesis/config"
	"genesis/db"
	"time"

	"github.com/ninja-software/terror"

	"github.com/gofrs/uuid"
	gocheckpasswd "github.com/ninja-software/go-check-passwd"
	"github.com/volatiletech/null"
)

func (r *queryResolver) VerifyResetToken(ctx context.Context, token string, email *null.String) (bool, error) {
	if len(token) <= 6 {
		// SMS reset token?
		if !email.Valid {
			return false, terror.New(ErrEmailInvalid, "")
		}

		user, err := r.UserStore.GetByEmail(email.String)
		if err != nil {
			return false, terror.New(ErrEmailInvalid, "")
		}

		if user.ResetToken != token {
			return false, terror.New(ErrTokenInvalid, "")
		}
		if user.ResetTokenExpires.Before(time.Now()) {
			return false, terror.New(ErrTokenExpired, "")
		}

		// Update verification status
		if !user.MobileVerified {
			user.MobileVerified = true

			_, err := r.UserStore.Update(user)
			if err != nil {
				return false, terror.New(err, "failed to update user mobile verified")
			}
		}

		return true, nil
	}

	user, err := r.UserStore.GetByResetToken(token)
	if err != nil {
		return false, terror.New(ErrTokenInvalid, "")
	}

	// Update verification status
	if !user.Verified {
		user.Verified = true
		_, err := r.UserStore.Update(user)
		if err != nil {
			return false, terror.New(err, "verify user")
		}
	}

	return true, nil
}

func (r *mutationResolver) CheckPasswordStrength(ctx context.Context, password string) error {
	// make sure password met minimum length
	if len(password) < int(config.PasswordMinimumLength) {
		return terror.New(ErrPasswordShort, "")
	}

	// make sure user not using common password
	if gocheckpasswd.IsCommon(password) {
		return terror.New(ErrPasswordCommon, "")
	}

	return nil
}

func (r *mutationResolver) ForgotPassword(ctx context.Context, email string, viaSms *bool) (bool, error) {
	user, err := r.UserStore.GetByEmail(email)
	if err != nil {
		return false, terror.New(ErrEmailInvalid, "")
	}

	if !user.Email.Valid {
		return false, terror.New(ErrEmailInvalid, "user has no email set")
	}

	sendViaSMS := viaSms != nil && *viaSms == true

	if sendViaSMS && !user.MobilePhone.Valid {
		return false, terror.New(ErrMobileNotSet, "")
	}

	// Generate new verify token
	if sendViaSMS {
		user.ResetToken = GenerateAlphanumericCode(6)
		user.ResetTokenExpires = time.Now().Add(time.Hour * time.Duration(1))
	} else {
		user.ResetToken = uuid.Must(uuid.NewV4()).String()
		user.ResetTokenExpires = time.Now().AddDate(0, 0, r.Config.UserAuth.ResetTokenExpiryDays)
	}
	_, err = r.UserStore.Update(user)
	if err != nil {
		return false, terror.New(err, "update user")
	}

	// Send SMS
	if sendViaSMS {
		err := r.SmsMessenger.SendToken(user.MobilePhone.String, user.ResetToken)
		if err != nil {
			return false, terror.New(err, "Failed to send SMS")
		}
		return true, nil
	}

	// Create email
	sender := r.Config.Email.Sender
	subject := "Genesis - Forgot Password"

	name := "User"
	if user.FirstName.Valid {
		name = user.FirstName.String
	}
	if user.LastName.Valid {
		name += " " + user.LastName.String
	}

	message := r.Mailer.NewMessage(sender, subject, "", user.Email.String)
	message.SetTemplate("portal_forgot_password")
	message.AddVariable("name", name)
	message.AddVariable("magic_link", fmt.Sprintf("%s/reset/%s", r.Config.API.AdminHost, user.ResetToken))

	// Send Email
	emailCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = r.Mailer.Send(emailCtx, message)
	if err != nil {
		return false, terror.New(err, "failed to send email")
	}

	return true, nil
}

func (r *mutationResolver) ResendEmailVerification(ctx context.Context, email string) (bool, error) {
	user, err := r.UserStore.GetByEmail(email)
	if err != nil {
		return false, terror.New(ErrEmailInvalid, "")
	}

	if !user.Email.Valid {
		return false, terror.New(ErrEmailInvalid, "user has no email set")
	}

	// Generate new verify token
	user.VerifyToken = uuid.Must(uuid.NewV4()).String()
	user.VerifyTokenExpires = time.Now().AddDate(0, 0, r.Config.UserAuth.TokenExpiryDays)
	_, err = r.UserStore.Update(user)
	if err != nil {
		return false, terror.New(err, "update user")
	}

	// Create email
	sender := r.Config.Email.Sender
	subject := "Genesis - Verify Email"

	name := "User"
	if user.FirstName.Valid {
		name = user.FirstName.String
	}
	if user.LastName.Valid {
		name += " " + user.LastName.String
	}

	message := r.Mailer.NewMessage(sender, subject, "", user.Email.String)
	message.SetTemplate("portal_confirm_email")
	message.AddVariable("name", name)
	message.AddVariable("email", email)
	message.AddVariable("new_account", true)
	message.AddVariable("magic_link", fmt.Sprintf("%s/verify/%s", r.Config.API.AdminHost, user.VerifyToken))

	// Send Email
	emailCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = r.Mailer.Send(emailCtx, message)
	if err != nil {
		return false, terror.New(err, "send email")
	}

	return true, nil
}

func (r *mutationResolver) ResetPassword(ctx context.Context, token string, password string, email *null.String) (bool, error) {
	var user *db.User
	var err error

	if len(token) < 6 {
		// SMS reset token?
		if !email.Valid {
			return false, terror.New(err, "Email is needed for SMS reset")
		}

		user, err = r.UserStore.GetByEmail(email.String)
		if err != nil {
			return false, terror.New(ErrEmailInvalid, "")
		}

		if user.ResetToken != token {
			return false, terror.New(ErrTokenInvalid, "")
		}
		if user.ResetTokenExpires.Before(time.Now()) {
			return false, terror.New(ErrTokenExpired, "")
		}
	} else {
		user, err = r.UserStore.GetByResetToken(token)
		if err != nil {
			return false, terror.New(ErrTokenInvalid, "")
		}
	}

	err = r.CheckPasswordStrength(ctx, password)
	if err != nil {
		return false, terror.New(err, "failed to check password strength")
	}

	// Change Password
	hashed := crypto.HashPassword(password)
	user.PasswordHash = hashed
	user.ResetToken = uuid.Must(uuid.NewV4()).String()
	user.ResetTokenExpires = time.Now()

	_, err = r.UserStore.Update(user)
	if err != nil {
		return false, terror.New(err, "update user")
	}

	// Blacklist all old tokens (JWT)
	err = r.Blacklister.BlacklistAll(user.ID)
	if err != nil {
		return false, terror.New(err, "blacklist old tokens")
	}

	return true, nil
}

// ChangePassword your password
func (r *mutationResolver) ChangePassword(ctx context.Context, oldPassword, password string) (bool, error) {
	u, err := r.Auther.UserFromContext(ctx)
	if err != nil {
		return false, terror.New(terror.ErrParse, "")
	}

	err = r.CheckPasswordStrength(ctx, password)
	if err != nil {
		return false, terror.New(err, "")
	}

	// Validate password
	if u.Email.Valid {
		err = r.Auther.ValidatePassword(ctx, u.Email.String, oldPassword)
		if err != nil {
			return false, terror.New(terror.ErrAuthWrongPassword, "")
		}
	}

	// Change password
	hashed := crypto.HashPassword(password)
	u.PasswordHash = hashed

	_, err = r.UserStore.Update(u)
	if err != nil {
		return false, terror.New(err, "update user")
	}

	return true, nil
}
