package keeper_test

import (
	"testing"
	"time"

	checkersapp "github.com/alice/checkers/app"
	"github.com/alice/checkers/x/checkers/keeper"
	"github.com/alice/checkers/x/checkers/testutil"
	"github.com/alice/checkers/x/checkers/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

const (
	alice = testutil.Alice
	bob   = testutil.Bob
	carol = testutil.Carol
)

const (
	balAlice = 50000000
	balBob   = 20000000
	balCarol = 10000000
)

// Define your test suite in a new keeper_integration_suite_test.go file. In a dedicated folder `tests/integration/checkers/keeper/`

type IntegrationTestSuite struct {
	suite.Suite

	app         *checkersapp.App
	msgServer   types.MsgServer
	ctx         sdk.Context
	queryClient types.QueryClient
}

var (
	checkersModuleAddress string
)

func TestCheckersKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

/**
This SetupTest function is like a `beforeEach` as found in other test libraries. With it, you always get a new `app`
in each test, without interference between them. Do not omit it unless you have specific reasons to do so.

- It collects your checkersModuleAddress for later use in tests that check events and balances
**/

func (suite *IntegrationTestSuite) SetupTest() {
	// setup the app
	app := checkersapp.Setup(false)
	// setup the context
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

	// set the auth to the account keeper
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	// set the auth to the bank keeper
	app.BankKeeper.SetParams(ctx, banktypes.DefaultParams())
	// getting the checkers module address from the account keeper by the module name.
	checkersModuleAddress = app.AccountKeeper.GetModuleAddress(types.ModuleName).String()

	// Something to deal with the queries.
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.CheckersKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	// Setting up the suite
	suite.app = app
	suite.msgServer = keeper.NewMsgServerImpl(app.CheckersKeeper)
	suite.ctx = ctx
	suite.queryClient = queryClient
}

// Make a bank genesis Balance type from primitives:
func makeBalance(address string, balance int64) banktypes.Balance {
	return banktypes.Balance{
		Address: address,
		Coins: sdk.Coins{
			sdk.Coin{
				Denom:  sdk.DefaultBondDenom,
				Amount: sdk.NewInt(balance),
			},
		},
	}
}

// Make your preferred bank genesis state:

func getBankGenesis() *banktypes.GenesisState {
	// array of coins of type banktypes.Balance
	// got the address and the aray of coin objects that thy have.
	coins := []banktypes.Balance{
		makeBalance(alice, balAlice),
		makeBalance(bob, balBob),
		makeBalance(carol, balCarol),
	}

	// supply is equal to the sum of all the balances.
	supply := banktypes.Supply{
		Total: coins[0].Coins.Add(coins[1].Coins...).Add(coins[2].Coins...),
	}

	// setting the state of the new genesis i assume because coins array  this just gets the genesis state lmao
	// has the balances and addresses associated with the balance this means that they are assigned that balance.
	state := banktypes.NewGenesisState(
		banktypes.DefaultParams(),
		coins,
		supply.Total,
		[]banktypes.Metadata{})

	return state
}

// Add a function to prepare your suite with your desired balances.
func (suite *IntegrationTestSuite) setupSuiteWithBalances() {
	suite.app.BankKeeper.InitGenesis(suite.ctx, getBankGenesis())
}

// Add a function to check balancs from primitives:

func (suite *IntegrationTestSuite) RequireBankBalance(expected int, atAddress string) {
	sdkAdd, err := sdk.AccAddressFromBech32(atAddress)
	suite.Require().Nil(err, "Failed to parse address: %s", atAddress)
	suite.Require().Equal(
		int64(expected),
		suite.app.BankKeeper.GetBalance(suite.ctx, sdkAdd, sdk.DefaultBondDenom).Amount.Int64())
}
