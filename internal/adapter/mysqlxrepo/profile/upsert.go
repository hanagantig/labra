package profile

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"labra/internal/apperror"
	"labra/internal/entity"
)

func (r *Repository) Upsert(ctx context.Context, patient entity.Profile) (entity.Profile, error) {
	const query = `INSERT INTO profiles(uuid, user_id, creator_user_id, f_name, l_name, birth_date, gender)
					VALUES (:uuid, :user_id, :creator_user_id, :f_name, :l_name, :birth_date, :gender);`

	if patient.Uuid == uuid.Nil {
		patient.Uuid = uuid.New()
	}

	model := NewFromEntity(patient)
	res, err := r.GetConn(ctx).NamedExecContext(ctx, query, model)
	if err != nil {
		return entity.Profile{}, apperror.ToAppError(err)
	}

	patientID, err := res.LastInsertId()
	if err != nil {
		return entity.Profile{}, apperror.ToAppError(err)
	}

	if patientID == 0 {
		return entity.Profile{}, apperror.ToAppError(errors.New("failed to insert user"))
	}

	model.ID = int(patientID)

	return model.buildEntity(), nil
}
