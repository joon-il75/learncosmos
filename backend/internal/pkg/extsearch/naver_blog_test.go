package extsearch

import "testing"

func TestStripNaverHTML(t *testing.T) {
	got := stripNaverHTML("<b>코바늘</b> 뜨개질 <b>기초</b> &amp; 시작")
	want := "코바늘 뜨개질 기초 & 시작"
	if got != want {
		t.Fatalf("stripNaverHTML mismatch: got %q want %q", got, want)
	}
}

func TestNaverBlogExternalID(t *testing.T) {
	got := naverBlogExternalID("https://blog.naver.com/songi6743/224230750637")
	want := "blog.naver.com/songi6743/224230750637"
	if got != want {
		t.Fatalf("naverBlogExternalID mismatch: got %q want %q", got, want)
	}
}
