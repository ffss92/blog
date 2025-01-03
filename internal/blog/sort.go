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
	res := strings.Compare(a.Date, b.Date)
	if res == 0 {
		return strings.Compare(a.Title, b.Title)
	}
	return -res
}
