package genesis

import (
	"context"
	"errors"
	"fmt"
	"genesis/db"
	"net/http"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofrs/uuid"
	"github.com/matishsiao/goInfo"
	"github.com/ninja-software/terror"
)

// SError is sentry error for JSON marshalling
type SError struct {
	Err             error                  `json:"-"`               // use for backend, original error, not send to web
	ErrStr          string                 `json:"error"`           // use for frontend, holds the error string
	CorrelationID   string                 `json:"correlationID"`   // short error id for user to tell sysadmin for reporting error
	FriendlyMessage string                 `json:"friendlyMessage"` // message for the user, more readable for human
	Tags            map[string]string      `json:"tags"`            // info such as mon error id
	Vars            map[string]interface{} `json:"variables"`       // variable set during the process of code, such as the shallow endpoint function
}

// NewCorrelationID generate short 8 character error id for user to quote easier
func NewCorrelationID() string {
	eidv4, _ := uuid.NewV4()
	eids := strings.Split(eidv4.String(), "-")
	if len(eids) < 1 || len(eids[0]) != 8 {
		// just incase
		eids = []string{"aaaaaaaa"}
	}
	eidc := eids[0]
	eid := fmt.Sprintf("%c%c%c%c-%c%c%c%c", eidc[0], eidc[1], eidc[2], eidc[3], eidc[4], eidc[5], eidc[6], eidc[7])

	return eid
}

// SentrySend send report to sentry
func SentrySend(ctx context.Context, user *db.User, r *http.Request, err error, stackTrace string) {
	friendlyMessage := err.Error()
	correlationID := NewCorrelationID()

	// if happen locally, put a blank to prevent crash
	if r == nil {
		r = &http.Request{
			RemoteAddr: "127.0.0.1",
		}
	}

	var bErr *terror.Error
	if errors.As(err, &bErr) {
		friendlyMessage = bErr.Message
	}
	if len(stackTrace) > 0 {
		friendlyMessage = stackTrace
	}

	level := sentry.LevelInfo
	if strings.Contains(friendlyMessage, "ERROR") {
		level = sentry.LevelError
	}
	if strings.Contains(friendlyMessage, "PANIC") {
		level = sentry.LevelFatal
	}

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.Clear()

		// user info
		if user != nil {
			scope.SetUser(sentry.User{
				ID:        user.ID,
				IPAddress: r.RemoteAddr,
				Email:     user.Email.String,
			})
		} else {
			scope.SetUser(sentry.User{
				IPAddress: r.RemoteAddr,
			})
		}

		// TODO extra details
		// Method Rest:    GET, POST
		// Method GraphQL: QUERY, MUTATION
		// TODO: re-add SdumpRequestData(r)?
		scope.SetRequest(r)

		// operating system info
		osInfo := goInfo.GetInfo()
		scope.SetExtras(map[string]interface{}{
			"sys.cpus":     osInfo.CPUs,
			"sys.core":     osInfo.Core,
			"sys.goos":     osInfo.GoOS,
			"sys.hostname": osInfo.Hostname,
			"sys.kernel":   osInfo.Kernel,
			"sys.os":       osInfo.OS,
			"sys.platform": osInfo.Platform,
		})

		// error info
		scope.SetExtra("error", err.Error())
		scope.SetExtra("correlation_id", correlationID)
		scope.SetLevel(level)

		// mon info
		if bErr != nil && bErr.Meta != nil {
			scope.SetExtras(map[string]interface{}{
				// "mon.tags": bErr.Tags,
				"mon.vars": bErr.Meta,
			})
		}
	})
	defer sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.Clear()
	})

	// sentry title
	sentry.CaptureMessage(friendlyMessage)

	sentry.Flush(2 * time.Second)
}
