func! Remove(string)

	silent! exe '%substitute/' . a:string . '//g'
endf

func! RemoveInline(string)

	silent! exe 'substitute/' . a:string . '//g'
endf

func! AbbrevPrimitives()

	" substitution table
	silent! %s/HERE/HRE/g
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
	silent! %s/OVER/OVR/g    
endf

func! MakeBoundedPattern(pat)
	return '\(^\|\s\)\zs' . a:pat . '\ze\(\s\|$\)'
endf

let g:word_id = 0
let g:varnames_used = ""

func! Mangle(word_def, prefix)


	if a:word_def =~ '^\s*$'

		echo "Empty string: " . a:word_def
		return
	endif

	" pattern for calls to that word
	let l:word_call = MakeBoundedPattern(a:word_def)

	" new name
	let l:generated_name = a:prefix . g:word_id

	while g:varnames_used =~ l:generated_name
		echo "Detected duplicate name, randomizing"

		let l:generated_name = a:prefix . (rand() % 100)
	endw

	let g:word_id = g:word_id + 1

	let g:varnames_used = g:varnames_used . l:generated_name . ';'

	exe '%sub/' . l:word_call . '/' . l:generated_name . '/g'
	exe 'normal!' . 'ggo// @' . line('.') . '/' . l:generated_name . '	->	' . a:word_def
endf

func! MangleWord() " this function must be called only on lines with a forth word definition

	" get the defined word
	let l:word_def = matchstr(getline('.'), '[0-9A-Za-z\-!@?]\+')

	" mangle it
	cal Mangle(l:word_def, 'w')
endf

func! MangleVar()

	" get the defined word
	let l:var_def = matchstr(getline('.'), '\(VARIABLE\|VALUE\)\s\+\zs[0-9A-Za-z\-!@?]\+')

	" mangle it
	cal Mangle(l:var_def, 'v')
endf

func! Unforth(outfile) abort

	let l:file_name = substitute(expand('%:t'), '\..*', '', 'g')

	" save to new file
	execute 'saveas!' a:outfile

	setlocal ft=none

	" create an empty line at the start
	normal! ggo

	" My shorthands:

	" 3-byte tagged buffer
	silent! %s/#3/TAG NOP NOP/g

	" 3-byte untagged buffer
	silent! %s/#_3/NOP NOP NOP/g

	" 4-byte tagged buffer
	silent! %s/#4/TAG NOP NOP NOP/g

	" 4-byte untagged buffer
	silent! %s/#_4/NOP NOP NOP NOP/g

	" actual forth stuff

	" remove quotes
	cal Remove('"')

	" remove comments in the form: \ something something
	cal Remove('\\.*')

	" remove comments in the form: ( something )
	cal Remove('(.*)')

	" remove blank lines
	g/^\s*$/d

	" remove trailing whitespace
	cal Remove('\s*$')

	" remove excess whitespace
	%s/\( \)\+/\1/g

	" Mangle word names to avoid conflicts
	g/^:/cal MangleWord()
	g/VARIABLE\|VALUE/cal MangleVar()

	cal AbbrevPrimitives()

 	" Text Decorations
	silent! exe '%s/' . MakeBoundedPattern('A')     . '/ 7 6 API /g'
	silent! exe '%s/' . MakeBoundedPattern('SPACE') . '/ 6 6 API /g'
	silent! exe '%s/' . MakeBoundedPattern('COLON') . '/ 8 6 API /g'

	" 'g;^\s*[^/];...': execute ... in uncommented lines

	" remove everything after 3 letters of each word
	g;^\s*[^/];cal RemoveInline(MakeBoundedPattern('[A-Za-z\-!@?*]\{1,3\}\zs[A-Za-z\-!@?*]*'))
	g;VAL\|VAR;cal RemoveInline(MakeBoundedPattern('\(VAR\|VAL\s\+[0-9A-Za-z\-!@?]\{1,3\}\)\zs[0-9A-Za-z\-!@?]*\*'))

	" quote everything
	g;^\s*[^/];normal! A" NL
	g;^\s*[^/];normal! I"

 	" write code variable
	normal! gg0O#define NL "\n"
	normal! gg0Oconst char code[] PROGMEM =

	" write guards
	let l:guard = '__' . l:file_name . '_FTH__'
	
	exe 'normal!' 'gg0O#define ' . l:guard
	exe 'normal!' 'gg0O#ifndef ' . l:guard

	" close variable/guard
	normal! G0o;
	normal! G0o#endif

	setlocal ft=c
	w
endf

setlocal noswapfile
command! -nargs=1 Unforth cal Unforth(<q-args>)

