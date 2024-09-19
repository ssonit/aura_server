package biz

import (
	"context"

	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/pin/models"
	"github.com/ssonit/aura_server/internal/pin/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AllPins = "all_pins"
)

type service struct {
	store utils.PinStore
}

func NewService(store utils.PinStore) *service {
	return &service{store: store}
}

func (s *service) UnLikePin(ctx context.Context, p *models.LikeDelete) error {

	userOId, err := primitive.ObjectIDFromHex(p.UserId)
	pinOId, err := primitive.ObjectIDFromHex(p.PinId)

	if err != nil {
		return utils.ErrFailedToDecodeObjID
	}

	err = s.store.UnlikePin(ctx, userOId, pinOId)

	if err != nil {
		return err
	}

	return nil

}

func (s *service) LikePin(ctx context.Context, p *models.LikeCreation) error {

	userOId, err := primitive.ObjectIDFromHex(p.UserId)
	pinOId, err := primitive.ObjectIDFromHex(p.PinId)

	if err != nil {
		return utils.ErrFailedToDecodeObjID
	}

	err = s.store.LikePin(ctx, userOId, pinOId)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) SaveBoardPin(ctx context.Context, p *models.BoardPinSave) (primitive.ObjectID, error) {
	userOId, _ := primitive.ObjectIDFromHex(p.UserId)
	pinOId, _ := primitive.ObjectIDFromHex(p.PinId)
	selectedBoardOId, _ := primitive.ObjectIDFromHex(p.BoardId)

	boardId, err := s.store.GetBoardByUserId(ctx, userOId, AllPins)

	if err != nil {
		return primitive.NilObjectID, err
	}

	if boardId.IsZero() {
		return primitive.NilObjectID, utils.ErrBoardNotFound
	}

	boardPinId, err := s.store.GetBoardPinItem(ctx, &models.BoardPinFilter{
		PinId:  pinOId,
		UserId: userOId,
	})

	if err != nil {
		return primitive.NilObjectID, err
	}

	pinInBoard, err := s.store.CheckIfPinExistsInBoard(ctx, boardId, pinOId)
	if err != nil {
		return primitive.NilObjectID, err
	}

	if !pinInBoard {
		// Create board pin for all pins
		_, err = s.store.CreateBoardPin(ctx, &models.BoardPinCreation{
			BoardId: boardId,
			PinId:   pinOId,
			UserId:  userOId,
		})

		if err != nil {
			return primitive.NilObjectID, err
		}
	}

	if boardPinId != nil {
		err = s.store.DeleteBoardPinById(ctx, boardPinId.ID)

		if err != nil {
			return primitive.NilObjectID, err
		}
	}

	// Create board pin for the pin
	_, err = s.store.CreateBoardPin(ctx, &models.BoardPinCreation{
		BoardId: selectedBoardOId,
		PinId:   pinOId,
		UserId:  userOId,
	})

	if err != nil {
		return primitive.NilObjectID, err
	}

	return primitive.NilObjectID, nil
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

	userOId, _ := primitive.ObjectIDFromHex(userId)

	err = s.store.DeleteBoardPinById(ctx, pin.BoardPinId)

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
	boardId, err := s.store.GetBoardByUserId(ctx, p.UserId, AllPins)

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

func (s *service) GetPinById(ctx context.Context, id, userId string) (*models.PinModel, error) {

	oID, err := primitive.ObjectIDFromHex(id)
	userOID, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return nil, utils.ErrFailedToDecodeObjID
	}

	data, err := s.store.GetItem(ctx, map[string]interface{}{"_id": oID, "user_id": userOID})

	if err != nil {
		return nil, err
	}

	return data, nil
}
