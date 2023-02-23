package keeper

import (
	"context"
	"strconv"

	"github.com/alice/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	// Getting the system info for the context gets the next game id because
	// the system info message holds the next id.
	systemInfo, found := k.Keeper.GetSystemInfo(ctx)
	if !found {
		panic("SystemInfo not found")
	}

	newIndex := strconv.FormatUint(systemInfo.NextId, 10)

	newGame := rules.New()

	storedGame := types.StoredGame{
		Index:       newIndex, // using the new index from system info here.
		Board:       newGame.String(),
		Turn:        rules.PieceStrings[newGame.Turn],
		Black:       msg.Black, // these come from the command line message. or grpc
		Red:         msg.Red,
		MoveCount:   0,
		BeforeIndex: types.NoFifoIndex,
		AfterIndex:  types.NoFifoIndex,
		Deadline:    types.FormatDeadline(types.GetNextDeadline(ctx)),
		Winner:      rules.PieceStrings[rules.NO_PLAYER],
		Wager:       msg.Wager,
	}

	// Confirm that the values in the object are correct by checking the validity of the players
	// addresses:

	err := storedGame.Validate()
	if err != nil {
		return nil, err
	}

	// Send the stored game to the tail.(because it is the most recent now)
	k.Keeper.SendToFifoTail(ctx, &storedGame, &systemInfo)

	//Save the storedGame object using the Keeper.SetStoredGame function created by the
	// ignite scaffold map storedGame command.
	k.Keeper.SetStoredGame(ctx, storedGame)

	// Prepare the ground work for the next game using Keeper.SetSystemInfo function
	// created by Ignite CLI
	systemInfo.NextId++
	k.Keeper.SetSystemInfo(ctx, systemInfo)

	// Consume the gas for creating the game.
	ctx.GasMeter().ConsumeGas(types.CreateGameGas, "Create game")

	// Emit the events when the game is created

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.GameCreatedEventType,
			sdk.NewAttribute(types.GameCreatedEventCreator, msg.Creator),
			sdk.NewAttribute(types.GameCreatedEventGameIndex, newIndex),
			sdk.NewAttribute(types.GameCreatedEventBlack, msg.Black),
			sdk.NewAttribute(types.GameCreatedEventRed, msg.Red),
			sdk.NewAttribute(types.GameCreatedEventWager, strconv.FormatUint(msg.Wager, 10)),
		),
	)
	// Return the newley creatd id for reference.
	return &types.MsgCreateGameResponse{
		GameIndex: newIndex,
	}, nil
}
