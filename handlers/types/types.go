package types

import "errors"

type UpdateCatRequest struct {
	Salary int64 `json:"salary"`
}

func (r *UpdateCatRequest) Validate() error {
	if r.Salary == 0 {
		return errors.New("salary should be more then 0")
	}
	return nil
}

type UpdateMissionRequest struct {
	Completed bool `json:"completed"`
}

type UpdateTargetRequest struct {
	Completed bool   `json:"completed"`
	Notes     string `json:"notes"`
}

type AssignToMissionRequest struct {
	CatID uint `json:"cat_id" `
}
