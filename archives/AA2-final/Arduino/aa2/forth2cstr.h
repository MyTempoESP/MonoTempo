#ifndef __aa2_FTH__
#define __aa2_FTH__
const char code[] PROGMEM =
#define NL "\n"
": w0" NL
// @381/v37	->	BAG
// @319/v36	->	DATA
// @299/v35	->	1-CODE
// @304/v34	->	2-CODE
// @278/v33	->	VERSION
// @378/w32	->	RST
// @376/w31	->	BA@
// @380/w30	->	do-button
// @369/w29	->	switch-screen
// @362/w28	->	confirm
// @325/w27	->	S-A
// @325/w26	->	lock-display
// @371/w25	->	restore
// @350/w24	->	redraw
// @365/w23	->	BAG!
// @367/w22	->	BAG@
// @342/w21	->	PROG-S-7
// @257/w20	->	S-7
// @339/w19	->	PROG-S-6
// @232/w18	->	S-6
// @336/w17	->	PROG-S-5
// @189/w16	->	S-5
// @333/w15	->	PROG-S-4
// @163/w14	->	S-4
// @330/w13	->	PROG-S-3
// @134/w12	->	S-3
// @327/w11	->	PROG-S-2
// @98/w10	->	S-2
// @356/w9	->	PROG-S-1
// @68/w8	->	S-1
// @245/w7	->	clean!
// @159/w6	->	data-Big!
// @300/w5	->	data-C!
// @259/w4	->	data!
// @52/w3	->	F!
// @48/w2	->	Encode!
// @350/w1	->	Dis
// @272/w0	->	Button-Data
	"NOP NOP NOP NOP" NL
	"NOP NOP NOP NOP" NL
	"NOP NOP NOP NOP" NL
	"NOP NOP NOP NOP" NL
	"NOP NOP NOP NOP" NL
	"NOP NOP NOP NOP" NL
	"NOP NOP NOP NOP" NL
	"NOP NOP NOP NOP" NL
	"NOP NOP NOP NOP" NL
	"NOP NOP NOP NOP" NL
	"NOP NOP NOP NOP" NL
	"NOP NOP NOP NOP" NL
	"NOP NOP NOP NOP" NL
";" NL
"501 VAL v33" NL
": w1" NL
	"NOP NOP NOP NOP NOP NOP NOP NOP NOP NOP NOP NOP NOP NOP NOP NOP" NL
	"0 5 API" NL
	"v33 4 API" NL
	 "2 API" NL
	"TAG NOP NOP NOP 2 API" NL
	"TAG NOP NOP NOP 2 API" NL
	"3 5 API NOP" NL
	"0 6 API" NL
	 "2 API" NL
	 "0 API" NL
";" NL
"VAL v34" NL
"VAL v35" NL
"' w1 VAL v36" NL
": w2" NL
	"$BF" NL
	"OVR C!" NL
	"1 +" NL
	 "!" NL
";" NL
": w3" NL
	"DUP >R" NL
	"w2" NL
	"R>" NL
	"3 +" NL
	 "C!" NL
";" NL
": w4" NL
	"v36 + w2" NL
";" NL
": w5 v36 + C! ;" NL
": w6 4 * v36 + w3 ;" NL
": w7" NL
	"$bbbb v36 !" NL
	"$bbbb v36 $2 + !" NL
	"$bbbb v36 $4 + !" NL
	"$bbbb v36 $6 + !" NL
	"$bbbb v36 $8 + !" NL
	"$bbbb v36 $C + !" NL
	"$bbbb v36 $A + !" NL
	"$bbbb v36 $E + !" NL
";" NL
": S-1*" NL
	"w7" NL
	 "$02" NL
	 "$0C" NL
	"w5" NL
	 "$02" NL
	"w6" NL
	 "$01" NL
	 "$04" NL
	"w5" NL
	 "$00" NL
	"w6" NL
";" NL
": w9" NL
	"0 0 0 0 S-1*" NL
	 "$05B4" NL
	 "$00" NL
	"v35 + !" NL
	 "$01B4" NL
	 "$02" NL
	"v35 + !" NL
 	 "$05B4" NL
	 "$00" NL
	"v34 + !" NL
	 "$01B4" NL
	 "$02" NL
	"v34 + !" NL
";" NL
": S-2*" NL
	"w7" NL
	 "$07" NL
	 "$0E" NL
	"w5" NL
	 "$02" NL
	"w4" NL
	 "$05" NL
	"w4" NL
	 "$08" NL
	"w4" NL
	 "$0B" NL
	"w4" NL
	 "$04" NL
	 "$01" NL
	"w5" NL
	 "$00" NL
	"w5" NL
";" NL
": w11" NL
	"5 0 0 0 0 S-2*" NL
	 "$05B4" NL
	 "$00" NL
	"v35 + !" NL
	 "$07B4" NL
	 "$02" NL
	"v35 + !" NL
 	 "$05B4" NL
	 "$00" NL
	"v34 + !" NL
	 "$06B4" NL
	 "$02" NL
	"v34 + !" NL
";" NL
": S-3*" NL
	"w7" NL
	 "$0C" NL
	"w5" NL
	 "$08" NL
	"w5" NL
	 "$0A" NL
	 "$04" NL
	"w5" NL
	 "$00" NL
	"w4" NL
";" NL
": w13" NL
	"0 2 6 S-3*" NL
	 "$05B4" NL
	 "$00" NL
	"v35 + !" NL
	 "$06B4" NL
	 "$02" NL
	"v35 + !" NL
 	 "$05B4" NL
	 "$00" NL
	"v34 + !" NL
	 "$03B4" NL
	 "$02" NL
	"v34 + !" NL
";" NL
": S-4*" NL
	"0 w6" NL
	"1 w6" NL
	"2 w6" NL
	"3 w6" NL
";" NL
": w15" NL
	"0 0" NL
	"0 0" NL
	"0 0" NL
	"0 0 S-4*" NL
	 "$08B4" NL
	 "$00" NL
	"v35 + !" NL
	 "$BBBB" NL
	 "$02" NL
	"v35 + !" NL
	 "$BBBB" NL
	 "$00" NL
	"v34 + !" NL
	 "$BBBB" NL
	 "$02" NL
	"v34 + !" NL
";" NL
": S-5*" NL
	 "$16" NL
	 "$08" NL
	"w5" NL
	 "$02" NL
	"w4" NL
	 "$01" NL
	"w5" NL
	 "$00" NL
	"w5" NL
	 "$DA7E" NL
	 "$05" NL
	"w4" NL
	 "$0B" NL
	 "$0F" NL
	"w5" NL
	 "$100" NL
	 "$0C" NL
	"w4" NL
	 "$0B" NL
	"w5" NL
	 "$0A" NL
	"w5" NL
	 "$09" NL
	"w5" NL
";" NL
": w17" NL
	"0 0 0 0 0 0 S-5*" NL
	 "$05B4" NL
	 "$00" NL
	"v35 + !" NL
	 "$07B4" NL
	 "$02" NL
	"v35 + !" NL
 	 "$05B4" NL
	 "$00" NL
	"v34 + !" NL
	 "$07B4" NL
	 "$02" NL
	"v34 + !" NL
";" NL
": S-6*" NL
	"w7" NL
	 "$0C" NL
	 "$04" NL
	"w5" NL
	 "$00" NL
	"w5" NL
";" NL
": w19" NL
	"2 S-6*" NL
	 "$05B4" NL
	 "$00" NL
	"v35 + !" NL
	 "$06B4" NL
	 "$02" NL
	"v35 + !" NL
 	 "$BBBB" NL
	 "$00" NL
	"v34 + !" NL
	 "$BBBB" NL
	 "$02" NL
	"v34 + !" NL
";" NL
": S-7*" NL
	"w7" NL
	 "$10" NL
	 "$05" NL
	"w5" NL
	 "$10" NL
	 "$04" NL
	"w5" NL
	 "$CA00 v33 1 RSH OR" NL
	 "$00" NL
	"w4" NL
	 "$12" NL
	 "$0D" NL
	"w5" NL
	 "$10" NL
	 "$0C" NL
	"w5" NL
	 "$08" NL
	"w4" NL
";" NL
": w21" NL
	"0 S-7*" NL
	 "$05B4" NL
	 "$00" NL
	"v35 + !" NL
	 "$01B4" NL
	 "$02" NL
	"v35 + !" NL
	 "$05B4" NL
	 "$00" NL
	"v34 + !" NL
	 "$01B4" NL
	 "$02" NL
	"v34 + !" NL
";" NL
"' w0 2 + VAL v37" NL
": w22 v37 + C@ ;" NL
": w23 v37 + C! ;" NL
"$00 $02 w23" NL
": w24" NL
	"v37 3 + @ C@" NL
	 "$02 *" NL
	 "$05 +" NL
	 "v37 + @ EXE" NL
";" NL
": w25" NL
	"v36 C@ $80 =" NL
	"IF" NL
		"w24" NL
	"THN" NL
	"v37 $13 + @ C@ $80 =" NL
	"IF" NL
		"$BB v37 $13 + @ C!" NL
	"THN" NL
";" NL
": w26" NL
	 "$80" NL
	 "$00" NL
	"w5" NL
	"100 DLY" NL
";" NL
": S-A" NL
 "w26" NL
	"DUP $14 === NOT IF" NL
		"$80 v37 $13 + @ C!" NL
	"THN" NL
	"2 API" NL
	"SWP 5 API 2 API" NL
	"DUP 0 = IF" NL
		"2 API DRP" NL
	"ELS" NL
		"5 API 2 API" NL
	"THN" NL
	"2 API" NL
	"0 API" NL
";" NL
": w28 $00 2 w23 ;" NL
"' w28 v37 3 + !" NL
"' w9 v37 $05 + !" NL
"' w11 v37 $07 + !" NL
"' w13 v37 $09 + !" NL
"' w15 v37 $0B + !" NL
"' w17 v37 $0D + !" NL
"' w19 v37 $0F + !" NL
"' w21 v37 $11 + !" NL
": w29" NL
	"v37 3 + @ C@ 1 + 7 MOD DUP" NL
	"v37 3 + @ C!" NL
	"$60 OR 2 w23" NL
	"w24" NL
";" NL
": w30" NL
	"TAG" NL
	"6 IN 0 = DUP" NL
	"7 IN 0 = DUP" NL
	"0 w22 NOT AND" NL
	"IF" NL
		"w28" NL
		"w25" NL
	"THN" NL
	"0 w23" NL
	"1 w22 NOT AND" NL
	"IF" NL
		"w29" NL
	"THN" NL
	"1 w23" NL
";" NL
"v37 $13 + !" NL
": BA@ w22 ;" NL
": RST w25 ;" NL
"w9" NL
"500 DLY" NL
"100 1 TMI w1" NL
"10 0 TMI w30" NL
"1 TME" NL
;
#endif
