package main

type Stack []int

func (s Stack) Peek() int {
	//fmt.Printf("Peek %d\n", s[len(s)-1])
	return s[len(s)-1]
}

func (s Stack) Push(i int) Stack {
	//fmt.Printf("Push %d\n", i)
	return append(s, i)
}

func (s Stack) IsEmpty() bool {
	return len(s) == 0
}

func (s Stack) Pop() (Stack, int) {

	i := s.Peek()
	s = s[:len(s)-1]

	//fmt.Printf("Pop %d\n", i)

	return s, i
}
