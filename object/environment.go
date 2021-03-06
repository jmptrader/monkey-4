package object

import (
	"fmt"
	"os"
)

// Environment stores our functions, variables, constants, etc.
type Environment struct {
	// store holds variables, including functions.
	store map[string]Object

	// outer holds any parent environment.  Our env. allows
	// nesting to implement scope.
	outer *Environment
}

// NewEnvironment creates new environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnvironment create new environment by outer parameter
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get object by name
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set object by name
func (e *Environment) Set(name string, val Object) Object {

	//
	// If a variable is constant then we don't allow it to be changed.
	//
	// But constants are scoped, they are not global, so we only need
	// to look in the current scope - not any parent.
	//
	// i.e. The parent-scope might have a constant-value, but
	// we just don't care.  Consider the following code:
	//
	//    const a = 3.13;
	//    function foo() {
	//       let a = 1976;
	//    };
	//
	// The variable inside the function _should_ not be constant.
	//
	cur := e.store[name]
	if cur != nil && (cur.Constant()) {
		fmt.Printf("Attempting to modify '%s' denied; it was defined as a constant.\n", name)
		os.Exit(3)
	}

	//
	// Store the (updated) value.
	//
	e.store[name] = val
	return val
}

// SetConst sets the value of a constant by name.
func (e *Environment) SetConst(name string, val Object) Object {

	// store the value
	e.store[name] = val

	// Mark this as a constant
	e.store[name].SetConstant(true)
	return val
}
