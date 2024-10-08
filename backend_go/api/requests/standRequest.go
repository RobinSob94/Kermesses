package requests

type StandRequest struct {
	Name         string `json:"name" binding:"required"`
	Type         string `son:"type" binding:"required"`
	JetonsRequis uint   `json:"jetons_requis" binding:"required"`
}
