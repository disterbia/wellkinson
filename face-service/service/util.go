// /exercise-service/service/util.go
package service

import (
	"common/util"
	"exercise-service/dto"
)

func validateExercise(exercise dto.ExerciseRequest) error {
	if err := util.ValidateDate(exercise.PlanStartAt); err != nil {
		return err
	}
	if err := util.ValidateDate(exercise.PlanEndAt); err != nil {
		return err
	}
	if err := util.ValidateTime(exercise.ExerciseStartAt); err != nil {
		return err
	}
	if err := util.ValidateTime(exercise.ExerciseEndAt); err != nil {
		return err
	}
	return nil
}
