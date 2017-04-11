package epub

import (
	"fmt"
	"strings"

	"github.com/ssor/epubgo/raw"
)

func NewNavigationPoint(src *raw.NavPoint, level, IndexInList int, tagPre string) *NavigationPoint {
	point := &NavigationPoint{
		Title: src.Title(),
		Level: level,
		URL:   src.URL(),
	}

	if sharpIndex := strings.Index(src.URL(), "#"); sharpIndex > 0 {
		point.URL = src.URL()[:sharpIndex]
	}

	if len(tagPre) <= 0 {
		point.Tag = fmt.Sprintf("%d", IndexInList)
	} else {
		point.Tag = fmt.Sprintf("%s.%d", tagPre, IndexInList)
	}

	return point
}

type NavigationPoint struct {
	Title               string `json:"title"`
	CharactorCountSelf  int    `json:"charactor_count"`
	CharactorCountTotal int    `json:"total_charactor_count"`
	Level               int    `json:"level"`
	URL                 string `json:"url"`
	Tag                 string `json:"tag"` //like 1.1 or 1.2.1
	Text                string `json:"text"`
}

type NavigationPointArray []*NavigationPoint

func (npa NavigationPointArray) SumSubLevelCharactorCount(tag string) int {
	if strings.HasSuffix(tag, ".") == false {
		tag = tag + "."
	}
	count := 0
	for _, np := range npa {
		if strings.HasPrefix(np.Tag, tag) == true {
			count += np.CharactorCountSelf
		}
	}

	return count
}

func (npa NavigationPointArray) Each(f func(*NavigationPoint) error) error {
	if f == nil {
		return nil
	}
	for _, np := range npa {
		err := f(np)
		if err != nil {
			return err
		}
	}
	return nil
}
