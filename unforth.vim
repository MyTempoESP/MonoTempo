function Unforth()
	silent! %s/"//g
	silent! %s/\/\//\\ /g
	silent! %s/NL//g

	%normal! 9<

	silent! %s/VAL/VALUE/g
	silent! %s/NXT/NEXT/g
	silent! %s/THN/THEN/g
	silent! %s/BGN/BEGIN/g
	silent! %s/RPT/REPEAT/g
	silent! %s/UTL/UNTIL/g
	silent! %s/WHL/WHILE/g
	silent! %s/ALO/ALLOT/g
	silent! %s/VAR/VARIABLE/g
	silent! %s/SWP/SWAP/g
endfunction

function Forth()
	silent! %s/\\.\+//g
	g/^ *$/d

	silent! %s/ *$//

	%normal! I"
	%normal! A"

	normal! gg0O#define NL "\n"
	normal! gg0Oconst char code[] PROGMEM =

  silent! %s/VALUE/VAL/g
	silent! %s/NEXT/NXT/g
	silent! %s/THEN/THN/g
	silent! %s/BEGIN/BGN/g
	silent! %s/REPEAT/RPT/g
	silent! %s/UNTIL/UTL/g
	silent! %s/WHILE/WHL/g
	silent! %s/ALLOT/ALO/g
	silent! %s/VARIABLE/VAR/g
	silent! %s/SWAP/SWP/g 

	silent! " Text Decorations
	silent! %s/ A / 7 6 API /g
	silent! %s/ SPACE / 6 6 API /g
	silent! %s/ COLON / 8 6 API /g

	silent! %s/ \(.\){3}//g

	execute 'w!' expand('%') . ".h"
	undo
endfunction

