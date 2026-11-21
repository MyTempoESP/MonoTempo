: lit  ( API_Call    -- )     API fwd ; \ Call an API function and proceed to the next line
: lbl  ( Label       -- ) 5   API     ; \ Display a named Label
: fwd  (             -- ) 2   API     ; \ Proceed to the next line
: num  ( Num         -- ) 4   lit     ; \ Number literal
: val  ( Val_ID      -- ) 6   lit     ; \ Predefined named values (e.g. "OK","CONECTAD","WEB")
: atn  ( a1 a2 a3 a4 -- ) 1   lit     ; \ Display each antenna state
: ip   ( o1 o2 o3 o4 -- ) 7   lit     ; \ Display each specified octet separated by '.'
: ms   ( Millis      -- ) 3   lit     ; \ Display Milliseconds in XXXms or Xs format
: hms  ( HH MM SS    -- ) 256 ip      ; \ Specifying an invalid ip to ip switches to HH:MM:SS format
: <cr> (             -- ) 999 lbl     ; \ if argument to lbl is > LABEL_COUNT then the line is cleared
: tim  (             -- ) 11  lbl     ; \ Shorthand for 11 lbl cuz lazy
