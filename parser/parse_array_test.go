package parser_test

import (
	"github.com/calico32/goose/ast"
	"github.com/calico32/goose/token"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("parseArrayLitOrInitializer", func() {
	It("should parse array literal", func() {
		src := `[1, 2, 3]`

		p := prepareParser(src)
		expr := p.ParseExpr()

		Expect(expr).To(BeAssignableToTypeOf(&ast.ArrayLiteral{}))
		arr := expr.(*ast.ArrayLiteral)

		Expect(arr.List).To(HaveLen(3))

		Expect(arr.Opening).ToNot(And(BeNil(), Equal(token.NoPos)))
		Expect(arr.Closing).ToNot(And(BeNil(), Equal(token.NoPos)))

		Expect(arr.List[0]).To(BeAssignableToTypeOf(&ast.Literal{}))
		Expect(arr.List[1]).To(BeAssignableToTypeOf(&ast.Literal{}))
		Expect(arr.List[2]).To(BeAssignableToTypeOf(&ast.Literal{}))
		Expect(arr.List[0].(*ast.Literal).Value).To(Equal("1"))
		Expect(arr.List[1].(*ast.Literal).Value).To(Equal("2"))
		Expect(arr.List[2].(*ast.Literal).Value).To(Equal("3"))
	})

	It("should parse sparse array literal", func() {
		src := `[1, 2, , 4, 5,]`

		p := prepareParser(src)
		expr := p.ParseExpr()

		Expect(expr).To(BeAssignableToTypeOf(&ast.ArrayLiteral{}))
		arr := expr.(*ast.ArrayLiteral)

		Expect(arr.List).To(HaveLen(5))

		Expect(arr.Opening).ToNot(And(BeNil(), Equal(token.NoPos)))
		Expect(arr.Closing).ToNot(And(BeNil(), Equal(token.NoPos)))

		Expect(arr.List[0]).To(BeAssignableToTypeOf(&ast.Literal{}))
		Expect(arr.List[1]).To(BeAssignableToTypeOf(&ast.Literal{}))
		Expect(arr.List[2]).To(BeAssignableToTypeOf(&ast.Literal{}))
		Expect(arr.List[3]).To(BeAssignableToTypeOf(&ast.Literal{}))
		Expect(arr.List[4]).To(BeAssignableToTypeOf(&ast.Literal{}))

		Expect(arr.List[0].(*ast.Literal).Value).To(Equal("1"))
		Expect(arr.List[1].(*ast.Literal).Value).To(Equal("2"))
		Expect(arr.List[2].(*ast.Literal).Kind).To(Equal(token.Null))
		Expect(arr.List[3].(*ast.Literal).Value).To(Equal("4"))
		Expect(arr.List[4].(*ast.Literal).Value).To(Equal("5"))
	})

	It("should parse array initializer", func() {
		src := `["foo"; 10]`

		p := prepareParser(src)
		expr := p.ParseExpr()

		Expect(expr).To(BeAssignableToTypeOf(&ast.ArrayInitializer{}))
		arr := expr.(*ast.ArrayInitializer)

		Expect(arr.Opening).ToNot(And(BeNil(), Equal(token.NoPos)))
		Expect(arr.Closing).ToNot(And(BeNil(), Equal(token.NoPos)))
		Expect(arr.Semi).ToNot(And(BeNil(), Equal(token.NoPos)))

		Expect(arr.Count).To(BeAssignableToTypeOf(&ast.Literal{}))
		Expect(arr.Count.(*ast.Literal).Value).To(Equal("10"))

		Expect(arr.Value).To(BeAssignableToTypeOf(&ast.StringLiteral{}))
		Expect(arr.Value.(*ast.StringLiteral).String()).To(Equal("foo"))
	})
})
