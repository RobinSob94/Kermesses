package requests

type AddStandRequest struct {
	StandIds []uint `json:"stand_ids" binding:"required,gt=0"`
}
