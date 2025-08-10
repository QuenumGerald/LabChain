package milestones

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"labchain/x/milestones/keeper"
	"labchain/x/milestones/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.Register(
		&types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *types.Module
	StoreService store.KVStoreService
	Cdc          codec.Codec
	AddressCodec address.Codec

	AuthKeeper  types.AuthKeeper
	BankKeeper  types.BankKeeper
	GovKeeper   *govkeeper.Keeper
	GroupKeeper groupkeeper.Keeper
}

type ModuleOutputs struct {
	depinject.Out

	MilestonesKeeper keeper.Keeper
	Module           appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(types.GovModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}
	k := keeper.NewKeeper(
		in.StoreService,
		in.Cdc,
		in.AddressCodec,
		authority,
		in.BankKeeper,
		in.GovKeeper,
		in.GroupKeeper,
	)
	m := NewAppModule(in.Cdc, k, in.AuthKeeper, in.BankKeeper)

	return ModuleOutputs{MilestonesKeeper: k, Module: m}
}
