package api

import "lavazares/content"

type lessonManager content.DefaultLessonManager

func InitAPI() {
	lessonManager := content.NewDefaultLessonManager()
}
