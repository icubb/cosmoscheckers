package keeper

import (
	"context"
	"fmt"

	"github.com/alice/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
- Differentiate between two types of errors:
	- Errors relating to the move, returning a reason.
	- Errors indicating that testing the move is impossible, returning an error.
**/

func (k Keeper) CanPlayMove(goCtx context.Context, req *types.QueryCanPlayMoveRequest) (*types.QueryCanPlayMoveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Process the query
	// The game needs to be fetched, if it does not exist return an error message.
	storedGame, found := k.GetStoredGame(ctx, req.GameIndex)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "%s", req.GameIndex)
	}

	// Has the game already been won?
	// If the winner in stored game is not an empty player.
	if storedGame.Winner != rules.PieceStrings[rules.NO_PLAYER] {
		return &types.QueryCanPlayMoveResponse{
			Possible: false,
			Reason:   types.ErrGameFinished.Error(),
		}, nil
	}

	// Is the player actually one of the game players?

	isBlack := rules.PieceStrings[rules.BLACK_PLAYER] == req.Player
	isRed := rules.PieceStrings[rules.RED_PLAYER] == req.Player
	var player rules.Player
	if isBlack && isRed {
		player = rules.StringPieces[storedGame.Turn].Player
	} else if isBlack {
		player = rules.BLACK_PLAYER
	} else if isRed {
		player = rules.RED_PLAYER
	} else {
		return &types.QueryCanPlayMoveResponse{
			Possible: false,
			Reason:   fmt.Sprintf("%s: %s", types.ErrCreatorNotPlayer.Error(), req.Player),
		}, nil
	}

	// Is it the players turn

	game, err := storedGame.ParseGame()
	if err != nil {
		return nil, err
	}
	if !game.TurnIs(player) {
		return &types.QueryCanPlayMoveResponse{
			Possible: false,
			Reason:   fmt.Sprintf("%s: %s", types.ErrNotPlayerTurn.Error(), player.Color),
		}, nil
	}

	// Attempt the move and report back

	_, moveErr := game.Move(
		rules.Pos{
			X: int(req.FromX),
			Y: int(req.FromY),
		},
		rules.Pos{
			X: int(req.FromX),
			Y: int(req.FromY),
		},
	)
	if moveErr != nil {
		return &types.QueryCanPlayMoveResponse{
			Possible: false,
			Reason:   fmt.Sprintf("%s: %s", types.ErrWrongMove.Error(), moveErr.Error()),
		}, nil
	}

	return &types.QueryCanPlayMoveResponse{
		Possible: true,
		Reason:   "ok",
	}, nil

	//return &types.QueryCanPlayMoveResponse{}, nil
}
