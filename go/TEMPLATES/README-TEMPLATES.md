# Colas Adventofcode templates

These are files to be copied in the day directories as needed, to provide some functionalities. So they are in the `main` package namespace. You can also find here files used to populate new day directories via the bash script `MAKEDAY`.

I did not make proper standalone packages because:
- I was learning Go, and it gave me the occasion to reinvent the wheel to do things by myself and better understand them
- There already exist better (but maybe more complex to use) packages on github, no need to publish them as standalone packages
- I may evolve them, and maybe later publish some of them as proper standalone packages, so I until then I can tweak them and change their API without fear of breaking past exercices relying on them, as they work with their own copy.
- Some are just bits copied from other authors (with attribution)
- I try to avoid dependencies, so i liked the idea of  having all the code necessary for a day problem contained it its directory.

## Used external packages
On a related note, I also used the external packages not in the Go standard library:

- `github.com/gammazero/deque`  fast FIFO & LIFO
- `github.com/fzipp/astar` simple and flexible Astar shortest path finder. I finally made my own version `astar.go`, even easier to use, using generics
- `github.com/emirpasic/gods/` plenty of useful data structures and algorithms in Go
- `github.com/deckarep/golang-set` a production-quality implementation of Sets
  - import by: `mapset "github.com/deckarep/golang-set/v2"`
