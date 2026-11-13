package lex

import ()

func Normalise(tokens []Token) []Token {
	n := normaliser{tokens, []Token{}}
	n.normalise()
	return n.out
}

type normaliser struct {
	in  []Token
	out []Token
}

func (n *normaliser) normalise() {
	current := n.next()
	next := n.next()

	for {
		if current.Kind == TokenEOF {
			return
		}
		if current.Kind == TokenIDENTIFIER && next.Kind == TokenOPEN_CURLY {
			n.emit(current)
			n.emit(Token{":", TokenCOLON}) //using empty string as value for synthetic tokens
		} else {
			n.emit(current)
		}

		current = next
		next = n.next()
	}
}

func (n *normaliser) next() Token {
	if len(n.in) == 0 {
		return EOFtoken
	}
	next := n.in[0]
	n.in = n.in[1:]
	return next
}

func (n *normaliser) emit(token Token) {
	n.out = append(n.out, token)
}
