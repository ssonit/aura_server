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

func (s *service) ListSoftDeletedPins(ctx context.Context, userId string) ([]models.PinModel, error) {
	userOId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return nil, utils.ErrFailedToDecodeObjID
	}

	data, err := s.store.ListSoftDeletedPins(ctx, userOId)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *service) RestorePin(ctx context.Context, id, userId string) error {
	oID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return utils.ErrFailedToDecodeObjID

	}

	data, err := s.store.GetItem(ctx, map[string]interface{}{"_id": oID})

	if err != nil {
		return err
	}

	if data.UserId.Hex() != userId {
		return utils.ErrUserNotPermitted

	}

	err = s.store.RestorePin(ctx, oID)

	if err != nil {
		return err
	}

	return nil

}

func (s *service) SoftDeletePin(ctx context.Context, id, userId string) error {
	oID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return utils.ErrFailedToDecodeObjID
	}

	data, err := s.store.GetItem(ctx, map[string]interface{}{"_id": oID})

	if err != nil {
		return err
	}

	if data.UserId.Hex() != userId {
		return utils.ErrUserNotPermitted
	}

	err = s.store.SoftDeletePin(ctx, oID)

	if err != nil {
		return err
	}

	return nil

}

func (s *service) UnSaveBoardPin(ctx context.Context, p *models.BoardPinUnSave) error {
	pinOId, err := primitive.ObjectIDFromHex(p.PinId)
	userOId, err := primitive.ObjectIDFromHex(p.UserId)
	boardPinId, err := primitive.ObjectIDFromHex(p.BoardPinId)

	if err != nil {
		return utils.ErrFailedToDecodeObjID
	}

	// Board of all pins by user id
	boardIdAllPins, err := s.store.GetBoardByUserId(ctx, userOId, AllPins)

	if err != nil {
		return err
	}

	if boardIdAllPins.IsZero() {
		return utils.ErrBoardNotFound
	}

	// delete board pin for all pins
	err = s.store.DeleteBoardPin(ctx, &models.BoardPinFilter{
		PinId:   pinOId,
		BoardId: boardIdAllPins,
		UserId:  userOId,
	})

	if err != nil {
		return err
	}

	// delete board pin for the pin
	err = s.store.DeleteBoardPinById(ctx, boardPinId)

	if err != nil {
		return err
	}

	return nil

}

func (s *service) DeleteComment(ctx context.Context, commentId, userId string) error {
	commentOID, err := primitive.ObjectIDFromHex(commentId)

	if err != nil {
		return utils.ErrFailedToDecodeObjID
	}

	comment, err := s.store.GetCommentById(ctx, commentOID)

	if err != nil {
		return err
	}

	if comment.UserId.Hex() != userId {
		return utils.ErrUserNotPermitted
	}

	err = s.store.DeleteComment(ctx, commentOID)

	if err != nil {
		return err
	}

	return nil

}

func (s *service) ListCommentsByPinId(ctx context.Context, pinId string, paging *common.Paging) ([]models.CommentModel, error) {
	oID, err := primitive.ObjectIDFromHex(pinId)

	if err != nil {
		return nil, utils.ErrFailedToDecodeObjID
	}

	data, err := s.store.ListCommentsByPinId(ctx, oID, paging)

	if err != nil {
		return nil, err
	}

	return data, nil

}

func (s *service) CreateComment(ctx context.Context, p *models.CommentCreation) (primitive.ObjectID, error) {

	userOId, err := primitive.ObjectIDFromHex(p.UserId)
	pinOId, err := primitive.ObjectIDFromHex(p.PinId)

	if err != nil {
		return primitive.NilObjectID, utils.ErrFailedToDecodeObjID
	}

	data, err := s.store.CreateComment(ctx, &models.CommentCreationStore{
		PinId:   pinOId,
		UserId:  userOId,
		Content: p.Content,
	})

	if err != nil {
		return primitive.NilObjectID, err
	}

	return data, nil

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

	tagIDs, err := s.store.CheckAndCreateTags(ctx, p.Tags)

	if err != nil {
		return primitive.NilObjectID, err
	}

	data, err := s.store.Create(ctx, p, tagIDs)

	if err != nil {
		return primitive.NilObjectID, err
	}

	if data.IsZero() {
		return primitive.NilObjectID, utils.ErrDataIsZero
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
