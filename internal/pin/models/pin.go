package models

import "github.com/ssonit/aura_server/common"

type Pin struct {
	common.CommonModel
	UserId      string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	BoardId     string `json:"board_id,omitempty" bson:"board_id,omitempty"`
	Title       string `json:"title,omitempty" bson:"title,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	MediaId     string `json:"media_id,omitempty" bson:"media_id,omitempty"`
	LinkUrl     string `json:"link_url,omitempty" bson:"link_url,omitempty"`
}

type PinCreation struct {
	UserId      string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	BoardId     string `json:"board_id,omitempty" bson:"board_id,omitempty"`
	Title       string `json:"title,omitempty" bson:"title,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	MediaId     string `json:"media_id,omitempty" bson:"media_id,omitempty"`
	LinkUrl     string `json:"link_url,omitempty" bson:"link_url,omitempty"`
}
