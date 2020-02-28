package main

import "fmt"

// Nodes for every step
type Node struct {
	i    float32
	deri float32
	derw float32

	previous *Node
	next     *Node
}

func (previousNode *Node) Plus(b float32) *Node {
	var current Node

	current.i = previousNode.i + b
	current.deri = 1
	current.derw = 1

	current.previous = previousNode
	current.next = nil
	return &current
}

func (previousNode *Node) Minus(b float32) *Node {
	var current Node

	current.i = previousNode.i - b
	current.deri = 1
	current.derw = 1

	current.previous = previousNode
	current.next = nil
	return &current
}

func (previousNode *Node) Multiply(w float32) *Node {
	var current Node

	current.i = previousNode.i * w
	current.deri = w
	current.derw = previousNode.i

	current.previous = previousNode
	current.next = nil
	return &current
}

func (previousNode *Node) Square() *Node {
	var current Node

	current.i = previousNode.i * previousNode.i
	current.deri = 2 * previousNode.i
	current.derw = 1

	current.previous = previousNode
	current.next = nil
	return &current
}

// Computatinal Graph
type Graph struct {
	length int
	head   *Node
}

func (g *Graph) AppendtoGraph(node Node) {
	//newNode := Node{1, 0, 0, nil, nil}
	if g.head == nil {
		g.head = &node
	} else {
		current := g.head
		for current.next != nil {
			current = current.next
		}
		current.next = &node
	}
	g.length++
	return

}

// traverse the linked list aka Graph
func (g *Graph) ComputationalGraph() error {
	if g.head == nil {
		return fmt.Errorf("Empty Graph")
	}
	current := g.head
	for current != nil {
		fmt.Println(current.i)
		current = current.next
	}
	return nil
}

// traverse the linked list aka Graph
func (g *Graph) BackwardPass() []float32 {
	if g.head == nil {
		println("ERROR: Graph is empty")
		return nil
	}
	current := g.head
	i := 0
	ders := make([]float32, g.length)
	for current != nil {
		//fmt.Println(current.i)
		ders[i] = current.deri
		i++
		current = current.next
	}
	return ders
}

func (g *Graph) B() []float32 { //for d(node)/d(weight)
	if g.head == nil {
		println("ERROR: Graph is empty")
		return nil
	}
	current := g.head
	i := 0
	ders := make([]float32, g.length)
	for current != nil {
		//fmt.Println(current.i)
		ders[i] = current.derw
		i++
		current = current.next
	}
	return ders
}

func initGraph() *Graph {
	return &Graph{}
}

func main() {

	var x, y float32 = 2, 9
	var a, b, c, d float32 = 1, 1, 1, 1

	graph := initGraph()

	X := Node{x, 0, 0, nil, nil}

	n1 := X.Multiply(a)
	graph.AppendtoGraph(*n1)
	n2 := n1.Plus(b)
	graph.AppendtoGraph(*n2)
	n3 := n2.Square()
	graph.AppendtoGraph(*n3)
	n4 := n3.Multiply(c)
	graph.AppendtoGraph(*n4)
	n5 := n4.Plus(d)
	graph.AppendtoGraph(*n5)
	n6 := n5.Square()
	graph.AppendtoGraph(*n6)
	n7 := n6.Minus(y)
	graph.AppendtoGraph(*n7)
	n8 := n7.Square()
	graph.AppendtoGraph(*n8)

	println("Forward Pass: ")
	graph.ComputationalGraph()

	//println(graph.length)
	fmt.Println("Backward Pass: ")

	der := graph.BackwardPass()

	for i := 0; i < len(der); i++ {
		fmt.Print(der[i])
		fmt.Print(" ")
	}

	fmt.Println()
	derw := graph.B()

	for i := 0; i < len(derw); i++ {
		fmt.Print(derw[i])
		fmt.Print(" ")
	}

	memoize := make([]float32, graph.length)
	// initializing memoize
	for i := graph.length - 1; i > 0; i-- {
		memoize[i] = 1
	}

	fmt.Println()
	//fmt.Println(len(der))
	//fmt.Println(len(memoize))

	// DP for BackPropagation
	for i := graph.length - 1; i > 0; i-- {
		memoize[i-1] = der[i] * memoize[i]
		//fmt.Print(memoize[i])
		//fmt.Print(" ")
	}

	fmt.Println(memoize)

	//	FOR DL/DW USE (N-1)

	// d at node 5
	dLdd := memoize[4] * derw[4]
	fmt.Println(dLdd)

	// c at node 4
	dLdc := memoize[3] * derw[3]
	fmt.Println(dLdc)

	// b at node 2
	dLdb := memoize[1] * derw[1]
	fmt.Println(dLdb)

	// a at node 1
	dLda := memoize[0] * derw[0]
	fmt.Println(dLda)

	// UPDATING WEIGHTS
	var alpha float32 = 0.1

	a = a - alpha*dLda
	b = b - alpha*dLdb
	c = c - alpha*dLdc
	d = d - alpha*dLdd

	fmt.Print("Updated Weights: ")
	fmt.Println(a, b, c, d)

	/*
		//d at node 5 ; dLdd
		var dNode int = 5
		var dLdd float32 = 1
		for i := len(der)-1; i>dNode-1; i-- {
			dLdd = dLdd*der[i]
		}

		dLdd = dLdd*derw[dNode-1]

		fmt.Println()
		fmt.Println(dLdd)

		//c at node 4 ; dLdc
		var cNode int = 4
		var dLdc float32 = 1
		for i := len(der)-1; i>cNode-1; i-- {
			dLdc = dLdc*der[i]
		}

		dLdc = dLdc*derw[cNode-1]

		fmt.Println()
		fmt.Println(dLdc)

		//b at node 2 ; dLdd
		var bNode int = 2
		var dLdb float32 = 1
		for i := len(der)-1; i>bNode-1; i-- {
			dLdb = dLdb*der[i]
		}

		dLdb = dLdb*derw[bNode-1]

		fmt.Println()
		fmt.Println(dLdb)

		//a at node 1 ; dLdd
		var aNode int = 1
		var dLda float32 = 1
		for i := len(der)-1; i>aNode-1; i-- {
			dLda = dLda*der[i]
		}

		dLda = dLda*derw[aNode-1]

		fmt.Println()
		fmt.Println(dLda)
	*/

}
