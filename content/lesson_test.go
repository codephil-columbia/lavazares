package content

import (
	"reflect"
	"sort"
	"testing"
)

func TestLessonsSort(t *testing.T) {
	cases := []struct {
		name     string
		lst      byLessonName
		expected byLessonName
	}{
		{
			"Should sort list of Lessons correctly",
			byLessonName{
				Lesson{
					LessonName: "Lesson 3: Third Lesson",
				},
				Lesson{
					LessonName: "Chapter 1 Test",
				},
				Lesson{
					LessonName: "Lesson 1: First Lesson",
				},
			},
			byLessonName{
				Lesson{
					LessonName: "Lesson 1: First Lesson",
				},
				Lesson{
					LessonName: "Lesson 3: Third Lesson",
				},
				Lesson{
					LessonName: "Chapter 1 Test",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sort.Sort(tc.lst)
			for i := range tc.lst {
				if tc.lst[i].LessonName != tc.expected[i].LessonName {
					t.Errorf("List was not sorted properly, got [%v] expected [%v]", tc.lst, tc.expected)
				}
			}
		})
	}
}

var (
	readOnlyLesson1 = Lesson{
		LessonName: "Lesson 1: First Lesson",
		ChapterID:  "1",
	}

	readOnlyLesson2 = Lesson{
		LessonName: "Lesson 3: Third Lesson",
		ChapterID:  "1",
	}
	readOnlyLesson3 = Lesson{
		LessonName: "Chapter 1 Test",
		ChapterID:  "1",
	}
	readOnlyLesson4 = Lesson{
		LessonName: "Chapter 1 Test",
		ChapterID:  "2",
	}
)

type lessonStoreMock struct{}

func (mock *lessonStoreMock) Query(id string) (*Lesson, error) {
	return &readOnlyLesson1, nil
}

func (mock *lessonStoreMock) QueryAll() ([]*Lesson, error) {
	return []*Lesson{
		&readOnlyLesson1,
		&readOnlyLesson2,
		&readOnlyLesson3,
		&readOnlyLesson4,
	}, nil
}

func TestGetLessonsInChapter(t *testing.T) {
	mockStore := lessonStoreMock{}
	cases := []struct {
		name      string
		store     lessonStoreMock
		chapterID string
		expected  []Lesson
	}{
		{
			"Should be able to query lessons by chapter",
			mockStore,
			"1",
			[]Lesson{
				readOnlyLesson1,
				readOnlyLesson2,
				readOnlyLesson3,
			},
		},
	}

	for _, tc := range cases {
		manager := newLessonManager(&tc.store)
		l, err := manager.GetLessonsInChapter(tc.chapterID)
		if err != nil {
			t.Errorf("%v", err)
		}

		if !reflect.DeepEqual(l, tc.expected) {
			t.Errorf("Expected %v, got %v", tc.expected, l)
		}
	}
}
