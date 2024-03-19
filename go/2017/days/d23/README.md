# Advent of code challenge 2017, in GO, day d23

Part 1 is a straughtforward subset of the day 19 problem, a virtual computer program running on a "tablet".

The Part 2 howver requires to examine the result and try to detect cycles as to predict a very distant end of the program.

The program in my input ends with the following intructions:
Note: I prefix with the line number, starting at zero.
```
19: jnz g -8
20: sub d -1
21: set g d
22: sub g b
23: jnz g -13
24: jnz f 2
25: sub h -1
26: set g b
27: sub g c
28: jnz g 2
29: jnz 1 3
30: sub b -17
31:jnz 1 -23
```
It acts as 2 barrers that must be opened for the program to terminate.
- first, line 19: is passed if g is `0`
- second, line 23: is passed if `d - b` is `0`, that is if `d` equals `b`, after incrementing `d` by one

We see also that:
- quickly, `b` has the value `108400` whichs then stays constant
- `d` starts at `2` at 19: as soon as step 9
- `g` starts at `-108397` at 19: as soon as step 17 and decrements by one on each pass through line 19.
- symetrically, `e` restarts at `3` to reach `108400`, incrementing by one
- the program reaches lines 19: every 8 steps after 17: 25, 33, 41, 39...

Looking at the registers when reaching line 19 we see:
```
...step   line      a      b      c d      e f       g h
  [17]     @19 jnz [1 108400 125400 2      3 1 -108397 0]
  [25]     @19 jnz [1 108400 125400 2      4 1 -108396 0]
  [33]     @19 jnz [1 108400 125400 2      5 1 -108395 0]
...
  [867194] @19 jnz [1 108400 125400 2 108400 0       0 0]
  [867207] @19 jnz [1 108400 125400 3      3 0 -108397 0]
  [867215] @19 jnz [1 108400 125400 3      4 0 -108396 0]

```
We can then deduce when at line 19 we will have `d == 108400` and `g == 0`, we will have the register values
```
  [?] @19 jnz [1 108400 125400 108400 108400 0       0 0]
```

But the proper solution was to see that the input program 
```
set b 84
set c b
jnz a 2
jnz 1 5
mul b 100
sub b -100000
set c b
sub c -17000
set f 1
set d 2
set e 2
set g d
mul g e
sub g b
jnz g 2
set f 0
sub e -1
set g e
sub g b
jnz g -8
sub d -1
set g d
sub g b
jnz g -13
jnz f 2
sub h -1
set g b
sub g c
jnz g 2
jnz 1 3
sub b -17
jnz 1 -23
```

was equivalent in pseudocode to:

```
c = b = 84

if a != 0:
	b = b * 100 + 100000  // b = 108400
	c = b + 17000         // c = 125400
	f = 1
	d = 2
		e = 2
			g = d * e - b
			
			if g == 0:
				f = 0
			
			
			e += 1
			g = e - b
			if g != 0: loop
		d += 1
		g = d - b
		if g != 0 loop
	
	if f == 0
		h += 1
	
	g = b - c
	if g == 0: break
	b += 17
	loop
```

And the following go code, that is still too slow for real use.
```
b := 108400
c := 125400
h := 0
for ; b <= c; b += 17 {
    f := 1
    for d := 2; d <= b; d++ {
        for e := 2, e <= b; e++ {
            if d * e == b {
                f = 0
            }
        }
    }
	if f == 0 {
		h++
    }
}
fmt.Print(h)

```

And thus, either implement this code directly, or recognizing the main loop a seach for prime numbers, counting all non-prime numbers between b and c (inclusive) by steps of 17, giving the much faster fast GO code:

```
b := 108400
c := 125400
h := 0
for ; b <= c; b += 17 {
	if ! IsPrime(b) {
		h++
    }
}
fmt.Print(h)
```

Note that ChatGPT (via Perplexity) was able to recognize this code. At the question:
```
Please explain what this Go code is doing:

for ; b <= c; b += 17 {
    f := 1
    for d := 2; d <= b; d++ {
        for e := 2, e <= b; e++ {
            if d * e == b {
                f = 0
            }
        }
    }
	if f == 0 {
		h++
    }
}
```

It answered:

This Go code snippet appears to be a program that iterates over a range of values from `b` to `c` with a step size of 17. Within each iteration, it checks if the value of `b` is a product of two numbers other than 1 and itself. Here is a breakdown of the code:

1. The outer loop iterates over values of `b` starting from the initial value of `b` until it reaches or exceeds the value of `c`, incrementing `b` by 17 in each iteration.
   
2. Inside this loop, a variable `f` is initialized to 1. This variable is used to track whether `b` is a prime number or not.

3. There is an inner loop that starts with `d` at 2 and increments it until it reaches the value of `b`. This loop is used to iterate over potential divisors of `b`.

4. Within the inner loop for `d`, there is another nested loop for `e`, starting at 2 and also iterating up to the value of `b`. This nested loop checks if the product of `d` and `e` equals `b`.

5. If a pair of numbers (`d`, `e`) is found such that their product equals `b`, then the variable `f` is set to 0, indicating that `b` is not a prime number.

6. After the inner loops complete, if the variable `f` remains 1 (indicating that no pair of numbers multiplied together equals `b`), then the counter variable `h` is incremented.

In summary, this code snippet seems to be counting how many numbers within the range from `b` to `c`, incremented by 17, are not prime numbers (i.e., can be expressed as a product of two numbers other than 1 and itself).
