package transaction

import "startup/user"

type GetCampaignTrasactionInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}
