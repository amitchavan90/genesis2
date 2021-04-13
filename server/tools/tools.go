// +build tools

package genesis

//go:generate go build -o ../../bin/mockery github.com/vektra/mockery/cmd/mockery
//go:generate go build -o ../../bin/caddy github.com/caddyserver/caddy/caddy
//go:generate go build -o ../../bin/migrate -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate
//go:generate go build -o ../../bin/realize github.com/oxequa/realize
//go:generate go build -o ../../bin/sqlboiler github.com/volatiletech/sqlboiler
//go:generate go build -o ../../bin/sqlboiler-psql github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql
//go:generate go build -o ../../bin/go-bindata github.com/kevinburke/go-bindata/go-bindata
//go:generate go build -o ../../bin/abigen github.com/ethereum/go-ethereum/cmd/abigen

import (
	_ "github.com/caddyserver/caddy"
	_ "github.com/ethereum/go-ethereum"
	_ "github.com/golang-migrate/migrate"
	_ "github.com/oxequa/realize"
	_ "github.com/vektra/mockery/cmd/mockery"
	_ "github.com/volatiletech/sqlboiler"
	_ "github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql"
)
