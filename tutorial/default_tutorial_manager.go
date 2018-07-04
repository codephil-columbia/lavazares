package tutorial

import (
	"lavazares/models"
)

type DefaultTutorialManager struct {
}

func (d DefaultTutorialManager) AddCompletedTutorial(l models.LessonsComplete) error {
	err := models.UserCompletedLesson(l)
	if err != nil {
		return err
	}

	return nil
}
