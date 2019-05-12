package content

import (
	"testing"
)

type mockChapterManager struct{}
type mockLessonManager struct{}

var (
	testLessons = []*Lesson{
		&Lesson{
			LessonName: "Lesson 1",
			LessonID:   "2",
			ChapterID:  "1",
		}, &Lesson{
			LessonName: "Lesson 2",
			LessonID:   "3",
			ChapterID:  "1",
		}, &Lesson{
			LessonName: "Lesson 1",
			LessonID:   "6",
			ChapterID:  "4",
		},
	}
	testLesson = Lesson{
		LessonName: "Lesson 1",
		LessonID:   "2",
		ChapterID:  "1",
	}

	testChapters = []*Chapter{
		&testFirstChapter,
		&testLastChapter,
	}
	testFirstChapter = Chapter{
		ChapterID:   "1",
		ChapterName: "Chapter 1",
	}
	testLastChapter = Chapter{
		ChapterID:   "4",
		ChapterName: "Chapter 2",
	}
)

func TestGetLessonsInChapter(t *testing.T) {
	contentManager := initWithMocks()

	cases := []struct {
		chapterID   string
		expectedErr bool
	}{
		// should be able to return lessons for a valid chapter
		// and correctly filter lessons belonging to the given
		// chapter id
		{
			"1",
			false,
		},
	}

	for _, tc := range cases {
		lessons, err := contentManager.GetLessonsInChapter(tc.chapterID)
		if err != nil {
			if !tc.expectedErr {
				t.Errorf("Unexpected error: [%v]", err)
				continue
			}
		}

		for _, lesson := range lessons {
			if lesson.ChapterID != tc.chapterID {
				t.Errorf(`Found lesson with incorrect chapter id, 
				wanted chapter id [%v], got [%v]`, tc.chapterID, lesson.ChapterID)
				continue
			}
		}
	}
}

func TestGetNextChapter(t *testing.T) {
	contentManger := initWithMocks()

	cases := []struct {
		name           string
		chapterID      string
		expectedNextID string
		expectedErr    bool
	}{
		{
			"Should get next chapter for chapter that is not the last one",
			"1",
			"4",
			false,
		},
		{
			"Should say that all chapters are completed if its the last chapter",
			"4",
			"",
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			chapter, err := contentManger.getNextChapter(tc.chapterID)
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("Unexpected error: [%v]", err)
				} else {
					if err != errCompletedAllChapters {
						t.Errorf("Expected [%v], got [%v]",
							errCompletedAllChapters, err,
						)
					}
				}
			} else {
				if chapter.ChapterID != tc.expectedNextID {
					t.Errorf("Expected Chapter with id: [%v], got: [%v]",
						chapter.ChapterID, tc.expectedNextID,
					)
				}
			}
		})
	}
}

func TestGetNextLesson(t *testing.T) {
	contentManager := initWithMocks()
	cases := []struct {
		name        string
		lessonID    string
		expectedID  string
		expectedErr bool
	}{
		{
			`Should get next lesson thats in the same chapter if its not the last one
			in its chapter`,
			"2",
			"3",
			false,
		},
		{
			"Should return first lesson in next chapter if its the last lesson in its chapter",
			"3",
			"6",
			false,
		},
		{
			"Should return correct error if it is the last lesson",
			"6",
			"",
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nextLesson, err := contentManager.GetNextLesson(tc.lessonID)
			if err != nil {
				if !tc.expectedErr || err != errCompletedAllLessons {
					t.Errorf("Unexpected error: [%v]", err)
				}
			} else {
				if nextLesson.LessonID != tc.expectedID {
					t.Errorf("Expected Lesson with id:[%v], got:[%v]",
						tc.expectedID, nextLesson.LessonID,
					)
				}
			}
		})
	}
}

func initWithMocks() *DefaultContentManager {
	return &DefaultContentManager{
		chapterManager: mockChapterManager{},
		lessonManager:  mockLessonManager{},
	}
}

// ChapterManager Mocks
func (manager mockChapterManager) GetChapters() ([]*Chapter, error) {
	return testChapters, nil
}

func (manager mockChapterManager) GetChapter(id string) (*Chapter, error) {
	return &testFirstChapter, nil
}

// LessonManager Mocks

func (manager mockLessonManager) GetLessons() ([]*Lesson, error) {
	return testLessons, nil
}

func (manager mockLessonManager) GetLesson(id string) (*Lesson, error) {
	switch id {
	case "2":
		return testLessons[0], nil
	case "3":
		return testLessons[1], nil
	case "6":
		return testLessons[2], nil
	default:
		return nil, nil
	}
}
