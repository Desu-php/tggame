package dto

type GameStartDto struct {
	Username string `json:"username" binding:"required,max=255"`
	TelegramId uint64 `json:"telegram_id" binding:"required,min=1"`
}