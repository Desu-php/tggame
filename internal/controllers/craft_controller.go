package controllers

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/internal/responses"
	"example.com/v2/internal/services"
	"example.com/v2/pkg/db"
	"example.com/v2/pkg/errs"
	"example.com/v2/pkg/transaction"
	"example.com/v2/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type CraftController struct {
	logger          *logrus.Logger
	db              *db.DB
	trx             transaction.TransactionManager
	userItemService *services.UserItemService
	userStatService *services.UserStatService
}

type ItemInput struct {
	ItemID uint `json:"item_id" binding:"required"`
	Count  uint `json:"count" binding:"required"`
}

type ItemInputs struct {
	Items []ItemInput `json:"items" binding:"required"`
}

type UserItemCount struct {
	ItemID uint `json:"item_id"`
	Count  uint `json:"count"`
}

func NewCraftController(
	logger *logrus.Logger,
	db *db.DB,
	trx transaction.TransactionManager,
	userItemService *services.UserItemService,
	userStatService *services.UserStatService,
) *CraftController {
	return &CraftController{
		logger,
		db,
		trx,
		userItemService,
		userStatService,
	}
}

func (cc *CraftController) Craft(c *gin.Context) {
	var itemInputs ItemInputs
	if err := c.ShouldBindJSON(&itemInputs); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	totalCount := 0
	itemIds := make([]uint, len(itemInputs.Items))
	for _, item := range itemInputs.Items {
		totalCount += int(item.Count)
		itemIds = append(itemIds, item.ItemID)
	}

	if totalCount != 4 {
		c.JSON(400, gin.H{"error": "Invalid count"})
		return
	}

	var items []models.Item

	err := cc.db.WithContext(c).Model(&models.Item{}).
		Where("id IN ?", itemIds).
		Preload("Rarity").
		Find(&items).Error

	if err != nil {
		cc.logger.WithError(err).Error("CraftController::Craft")
		responses.ServerErrorResponse(c)
		return
	}

	if len(items) == 0 {
		c.JSON(400, gin.H{"error": "items is required"})
		return
	}

	rarityId := items[0].RarityID

	for _, item := range items {
		if rarityId != item.RarityID {
			c.JSON(400, gin.H{"error": "You can only craft items of the same rarity", "code": errs.CraftingItemsDifferentRarityCode})
			return
		}
	}

	user, ok := utils.GetUser(c)

	if !ok {
		cc.logger.WithError(err).Error("CraftController::Craft")
		responses.ServerErrorResponse(c)
		return
	}

	var userItems []UserItemCount

	err = cc.db.WithContext(c).Model(models.UserItem{}).
		Select("item_id, COUNT(item_id) as count").
		Where("user_items.user_id = ?", user.ID).
		Where("item_id IN ?", itemIds).
		Group("item_id").
		Scan(&userItems).Error

	if err != nil {
		cc.logger.WithError(err).Error("CraftController::Craft")
		responses.ServerErrorResponse(c)
		return
	}

	var itemsMap = make(map[uint]uint)
	for _, input := range itemInputs.Items {
		hasEnough := false
		for _, userItem := range userItems {
			if userItem.ItemID == input.ItemID {
				if userItem.Count >= input.Count {
					itemsMap[input.ItemID] = userItem.Count - input.Count
					hasEnough = true
				}
				break
			}
		}

		if !hasEnough {
			c.JSON(400, gin.H{"error": "Not enough items for crafting", "code": errs.CraftingNotEnoughItemsCode})
			return
		}
	}

	var rarity models.Rarity

	err = cc.db.WithContext(c).Model(models.Rarity{}).
		Where("sort > ?", items[0].Rarity.Sort).
		Order("sort ASC").
		First(&rarity).Error

	if err != nil {
		cc.logger.WithError(err).Error("CraftController::Craft")
		responses.ServerErrorResponse(c)
		return
	}

	itemRarity := craftRandomizer(&rarity, &items[0].Rarity)

	var randomItem models.Item

	err = cc.db.WithContext(c).
		Model(&models.Item{}).
		Where("rarity_id = ?", itemRarity.ID).
		Order("RANDOM()").
		Preload("Rarity").
		Preload("Type").
		Limit(1).
		First(&randomItem).Error

	if err != nil {
		cc.logger.WithError(err).Error("CraftController::Craft")
		responses.ServerErrorResponse(c)
		return
	}

	err = cc.trx.RunInTransaction(c, func(ctx context.Context) error {
		for _, item := range itemInputs.Items {
			var ids []uint
			err = cc.db.WithContext(ctx).Model(&models.UserItem{}).
				Where("user_id = ?", user.ID).
				Where("item_id = ?", item.ItemID).
				Limit(int(item.Count)).
				Pluck("id", &ids).Error

			if err != nil {
				return err
			}

			if len(ids) > 0 {
				err = cc.db.WithContext(ctx).Delete(&models.UserItem{}, ids).Error
				if err != nil {
					return err
				}
			}

			if itemsMap[item.ItemID] == 0 {
				var downgradeItem models.Item

				err = cc.db.WithContext(ctx).Model(&models.Item{}).First(&downgradeItem, item.ItemID).Error

				if err != nil {
					return err
				}

				err = cc.userStatService.Downgrade(ctx, &services.UserStatDowngradeDto{
					User:         user,
					Attributable: &downgradeItem,
				})

				if err != nil {
					return err
				}
			}
		}

		_, err = cc.userItemService.SetUserItem(ctx, user, &randomItem, nil)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		cc.logger.WithError(err).Error("CraftController::Craft")
		responses.ServerErrorResponse(c)
		return
	}

	c.JSON(200, gin.H{"data": randomItem})
}

func craftRandomizer(betterRarity *models.Rarity, sameRarity *models.Rarity) *models.Rarity {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomChance := r.Float64() * 100

	if randomChance < betterRarity.CraftChance {
		return betterRarity
	}

	return sameRarity
}
