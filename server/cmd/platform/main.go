package main

import (
	"crypto/ecdsa"
	"encoding/base64"
	"flag"
	"genesis"
	"genesis/api"
	"genesis/bindata"
	"genesis/blockchain"
	"genesis/config"
	"genesis/scheduler"
	"genesis/seed"
	"genesis/sms"
	"genesis/store"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/getsentry/sentry-go"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/volatiletech/sqlboiler/boil"

	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	migrate_bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/oklog/run"
)

const version = "v1.0.29"
const envPrefix = "GENESIS"

// const set by sed during fieldapp build, please dont remove
const fieldappVersion = "0.1.929"

func main() {
	var (
		rootFlagSet  = flag.NewFlagSet("genesis", flag.ExitOnError)
		dbFlagSet    = flag.NewFlagSet("db", flag.ExitOnError)
		serveFlagSet = flag.NewFlagSet("serve", flag.ExitOnError)

		apiAddr         = serveFlagSet.String("api_addr", "localhost:8081", "host:port to run the API")
		apiBlobbaseurl  = serveFlagSet.String("api_blobbaseurl", "/api/files/", "the sub url for downloading blobs")
		apiConsumerhost = serveFlagSet.String("api_consumerhost", "https://consumer.genesis.staging.theninja.life", "the url for the consumer site")
		apiAdminHost    = serveFlagSet.String("api_adminhost", "https://admin.genesis.staging.theninja.life", "the url for the admin site")

		loadbalancerAdminaddr         = serveFlagSet.String("loadbalancer_adminaddr", ":8080", "host:port to run caddy for admin")
		loadbalancerAdminindexfile    = serveFlagSet.String("loadbalancer_adminindexfile", "admin.html", "index file name for admin")
		loadbalancerAdminrootpath     = serveFlagSet.String("loadbalancer_adminrootpath", "../web/dist", "folder path of admin.html")
		loadbalancerConsumeraddr      = serveFlagSet.String("loadbalancer_consumeraddr", ":8082", "host:port to run caddy for consumer")
		loadbalancerConsumerindexfile = serveFlagSet.String("loadbalancer_consumerindexfile", "consumer.html", "index file name for consumer")
		loadbalancerConsumerrootpath  = serveFlagSet.String("loadbalancer_consumerrootpath", "../web/dist", "folder path of consumer.html")

		userauthJwtsecret                  = serveFlagSet.String("userauth_jwtsecret", "******", "JWT secret")
		userauthTokenexpirydays            = serveFlagSet.Int("userauth_tokenexpirydays", 30, "How many days before the token expires")
		userauthResettokenexpirydays       = serveFlagSet.Int("userauth_resettokenexpirydays", 1, "How many days before the reset token expires")
		userauthBlacklistrefreshhours      = serveFlagSet.Int("userauth_blacklistrefreshhours", 1, "How often should the issued_tokens list be cleared of expired tokens in hours")
		emailDomain                        = serveFlagSet.String("email_domain", "njs.dev", "Mailer domain")
		emailSender                        = serveFlagSet.String("email_sender", "Ninja Software <noreply@njs.dev>", "Default email address to send emails from")
		emailApikey                        = serveFlagSet.String("email_apikey", "SAMPLE KEY", "MailGun API key")
		smsEnabled                         = serveFlagSet.Bool("sms_enabled", false, "Should SMS sending be enabled. If false, will use mock sms, sending to console only.")
		smsSid                             = serveFlagSet.String("sms_sid", "", "Twilio Auth SID. Required if SMS Enabled")
		smsToken                           = serveFlagSet.String("sms_token", "", "Twilio Auth Token. Required if SMS Enabled")
		skipUserverification               = serveFlagSet.Bool("skip_userverification", false, "Should Skip Verification sending be enabled. If true, verification is not required for users.")
		blockchainPrivatekeybytes          = serveFlagSet.String("blockchain_privatekeybytes", "", "The private key JSON in base 64 encoding")
		blockchainPrivatekeypassword       = serveFlagSet.String("blockchain_privatekeypassword", "tiger", "Password to unlock the private key")
		blockchainEthereumhost             = serveFlagSet.String("blockchain_ethereumhost", "http://127.0.0.1:8545", "Host of the geth node")
		blockchainEtherscanhost            = serveFlagSet.String("blockchain_etherscanhost", "https://rinkeby.etherscan.io", "Etherscan host for the transaction links")
		blockchainTestprivatekey           = serveFlagSet.String("blockchain_testprivatekey", "", "Use a single private key (for testing) instead of the private key file + pass")
		blockchainCreateContractOnFirstRun = serveFlagSet.Bool("blockchain_createcontractonfirstrun", false, "Whether or not ")
		blockchainFlushPendingInterval     = serveFlagSet.Int("blockchain_flushpendinginterval", 24, "Interval between flushing pending transactions (in hours)")
		blockchainTxConfirmChkIntervalMin  = serveFlagSet.Int("blockchain_checkcommitinterval", 1, "Interval between checking unconfirmed transactions (in minutes)")
		blockchainEthLowNotifyAmount       = serveFlagSet.Float64("blockchain_eth_low_notify_amount", 0.2, "(in ETH) An notification email will be sent if the account balance is equal or less than this amount")
		blockchainEthLowNotifyEmail        = serveFlagSet.String("blockchain_eth_low_notify_email", "", "Email to send ETH low notification emails")

		fieldappVersion = serveFlagSet.String("fieldapp_version", fieldappVersion, "Current Fieldapp version")

		sentryDSN         = serveFlagSet.String("sentry_dsn", "", "Sends error to remote server. If set, it will send error.")
		sentryServerName  = serveFlagSet.String("sentry_server_name", "dev-pc", "The machine name that this program is running on")
		sentryEnvironment = serveFlagSet.String("sentry_environment", "development", "This program environment")

		databaseUser = rootFlagSet.String("database_user", "genesis", "Postgres username")
		databasePass = rootFlagSet.String("database_pass", "dev", "Postgres password")
		databaseHost = rootFlagSet.String("database_host", "localhost", "Postgres host")
		databasePort = rootFlagSet.String("database_port", "5438", "Postgres port")
		databaseName = rootFlagSet.String("database_name", "genesis", "Postgres database name")

		dbDrop     = dbFlagSet.Bool("db_drop", false, "drop the database")
		dbMigrate  = dbFlagSet.Bool("db_migrate", false, "migrate the database")
		dbSeed     = dbFlagSet.Bool("db_seed", false, "seed the database")
		dbSeedProd = dbFlagSet.Bool("db_seed_prod", false, "seed the production database")
		dbVersion  = dbFlagSet.Bool("db_version", false, "version of the database")
	)
	dbCmd := &ffcli.Command{
		Name:       "db",
		Options:    []ff.Option{ff.WithEnvVarPrefix(envPrefix)},
		ShortUsage: "genesis db [flags]",
		ShortHelp:  "Run database commands.",
		FlagSet:    dbFlagSet,
		Exec: func(_ context.Context, args []string) error {
			if !*dbDrop && !*dbMigrate && !*dbSeed && !*dbSeedProd && !*dbVersion {
				return errors.New("-db_drop, -db_migrate, -db_version or -db_seed is required but not provided ")
			}
			conn := connect(
				*databaseUser,
				*databasePass,
				*databaseHost,
				*databasePort,
				*databaseName,
			)

			if *dbDrop {
				m, err := newMigrateInstance(conn)
				if err != nil {
					return err
				}
				err = m.Drop()
				if err != nil {
					return err
				}
				fmt.Println("database dropped successfully")
			}
			if *dbMigrate {
				m, err := newMigrateInstance(conn)
				if err != nil {
					return err
				}
				err = m.Up()
				if err != nil && err != migrate.ErrNoChange {
					return err
				}
				fmt.Println("database migrated successfully")
			}
			if *dbSeed {
				err := seed.Run(conn)
				if err != nil {
					return err
				}
				fmt.Println("database seeded successfully")
			}
			if *dbSeedProd {
				err := seed.RunProd(conn)
				if err != nil {
					return err
				}
				fmt.Println("database seeded (production) successfully")
			}
			if *dbVersion {
				m, err := newMigrateInstance(conn)
				if err != nil {
					return err
				}
				dbVersion, dirty, err := m.Version()
				if err != nil {
					return err
				}
				fmt.Println("version:", dbVersion, "dirty", dirty)
			}

			return nil
		},
	}

	serveCmd := &ffcli.Command{
		Name:       "serve",
		Options:    []ff.Option{ff.WithEnvVarPrefix(envPrefix)},
		ShortUsage: "genesis serve [flags]",
		ShortHelp:  "Start the server.",
		FlagSet:    serveFlagSet,
		Exec: func(_ context.Context, args []string) error {
			var err error
			// Manual hydration for now, until we refactor everything to not use config structs
			_config := config.NewConfig()
			_config.API = &config.API{
				Addr:         *apiAddr,
				BlobBaseURL:  *apiBlobbaseurl,
				ConsumerHost: *apiConsumerhost,
				AdminHost:    *apiAdminHost,
			}
			_config.LoadBalancer = &config.LoadBalancer{
				AdminAddr:         *loadbalancerAdminaddr,
				AdminIndexFile:    *loadbalancerAdminindexfile,
				AdminRootPath:     *loadbalancerAdminrootpath,
				ConsumerAddr:      *loadbalancerConsumeraddr,
				ConsumerIndexFile: *loadbalancerConsumerindexfile,
				ConsumerRootPath:  *loadbalancerConsumerrootpath,
			}
			_config.Database = &config.Database{
				User: *databaseUser,
				Pass: *databasePass,
				Host: *databaseHost,
				Port: *databasePort,
				Name: *databaseName,
			}
			_config.UserAuth = &config.UserAuth{
				JWTSecret:             *userauthJwtsecret,
				TokenExpiryDays:       *userauthTokenexpirydays,
				ResetTokenExpiryDays:  *userauthResettokenexpirydays,
				BlacklistRefreshHours: *userauthBlacklistrefreshhours,
			}
			_config.Email = &config.Email{
				Domain: *emailDomain,
				Sender: *emailSender,
				APIKey: *emailApikey,
			}
			_config.SMS = &config.SMS{
				Enabled: *smsEnabled,
				SID:     *smsSid,
				Token:   *smsToken,
			}
			_config.Skip = &config.Skip{
				UserVerification: *skipUserverification,
			}
			_config.Blockchain = &config.Blockchain{
				PrivateKeyBytes:             *blockchainPrivatekeybytes,
				PrivateKeyPassword:          *blockchainPrivatekeypassword,
				EthereumHost:                *blockchainEthereumhost,
				EtherscanHost:               *blockchainEtherscanhost,
				TestPrivateKey:              *blockchainTestprivatekey,
				CreateContractOnFirstRun:    *blockchainCreateContractOnFirstRun,
				FlushPendingActionsInterval: *blockchainFlushPendingInterval,
				EthLowNotifyAmount:          *blockchainEthLowNotifyAmount,
				EthLowNotifyEmail:           *blockchainEthLowNotifyEmail,
			}
			_config.Fieldapp = &config.Fieldapp{
				Version: *fieldappVersion,
			}
			_config.Sentry = &config.Sentry{
				DSN:         *sentryDSN,
				ServerName:  *sentryServerName,
				Environment: *sentryEnvironment,
			}

			conn := connect(
				*databaseUser,
				*databasePass,
				*databaseHost,
				*databasePort,
				*databaseName,
			)
			g := &run.Group{}
			ctx, cancel := context.WithCancel(context.Background())

			if len(_config.Sentry.DSN) > 0 {
				err = sentry.Init(sentry.ClientOptions{
					Dsn:         _config.Sentry.DSN,
					ServerName:  _config.Sentry.ServerName,
					Release:     version,
					Environment: _config.Sentry.Environment,
				})
				if err != nil {
					fmt.Printf("Sentry init failed: %v\n", err)
				}
			} else {
				fmt.Println("Sentry init skipped...")
			}

			g.Add(func() error {
				l := genesis.NewLogToStdOut("api", version, false)
				skuStore := store.NewSKUStore(conn)
				orderStore := store.NewOrderStore(conn)
				containerStore := store.NewContainerStore(conn)
				palletStore := store.NewPalletStore(conn)
				cartonStore := store.NewCartonStore(conn)
				productStore := store.NewProductStore(conn)
				organisationStore := store.NewOrganisationStore(conn)
				userStore := store.NewUserStore(conn)
				referralStore := store.NewReferralStore(conn)
				taskStore := store.NewTaskStore(conn)
				userTaskStore := store.NewUserTaskStore(conn)
				roleStore := store.NewRoleStore(conn)
				blobStore := store.NewBlobStore(conn)
				loyaltyStore := store.NewLoyaltyStore(conn)
				contractStore := store.NewContractStore(conn)
				distributorStore := store.NewDistributorStore(conn)
				transactionsStore := store.NewTransactionStore(conn)
				manifestStore := store.NewManifestStore(conn)
				trackActionStore := store.NewTrackActionStore(conn)
				userActivityStore := store.NewUserActivityStore(conn)
				tokenStore := store.NewTokenStore(conn)
				blacklistRefreshHours := _config.UserAuth.BlacklistRefreshHours
				blacklistProvider := genesis.NewBlacklister(l, tokenStore, blacklistRefreshHours)

				var smsMessenger genesis.Messenger = sms.NewMock()
				if _config.SMS.Enabled {
					fmt.Println("Using Real SMS")
					if _config.SMS.SID == "" || _config.SMS.Token == "" {
						panic("Twilio SID and Token environment variables must be set when SMS is enabled.")
					}
					smsMessenger = sms.New(_config.SMS.SID, _config.SMS.Token)
				} else {
					fmt.Println("Using Mock SMS")
				}

				mailer := mailgun.NewMailgun(_config.Email.Domain, _config.Email.APIKey)

				jwtSecret := _config.UserAuth.JWTSecret
				auther := genesis.NewAuther(jwtSecret,
					userStore,
					blacklistProvider,
					tokenStore,
					roleStore,
					_config.UserAuth.TokenExpiryDays,
				)

				// Setup blockchain
				blk := &blockchain.Service{}
				var privateKey *ecdsa.PrivateKey
				if _config.Blockchain.TestPrivateKey != "" || _config.Blockchain.PrivateKeyBytes != "" {
					if _config.Blockchain.TestPrivateKey != "" {
						privateKey, err = crypto.HexToECDSA(_config.Blockchain.TestPrivateKey)
						if err != nil {
							return fmt.Errorf("convert hex to ecdsa: %w", err)
						}
					} else {
						if _config.Blockchain.PrivateKeyBytes == "" {
							return fmt.Errorf("Private key base64 bytes not provided")
						}
						if _config.Blockchain.PrivateKeyPassword == "" {
							return fmt.Errorf("Private key pass not provided")
						}

						b, err := base64.StdEncoding.DecodeString(_config.Blockchain.PrivateKeyBytes)
						if err != nil {
							return fmt.Errorf("decode base64: %w", err)
						}
						key, err := keystore.DecryptKey(b, _config.Blockchain.PrivateKeyPassword)
						if err != nil {
							return fmt.Errorf("decrypt key: %w", err)
						}

						privateKey = key.PrivateKey
					}

					blk, err = blockchain.New(
						ctx,
						_config.Blockchain.EthereumHost,
						privateKey,
						false,
						mailer,
						_config,
					)
					if err != nil {
						return fmt.Errorf("create new blockchain client: %w", err)
					}
				} else {
					blk, err = blockchain.New(
						ctx,
						_config.Blockchain.EthereumHost,
						privateKey,
						true,
						mailer,
						_config,
					)
					if err != nil {
						return fmt.Errorf("create new blockchain client: %w", err)
					}
				}

				systemTicker := genesis.NewTicker(
					ctx,
					_config,
					transactionsStore,
					manifestStore,
					blk,
				)

				// create and register a ticker and start
				// for blockchain transaction confirmation check
				ss := scheduler.New("BlockchainTransactionConfirmationCheckDaemon", *blockchainTxConfirmChkIntervalMin*60, true)
				bcTxChk := scheduler.BlockchainTransactionConfirmCheck{
					Conn:             conn,
					TransactionStore: transactionsStore,
					ManifestStore:    manifestStore,
					Blk:              blk,
				}
				ss.TaskRegisterAndStart(bcTxChk.Runner)

				APIController := api.NewAPIController(
					l,
					_config,
					skuStore,
					orderStore,
					containerStore,
					palletStore,
					cartonStore,
					productStore,
					organisationStore,
					userStore,
					referralStore,
					taskStore,
					userTaskStore,
					roleStore,
					blobStore,
					loyaltyStore,
					contractStore,
					distributorStore,
					transactionsStore,
					manifestStore,
					trackActionStore,
					userActivityStore,
					blacklistProvider,
					tokenStore,
					jwtSecret,
					auther,
					mailer,
					smsMessenger,
					blk,
					systemTicker,
				)

				server := &genesis.APIService{
					Log:  l,
					Addr: _config.API.Addr,
				}
				return server.Run(ctx, APIController)
			}, func(err error) {
				fmt.Println(err)
				cancel()
			})
			g.Add(func() error {
				l := genesis.NewLogToStdOut("loadbalancer", version, false)
				lb := genesis.LoadbalancerService{Log: l}
				return lb.Run(ctx, _config.API.Addr, _config.LoadBalancer, version)
			}, func(error) {
				fmt.Println(err)
				cancel()
				return
			})
			g.Add(func() error {
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt)

				select {
				case <-c:
					return errors.New("ctrl-c caught, exiting gracefully")
				}

			}, func(error) {
				fmt.Println(err)
				cancel()
				return
			})

			err = g.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}
	configCmd := &ffcli.Command{
		Name:        "config",
		Options:     []ff.Option{ff.WithEnvVarPrefix(envPrefix)},
		ShortUsage:  "genesis config",
		ShortHelp:   "Show environment variables.",
		Subcommands: []*ffcli.Command{dbCmd, serveCmd},
		Exec: func(ctx context.Context, args []string) error {
			var envVarReplacer = strings.NewReplacer(
				"-", "_",
				".", "_",
				"/", "_",
			)
			rootFlagSet.VisitAll(func(f *flag.Flag) {
				var key string
				key = strings.ToUpper(f.Name)
				key = envVarReplacer.Replace(key)
				key = strings.ToUpper(envPrefix) + "_" + key
				fmt.Println(key, f.Value.String())
			})
			dbFlagSet.VisitAll(func(f *flag.Flag) {
				var key string
				key = strings.ToUpper(f.Name)
				key = envVarReplacer.Replace(key)
				key = strings.ToUpper(envPrefix) + "_" + key
				fmt.Println(key, f.Value.String())
			})
			serveFlagSet.VisitAll(func(f *flag.Flag) {
				var key string
				key = strings.ToUpper(f.Name)
				key = envVarReplacer.Replace(key)
				key = strings.ToUpper(envPrefix) + "_" + key
				fmt.Println(key, f.Value.String())
			})
			return nil
		},
	}

	rootCmd := &ffcli.Command{
		Options:     []ff.Option{ff.WithEnvVarPrefix(envPrefix)},
		ShortUsage:  "genesis [flags]",
		ShortHelp:   "Ninja Genesis",
		Subcommands: []*ffcli.Command{dbCmd, serveCmd, configCmd},
		FlagSet:     rootFlagSet,
		Exec: func(ctx context.Context, args []string) error {
			var envVarReplacer = strings.NewReplacer(
				"-", "_",
				".", "_",
				"/", "_",
			)
			rootFlagSet.VisitAll(func(f *flag.Flag) {
				var key string
				key = strings.ToUpper(f.Name)
				key = envVarReplacer.Replace(key)
				key = strings.ToUpper(envPrefix) + "_" + key
				fmt.Println(key, f.Value.String())
			})
			dbFlagSet.VisitAll(func(f *flag.Flag) {
				var key string
				key = strings.ToUpper(f.Name)
				key = envVarReplacer.Replace(key)
				key = strings.ToUpper(envPrefix) + "_" + key
				fmt.Println(key, f.Value.String())
			})
			serveFlagSet.VisitAll(func(f *flag.Flag) {
				var key string
				key = strings.ToUpper(f.Name)
				key = envVarReplacer.Replace(key)
				key = strings.ToUpper(envPrefix) + "_" + key
				fmt.Println(key, f.Value.String())
			})
			return flag.ErrHelp
		},
	}
	if err := rootCmd.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := rootCmd.Run(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newMigrateInstance(conn *sqlx.DB) (*migrate.Migrate, error) {
	s := migrate_bindata.Resource(bindata.AssetNames(),
		func(name string) ([]byte, error) {
			return bindata.Asset(name)
		})
	d, err := migrate_bindata.WithInstance(s)
	if err != nil {
		return nil, fmt.Errorf("bindata instance: %w", err)
	}
	dbDriver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("db instance: %w", err)
	}
	m, err := migrate.NewWithInstance("go-bindata", d, "postgres", dbDriver)
	if err != nil {
		return nil, fmt.Errorf("migrate instance: %w", err)
	}
	return m, nil
}

func connect(
	DatabaseUser string,
	DatabasePass string,
	DatabaseHost string,
	DatabasePort string,
	DatabaseName string,

) *sqlx.DB {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DatabaseUser,
		DatabasePass,
		DatabaseHost,
		DatabasePort,
		DatabaseName,
	)
	conn, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatal("could not initialise database:", err)
	}
	if conn == nil {
		panic("conn is nil")
	}

	boil.SetDB(conn)
	return conn
}
