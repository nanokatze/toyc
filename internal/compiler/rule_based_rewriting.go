package compiler

type RewriteRule struct {
	Name    string
	Pattern *Pattern
	Replace func(*Sea, *Value) *Value // TODO: also sneak in user data somehow
}
