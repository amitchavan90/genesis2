package config

// NewConfig returns a new PlatformConfig
func NewConfig() *PlatformConfig {
	c := &PlatformConfig{
		API:          &API{},
		LoadBalancer: &LoadBalancer{},
		Database:     &Database{},
		UserAuth:     &UserAuth{},
		Email:        &Email{},
		SMS:          &SMS{},
		Skip:         &Skip{},
		Blockchain:   &Blockchain{},
		Fieldapp:     &Fieldapp{},
		Sentry:       &Sentry{},
	}
	return c
}

// PlatformConfig for the Platform
type PlatformConfig struct {
	API          *API
	LoadBalancer *LoadBalancer
	Database     *Database
	UserAuth     *UserAuth
	Email        *Email
	SMS          *SMS
	Skip         *Skip
	Blockchain   *Blockchain
	Fieldapp     *Fieldapp
	Sentry       *Sentry
}

// API for the API service
type API struct {
	Addr         string
	BlobBaseURL  string
	ConsumerHost string
	AdminHost    string
}

// LoadBalancer for the LoadBalancer service
type LoadBalancer struct {
	AdminAddr         string
	AdminIndexFile    string
	AdminRootPath     string
	ConsumerAddr      string
	ConsumerIndexFile string
	ConsumerRootPath  string
}

// Database for the Database service
type Database struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

// UserAuth holds variables for user auth config, such as token expiry
type UserAuth struct {
	JWTSecret             string
	TokenExpiryDays       int
	ResetTokenExpiryDays  int
	BlacklistRefreshHours int
}

// Email holds variables for the mailer
type Email struct {
	Domain string
	Sender string
	APIKey string
}

// SMS holds config variables for SMS messaging functions.
type SMS struct {
	Enabled bool
	SID     string
	Token   string
}

// Skip holds config variables that allows skipping certain things
type Skip struct {
	UserVerification bool
}

// Blockchain for the Blockchain service
type Blockchain struct {
	PrivateKeyBytes             string
	PrivateKeyPassword          string
	EthereumHost                string
	EtherscanHost               string
	TestPrivateKey              string
	CreateContractOnFirstRun    bool
	FlushPendingActionsInterval int
	EthLowNotifyAmount          float64
	EthLowNotifyEmail           string
}

// Fieldapp info config
type Fieldapp struct {
	Version string
}

// Sentry holds config variables that allows skipping certain things
type Sentry struct {
	DSN         string `desc:"Sends error to remote server. If set, it will send error." default:""`
	ServerName  string `desc:"The machine name that this program is running on" default:"dev-pc"`
	Environment string `desc:"This program environment" default:"development"`
}

// PasswordMinimumLength global constant of password minimum length requirement
const PasswordMinimumLength uint = 8
