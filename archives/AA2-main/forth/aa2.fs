\ __AA2_FTH__

\ General overview

\ stack model:
\           __ TOP
\         /
\        v
\  ( a b c )

\ Tagged addresses are denoted by [ADDR_NAME]
\ the word for them is TAG, which puts the address inside the function in the ~~return~~ normal stack. Values greater than $7F have to be 'encoded' using the 16-bit-encode word

\ type BigNumber -> 32-bit, format: $BF(8-bit) prefixed 16-bit value + 8-bit magnitude
\                   like: 0xBF[vvvv][mm], example: 0xBF004103 [65K]

\ #3  -> TAG NOP NOP    , 3-byte tagged   buffer, this is usually a placeholder for a CALL instruction;
\ #4  -> TAG NOP NOP NOP, 4-byte tagged   buffer, this is usually a placeholder for code, allowing 2 calls;
\ #_3 -> NOP NOP NOP    , 3-byte tagged   buffer.
\ #_4 -> NOP NOP NOP NOP, 4-byte untagged buffer, this is usually a placeholder for raw data.

\ Word names ending with a ! mean they WRITE something to memory.
\ Word names ending with a @ mean they FETCH something from memory.

\ functions marked with an asterisk(*) need to be double checked for duplicate names, consider only the first 3 chars!!
\ they are also EXTERNAL, DO NOT CALL THEM IN THIS SOURCE CODE!

\ the words 'alignment' and other technical terms are, for historical reasons, severely misused
\ in this source code, so refer to the following definitions:

\	'alignment'/'align'/'aligned'	generally means something is offset according to a given byte size
\					... e.g. BigNumbers are 4-byte values, therefore anything written after
\					... a BigNumber has to be 'aligned' by 4 bytes, i could say 'corrected' or
\					... 'skip 4 bytes', but aligned sounds smarter. Fake it til' you make it!

\	'instruction'			an instruction is simply an opcode | whatever arguments it may have, it does not
\					... exceed 3 bytes.

\	'encode'/'encoded'/'encoding'	this one is funny, because of how wrong it is.
\					... Unfortunately, we can't just shove a random value in memory,
\					... instead, we build an instruction which tells the computer:
\					... 'push this random number to the stack'. Then we put THAT in memory.
\					... We called that technique 'encoding'(TM)

\									[useless technical details]
\					there's something in forth called a LIT instruction, which, if you don't know forth,
\					... simply pushes a number to the stack. In nanoFORTH, LIT has the opcode 0xBF, followed
\					... by the actual 16-bit number.

\ ======================

\ C-API Drawing words

: Label    5   API ;
: Forward  2   API ;
: BigNum   1   API ; ( ; biblically accurate dump of 65K in BigNumber format: 0xBF004103, precision is for nerds )
: Number   4   API ;
: Value    6   API ;
: Ip       7   API ;
: Ms       3   API ;
: helper-Antenna
	A 1 Number COLON BigNum SPACE
	A 2 Number COLON BigNum Forward
	A 3 Number COLON BigNum SPACE
	A 4 Number COLON BigNum Forward
;

\ ======================

501 VALUE VERSION

\ The screen

: Dis			( -- ; Screen buffer )
 
	\ Memory layout:

	\               080: 3820       :

	\         REL   ADDR
        \       DATA+00 082: bbbbbbbb   DATA Row 1
        \       DATA+04 086: bbbbbbbb   DATA Row 2
        \       DATA+08 08A: bbbbbbbb   DATA Row 3
        \       DATA+0C 08E: bbbbbbbb   DATA Row 4

        \               092:   00       First label ( Portal myXXX )
        \               093: c005
        \               095: c079
        \               097: c01d
        \               099: c00d
        \       
        \       COD1+00 09B: bcbb       Code-1 - Call slot 1
        \       COD1+02 09D: bbbb       Code-1 - Call slot 2
        \       COD1+04 09F: c00d       Code-1 - Forward (END)
        \       
        \       COD2+00 0A1: bcbb       Code-2 - Call slot 1
        \       COD2+02 0A3: bbbb       Code-2 - Call slot 2
        \       COD2+04 0A5: c00d       Code-2 - Forward (END)
        \
        \               0A7:   03       Last label ( Comunicando WEB )
        \               0A8: c005
        \               0AA: bb00
        \               0AC: c025
        \               0AE: c00d
        
	\		0B0: 00b4       Draw()
        \       	0B2:   80       ;
         

	#_4 #_4 #_4 #_4		( COMPILE TIME: data-end-addr ; 16 byte data buffer ) 

	0       Label		( label-id ; Main label, shows device name )
	VERSION Number		( device-version )
	        Forward

	#4      Forward		( COMPILE TIME: tagged-addr )
	#4      Forward		( COMPILE TIME: tagged-addr tagged-addr-2 )

	3       Label NOP	( label-id ; Main label #2, shows device connectivity )
	0       Value		( value-id )
	        Forward

	0 API			( 0 ; C-Api call #0: Draw screen )
;

VALUE 2-CODE     ( data-end-addr tagged-addr tagged-addr-2 ; Address of the second screen code buffer )
VALUE 1-CODE     ( data-end-addr tagged-addr               ; Address of the first screen code buffer )

' Dis VALUE DATA ( ; Address of beginning of screen data buffer )

\ ======================

\ Utility functions for memory access

: 16-bit-encode!	( int16 addr -- )
	$BF			( int16 addr Lit-opcode      ; $BF is the opcode for pushing 16-bit integers to the stack )
	OVER C!			( int16 addr Lit-opcode addr ; writes a single byte LIT instruction to addr )
	1 +			( int16 addr +1              ; offsets addr by 1, right next to the LIT opcode )
	      !			( int16 addr'                ; writes a 2-byte, i.e. 16-bit, signed integer to addr' )
;

: F!			( m v iaddr -- )
	DUP >R			( m v iaddr       ; saves a copy of indexed-addr for later )
	16-bit-encode!		( m v iaddr       ; encodes value to specified DATA buffer index )
	R>			( m   iaddr       ; restores indexed-addr from Rstack )
	3 +			( m   iaddr +3    ; offsets the indexed address by 3 bytes )
	     C!			( m   iaddr'      ; writes the single-byte magnitude right next to the encoded value )
;

: calc-align		( idx align addr -- iaddr ; calculates aligned offset of the specified index )
	>R			( idx align addr )
	 *			( idx align      ; calculates ALIGN-byte *offset* )
	R>			( offset         )
	 +			( offset addr    ; calculates base address + offset to get the indexed-address )
;

\ ======================

\ Aligned writes

: data!			( value idx -- ; n-byte aligned 16-bit write to DATA buffer )
	DATA + 16-bit-encode!
;
: data-C!		( value idx -- ; n-byte aligned 8-bit write to DATA buffer )
	DATA + C!
;

: data-Big!		( m v idx -- ; automatic 4-byte aligned BigNumber write to DATA buffer, according to the BigNumber type )
	4			( value idx align ; this one has fixed alignment because 4 is the max size )
	DATA calc-align         ( value idx align ) 
	            F!		( value iaddr     ; writes a BigNumber to the specified DATA index )
;

\ ======================

\ Calling Words

: prepare-call		( word-addr target-addr -- ; this function is specific to the [encode-]call! functions )
	>R			( word-addr        target-addr ; puts target-addr in a pokeball )
	$C000 OR		( word-addr        call-opcode ; creates a call instruction, C000 | Address )
	R>			( call-instruction             ; goes like "target-addr, i choose you!" )
;

: encode-call!  	( word-addr target-addr -- ; useful in code modifying words, encode a literal call instruction )
	prepare-call		( word-addr        target-addr ; )
	16-bit-encode!		( call-instruction target-addr ; encodes the created call instruction as a value to target-addr )
;

: call! 		( word-addr target-addr -- ; writes a call instruction in target-addr )
	prepare-call		( word-addr        target-addr ; )
	!			( call-instruction target-addr ; writes the created call instruction, modifying target-addr )
;

: call-idx!		( word-addr idx code-addr -- ; writes a call instruction to aligned offset )
	>R			( word-addr idx align code-addr ; storing code-addr )
	2 R>			( word-addr idx align           ; retrieving code-addr, to reorder )
	calc-align		( word-addr idx align code-addr ; 2-byte call )
	call!			( word-addr indexed-addr )
;

\ ======================

\ Screen function addresses

' Label          VALUE addr-Label
' Forward        VALUE addr-Forward
' BigNum         VALUE addr-BigNum
' Number         VALUE addr-Number
' Value          VALUE addr-Value
' Ip             VALUE addr-Ip
' Ms             VALUE addr-Ms
' helper-Antenna VALUE addr-Antenna  

\ Screen functions

( ADDR IDX N-CODE call-idx! )
: program-calls!	( addr-W1 addr-W2 n-CODE -- ; programs the passed words into the specified CODE section )
	DUP >R >R		( addr-W1 addr-W2 n-CODE -- addr-W1 addr-W2 ; puts n-CODE in the return stack for 2 uses )
	1      R> call-idx!	( addr-W1 addr-W2 ; programs W2 to n-CODE index 1 )
	0   R>    call-idx!	( addr-W1         ; programs W1 to n-CODE index 0 )
;

\ Fixed screen graphic codes:
\
\ 			     LABELS                 VALUES
\                         _____________          _____________
\ 			$00  PORTAL   My       $00  WEB        
\ 			$01  ATLETAS           $01  CONECTAD   
\ 			$02  REGIST.           $02  DESLIGAD   
\ 			$03  COMUNICANDO       $03  AUTOMATIC  
\ 			$04  LEITOR            $04  OK         
\ 			$05  LTE/4G:           $05  X          
\ 			$06  WIFI:             $06  A          
\ 			$07  IP:               $07  COLON      
\ 			$08  LOCAL:
\ 			$09  PROVA:
\ 			$0A  PING:
\ 			$0B  HORA:
\ 			$0C  USB:
\ 			$0D  AGUARDE...
\ 			$0E  ERRO TENTAR,
\ 			$0F    NOVAMENTE
\
\ 			$10  RFID  -
\ 			$11  SERIE:
\ 			$12  SIST.
\
\ 			$13  PRESSIONE,
\ 			$14  PARA CONFIRMAR
\
\ 			$15  OFFLINE
\ 			$16  DATA:
 
: clean!
	\ DATA:
	\         REL   ADDR
        \       DATA+00 082: bbbbbbbb   DATA Row 1
        \       DATA+04 086: bbbbbbbb   DATA Row 2
        \       DATA+08 08A: bbbbbbbb   DATA Row 3
        \       DATA+0C 08E: bbbbbbbb   DATA Row 4

	\	16:22:21.120 -> 080: 3820bbbb bbbbbbbb bbbbbbbb bbbbbbbb  8 ;;;;;;;;;;;;;;
	\ 	16:22:21.153 -> 090: bbbb00c0 05c079c0 1dc00dc0 05c02dc0  ;;_@_@y@_@_@_@-@
	\ 	16:22:21.153 -> 0a0: 0dc005c0 25c00d03 c005bb00 c025c00d  _@_@%@__@_;_@%@_
	\ 	16:22:21.153 -> 0b0: 00b48000 7d763333 bf00a180 00b37633  _4__}v33?_!__3v3

	$bbbb DATA      !
	$bbbb DATA $2 + !
	$bbbb DATA $4 + !
	$bbbb DATA $6 + !
	$bbbb DATA $8 + !
	$bbbb DATA $C + !
	$bbbb DATA $A + !
	$bbbb DATA $E + !
;

: code-clean!
	$bbbb 1-CODE      !
	$bbbb 1-CODE $2 + !

	$bbbb 2-CODE      !
	$bbbb 2-CODE $2 + !
;

: S-1* ( m2 vvv2 m1 vvv1 -- ; extern, Tags+Unicas )

	clean!

	\ _ -> bb ( NOP )

	\ Memory layout:

	\ DATA:
	\         REL   ADDR
        \       DATA+00 082: bfvvv2m2   DATA Row 1
        \       DATA+04 086: 01______   DATA Row 2
        \       DATA+08 08A: bfvvv1m1   DATA Row 3
        \       DATA+0C 08E: 02______   DATA Row 4

	\ DMP should look something like this ( let m2 = 3, m2 = 10, m1 = 3, and m1 = 10 )
	\	15:40:40.516 -> 080: 3820bf00 0a0301bb bbbbbf00 0a0302bb  8 ?____;;;?____;
	\ 	15:40:40.516 -> 090: bbbb00c0 05c079c0 1dc00dc0 05c015c0  ;;_@_@y@_@_@_@_@
	\ 	15:40:40.516 -> 0a0: 0dc005c0 15c00d03 c005bb00 c025c00d  _@_@_@__@_;_@%@_
	\	15:40:40.516 -> 0b0: 00b480

	\ Data for 1-CODE
	\ we reserve the two LAST rows for that data.

	( LABEL: REGIST. )
	( value ) $02
	( index ) $0C
	data-C!

	( LITERAL: Tags[bfvvv1m1] )
	( value ) ( m2 vvv2 m1 vvv1 )
	( index )      $02
	( align )     ( 4 )
	data-Big!

	\ Data for 2-CODE
	\ we reserve the two FIRST rows for that data

	( LABEL: ATLETAS )
	( value ) $01
	( index ) $04
	data-C!

	( LITERAL: Atletas[bfvvv2m2] )
	( value ) ( m2 vvv2 )
	( index )      $00
	( align )     ( 4 )
	data-Big!

	addr-Label
	addr-BigNum 1-CODE program-calls!

	addr-Label
	addr-BigNum 2-CODE program-calls!
;

: S-2* ( st o1 o2 o3 o4 -- ; extern, IP+State )

	clean!
	 
	\ st        -> state
	\ o1 ... o4 -> octets
	\ _         -> bb ( NOP )

	\ Memory layout:

	\ DATA:
	\         REL   ADDR
        \       DATA+00 082: st04____   DATA Row 1
        \       DATA+04 086: ________   DATA Row 2
        \       DATA+08 08A: o1o2o3o4   DATA Row 3
        \       DATA+0C 08E: 07______   DATA Row 4

	\ Data for 1-CODE
	\ we reserve the two LAST rows for that data.

	( LABEL: IP: )
	( value )  $07
	( index )  $0C
	data-C!

	( value ) ( st o1 o2 o3 o4 )
	( index )  $0B
	data-C!
	( value ) ( st o1 o2 o3 )
	( index )  $0A
	data-C!
	( value ) ( st o1 o2 )
	( index )  $09
	data-C!
	( value ) ( st o1 )
	( index )  $08
	data-C!

	\ Data for 2-CODE
	\ we reserve the two FIRST rows for that data

	( LABEL: LEITOR )
	( value ) $04
	( index ) $01
	data-C!
	
	( value ) ( st )
	( index )  $00
	data-C!

	addr-Label
	addr-Ip 1-CODE program-calls!

	addr-Label
	addr-Value 2-CODE program-calls!
;

: S-3* ( ping state lte-or-wifi[$05/$06] -- ; extern, Wifi/Lte )

	clean!
	 
	\ st        -> state
	\ pvvv      -> ping-v
	\ _         -> bb ( NOP )

	\ Memory layout:

	\ DATA:
	\         REL   ADDR
        \       DATA+00 082: bfpvvv__   DATA Row 1
        \       DATA+04 086: 0A______   DATA Row 2
        \       DATA+08 08A: st______   DATA Row 3
        \       DATA+0C 08E: 06______   DATA Row 4

	\ let ping = 0x1f4, state = 0x01
	\	16:43:33.717 -> 080: 3820bf01 f4bb0abb bbbb01bb bbbb06bb  8 ?_t;_;;;_;;;_;
	\	16:43:33.717 -> 090: bbbb00c0 05c079c0 1dc00dc0 05c025c0  ;;_@_@y@_@_@_@%@
	\	16:43:33.717 -> 0a0: 0dc005c0 35c00d03 c005bb00 c025c00d  _@_@5@__@_;_@%@_
	\	16:43:33.717 -> 0b0: 00b48000 7d763334 bf00a180 00b37633  _4__}v34?_!__3v3

	\ Data for 1-CODE
	\ we reserve the two LAST rows for that data.

	( LABEL: WIFI: or LTE/4G: )
	( value )  ( lte-or-wifi )
	( index )  $0C
	data-C!

	( value ) ( wifi-state )
	( index )  $08
	data-C!

	\ Data for 2-CODE
	\ we reserve the two FIRST rows for that data

	( LABEL: PING: )
	( value ) $0A
	( index ) $04
	data-C!
	
	( value ) ( ping )
	( index )  $00
	data!

	addr-Label
	addr-Value 1-CODE program-calls!

	addr-Label
	addr-Ms 2-CODE program-calls!
;

: S-4* ( m1 v1 m2 v2 m3 v3 m4 v4 -- ; extern, Antennas )

	\ clean!		( we do   use all slots )
	code-clean!		( we dont use all slots )
	 
	\ _         -> bb ( NOP )

	\ Memory layout:

	\ DATA:
	\         REL   ADDR
        \       DATA+00 082: bfvvv1m1   DATA Row 1
        \       DATA+04 086: bfvvv2m2   DATA Row 2
        \       DATA+08 08A: bfvvv3m3   DATA Row 3
        \       DATA+0C 08E: bfvvv4m4   DATA Row 4

	0 data-Big!		( m1 v1 m2 v2 m3 v3 m4 v4 idx )
	1 data-Big!		( m1 v1 m2 v2 m3 v3       idx )
	2 data-Big!		( m1 v1 m2 v2             idx )
	3 data-Big!		( m1 v1                   idx )

	addr-Antenna 0 1-CODE call-idx!
;

: S-5*	( hour min sec day mon year -- ; extern, DateTime, this one is a mess )
 	\ DATA:

	\ * capitalized hexadecimal values for clarity

	\         REL   ADDR
        \       DATA+00 082: ddmmBFyy   DATA Row 1
        \       DATA+04 086: yyBFDA7E   DATA Row 2
        \       DATA+08 08A: 16hhmmss   DATA Row 3
        \       DATA+0C 08E: BF00FF0B   DATA Row 4

	\ clean!		( ; we use all data slots )

	( LABEL: DATA: )
	( value ) $16
	( index ) $08
	data-C!

	( value ) ( year )
	( index )  $02
	data!

	( value ) ( mon )
	( index )  $01
	data-C!

	( value ) ( day )
	( index )  $00
	data-C!

	( value ) $DA7E
	( index )   $05
	data!

	( LABEL: HORA: )
	( value )  $0B
	( index )  $0F
	data-C!

	( value ) $100
	( index )  $0C
	data!

	( value ) ( sec )
	( index )  $0B
	data-C!

	( value ) ( min )
	( index )  $0A
	data-C!

	( value ) ( hour )
	( index )  $09
	data-C!

	addr-Label
	addr-Ip 1-CODE program-calls!

	addr-Label
	addr-Ip 2-CODE program-calls!
;

: S-6*	( state ; extern, USB )
  	\ DATA:

	\         REL   ADDR
        \       DATA+00 082: st______   DATA Row 1
        \       DATA+04 086: 0c______   DATA Row 2
        \       DATA+08 08A: ________   DATA Row 3
        \       DATA+0C 08E: ________   DATA Row 4

	clean!
	code-clean!		( ; we do not use all call slots )

	( LABEL: USB: )
	( value ) $0C
	( index ) $04
	data-C!

	( value ) ( state )
	( index )  $00
	data-C!

	addr-Label
	addr-Value 1-CODE program-calls!
;

: S-7*	( version-hash ; extern, System info )
  	\ DATA:

	\         REL   ADDR
        \       DATA+00 082: BFrfid__   DATA Row 1
        \       DATA+04 086: 1010____   DATA Row 2
        \       DATA+08 08A: BFhash__   DATA Row 3
        \       DATA+0C 08E: 1012____   DATA Row 4

	clean!

	( LABEL: RFID - )
	( value ) $10
	( index ) $05
	data-C!

	\ Hexadecimal prefix
	( value ) $10
	( index ) $04
	data-C!

	( value ) $CA00 VERSION 1 RSHIFT OR
	( index ) $00
	data!

	( LABEL: SIST. )
	( value ) $12
	( index ) $0D
	data-C!

	\ Hexadecimal prefix
	( value ) $10
	( index ) $0C
	data-C!

	( value ) ( version-hash )
	( index ) $08
	data!

	addr-Label
	addr-BigNum 1-CODE program-calls!       

	addr-Label
	addr-BigNum 2-CODE program-calls!
;

\ ======================

\ Initialization

1500 DLY
50 0 TMI Dis
\ 10 1 TMI do-button
1 TME

\ ======================

