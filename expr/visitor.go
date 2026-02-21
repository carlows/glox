package expr

type Visitor interface {
	VisitBinary(*Binary) any
	VisitGrouping(*Grouping) any
	VisitLiteral(*Literal) any
	VisitUnary(*Unary) any
}
