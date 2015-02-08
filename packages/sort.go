package packages

import "errors"

func (ps *Packages) Sort() error {
	// this is a bruteforce topological sorting, because I'm lazy
	order := []string{}
	defer func() { ps.Order = order }()

	sorted := make(map[string]bool)

	// try iterate 100x
	for pass := 0; pass < 100; pass++ {
		changes := false

	nextpkg:
		for _, p := range ps.List {
			if sorted[p.Name] {
				continue
			}

			for _, dep := range p.Deps {
				if !sorted[dep] {
					continue nextpkg
				}
			}

			changes = true
			sorted[p.Name] = true
			order = append(order, p.Name)
		}

		if len(sorted) == len(ps.List) {
			return nil
		}

		if !changes {
			return errors.New("packages contain a circular dependency")
		}
	}

	return errors.New("unable to find order in 100 passes")
}
