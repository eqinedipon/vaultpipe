package env

// Chain applies a sequence of transform functions to a map, threading the
// output of each step into the input of the next. Each step may return an
// error, which aborts the chain immediately.
//
// Chain is useful when you want to compose sanitize, coerce, truncate, and
// validate passes into a single reusable pipeline without intermediate
// variables.
type Chain struct {
	steps []ChainStep
}

// ChainStep is a single transformation in a Chain.
type ChainStep func(map[string]string) (map[string]string, error)

// NewChain creates a Chain with the provided steps.
func NewChain(steps ...ChainStep) *Chain {
	return &Chain{steps: steps}
}

// Run executes each step in order, returning the final map or the first error
// encountered.
func (c *Chain) Run(src map[string]string) (map[string]string, error) {
	current := copyMap(src)
	for _, step := range c.steps {
		var err error
		current, err = step(current)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}

// Append returns a new Chain with the given steps added at the end.
func (c *Chain) Append(steps ...ChainStep) *Chain {
	newSteps := make([]ChainStep, len(c.steps)+len(steps))
	copy(newSteps, c.steps)
	copy(newSteps[len(c.steps):], steps)
	return &Chain{steps: newSteps}
}

// WrapTransformer adapts a *Transformer into a ChainStep.
func WrapTransformer(t *Transformer) ChainStep {
	return func(m map[string]string) (map[string]string, error) {
		return t.Apply(m)
	}
}

// WrapValidation adapts a Validate call into a ChainStep that passes the map
// through unchanged when validation succeeds.
func WrapValidation(opts ...ValidationOption) ChainStep {
	return func(m map[string]string) (map[string]string, error) {
		if err := Validate(m, opts...); err != nil {
			return nil, err
		}
		return m, nil
	}
}

func copyMap(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
