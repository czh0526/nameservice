package keeper

import (
	"context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"nameservice/x/nameservice/types"
)

func (k msgServer) BuyName(goCtx context.Context, msg *types.MsgBuyName) (*types.MsgBuyNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	whois, isFound := k.GetWhois(ctx, msg.Name)
	minPrice := sdk.Coins{
		sdk.NewInt64Coin("token", 10),
	}
	price, _ := sdk.ParseCoinsNormalized(whois.Price)
	bid, _ := sdk.ParseCoinsNormalized(msg.Bid)

	owner, _ := sdk.AccAddressFromBech32(whois.Owner)
	buyer, _ := sdk.AccAddressFromBech32(msg.Creator)

	if isFound {
		if price.IsAllGT(bid) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid is not high enough")
		}
		k.bankKeeper.SendCoins(ctx, buyer, owner, bid)
	} else {
		if minPrice.IsAllGT(bid) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid is less than min amount")
		}
		k.bankKeeper.SendCoinsFromAccountToModule(ctx, buyer, types.ModuleName, bid)
	}

	newWhois := types.Whois{
		Index: msg.Name,
		Name:  msg.Name,
		Value: whois.Value,
		Price: bid.String(),
		Owner: buyer.String(),
	}

	k.SetWhois(ctx, newWhois)

	return &types.MsgBuyNameResponse{}, nil
}
