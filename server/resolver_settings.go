package genesis

import (
	"context"
	"genesis/db"
	"genesis/graphql"

	"github.com/ninja-software/terror"
)

////////////////
//  Resolver  //
////////////////

// Settings resolver
func (r *Resolver) Settings() graphql.SettingsResolver {
	return &settingsResolver{r}
}

type settingsResolver struct{ *Resolver }

func (r *queryResolver) Settings(ctx context.Context) (*db.Setting, error) {
	settings, err := r.TransactionStore.GetSettings()
	if err != nil {
		return nil, terror.New(err, "get settings (1)")
	}
	return settings, nil
}

func (r *settingsResolver) ConsumerHost(ctx context.Context, obj *db.Setting) (string, error) {
	return r.Config.API.ConsumerHost, nil
}
func (r *settingsResolver) AdminHost(ctx context.Context, obj *db.Setting) (string, error) {
	return r.Config.API.AdminHost, nil
}
func (r *settingsResolver) EtherscanHost(ctx context.Context, obj *db.Setting) (string, error) {
	return r.Config.Blockchain.EtherscanHost, nil
}
func (r *settingsResolver) FieldappVersion(ctx context.Context, obj *db.Setting) (string, error) {
	return r.Config.Fieldapp.Version, nil
}
