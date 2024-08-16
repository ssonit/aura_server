package models

import "github.com/ssonit/aura_server/common"

type Media struct {
	common.CommonModel
	Url       string `json:"url,omitempty" bson:"url,omitempty"`
	SecureUrl string `json:"secure_url,omitempty" bson:"secure_url,omitempty"`
	PublicId  string `json:"public_id,omitempty" bson:"public_id,omitempty"`
	Format    string `json:"format,omitempty" bson:"format,omitempty"`
	Width     int    `json:"width,omitempty" bson:"width,omitempty"`
	Height    int    `json:"height,omitempty" bson:"height,omitempty"`
}
