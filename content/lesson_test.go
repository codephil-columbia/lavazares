package content

import (
	"sort"
	"testing"
)

func TestLessonsSort(t *testing.T) {
	cases := []struct {
		name     string
		lst      lessons
		expected lessons
	}{
		{
			"Should sort list of Lessons correctly",
			lessons{
				&Lesson{
					LessonName: "Lesson 3: Third Lesson",
				},
				&Lesson{
					LessonName: "Chapter 1 Test",
				},
				&Lesson{
					LessonName: "Lesson 1: First Lesson",
				},
			},
			lessons{
				&Lesson{
					LessonName: "Lesson 1: First Lesson",
				},
				&Lesson{
					LessonName: "Lesson 3: Third Lesson",
				},
				&Lesson{
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
					t.FailNow()
				}
			}
		})
	}
}
