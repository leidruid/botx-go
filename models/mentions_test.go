package models

import "testing"

func TestFindAndReplaceEmbedMentions(t *testing.T) {
	body := "hi <embed_mention>user:123e4567-e89b-12d3-a456-426655440000:Bob</embed_mention>"
	newBody, mentions := FindAndReplaceEmbedMentions(body)
	if newBody == body {
		t.Fatalf("expected body to change")
	}
	if len(mentions) != 1 {
		t.Fatalf("expected 1 mention")
	}
	if mentions[0].MentionType != string(MentionUser) {
		t.Fatalf("unexpected mention type")
	}
}
