func Eval(n Node) (Node, error) {
	switch nt := n.(type) {
	case None:
		return nt, nil
	case Literal:
		return nt, nil
	case BinOp:
		var result int
		for i, n := range []Node{nt.Left, nt.Right} {
			l, err := Eval(n)
			if err != nil {
				return nil, fmt.Errorf("error while evaluating operand # %d: %s", i, err)
			}

			lit, ok := l.(Literal)
			if !ok {
				return nil, fmt.Errorf("operand %d is not a literal: %+v", i, lit)
			}

			val := lit.Contents

			if i == 0 {
				result = val
				continue
			}

			switch nt.Op {
			case Plus:
				result += val
			case Minus:
				result -= val
			case Times:
				result *= val
			}
		}

		return Literal{result}, nil

	default:
		return nil, fmt.Errorf("unhandled node type: %+v", nt)
	}

}
