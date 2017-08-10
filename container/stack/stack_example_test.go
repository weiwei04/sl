package stack

import "fmt"

// ExampleStack use stack to transform a arithmetic expression
// from infix to postfix, from Algorithm 4th 1.3.10
func ExampleStack_infixToPostFix() {
	// suppose we don't contain ' ' and the parenthsees is valid
	var infixExp = "((1+2)*((3-4)*(5-6)))"
	var postExp = []byte{}
	s := NewSliceStack()
	for _, token := range []byte(infixExp) {
		if token == '(' {
			continue
		} else if token == '+' || token == '-' || token == '*' || token == '/' {
			s.Push(token)
		} else if token == ')' {
			postExp = append(postExp, s.Pop().(byte))
		} else {
			postExp = append(postExp, token)
		}
	}
	for s.Size() != 0 {
		postExp = append(postExp, s.Pop().(byte))
	}
	fmt.Println(string(postExp))
	// Output: 12+34-56-**
}
