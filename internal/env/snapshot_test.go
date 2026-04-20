package env

import (
	"testing"
)

func baseMap() map[string]string {
	return map[string]string{
		"HOME":  "/home/user",
		"PATH":  "/usr/bin",
		"TOKEN": "abc123",
	}
}

func TestSnapshot_Get(t *testing.T) {
	s := NewSnapshotFromMap(baseMap())

	v, ok := s.Get("HOME")
	if !ok || v != "/home/user" {
		t.Fatalf("expected /home/user, got %q (ok=%v)", v, ok)
	}

	_, ok = s.Get("MISSING")
	if ok {
		t.Fatal("expected MISSING key to be absent")
	}
}

func TestSnapshot_Keys(t *testing.T) {
	s := NewSnapshotFromMap(baseMap())
	keys := s.Keys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
}

func TestSnapshot_Diff_DetectsAddedAndChanged(t *testing.T) {
	before := NewSnapshotFromMap(map[string]string{
		"A": "1",
		"B": "2",
	})
	after := NewSnapshotFromMap(map[string]string{
		"A": "1",
		"B": "changed",
		"C": "new",
	})

	diff := before.Diff(after)

	if diff["B"] != "changed" {
		t.Errorf("expected B=changed, got %q", diff["B"])
	}
	if diff["C"] != "new" {
		t.Errorf("expected C=new, got %q", diff["C"])
	}
	if _, ok := diff["A"]; ok {
		t.Error("A should not appear in diff (unchanged)")
	}
}

func TestSnapshot_Diff_EmptyWhenIdentical(t *testing.T) {
	s := NewSnapshotFromMap(baseMap())
	diff := s.Diff(NewSnapshotFromMap(baseMap()))
	if len(diff) != 0 {
		t.Fatalf("expected empty diff, got %v", diff)
	}
}

func TestSnapshot_Environ_RoundTrip(t *testing.T) {
	orig := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := NewSnapshotFromMap(orig)
	env := s.Environ()

	if len(env) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(env))
	}
	rebuilt := make(map[string]string)
	for _, e := range env {
		var k, v string
		for i, c := range e {
			if c == '=' {
				k = e[:i]
				v = e[i+1:]
				break
			}
		}
		rebuilt[k] = v
	}
	for k, v := range orig {
		if rebuilt[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, rebuilt[k])
		}
	}
}

func TestTakeSnapshot_CapturesOS(t *testing.T) {
	s := TakeSnapshot()
	if len(s.Keys()) == 0 {
		t.Fatal("expected at least one env var from os.Environ")
	}
}
