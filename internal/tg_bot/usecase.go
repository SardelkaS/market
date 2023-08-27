package tg_bot

import tg_bot_model "market_auth/internal/tg_bot/model"

type UC interface {
	NotifyNewOrder(input tg_bot_model.NotifyNewOrderLogicInput)
}
