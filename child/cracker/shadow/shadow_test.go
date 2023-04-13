package shadow_test

import (
	"amogus/child/cracker/shadow"
	"testing"
)

func TestExtractSalt(t *testing.T) {
	line := "$6$ggQ/IFh.2tfUYstz$9EwU55NwWjn283GqawXvQ.3KQNxrxDB58Pwc3imX8hejt16ATxsXyyhu8LRrF/SONGzXPpMwbRhwoAm9963KF1"

	salt := shadow.ExtractSha512Salt(line)
	expectedSalt := "$6$ggQ/IFh.2tfUYstz$"

	if salt != expectedSalt {
		t.Errorf("expected %s, got %s", expectedSalt, salt)
	}
}

func TestSha512Crypt(t *testing.T) {
	line := "$6$ggQ/IFh.2tfUYstz$9EwU55NwWjn283GqawXvQ.3KQNxrxDB58Pwc3imX8hejt16ATxsXyyhu8LRrF/SONGzXPpMwbRhwoAm9963KF1"
	pwd := "Wiadro123"
	salt := "$6$ggQ/IFh.2tfUYstz$"
	crypter := *shadow.GetSaltySha512Crypter(salt)

	got := crypter.Crypt([]byte(pwd)).String()
	expected := line

	if expected != got {
		t.Errorf("expected %s, got %s", expected, got)
	}
}
