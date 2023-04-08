package transaction

import (
	"errors"
	"startup/campaign"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTrasactionInput) ([]Transaction, error)
}

func NewServiceTransaction(repository Repository, campaignRepository campaign.Repository ) *service {
	return &service{repository: repository, campaignRepository: campaignRepository}
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTrasactionInput) ([]Transaction, error) {

	// get campaign -> check campaign.userid != user_id yang melakukan request
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}
	
	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("Not a owner of the campign")
	}
	 

	transaction, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
