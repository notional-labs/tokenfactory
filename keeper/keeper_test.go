package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/tokenfactory/app/apptesting"
	"github.com/osmosis-labs/tokenfactory/keeper"
	"github.com/stretchr/testify/suite"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/osmosis-labs/tokenfactory/types"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	queryClient     types.QueryClient
	msgServer       types.MsgServer
	bankQueryClient banktypes.QueryClient

	// defaultDenom is on the suite, as it depends on the creator test address.
	defaultDenom string
}

var (
	baseTestAccts        = []sdk.AccAddress{}
	defaultTestStartTime = time.Now().UTC()
)

func init() {
	baseTestAccts = CreateRandomAccounts(3)
}

// CreateRandomAccounts is a function return a list of randomly generated AccAddresses
func CreateRandomAccounts(numAccts int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, numAccts)
	for i := 0; i < numAccts; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.Setup()

	// Fund every TestAcc with two denoms, one of which is the denom creation fee
	fundAccsAmount := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().DenomCreationFee[0].Denom, types.DefaultParams().DenomCreationFee[0].Amount.MulRaw(100)), sdk.NewCoin("utwo", sdkmath.NewInt(100000000)))
	for _, acc := range suite.TestAccs {
		suite.FundAcc(acc, fundAccsAmount)
	}

	suite.queryClient = types.NewQueryClient(suite.QueryHelper)
	suite.bankQueryClient = banktypes.NewQueryClient(suite.QueryHelper)
	suite.msgServer = keeper.NewMsgServerImpl(suite.App.TokenFactoryKeeper)
}

// FundAcc funds target address with specified amount.
func (suite *KeeperTestSuite) CreateDefaultDenom() {
	res, _ := suite.msgServer.CreateDenom(suite.Ctx, types.NewMsgCreateDenom(suite.TestAccs[0].String(), "bitcoin"))
	suite.defaultDenom = res.GetNewTokenDenom()
}

func (suite *KeeperTestSuite) OverrideMsgServer(newKeeper keeper.Keeper) {
	suite.msgServer = keeper.NewMsgServerImpl(newKeeper)
}
