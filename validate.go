package raml

func (r *RAML) ValidateShapes() error {
	// Unwrap cache stores the mapping of original IDs to unwrapped shapes
	// to ensure the original references (aliases and links) match.
	unwrapCache := make(map[string]Shape)

	for _, frag := range r.fragmentsCache {
		switch f := frag.(type) {
		case *Library:
			for pair := f.AnnotationTypes.Oldest(); pair != nil; pair = pair.Next() {
				shape := pair.Value
				s := *shape
				if !s.Base().unwrapped {
					us, err := r.UnwrapShape(shape, make([]Shape, 0))
					if err != nil {
						return NewWrappedError("unwrap shape", err, s.Base().Location, WithPosition(&s.Base().Position))
					}
					unwrapCache[s.Base().Id] = s
					s = us
				}
				if err := s.Check(); err != nil {
					return NewWrappedError("check annotation type", err, s.Base().Location, WithPosition(&s.Base().Position))
				}
				if err := r.validateShapeCommons(s); err != nil {
					return NewWrappedError("validate shape commons", err, s.Base().Location, WithPosition(&s.Base().Position))
				}
			}
			for pair := f.Types.Oldest(); pair != nil; pair = pair.Next() {
				shape := pair.Value
				s := *shape
				if !s.Base().unwrapped {
					us, err := r.UnwrapShape(shape, make([]Shape, 0))
					if err != nil {
						return NewWrappedError("unwrap shape", err, s.Base().Location, WithPosition(&s.Base().Position))
					}
					unwrapCache[s.Base().Id] = s
					s = us
				}
				if err := s.Check(); err != nil {
					return NewWrappedError("check type", err, s.Base().Location, WithPosition(&s.Base().Position))
				}
				if err := r.validateShapeCommons(s); err != nil {
					return NewWrappedError("validate shape commons", err, s.Base().Location, WithPosition(&s.Base().Position))
				}
			}
		case *DataType:
			s := *f.Shape
			if !s.Base().unwrapped {
				us, err := r.UnwrapShape(f.Shape, make([]Shape, 0))
				if err != nil {
					return NewWrappedError("unwrap shape", err, s.Base().Location, WithPosition(&s.Base().Position))
				}
				unwrapCache[s.Base().Id] = s
				s = us
			}
			if err := s.Check(); err != nil {
				return NewWrappedError("check data type", err, s.Base().Location, WithPosition(&s.Base().Position))
			}
			if err := r.validateShapeCommons(s); err != nil {
				return NewWrappedError("validate shape commons", err, s.Base().Location, WithPosition(&s.Base().Position))
			}
		}
	}
	for _, item := range r.domainExtensions {
		db := *item.DefinedBy
		if !db.Base().unwrapped {
			us, ok := unwrapCache[db.Base().Id]
			if !ok {
				return NewError("unwrapped shape not found", db.Base().Location, WithPosition(&db.Base().Position))
			}
			db = us
		}
		if err := db.Validate(item.Extension.Value, "$"); err != nil {
			return NewWrappedError("check domain extension", err, item.Extension.Location, WithPosition(&item.Extension.Position))
		}
	}
	return nil
}

func (r *RAML) validateShapeCommons(s Shape) error {
	if err := r.validateShapeFacets(s); err != nil {
		return err
	}
	if err := r.validateExamples(s); err != nil {
		return err
	}

	switch s := s.(type) {
	case *ObjectShape:
		if s.Properties != nil {
			for pair := s.Properties.Oldest(); pair != nil; pair = pair.Next() {
				k, prop := pair.Key, pair.Value
				s := *prop.Shape
				if err := r.validateShapeCommons(s); err != nil {
					return NewWrappedError("validate property", err, s.Base().Location, WithPosition(&s.Base().Position), WithInfo("property", k))
				}
			}
			for pair := s.PatternProperties.Oldest(); pair != nil; pair = pair.Next() {
				k, prop := pair.Key, pair.Value
				s := *prop.Shape
				if err := r.validateShapeCommons(s); err != nil {
					return NewWrappedError("validate pattern property", err, s.Base().Location, WithPosition(&s.Base().Position), WithInfo("property", k))
				}
			}
		}
	case *ArrayShape:
		if s.Items != nil {
			if err := r.validateShapeCommons(*s.Items); err != nil {
				return NewWrappedError("validate items", err, s.Base().Location, WithPosition(&s.Base().Position))
			}
		}
	case *UnionShape:
		for _, item := range s.AnyOf {
			if err := r.validateShapeCommons(*item); err != nil {
				return NewWrappedError("validate union item", err, s.Base().Location, WithPosition(&s.Base().Position))
			}
		}
	}
	return nil
}

func (r *RAML) validateExamples(s Shape) error {
	base := s.Base()
	if base.Example != nil {
		if err := s.Validate(base.Example.Data.Value, "$"); err != nil {
			return NewWrappedError("validate example", err, base.Example.Location, WithPosition(&base.Example.Position))
		}
	}
	if base.Examples != nil {
		for pair := base.Examples.Map.Oldest(); pair != nil; pair = pair.Next() {
			ex := pair.Value
			if err := s.Validate(ex.Data.Value, "$"); err != nil {
				return NewWrappedError("validate example", err, ex.Location, WithPosition(&ex.Position))
			}
		}
	}
	if base.Default != nil {
		if err := s.Validate(base.Default.Value, "$"); err != nil {
			return NewWrappedError("validate default", err, base.Default.Location, WithPosition(&base.Default.Position))
		}
	}
	return nil
}

func (r *RAML) validateShapeFacets(s Shape) error {
	// TODO: Doesn't support multiple inheritance.
	base := s.Base()
	inherits := base.Inherits
	shapeFacetDefs := base.CustomShapeFacetDefinitions
	validationFacetDefs := make(map[string]Property)
	for {
		if len(inherits) == 0 {
			break
		}
		parent := *inherits[0]
		for pair := parent.Base().CustomShapeFacetDefinitions.Oldest(); pair != nil; pair = pair.Next() {
			f := pair.Value
			if _, ok := shapeFacetDefs.Get(f.Name); ok {
				base := (*f.Shape).Base()
				return NewError("duplicate custom facet", base.Location, WithPosition(&base.Position), WithInfo("facet", f.Name))
			}
			validationFacetDefs[f.Name] = f
		}
		inherits = parent.Base().Inherits
	}

	shapeFacets := base.CustomShapeFacets
	for k, facetDef := range validationFacetDefs {
		f, ok := shapeFacets.Get(k)
		if !ok {
			if facetDef.Required {
				return NewError("required custom facet is missing", base.Location, WithPosition(&base.Position), WithInfo("facet", k))
			}
			continue
		}
		if err := (*facetDef.Shape).Validate(f.Value, "$"); err != nil {
			return NewWrappedError("validate custom facet", err, f.Location, WithPosition(&f.Position), WithInfo("facet", k))
		}
	}

	for pair := shapeFacets.Oldest(); pair != nil; pair = pair.Next() {
		k, f := pair.Key, pair.Value
		if _, ok := validationFacetDefs[k]; !ok {
			return NewError("unknown facet", f.Location, WithPosition(&f.Position), WithInfo("facet", k))
		}
	}
	return nil
}
