package env

import (
	"errors"
	"testing"
)

func TestChain_EmptySteps_ReturnsInputCopy(t *testing.T) {
	src := map[string]string{"KEY": "value"}
	chain := NewChain()
	out, err := chain.Run(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", out["KEY"])
	}
}

func TestChain_StepsRunInOrder(t *testing.T) {
	var order []int
	makeStep := func(n int) ChainStep {
		return func(m map[string]string) (map[string]string, error) {
			order = append(order, n)
			return m, nil
		}
	}
	chain := NewChain(makeStep(1), makeStep(2), makeStep(3))
	_, err := chain.Run(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 3 || order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Errorf("unexpected order: %v", order)
	}
}

func TestChain_ErrorAbortsEarly(t *testing.T) {
	called := false
	failStep := func(m map[string]string) (map[string]string, error) {
		return nil, errors.New("step failed")
	}
	neverStep := func(m map[string]string) (map[string]string, error) {
		called = true
		return m, nil
	}
	chain := NewChain(failStep, neverStep)
	_, err := chain.Run(map[string]string{"A": "1"})
	if err == nil {
		t.Fatal("expected error")
	}
	if called {
		t.Error("second step should not have been called")
	}
}

func TestChain_DoesNotMutateInput(t *testing.T) {
	src := map[string]string{"K": "v"}
	chain := NewChain(func(m map[string]string) (map[string]string, error) {
		m["ADDED"] = "yes"
		return m, nil
	})
	_, err := chain.Run(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := src["ADDED"]; ok {
		t.Error("original map was mutated")
	}
}

func TestChain_Append_CreatesNewChain(t *testing.T) {
	base := NewChain()
	extended := base.Append(func(m map[string]string) (map[string]string, error) {
		m["EXTRA"] = "1"
		return m, nil
	})
	out, err := extended.Run(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["EXTRA"] != "1" {
		t.Errorf("expected EXTRA=1, got %q", out["EXTRA"])
	}
	if len(base.steps) != 0 {
		t.Error("base chain was mutated by Append")
	}
}

func TestWrapTransformer_IntegratesInChain(t *testing.T) {
	t.Parallel()
	tr := NewTransformer(TrimSpaceTransform)
	chain := NewChain(WrapTransformer(tr))
	out, err := chain.Run(map[string]string{"KEY": "  hello  "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["KEY"])
	}
}

func TestWrapValidation_PassesOnValid(t *testing.T) {
	chain := NewChain(WrapValidation(RequireKeys("A")))
	out, err := chain.Run(map[string]string{"A": "present"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "present" {
		t.Errorf("expected 'present', got %q", out["A"])
	}
}

func TestWrapValidation_FailsOnMissingKey(t *testing.T) {
	chain := NewChain(WrapValidation(RequireKeys("MUST_EXIST")))
	_, err := chain.Run(map[string]string{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
