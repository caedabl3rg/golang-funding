package campaign

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	FindALl() ([]Campaign, error)
	FindByUserID(userID int) ([]Campaign, error)
}
type repository struct {
	db *gorm.DB
}

func NewRepositoryCampaign(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) FindALl() ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *repository) FindByUserID(userID int) ([]Campaign, error) {
	var campaigns []Campaign

	err := r.db.Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	fmt.Println("sampai sini", campaigns)
	return campaigns, nil
}
