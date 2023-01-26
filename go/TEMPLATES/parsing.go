// Simple parsing with Scanf, example:

	for lineno := 0; lineno < len(lines); lineno++ {
		line := lines[lineno]
		if nf, _ := fmt.Sscanf(line, "cpy %d %1s", &n, &regname); nf == 2 {

		} else if nf, _ := fmt.Sscanf(line, "cpy %1s %1s", &regname2, &regname); nf == 2 {

		} else {
			panic(fmt.Sprintf("Syntax error line %d: %s\n", lineno, line))
		}
	}
	
