package biz

import (
	"context"

	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/pin/models"
	"github.com/ssonit/aura_server/internal/pin/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	store utils.PinStore
}

func NewService(store utils.PinStore) *service {
	return &service{store: store}
}

func (s *service) GetBoardPinItem(ctx context.Context, filter *models.BoardPinFilter) (*models.BoardPinModel, error) {
	data, err := s.store.GetBoardPinItem(ctx, filter)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *service) ListBoardPinItem(ctx context.Context, filter *models.BoardPinFilter, paging *common.Paging) ([]models.BoardPinModel, error) {
	data, err := s.store.ListBoardPinItem(ctx, filter, paging)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *service) UpdatePin(ctx context.Context, id string, pin *models.PinUpdate, userId string) error {
	oID, _ := primitive.ObjectIDFromHex(id)

	data, err := s.store.GetItem(ctx, map[string]interface{}{"_id": oID})

	if err != nil {
		return err
	}

	if data.UserId.Hex() != userId {
		return utils.ErrUserNotPermitted
	}

	err = s.store.UpdatePin(ctx, id, pin)

	if err != nil {
		return err
	}

	var filter models.BoardPinFilter

	userOId, _ := primitive.ObjectIDFromHex(userId)

	filter.PinId = oID
	filter.UserId = userOId

	err = s.store.DeleteBoardPin(ctx, &filter)

	if err != nil {
		return err
	}

	_, err = s.store.CreateBoardPin(ctx, &models.BoardPinCreation{
		BoardId: pin.BoardId,
		PinId:   oID,
		UserId:  userOId,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *service) CreatePin(ctx context.Context, p *models.PinCreation) (primitive.ObjectID, error) {
	boardId, err := s.store.GetBoardByUserId(ctx, p.UserId)

	if err != nil {
		return primitive.NilObjectID, err
	}

	if boardId.IsZero() {
		return primitive.NilObjectID, utils.ErrBoardNotFound
	}

	data, err := s.store.Create(ctx, p)

	if err != nil {
		return primitive.NilObjectID, err
	}

	if data.IsZero() {
		return primitive.NilObjectID, utils.ErrDataIsZero
	}

	if err != nil {
		return primitive.NilObjectID, err
	}

	// Create board pin for all pins
	_, err = s.store.CreateBoardPin(ctx, &models.BoardPinCreation{
		BoardId: boardId,
		PinId:   data,
		UserId:  p.UserId,
	})

	// Create board pin for the pin
	_, err = s.store.CreateBoardPin(ctx, &models.BoardPinCreation{
		BoardId: p.BoardId,
		PinId:   data,
		UserId:  p.UserId,
	})

	if err != nil {
		return primitive.NilObjectID, err
	}

	return data, nil

}

func (s *service) ListPinItem(ctx context.Context, filter *models.Filter, paging *common.Paging) ([]models.PinModel, error) {

	data, err := s.store.ListItem(ctx, filter, paging)

	if err != nil {
		return nil, err
	}

	return data, nil

}

func (s *service) GetPinById(ctx context.Context, id string) (*models.PinModel, error) {

	oID, _ := primitive.ObjectIDFromHex(id)

	data, err := s.store.GetItem(ctx, map[string]interface{}{"_id": oID})

	if err != nil {
		return nil, err
	}

	return data, nil
}
