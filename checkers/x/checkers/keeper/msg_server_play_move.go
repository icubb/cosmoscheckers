package keeper

import (
	"context"

	"strconv"

	"github.com/alice/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) PlayMove(goCtx context.Context, msg *types.MsgPlayMove) (*types.MsgPlayMoveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	// get the stored game information from the keeper
	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.GameIndex)
	// you return an error since this is a player mistake
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "%s", msg.GameIndex)
	}

	// if the stored game winner is not a no player means that the game has ended. with a winner.
	if storedGame.Winner != rules.PieceStrings[rules.NO_PLAYER] {
		return nil, types.ErrGameFinished
	}

	// WHat this is doing is checking if the player is legitimate in this game to play
	// if they are black or red they will be the same
	isBlack := storedGame.Black == msg.Creator
	isRed := storedGame.Red == msg.Creator
	var player rules.Player
	if !isBlack && !isRed {
		return nil, sdkerrors.Wrapf(types.ErrCreatorNotPlayer, "%s", msg.Creator)
	} else if isBlack && isRed { // i dont understand this it sets the player to the player who's turn it is??
		player = rules.StringPieces[storedGame.Turn].Player
	} else if isBlack { // if the player is black then player is set to black.
		player = rules.BLACK_PLAYER
	} else { // if they are red they are set to red.
		player = rules.RED_PLAYER
	}

	// Instantiate the board in order to implement the rules:

	game, err := storedGame.ParseGame()

	if err != nil {
		panic(err.Error())
	}

	// Is it the players turn? Check using the rules file's own TurnIs function:

	if !game.TurnIs(player) {
		return nil, sdkerrors.Wrapf(types.ErrNotPlayerTurn, "%s", player)
	}

	// Collect the wager from the player.
	err = k.Keeper.CollectWager(ctx, &storedGame)
	if err != nil {
		return nil, err
	}

	// Properly conduct the move using the rules move function:

	captured, moveErr := game.Move(
		rules.Pos{
			X: int(msg.FromX),
			Y: int(msg.FromY),
		},
		rules.Pos{
			X: int(msg.ToX),
			Y: int(msg.ToY),
		},
	)

	// if the move is an error state it so.
	if moveErr != nil {
		return nil, sdkerrors.Wrapf(types.ErrWrongMove, moveErr.Error())
	}

	// Update the winner field, which remains neutral if there is no winner yet:
	storedGame.Winner = rules.PieceStrings[game.Winner()]

	systemInfo, found := k.Keeper.GetSystemInfo(ctx)
	if !found {
		panic("SystemInfo not found")
	}

	lastBoard := game.String()
	if storedGame.Winner == rules.PieceStrings[rules.NO_PLAYER] {
		k.Keeper.SendToFifoTail(ctx, &storedGame, &systemInfo)
		storedGame.Board = lastBoard
	} else {
		k.Keeper.RemoveFromFifo(ctx, &storedGame, &systemInfo)
		storedGame.Board = ""
		k.Keeper.MustPayWinnings(ctx, &storedGame)
	}

	//k.Keeper.SendToFifoTail(ctx, &storedGame, &systemInfo)

	// update the move count for the game.
	storedGame.MoveCount++
	storedGame.Deadline = types.FormatDeadline(types.GetNextDeadline(ctx))
	// Prepare the updated board to be stored and store the information:
	//storedGame.Board = game.String()
	storedGame.Turn = rules.PieceStrings[game.Turn]
	// Sets the stored and system info that changed in the send to fifo tail section.
	k.Keeper.SetStoredGame(ctx, storedGame)
	k.Keeper.SetSystemInfo(ctx, systemInfo)

	// Consume the gas for playing a move
	ctx.GasMeter().ConsumeGas(types.PlayMoveGas, "Play a move")

	// This updates the fields that were modified using the Keeper.SetStoredGame
	// Function, as when you created and saved the game.

	// Emit the events

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.MovePlayedEventType,
			sdk.NewAttribute(types.MovePlayedEventCreator, msg.Creator),
			sdk.NewAttribute(types.MovePlayedEventGameIndex, msg.GameIndex),
			sdk.NewAttribute(types.MovePlayedEventCapturedX, strconv.FormatInt(int64(captured.X), 10)),
			sdk.NewAttribute(types.MovePlayedEventCapturedY, strconv.FormatInt(int64(captured.Y), 10)),
			sdk.NewAttribute(types.MovePlayedEventWinner, rules.PieceStrings[game.Winner()]),
			sdk.NewAttribute(types.MovePlayedEventBoard, lastBoard),
		),
	)

	// Return relevant information regarding the move's result:

	return &types.MsgPlayMoveResponse{
		CapturedX: int32(captured.X),
		CapturedY: int32(captured.Y),
		Winner:    rules.PieceStrings[game.Winner()],
	}, nil

	// The captured and winner information would be lost if you did not get it out of the function
	// More accurately, one would have to replay the transaction to discover the values. It is best
	// to make this information easily accessible.

}
