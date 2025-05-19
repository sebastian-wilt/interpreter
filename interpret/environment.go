package interpret

type Environment struct {
	values map[string]Value
	types  map[string]Type
	parent *Environment
}

var outer *Environment

func NewEnvironment() *Environment {
	if outer == nil {
		outer = &Environment{
			values: map[string]Value{},
			types:  getInbuilts(),
			parent: nil,
		}
	}

	return outer
}

func NewEnvironmentWithParent(parent *Environment) *Environment {
	env := &Environment{
		values: map[string]Value{},
		types:  map[string]Type{},
		parent: parent,
	}

	return env
}

// Define binding from name to value in current environment
func (env *Environment) define(name string, value Value) {
	env.values[name] = value
}

// Assign value to binding with identifier name in closest
// enclosing scope where name is defined
func (env *Environment) assign(name string, value Value) {
	_, ok := env.values[name]

	if ok {
		env.values[name] = value
		return
	}

	env.parent.assign(name, value)
}

// Lookup binding with name
// Look through all enclosing scopes
func (env *Environment) lookup(name string) Value {
	if val, ok := env.values[name]; ok {
		return val
	}

	// Look outwards in lexical scope if not defined here
	if env.parent != nil {
		return env.parent.lookup(name)
	}

	return nil
}

// Lookup type
func (env *Environment) lookupType(name string) Type {
	if t, ok := env.types[name]; ok {
		return t
	}

	if env.parent != nil {
		return env.parent.lookupType(name)
	}

	return nil
}
