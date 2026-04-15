package scanner

// Filter defines criteria for including or excluding ports from scan results.
type Filter struct {
	// Protocols limits results to the given protocols (e.g. "tcp", "udp").
	// An empty slice means all protocols are accepted.
	Protocols []string

	// PortRange, if non-zero, restricts results to ports within [Min, Max].
	PortRange *RangeSpec
}

// RangeSpec defines an inclusive port number range.
type RangeSpec struct {
	Min uint16
	Max uint16
}

// NewFilter returns a Filter that accepts all ports.
func NewFilter() *Filter {
	return &Filter{}
}

// WithProtocols restricts the filter to the provided protocols.
func (f *Filter) WithProtocols(protos ...string) *Filter {
	f.Protocols = protos
	return f
}

// WithPortRange restricts the filter to ports in [min, max].
func (f *Filter) WithPortRange(min, max uint16) *Filter {
	f.PortRange = &RangeSpec{Min: min, Max: max}
	return f
}

// Accept returns true when the given PortState passes the filter criteria.
func (f *Filter) Accept(p PortState) bool {
	if len(f.Protocols) > 0 {
		matched := false
		for _, proto := range f.Protocols {
			if p.Protocol == proto {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	if f.PortRange != nil {
		if p.Port < f.PortRange.Min || p.Port > f.PortRange.Max {
			return false
		}
	}

	return true
}

// Apply returns only those PortStates that satisfy the filter.
func (f *Filter) Apply(ports []PortState) []PortState {
	out := make([]PortState, 0, len(ports))
	for _, p := range ports {
		if f.Accept(p) {
			out = append(out, p)
		}
	}
	return out
}
