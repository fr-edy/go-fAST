package parser

import (
	"errors"
	"fmt"

	"github.com/t14raptor/go-fast/token"
)

const errUnexpectedToken = "Unexpected token %v"

var (
	errUnexpectedEndOfInput                  = errors.New("Unexpected end of input")
	errUnexpectedIdentifier                  = errors.New("Unexpected identifier")
	errUnexpectedReservedWord                = errors.New("Unexpected reserved word")
	errKeywordEscapedChars                   = errors.New("Keyword must not contain escaped characters")
	errUnexpectedNumber                      = errors.New("Unexpected number")
	errUnexpectedString                      = errors.New("Unexpected string")
	errMissingCatchOrFinally                 = errors.New("Missing catch or finally after try")
	errStaticPrototype                       = errors.New("Classes may not have a static property named 'prototype'")
	errConstructorAccessor                   = errors.New("Class constructor may not be an accessor")
	errConstructorAsync                      = errors.New("Class constructor may not be an async method")
	errConstructorGenerator                  = errors.New("Class constructor may not be a generator")
	errConstructorPrivate                    = errors.New("Class constructor may not be a private method")
	errFieldConstructor                      = errors.New("Classes may not have a field named 'constructor'")
	errIllegalReturn                         = errors.New("Illegal return statement")
	errIllegalNewlineAfterThrow              = errors.New("Illegal newline after throw")
	errDuplicateDefault                      = errors.New("Already saw a default in switch")
	errForInInitializer                      = errors.New("for-in loop variable declaration may not have an initializer")
	errInvalidLHSForIn                       = errors.New("Invalid left-hand side in for-in or for-of")
	errMissingDestructuringInit              = errors.New("Missing initializer in destructuring declaration")
	errLexicalSingleStatement                = errors.New("Lexical declaration cannot appear in a single-statement context")
	errIllegalBreak                          = errors.New("Illegal break statement")
	errIllegalContinue                       = errors.New("Illegal continue statement")
	errSuperUnexpected                       = errors.New("'super' keyword unexpected here")
	errRestParamLast                         = errors.New("Rest parameter must be last formal parameter")
	errGetterParams                          = errors.New("Getter must not have any formal parameters.")
	errSetterParams                          = errors.New("Setter must have exactly one formal parameter.")
	errInvalidTemplateLiteralOnOptionalChain = errors.New("Invalid template literal on optional chain")
	errInvalidLHSAssignment                  = errors.New("Invalid left-hand side in assignment")
	errIllegalAwaitInParams                  = errors.New("Illegal await-expression in formal parameters of async function")
	errMixedLogicalCoalesce                  = errors.New("Logical expressions and coalesce expressions cannot be mixed. Wrap either by parentheses")
	errExponentUnary                         = errors.New("Unary operator used immediately before exponentiation expression. Parenthesis must be used to disambiguate operator precedence")
	errMalformedArrowParams                  = errors.New("Malformed arrow function parameter list")
	errYieldInParams                         = errors.New("Yield expression not allowed in formal parameter")
	errCommaNotAllowed                       = errors.New("Comma is not allowed here")
	errRestElementLast                       = errors.New("Rest element must be last element")
	errInvalidDestructuringBinding           = errors.New("Invalid destructuring binding target")
	errInvalidDestructuringAssignment        = errors.New("Invalid destructuring assignment target")
	errInvalidBindingRest                    = errors.New("Invalid binding rest")
)

func (p *parser) error(err error) error {
	p.errs = append(p.errs, err)
	return err
}

func (p *parser) errorf(msg string, msgValues ...any) error {
	err := fmt.Errorf(msg, msgValues...)
	p.errs = append(p.errs, err)
	return err
}

func (p *parser) errorUnexpected(chr rune) error {
	if chr == -1 {
		return p.error(errUnexpectedEndOfInput)
	}
	return p.errorf(errUnexpectedToken, token.Illegal)
}

func (p *parser) errorUnexpectedToken(tkn token.Token) error {
	switch tkn {
	case token.Eof:
		return p.error(errUnexpectedEndOfInput)
	case token.Boolean, token.Null:
	case token.Identifier:
		return p.error(errUnexpectedIdentifier)
	case token.Keyword:
		return p.error(errUnexpectedReservedWord)
	case token.EscapedReservedWord:
		return p.error(errKeywordEscapedChars)
	case token.Number:
		return p.error(errUnexpectedNumber)
	case token.String:
		return p.error(errUnexpectedString)
	}
	return p.errorf(errUnexpectedToken, tkn.String())
}
