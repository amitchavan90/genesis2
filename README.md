```txt
   ___ ___ _  _ ___ ___ ___ ___
  / __| __| \| | __/ __|_ _/ __|
 | (_ | _|| .` | _|\__ \| |\__ \
  \___|___|_|\_|___|___/___|___/

```

_Dependencies_

- [go](https://golang.org/)
- [node](https://nodejs.org/en/)
- [postgres](https://www.postgresql.org/)
- [docker](https://docs.docker.com/install/linux/docker-ce/ubuntu/)
- [docker-compose](https://docs.docker.com/compose/install/)
- [solidity](https://github.com/ethereum/solidity)
- [geth](https://geth.ethereum.org)

_Included dependent binaries_

- [sqlboiler](https://github.com/volatiletech/sqlboiler)
- [migrate](https://github.com/golang-migrate/migrate)
- [mockery](https://github.com/vektra/mockery)
- [realize](https://github.com/oxequa/realize)
- [go-ethereum](https://github.com/ethereum/go-ethereum)

Download [solidity](https://github.com/ethereum/solidity/releases) v0.6.x to the \$PATH/bin folder, set execute bit

Note:

Project use on some libraries which does not work on go v1.15

## Development

### Database

```bash
docker run -d -p 5438:5432 \
--name genesis-db \
-e POSTGRES_USER=genesis \
-e POSTGRES_PASSWORD=dev \
-e POSTGRES_DB=genesis \
postgres:11-alpine
```

```bash
docker exec -it genesis-db psql -U genesis
```

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
\q
```

### Geth is removed

Due to requiring hosting large db to work, as lite mote barely get any host.

### Caddy is removed

Due to too much magic. Replace by nginx.

### Update Tools (if required)

```bash
cd server
go generate -tags tools ./tools/...
```

### Web

```bash
cd web
npm install
npm start
```

### Env

Make sure these environment variables are set:

```bash
GENESIS_BLOCKCHAIN_PRIVATEKEYBYTES
GENESIS_EMAIL_APIKEY
GENESIS_SENTRY_DSN
FONTAWESOME_TOKEN
```

(Find in Bitwarden, Mailgun and [Sentry](https://sentry.theninja.life/settings/sentry/projects/genesis-backend/keys/))

### Server

```bash
cd server
../scripts/db-prepare.sh
go generate ./...

go run cmd/platform/main.go db -db_drop
go run cmd/platform/main.go db -db_migrate
go run cmd/platform/main.go db -db_seed
go run cmd/platform/main.go serve
```

### Build Fieldapp

### Windows Dev

1. follow instruction on https://flutter.dev/docs/get-started/install/windows
   a. unpack windows flutter, try c:\src\flutter
2. create a virtual phone using Android Studio AVD
3. run build below using git-bash, this will take awhile
   a. if stuck on gradle, try cancel and run again
   b. may take a while
4. for android phone dev
   a. in settings  
    i. tap System, tapping About phone
   ii. keep tapping build number to enable developer mode
   iii. unplug and plug in usb cable
   iv. on phone, allow permission for pc to debug
   b. in developer options, enable USB debugging
   c. unplug and plug in usb cable

```bash
cd fieldapp
android/build.sh
```

apk located at `build/app/outputs/apk/release`

## Packaging

```bash
./scripts/build-docker.sh
```

## Deployment

### development

Add this in CapRover nginx setting to hardcode WeChat url auth response.

```txt
server {

    location /api/steak/ZwLGFBAL10.txt {
        return 200 '8b8a8b3ba30eb7b687a83bc55e446db4';
    }
```

### server side

```bash
# as root

# install nginx (previously caddy)
apt install nginx

# prod dir
mkdir /usr/share/latitude28/genesis

# after pushed genesis build to server
cd /etc/systemd/system
ln -s /usr/share/latitude28/genesis/init/genesis.service
systemctl enable genesis
systemctl start genesis

# reload systemd service
systemctl daemon-reload

# edit config
cd /usr/share/latitude28/genesis/init
cp genesis-prod.env.sample genesis-prod.env
nano genesis-prod.env


# as postgres user
su postgres -
psql
# change password
 > \password
# create db
 > create database genesis;
# select db
 > \c genesis
# create sql extension as above
 > ...
```

### dev side

```bash
export PROD_HOST="1.2.3.456"
make all
make deploy-prod-full
```

### tools

#### calculate and check hashes

https://passwordsgenerator.net/sha256-hash-generator/
https://www.fileformat.info/tool/hash.htm
https://www.pelock.com/products/hash-calculator

#### doco

codes in server/cmd/lab are used for tests and examples
https://godoc.org/github.com/ethereum/go-ethereum/ethclient
https://goethereumbook.org/smart-contract-compile/

production server has two folders

- /usr/share/latitude28/genesis-online (actual running and alias to genesis dir without -online)
- /usr/share/latitude28/genesis-upload (make deploy)

#### ssl certificate

Create combined certificate, using nginx

```bash
apt install certbot
certbot certonly -d admin.gn.latitude28.cn -d consumer.gn.latitude28.cn -d admin.gn.l28produce.com.au -d consumer.gn.l28produce.com.au
```
