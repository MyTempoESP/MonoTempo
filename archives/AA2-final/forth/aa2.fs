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

\ #2  -> TAG NOP        , 2-byte tagged   buffer, this is usually a placeholder for a single call;
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

\ SYMBOLIC ADDRESSES (they are off by about +0x3D)

\ C-API Drawing words

: Button-Data	( ; storage for button presses )

	\ 000: ffff7730 20bbbb00 00000365 01350180
	\ 010: 01c201ff 025b0298 02e7bbbb bbbbbbbb
	\ 020: bbbbbbbb bbbbbbbb bbbbbbbb bbbbbbbb
	\ 030: bbbbbbbb bbbbbbbb bb80

	\ 008->OK-button
	\ 009->Arrow-button
	\ 00A->Action

	\ 00B:039->Free memory for storage

	#_4
	#_4
	#_4
	#_4
	#_4
	#_4
	#_4
	#_4
	#_4
	#_4
	#_4
	#_4
	#_4
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

	0       5 API		( label-id ; LABEL, device name )
	VERSION 4 API 		(          ; NUMBER device-version )
	        2 API		(          ; FORWARD )

	#4      2 API  		( COMPILE TIME: tagged-addr )
	#4      2 API  		( COMPILE TIME: tagged-addr tagged-addr-2 )

	3       5 API NOP	( label-id ; LABEL, Main label #2, shows device connectivity )
	0       6 API		( value-id ; VALUE )
	        2 API

	        0 API		( 0        ; C-Api call #0: Draw screen )
;

VALUE 2-CODE     ( data-end-addr tagged-addr tagged-addr-2 ; Address of the second screen code buffer )
VALUE 1-CODE     ( data-end-addr tagged-addr               ; Address of the first screen code buffer )

' Dis VALUE DATA ( ; Address of beginning of screen data buffer )

\ ======================

\ Utility functions for memory access

: Encode!	( int16 addr -- )
	$BF		( int16 addr Lit-opcode      ; $BF is the opcode for pushing 16-bit integers to the stack )
	OVER C!		( int16 addr Lit-opcode addr ; writes a single byte LIT instruction to addr )
	1 +		( int16 addr +1              ; offsets addr by 1, right next to the LIT opcode )
	      !		( int16 addr'                ; writes a 2-byte, i.e. 16-bit, signed integer to addr' )
;

: F!		( m v iaddr -- )
	DUP >R		( m v iaddr       ; saves a copy of indexed-addr for later )
	Encode!		( m v iaddr       ; encodes value to specified DATA buffer index )
	R>		( m   iaddr       ; restores indexed-addr from Rstack )
	3 +		( m   iaddr +3    ; offsets the indexed address by 3 bytes )
	     C!		( m   iaddr'      ; writes the single-byte magnitude right next to the encoded value )
;

\ ======================

\ writes based on bitsize

: data!	    (    v16 i -- )
	DATA + Encode!
;
: data-C!   (    v8  i -- )     DATA + C!      ;
: data-Big! ( m7 v16 i -- ) 4 * DATA + F!      ;

\ ======================

\ Screen functions

\
\ 	    	     LABELS                 VALUES            API
\                 _____________          _____________     _________
\ 	        $00  PORTAL   My       $00  WEB           Label    5
\ 	        $01  ATLETAS           $01  CONECTAD      Forward  2
\ 	        $02  REGIST.           $02  DESLIGAD      BigNum   1
\ 	        $03  COMUNICANDO       $03  AUTOMATIC     Number   4
\ 	        $04  LEITOR            $04  OK            Value    6
\ 	        $05  LTE/4G:           $05  X             Ip       7
\ 	        $06  WIFI:             $06  A             Ms       3
\ 	        $07  IP:               $07  COLON      
\ 	        $08  LOCAL:
\ 	        $09  PROVA:
\ 	        $0A  PING:
\ 	        $0B  HORA:
\ 	        $0C  USB:
\ 	        $0D  AGUARDE...
\ 	        $0E  ERRO TENTAR,
\ 	        $0F    NOVAMENTE
\
\ 	        $10  RFID  -
\ 	        $11  SERIE:
\ 	        $12  SIST.
\
\ 	        $13  PRESSIONE,
\ 	        $14  PARA CONFIRMAR
\
\ 	        $15  OFFLINE
\ 	        $16  DATA:
 
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

;
: PROG-S-1
	0 0 0 0 S-1*

	\ Code

	( id for label )
	( value  ) $05B4
	( offset ) $00
	1-CODE + !

	( id for BigNum)
	( value  ) $01B4
	( offset ) $02
	1-CODE + !

	( id for label )
 	( value  ) $05B4
	( offset ) $00
	2-CODE + !

	( id for BigNum)
	( value  ) $01B4
	( offset ) $02
	2-CODE + !
;

: S-2* ( st o1 o2 o3 o4 -- ; extern, IP+State )

	clean!
	 
	\ st        -> state
	\ o1 ... o4 -> octets
	\ _         -> bb ( NOP )

	\ Memory layout:

	\ DATA:
	\         REL   ADDR
        \       DATA+00 082: st04BFo1   DATA Row 1
        \       DATA+04 086: o1BFo2o2   DATA Row 2
        \       DATA+08 08A: BFo3o3BF   DATA Row 3
        \       DATA+0C 08E: o4o407__   DATA Row 4

	\ Data for 1-CODE
	\ we reserve the two LAST rows for that data.

	( LABEL: IP: )
	( value )  $07
	( index )  $0E
	data-C!

	( value ) ( st o1 o2 o3 o4 )
	( index )  $02
	data!
	( value ) ( st o1 o2 o3 )
	( index )  $05
	data!
	( value ) ( st o1 o2 )
	( index )  $08
	data!
	( value ) ( st o1 )
	( index )  $0B
	data!

	\ Data for 2-CODE
	\ we reserve the two FIRST rows for that data

	( LABEL: LEITOR )
	( value ) $04
	( index ) $01
	data-C!
	
	( value ) ( st )
	( index )  $00
	data-C!
;
: PROG-S-2
	5 0 0 0 0 S-2*

	\ Code

	( id for label )
	( value  ) $05B4
	( offset ) $00
	1-CODE + !

	( id for IP )
	( value  ) $07B4
	( offset ) $02
	1-CODE + !

	( id for label )
 	( value  ) $05B4
	( offset ) $00
	2-CODE + !
	
	( id for value )
	( value  ) $06B4
	( offset ) $02
	2-CODE + !
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

;
: PROG-S-3
	0 2 6 S-3*

	\ Code

	( id for label )
	( value  ) $05B4
	( offset ) $00
	1-CODE + !

	( id for value )
	( value  ) $06B4
	( offset ) $02
	1-CODE + !

	( id for label )
 	( value  ) $05B4
	( offset ) $00
	2-CODE + !

	( id for MS )
	( value  ) $03B4
	( offset ) $02
	2-CODE + !
;

: S-4* ( m1 v1 m2 v2 m3 v3 m4 v4 -- ; extern, Antennas )

	\ clean!		( we do   use all slots )
	 
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
;
: PROG-S-4
	0 0
	0 0
	0 0
	0 0 S-4*

	\ antenna API call

	( value  ) $08B4
	( offset ) $00
	1-CODE + !

	( value  ) $BBBB
	( offset ) $02
	1-CODE + !

	( value  ) $BBBB
	( offset ) $00
	2-CODE + !

	( value  ) $BBBB
	( offset ) $02
	2-CODE + !
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
;
: PROG-S-5
	0 0 0 0 0 0 S-5*

	\ Code

	( id for label )
	( value  ) $05B4
	( offset ) $00
	1-CODE + !

	( id for IP )
	( value  ) $07B4
	( offset ) $02
	1-CODE + !

	( id for label )
 	( value  ) $05B4
	( offset ) $00
	2-CODE + !

	( id for IP )
	( value  ) $07B4
	( offset ) $02
	2-CODE + !
;

: S-6*	( state ; extern, USB )
  	\ DATA:

	\         REL   ADDR
        \       DATA+00 082: st______   DATA Row 1
        \       DATA+04 086: 0c______   DATA Row 2
        \       DATA+08 08A: ________   DATA Row 3
        \       DATA+0C 08E: ________   DATA Row 4

	clean!

	( LABEL: USB: )
	( value ) $0C
	( index ) $04
	data-C!

	( value ) ( state )
	( index ) $00
	data-C!
;
: PROG-S-6
	2 S-6*

	\ Code

	( id for label )
	( value  ) $05B4
	( offset ) $00
	1-CODE + !

	( id for value )
	( value  ) $06B4
	( offset ) $02
	1-CODE + !

 	( value  ) $BBBB
	( offset ) $00
	2-CODE + !

	( value  ) $BBBB
	( offset ) $02
	2-CODE + !
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
;
: PROG-S-7
	0 S-7*

	\ Code

	( id for label )
	( value  ) $05B4
	( offset ) $00
	1-CODE + !

	( id for BigNum )
	( value  ) $01B4
	( offset ) $02
	1-CODE + !

	( id for label )
	( value  ) $05B4
	( offset ) $00
	2-CODE + !

	( id for BigNum )
	( value  ) $01B4
	( offset ) $02
	2-CODE + !
;

\ TODO redesign waiting, should it be here or in the pc?

\ ======================

\ User input: ( 💀 )

\ +$00     -> 
\ +$01     -> 
\ +$02     -> Action
\ +$03     -> confirm-addr CONST
\ +$05:$11 -> jmp table    CONST
\ +$13     -> Lock input
\ +$14     -> ok-down?
\ +$16     -> arrow-down?

' Button-Data 2 + VALUE BAG

: BAG@ ( slot -- value ) BAG + C@ ;
: BAG! ( value slot -- ) BAG + C! ;

$00 $02 BAG! ( ; action )

: redraw
	BAG 3 + @ C@

	( align  ) $02 *
	( offset ) $05 +
	( calc   ) BAG + @ EXECUTE
;

: restore	( ; restores screen/input when we block them )
	DATA C@ $80 =
	IF
		redraw
	THEN

	BAG $13 + @ C@ $80 =
	IF
		$BB BAG $13 + @ C!
	THEN
; 

: lock-display
	( EXIT opcode )
	( value ) $80
	( index ) $00
	data-C!

	100 DLY
;

: S-A*  ( l1 l2 -- ; extern, Error[$0E-$0F]/Loading[$0D]/Confirmation screen[$13-$14] )

        lock-display

	DUP $14 === NOT IF
		$80 BAG $13 + @ C! \ Block input
	THEN

	\ wait-key

	2 API

	SWAP 5 API 2 API	( l1 l2 )

	DUP 0 = IF		( l2 )
		2 API DROP
	ELSE
		5 API 2 API	( l2 )
	THEN

	2 API

	0 API
; 

: confirm $00 2 BAG! ;
' confirm BAG 3 + !

\ programmers JMP table
' PROG-S-1 BAG $05 + !
' PROG-S-2 BAG $07 + !
' PROG-S-3 BAG $09 + !
' PROG-S-4 BAG $0B + !
' PROG-S-5 BAG $0D + !
' PROG-S-6 BAG $0F + !
' PROG-S-7 BAG $11 + !

: switch-screen
	BAG 3 + @ C@ 1 + 7 MOD DUP	( current-screen +1 MOD7 )
	BAG 3 + @ C!
	$60 OR 2 BAG!
	redraw
;

: do-button	( ; processes button input )
	TAG			( ; LOCK )

	6 IN 0 = DUP    	(               ; detect arrow button )
	7 IN 0 = DUP    	( arrow-down?x2 ; detect ok    button )

	0 BAG@ NOT AND		( arrow-down?x2 ok-down?x2 )
	IF			( arrow-down?x2 ok-down? ok-pressed? )
		confirm
		restore
	THEN
	0 BAG!			( arrow-down?x2 ok-down? )

	1 BAG@ NOT AND		( arrow-down?x2 )
	IF			( arrow-down? arrow-pressed? )
		switch-screen
	THEN
	1 BAG!			( arrow-down? )
;
BAG $13 + !	( [LOCKTAG] ; storing the lock )

\ ======================

\ Extern

: BA@* BAG@    ;
: RST* restore ;

\ ======================

\ Initialization

PROG-S-1	( ; set default screen )

500 DLY

100  1 TMI Dis
10   0 TMI do-button

1 TME

\ ======================

