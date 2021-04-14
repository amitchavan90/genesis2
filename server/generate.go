// package genesis holds platform wide source code

//go:generate ../bin/go-bindata -prefix migrations/ -pkg bindata -nocompress -o ./bindata/bindata.go migrations
//go:generate ../bin/sqlboiler ../bin/sqlboiler-psql --wipe --tag db --config ./sqlboiler.toml --output ./db

package genesis
