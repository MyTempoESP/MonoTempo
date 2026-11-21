#ifndef __aa2_FTH__
#define __aa2_FTH__
const char code[] PROGMEM =
#define NL "\n"
": w0 5 API ;" NL
// @210/v40	->	addr-Antenna
// @201/v39	->	addr-Ms
// @238/v38	->	addr-Ip
// @248/v37	->	addr-Value
// @118/v36	->	addr-Number
// @270/v35	->	addr-BigNum
// @114/v34	->	addr-Forward
// @267/v33	->	addr-Label
// @132/v32	->	DATA
// @264/v31	->	1-CODE
// @265/v30	->	2-CODE
// @250/v29	->	VERSION
// @241/w28	->	S-7
// @229/w27	->	S-6
// @198/w26	->	S-5
// @189/w25	->	S-4
// @172/w24	->	S-3
// @148/w23	->	S-2
// @130/w22	->	S-1
// @225/w21	->	code-clean!
// @234/w20	->	clean!
// @254/w19	->	program-calls!
// @188/w18	->	call-idx!
// @94/w17	->	call!
// @81/w16	->	encode-call!
// @85/w15	->	prepare-call
// @183/w14	->	data-Big!
// @242/w13	->	data-C!
// @243/w12	->	data!
// @87/w11	->	calc-align
// @68/w10	->	F!
// @76/w9	->	16-bit-encode!
// @246/w8	->	Dis
// @93/w7	->	helper-Antenna
// @91/w6	->	Ms
// @89/w5	->	Ip
// @87/w4	->	Value
// @85/w3	->	Number
// @83/w2	->	BigNum
// @81/w1	->	Forward
// @79/w0	->	Label
": w1 2 API ;" NL
": w2 1 API ;" NL
": w3 4 API ;" NL
": w4 6 API ;" NL
": w5 7 API ;" NL
": w6 3 API ;" NL
": w7" NL
	 "7 6 API  1 w3  8 6 API  w2  6 6 API " NL
	 "7 6 API  2 w3  8 6 API  w2 w1" NL
	 "7 6 API  3 w3  8 6 API  w2  6 6 API " NL
	 "7 6 API  4 w3  8 6 API  w2 w1" NL
";" NL
"501 VAL v29" NL
": w8" NL
	"NOP NOP NOP NOP NOP NOP NOP NOP NOP NOP NOP NOP NOP NOP NOP NOP" NL
	"0 w0" NL
	"v29 w3" NL
	 "w1" NL
	"TAG NOP NOP NOP w1" NL
	"TAG NOP NOP NOP w1" NL
	"3 w0 NOP" NL
	"0 w4" NL
	 "w1" NL
	"0 API" NL
";" NL
"VAL v30" NL
"VAL v31" NL
"' w8 VAL v32" NL
": w9" NL
	"$BF" NL
	"OVR C!" NL
	"1 +" NL
	 "!" NL
";" NL
": w10" NL
	"DUP >R" NL
	"w9" NL
	"R>" NL
	"3 +" NL
	 "C!" NL
";" NL
": w11" NL
	">R" NL
	 "*" NL
	"R>" NL
	 "+" NL
";" NL
": w12" NL
	"v32 + w9" NL
";" NL
": w13" NL
	"v32 + C!" NL
";" NL
": w14" NL
	"4" NL
	"v32 w11" NL
	 "w10" NL
";" NL
": w15" NL
	">R" NL
	"$C000 OR" NL
	"R>" NL
";" NL
": w16" NL
	"w15" NL
	"w9" NL
";" NL
": w17" NL
	"w15" NL
	"!" NL
";" NL
": w18" NL
	">R" NL
	"2 R>" NL
	"w11" NL
	"w17" NL
";" NL
"' w0 VAL v33" NL
"' w1 VAL v34" NL
"' w2 VAL v35" NL
"' w3 VAL v36" NL
"' w4 VAL v37" NL
"' w5 VAL v38" NL
"' w6 VAL v39" NL
"' w7 VAL v40" NL
": w19" NL
	"DUP >R >R" NL
	"1 R> w18" NL
	"0 R> w18" NL
";" NL
": w20" NL
	"$bbbb v32 !" NL
	"$bbbb v32 $2 + !" NL
	"$bbbb v32 $4 + !" NL
	"$bbbb v32 $6 + !" NL
	"$bbbb v32 $8 + !" NL
	"$bbbb v32 $C + !" NL
	"$bbbb v32 $A + !" NL
	"$bbbb v32 $E + !" NL
";" NL
": w21" NL
	"$bbbb v31 !" NL
	"$bbbb v31 $2 + !" NL
	"$bbbb v30 !" NL
	"$bbbb v30 $2 + !" NL
";" NL
": S-1*" NL
	"w20" NL
	 "$02" NL
	 "$0C" NL
	"w13" NL
	 "$02" NL
	"w14" NL
	 "$01" NL
	 "$04" NL
	"w13" NL
	 "$00" NL
	"w14" NL
	"v33" NL
	"v35 v31 w19" NL
	"v33" NL
	"v35 v30 w19" NL
";" NL
": S-2*" NL
	"w20" NL
	 "$07" NL
	 "$0C" NL
	"w13" NL
	 "$0B" NL
	"w13" NL
	 "$0A" NL
	"w13" NL
	 "$09" NL
	"w13" NL
	 "$08" NL
	"w13" NL
	 "$04" NL
	 "$01" NL
	"w13" NL
	 "$00" NL
	"w13" NL
	"v33" NL
	"v38 v31 w19" NL
	"v33" NL
	"v37 v30 w19" NL
";" NL
": S-3*" NL
	"w20" NL
	 "$0C" NL
	"w13" NL
	 "$08" NL
	"w13" NL
	 "$0A" NL
	 "$04" NL
	"w13" NL
	 "$00" NL
	"w12" NL
	"v33" NL
	"v37 v31 w19" NL
	"v33" NL
	"v39 v30 w19" NL
";" NL
": S-4*" NL
	"w21" NL
	"0 w14" NL
	"1 w14" NL
	"2 w14" NL
	"3 w14" NL
	"v40 0 v31 w18" NL
";" NL
": S-5*" NL
	 "$16" NL
	 "$08" NL
	"w13" NL
	 "$02" NL
	"w12" NL
	 "$01" NL
	"w13" NL
	 "$00" NL
	"w13" NL
	 "$DA7E" NL
	 "$05" NL
	"w12" NL
	 "$0B" NL
	 "$0F" NL
	"w13" NL
	 "$100" NL
	 "$0C" NL
	"w12" NL
	 "$0B" NL
	"w13" NL
	 "$0A" NL
	"w13" NL
	 "$09" NL
	"w13" NL
	"v33" NL
	"v38 v31 w19" NL
	"v33" NL
	"v38 v30 w19" NL
";" NL
": S-6*" NL
	"w20" NL
	"w21" NL
	 "$0C" NL
	 "$04" NL
	"w13" NL
	 "$00" NL
	"w13" NL
	"v33" NL
	"v37 v31 w19" NL
";" NL
": S-7*" NL
	"w20" NL
	 "$10" NL
	 "$05" NL
	"w13" NL
	 "$10" NL
	 "$04" NL
	"w13" NL
	 "$CA00 v29 1 RSH OR" NL
	 "$00" NL
	"w12" NL
	 "$12" NL
	 "$0D" NL
	"w13" NL
	 "$10" NL
	 "$0C" NL
	"w13" NL
	 "$08" NL
	"w12" NL
	"v33" NL
	"v35 v31 w19" NL
	"v33" NL
	"v35 v30 w19" NL
";" NL
"1500 DLY" NL
"50 0 TMI w8" NL
"1 TME" NL
;
#endif
