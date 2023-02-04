package keeper

import (
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) RemoveFromFifo(ctx sdk.Context, game *types.StoredGame, info *types.SystemInfo) {
	// Does it have a predecessor?
	if game.BeforeIndex != types.NoFifoIndex {
		// Get game at the index prior to the current game
		beforeElement, found := k.GetStoredGame(ctx, game.BeforeIndex)
		// if the stored game is not found then panic. because it should be reasonably.
		if !found {
			panic("Element before in Fifo was not found")
		}
		// set the before games after index to equal the currents game after index.
		beforeElement.AfterIndex = game.AfterIndex
		// sets the before element.
		k.SetStoredGame(ctx, beforeElement)
		// if the after index is equal to -1 then set the tail index to the before element index.
		// if the after index was at the tail then it sets the bfore element index to the tail.
		if game.AfterIndex == types.NoFifoIndex {
			info.FifoTailIndex = beforeElement.Index
		}
		// Is it at the FIFO head?
	} else if info.FifoHeadIndex == game.Index {
		info.FifoHeadIndex = game.AfterIndex
	}
	// Does it have a successor?
	if game.AfterIndex != types.NoFifoIndex {
		afterElement, found := k.GetStoredGame(ctx, game.AfterIndex)
		if !found {
			panic("Element after in Fifo was not found")
		}
		afterElement.BeforeIndex = game.BeforeIndex
		k.SetStoredGame(ctx, afterElement)
		if game.BeforeIndex == types.NoFifoIndex {
			info.FifoHeadIndex = afterElement.Index
		}
		// is it at the FIFO tail?
	} else if info.FifoTailIndex == game.Index {
		info.FifoTailIndex = game.BeforeIndex
	}
	game.BeforeIndex = types.NoFifoIndex
	game.AfterIndex = types.NoFifoIndex
}

func (k Keeper) SendToFifoTail(ctx sdk.Context, game *types.StoredGame, info *types.SystemInfo) {
	// Essentially if the head and tail don't exist yet, list is empty.
	if info.FifoHeadIndex == types.NoFifoIndex && info.FifoTailIndex == types.NoFifoIndex {
		game.BeforeIndex = types.NoFifoIndex
		game.AfterIndex = types.NoFifoIndex
		info.FifoHeadIndex = game.Index
		info.FifoTailIndex = game.Index
		// if the head or tail are empty.
	} else if info.FifoHeadIndex == types.NoFifoIndex || info.FifoTailIndex == types.NoFifoIndex {
		panic("Fifo should have both head and tail or none")
	} else if info.FifoTailIndex == game.Index {
		// Nothing to do, already at tail
	} else {
		// Snip game out
		k.RemoveFromFifo(ctx, game, info)

		// Now add to tail
		currentTail, found := k.GetStoredGame(ctx, info.FifoTailIndex)
		if !found {
			panic("Current Fifo tail was not found")
		}
		currentTail.AfterIndex = game.Index
		k.SetStoredGame(ctx, currentTail)

		game.BeforeIndex = currentTail.Index
		info.FifoTailIndex = game.Index
	}
}
