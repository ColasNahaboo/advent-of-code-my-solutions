package main

import (
	//"reflect" // reflect.DeepEqual(got, expected) 
	"testing"
)

// example.map
//   1        
// 234        
//   56       

//\x          11111111
//y  12345678901234567
//1          ...#     1
//2          .#..     2
//3          #...     3
//4          ....     4
//5  ...#.......#     5
//6  ........#...     6
//7  ..#....#....     7
//8  ..........#.     8
//9          ...#.... 9 
//10         .....#.. 10
//11         .#...... 11
//12         ......#. 12
//   12345678901234567
//            11111111


// reflect.DeepEqual(got, expected)
func Test_wrapCube(t *testing.T) {
	lines := fileToLines("example.txt")
	parse(lines)
	parse3Dmap("example.map")

	t_wc("1lt", t,  9, 1, 2,  5,  5, 1)
	t_wc("1lb", t,  9, 4, 2,  8,  5, 1)
	t_wc("1rt", t, 12, 1, 0, 16, 12, 2)
	t_wc("1rb", t, 12, 4, 0, 16,  9, 2)
	t_wc("1ul", t,  9, 1, 3,  4,  5, 1)
	t_wc("1ur", t, 12, 1, 3,  1,  5, 1)

	t_wc("2ul", t,  1, 5, 3, 12,  1, 1) 
	t_wc("2ur", t,  4, 5, 3,  9,  1, 1) 
	t_wc("2dl", t,  1, 8, 1, 12, 12, 3) 
	t_wc("2dr", t,  4, 8, 1,  9, 12, 3) 

	t_wc("3ul", t,  5, 5, 3,  9,  1, 0) 
	t_wc("3ur", t,  8, 5, 3,  9,  4, 0) 
	t_wc("3dl", t,  5, 8, 1,  9, 12, 0) 
	t_wc("3dr", t,  8, 8, 1,  9,  9, 0) 

 	t_wc("4rt", t, 12, 5, 0, 16,  9, 1)
	t_wc("4rb", t, 12, 8, 0, 13,  9, 1)

	t_wc("5lt", t, 9, 9, 2, 8, 8, 3)
	t_wc("5lb", t, 9,12, 2, 5, 8, 3)
	t_wc("5dl", t, 9,12, 1, 4, 8, 3)
	t_wc("5fr", t,12,12, 1, 1, 8, 3)

	t_wc("6ul", t, 13, 9, 3, 12,  8, 2)
	t_wc("6ur", t, 16, 9, 3, 12,  5, 2)
	t_wc("6rt", t, 16, 9, 0, 12,  4, 2)
	t_wc("6rb", t, 16,12, 0, 12,  1, 2)
	t_wc("6dl", t, 13,12, 1,  1,  8, 0)
	t_wc("6dr", t, 16,12, 1,  1,  5, 0)
}

func t_wc(label string, t *testing.T, x1, y1, d1, x2, y2, d2 int) {
	p1 := x1 + y1*gw
	p2 := x2 + y2*gw
	t.Run(label, func(t *testing.T) {
		p, d := wrapCube(p1, d1)
		if p != p2 || d != d2 {
			t.Errorf("expected %d:%d (%d,%d) but got %d:%d (%d,%d)", p2, d2, x2, y2, p, d, p % gw, p / gw)
		}
	})
}
