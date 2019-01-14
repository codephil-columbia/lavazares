package content

import (
	"reflect"
	"sort"
	"testing"
)

func TestChaptersSort(t *testing.T) {
	cases := []struct {
		name string
		chapters chapters
		expected chapters
	}{
		{
			"Should be able to sort list of Chapters correctly",
			chapters{
				&Chapter{
					ChapterName: "Chapter 1",
				},
				&Chapter{
					ChapterName: "Chapter 6",
				},
				&Chapter{
					ChapterName: "Chapter 3",
				},
			},
			chapters{
				&Chapter{
					ChapterName: "Chapter 1",
				},
				&Chapter{
					ChapterName: "Chapter 3",
				},
				&Chapter{
					ChapterName: "Chapter 6",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sort.Sort(tc.chapters)
			if !reflect.DeepEqual(tc.chapters, tc.expected) {
				t.Errorf("Chapters where sorted incorrectly, expected [%v]\n, got [%v]",
					tc.expected, tc.chapters,
				)
			}
		})
	}
}
