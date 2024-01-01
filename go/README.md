# Advent of code solutions in GO
Here are my solutions to the "Advent of code" challenges implemented in GO. See https://adventofcode.com/

I discovered "advent of code" in 2021, and did it in [bash](../) for the challenge. But since I started learning GO, I decided to code the previous AoC years in GO as a mean to practice it. So, although my bash code can be useful to see some tricks of a seasoned bash programmer, these GO solutions must be considered as "student code".

The code is in GO, with some housekeeping scripts in bash.

I tried to keep each day directory as standalone as possible. For instance, instead of factorizing code that is often used acrosss days, I store them in the `TEMPLATES/` directory, and I copy them in the day dir if needed, instead of making separate packages. This is because I keep tweaing them to experiment and do not want to manage backwards compatibility, and ensure you can re-compile my exemple by downloading only the directory of single day without strings attached. Of course, I may publish as proper separate packages on GitHub these libraries once I am confident enough their API is stabilized.

More details in each of the per-year subdirectories.
