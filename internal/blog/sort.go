package blog

import "strings"

func popularSort(a, b *Article) int {
	if a.PageViews > b.PageViews {
		return -1
	}
	if a.PageViews < b.PageViews {
		return 1
	}
	return 0
}

func dateSort(a, b *Article) int {
	return -strings.Compare(a.Date, b.Date)
}
