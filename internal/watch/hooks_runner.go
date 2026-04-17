package watch

// HookedRunner wraps a ScanFunc and fires lifecycle hooks around each invocation.
type HookedRunner struct {
	inner    func() error
	registry *HookRegistry
}

// NewHookedRunner returns a HookedRunner that fires hooks around fn.
func NewHookedRunner(fn func() error, reg *HookRegistry) *HookedRunner {
	if reg == nil {
		reg = NewHookRegistry()
	}
	return &HookedRunner{inner: fn, registry: reg}
}

// Run fires BeforeScan, calls the inner function, then fires AfterScan or OnError.
func (h *HookedRunner) Run() error {
	h.registry.Fire(HookBeforeScan, nil)
	err := h.inner()
	if err != nil {
		h.registry.Fire(HookOnError, map[string]any{"err": err.Error()})
		return err
	}
	h.registry.Fire(HookAfterScan, nil)
	return nil
}
