package tutorial

import "lavazares/models"

type TutorialManager interface {
	AddCompletedTutorial(u models.User, l models.Lesson) (e error)
}
