package expr

import (
	"fmt"
	"strconv"
)

const greatestOpWeight = 99

func (n nilNode) String() string {
	return "nil"
}

func (n identifierNode) String() string {
	return n.value
}

func (n numberNode) String() string {
	return fmt.Sprintf("%v", n.value)
}

func (n boolNode) String() string {
	if n.value {
		return "true"
	}
	return "false"
}

func (n textNode) String() string {
	return strconv.Quote(n.value)
}

func (n nameNode) String() string {
	return n.name
}

func (n unaryNode) String() string {
	switch n.operator {
	case "!", "not":
		return fmt.Sprintf("%v %v", n.operator, n.node)
	}
	return fmt.Sprintf("(%v%v)", n.operator, n.node)
}

func (n binaryNode) String() string {
	leftBinary := false
	rightBinary := false
	switch n.left.(type) {
	case binaryNode:
		leftBinary = true
	}
	switch n.right.(type) {
	case binaryNode:
		rightBinary = true
	}

	if leftBinary && !rightBinary {
		if n.left.(binaryNode).findLowestOpWeight(greatestOpWeight) < n.opWeight() {
			return fmt.Sprintf("(%v) %v %v", n.left, n.operator, n.right)
		} else {
			return fmt.Sprintf("%v %v %v", n.left, n.operator, n.right)
		}
	}
	if !leftBinary && rightBinary {
		if n.right.(binaryNode).findLowestOpWeight(greatestOpWeight) < n.opWeight() {
			return fmt.Sprintf("%v %v (%v)", n.left, n.operator, n.right)
		} else {
			return fmt.Sprintf("%v %v %v", n.left, n.operator, n.right)
		}
	}

	if leftBinary && rightBinary {
		var l, r string
		if n.left.(binaryNode).findLowestOpWeight(greatestOpWeight) < n.opWeight() {
			l = fmt.Sprintf("(%v)", n.left)
		} else {
			l = fmt.Sprintf("%v", n.left)
		}
		if n.right.(binaryNode).findLowestOpWeight(greatestOpWeight) < n.opWeight() {
			r = fmt.Sprintf("(%v)", n.right)
		} else {
			r = fmt.Sprintf("%v", n.right)
		}
		return fmt.Sprintf("%v %v %v", l, n.operator, r)

	}

	return fmt.Sprintf("%v %v %v", n.left, n.operator, n.right)
}

func (n binaryNode) findLowestOpWeight(min int) int {
	if n.opWeight() < min {
		min = n.opWeight()
	}
	switch n.left.(type) {
	case binaryNode:
		return n.left.(binaryNode).findLowestOpWeight(min)
	}
	switch n.right.(type) {
	case binaryNode:
		return n.right.(binaryNode).findLowestOpWeight(min)
	}
	return n.opWeight()
}

func (n binaryNode) opWeight() int {
	op := map[string]int{
		"+":   2,
		"-":   2,
		"*":   3,
		"/":   3,
		"and": 11,
		"&&":  11,
		"or":  12,
		"||":  12,
	}
	if val, ok := op[n.operator]; ok {
		return val
	}
	return greatestOpWeight
}

func (n matchesNode) String() string {
	return fmt.Sprintf("(%v matches %v)", n.left, n.right)
}

func (n propertyNode) String() string {
	return fmt.Sprintf("%v.%v", n.node, n.property)
}

func (n indexNode) String() string {
	return fmt.Sprintf("%v[%v]", n.node, n.index)
}

func (n methodNode) String() string {
	s := fmt.Sprintf("%v.%v(", n.node, n.method)
	for i, a := range n.arguments {
		if i != 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", a)
	}
	return s + ")"
}

func (n builtinNode) String() string {
	s := fmt.Sprintf("%v(", n.name)
	for i, a := range n.arguments {
		if i != 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", a)
	}
	return s + ")"
}

func (n functionNode) String() string {
	s := fmt.Sprintf("%v(", n.name)
	for i, a := range n.arguments {
		if i != 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", a)
	}
	return s + ")"
}

func (n conditionalNode) String() string {
	return fmt.Sprintf("%v ? %v : %v", n.cond, n.exp1, n.exp2)
}

func (n arrayNode) String() string {
	s := "["
	for i, n := range n.nodes {
		if i != 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", n)
	}
	return s + "]"
}

func (n mapNode) String() string {
	s := "{"
	for i, p := range n.pairs {
		if i != 0 {
			s += ", "
		}
		s += fmt.Sprintf("%v", p)
	}
	return s + "}"
}

func (n pairNode) String() string {
	switch n.key.(type) {
	case binaryNode, unaryNode:
		return fmt.Sprintf("%v: %v", n.key, n.value)
	}
	return fmt.Sprintf("%q: %v", n.key, n.value)
}
