package responses

import "time"

type RarityResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type ChestResponse struct {
	ID     uint            `json:"id"`
	Name   string          `json:"name"`
	Rarity *RarityResponse `json:"rarity"`
	Image  string          `json:"image"`
}

type UserChestResponse struct {
	ID            uint           `json:"id"`
	Health        uint           `json:"health"`
	CurrentHealth int            `json:"current_health"`
	Level         int            `json:"level"`
	Amount        uint32         `json:"amount"`
	Chest         *ChestResponse `json:"chest"`
}

type UserResponse struct {
	ID         uint               `json:"id"`
	Username   string             `json:"username"`
	TelegramId uint64             `json:"telegram_id"`
	Session    string             `json:"session"`
	UserChest  *UserChestResponse `json:"user_chest"`
}

type ItemTypeResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ItemResponse struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Rarity      *RarityResponse   `json:"rarity"`
	Image       string            `json:"image"`
	Type        *ItemTypeResponse `json:"type"`
	Description string            `json:"description"`
}

type UserItemResponse struct {
	ID        uint          `json:"id"`
	Item      *ItemResponse `json:"item"`
	CreatedAt time.Time     `json:"created_at"`
}

type GroupedUserItemResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Count  int    `json:"count"`
	Type   string `json:"type"`
	Rarity string `json:"rarity"`
	Image  string `json:"image"`
}

type ReferralUserResponse struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type AspectResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type AspectWithStatsResponse struct {
	AspectResponse
	UserLevel          uint    `json:"user_level,omitempty"`
	AspectStatId       uint    `json:"aspect_stat_id,omitempty"`
	Damage             uint    `json:"damage,omitempty"`
	CriticalDamage     uint    `json:"critical_damage,omitempty"`
	CriticalChance     float64 `json:"critical_change,omitempty"`
	GoldMultiplier     float64 `json:"gold_multiplier,omitempty"`
	Amount             uint    `json:"amount,omitempty"`
	AmountGrowthFactor float64 `json:"-"`
	PassiveDamage      uint    `json:"passive_damage,omitempty"`
}
