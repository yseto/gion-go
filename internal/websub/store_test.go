package websub

import "testing"

func TestStore(t *testing.T) {
	if topicURLs.insert("abc") != true {
		t.Error("invalid")
	}
	if topicURLs.insert("abc") != false {
		t.Error("invalid")
	}

	if topicURLs.contains("abc") != true {
		t.Error("invalid")
	}

	if topicURLs.delete("cde") != false {
		t.Error("invalid")
	}
	if topicURLs.contains("abc") != true {
		t.Error("invalid")
	}

	if topicURLs.delete("abc") != true {
		t.Error("invalid")
	}
	if topicURLs.contains("abc") != false {
		t.Error("invalid")
	}
}
