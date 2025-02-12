 \  Globals.fth

  501 VALUE VERSION   \  version
  
  \  Generic button interface definition
  
  \   IFC b
  \     HAS AC 
  \          ST 
  \          RL 
  \          Pr 
  \          IN ; 

  \   b m         
  \   b a         

  \  Class consts ( they use property capitalization )
  6 VALUE mIN       
  7 VALUE aIN       

  VARIABLE Scr          \  Screen Number

\  Button.fth

  \  Properties                 CLASS  PROP     DESC
  VARIABLE mRL          \  m ->  MAIN , m.RL, MAIN RELEASE
  VARIABLE Mac          \  (EXP) MAIN , M.ac, MAIN ACTION
  VARIABLE mST          \        MAIN , m.ST, MAIN STATE

  VARIABLE Aac          \  a ->  ALT  , a.AC, ALT  ACTION
  VARIABLE aST          \        ALT  , a.ST, ALT  STATE

  \  Methods                    CLASS  METHOD   DESC
  : mPr              \  M.PR(--)
     mIN IN 0 = ;  \        MAIN , m.Pr, MAIN PRESSED
  : aPr              \  A.PR(--)
     aIN IN 0 = ;  \        ALT  , a.Pr, ALT  PRESSED

  : b?               \  B.? ( b.Pr b.ST -- b.RL b.ST' )
     DUP ROT SWAP     \  Desc: returns the truth
     NOT AND ;     \  value of b.RL, i.e. checks if
  \                       the button has been PR + RL
  \                       (RELEASED), from the current
  \                       state vs the previous one, i.e.
  \                       b.Pr vs b.ST.
  \                         - Also returns the new value
  \                       for b.ST, which is the return
  \                       from b.PR.

  : sWi              \  SWI ( m.RL -- )
     IF              \  Desc: Switches screens if
       Scr @ 1 +     \  m.RL is set.
       7 MOD
       Scr ! THEN ;

\  Screen.fth

  : lBl  5   API ;
  : fWd  2   API ;
  : fNm  1   API ;
  : nUm  4   API ;
  : vAl  6   API ;
  : iP   7   API ;
  : mS   3   API ;
  \ : hMs  256 iP  ;
  \ : uSb  12  lBl ;
  \ : tIm  11  lBl ;
  
  \  Text Decorations
#define A      7 6 API 
#define SPACE  6 6 API 
#define COLON  8 6 API 

  \  escovando bit
  \  Antenna Data
  : aTn \  ( N Mag N Mag N Mag N Mag -- )
    A 1 nUm COLON fNm SPACE
    A 2 nUm COLON fNm fWd
  /**/              
    A 3 nUm COLON fNm SPACE
    A 4 nUm COLON fNm fWd
   ;              

  \  Display memory
  : Dis

    \  Data: 16 bytes
     NOP NOP NOP NOP \  Fits a 16-bit number
     NOP NOP NOP NOP \  in the form $bf[nnnn].
     NOP NOP NOP NOP
     NOP NOP NOP NOP
    /**/            

    \  Heading: 9 bytes
     0 lBl VER nUm fWd \  fWd ($8D)

    \  Screen code: 7 bytes each
     TAG NOP NOP NOP NOP fWd
     TAG NOP NOP NOP NOP fWd

    \  Heading: 9 bytes
     3 lBl 0 vAl NOP fWd

     0 API
   ;              

   ' Dis VALUE DAT  
     R>  VALUE LB2  
     R>  VALUE LB1  

  : enc              \  ENC ( addr n -- )
     $BF OVR C! 1 + !\  encodes a number to a LIT instruction
  \                       in the specified address.
   ;              

  : scr!
   DAT 4 * + DUP enc
   ;              

  : Atn
     1 - FOR I scr! NEXT DAT enc
     TAG NOP NOP   \  $bf[0000], 16-bit literal
     LB1 !           \  Placeholder for addr of aTn + CALL
   ;              
  
  $C000 ' aTn OR R>  \  make a Call instruction: $c0|[ADDR],
   enc            

\  Timers.fth
  500 DLY          \  Wait until everything is loaded
  \ 10 0 TMI bUp     \  Init button checking
  50 0 TMI Dis     \  Init display
  1 TME            \  Init timers
;
