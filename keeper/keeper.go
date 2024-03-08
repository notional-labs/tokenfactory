package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	store "cosmossdk.io/store/types"
	storetypes "cosmossdk.io/store/types"
	"github.com/osmosis-labs/tokenfactory/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey

		accountKeeper  types.AccountKeeper
		bankKeeper     types.BankKeeper
		contractKeeper types.ContractKeeper

		distrKeeper types.DistrKeeper

		authority string
	}
)

// NewKeeper returns a new instance of the x/tokenfactory keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distrKeeper types.DistrKeeper,
	authority string,
) Keeper {

	return Keeper{
		storeKey:      storeKey,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrKeeper:   distrKeeper,
		cdc:           cdc,
		authority:     authority,
	}
}

// Logger returns a logger for the x/tokenfactory module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetDenomPrefixStore returns the substore for a specific denom
func (k Keeper) GetDenomPrefixStore(ctx sdk.Context, denom string) store.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.GetDenomPrefixStore(denom))
}

// GetCreatorPrefixStore returns the substore for a specific creator address
func (k Keeper) GetCreatorPrefixStore(ctx sdk.Context, creator string) store.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.GetCreatorPrefix(creator))
}

// GetCreatorsPrefixStore returns the substore that contains a list of creators
func (k Keeper) GetCreatorsPrefixStore(ctx sdk.Context) store.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.GetCreatorsPrefix())
}

// Set the wasm keeper.
func (k *Keeper) SetContractKeeper(contractKeeper types.ContractKeeper) {
	k.contractKeeper = contractKeeper
}
