# Documentator
create a program that generates documentation from C/C++ comments
(recommended language: Haskell)

## format:

```
    % too permissive?
    % word -> [A-Za-z0-9?$#@!'"-+<>,.;:]+
    
    DOC ::= "~@"

    plaintext ::= {word}

    pexp ::= `(` {expr} `)`
    expr ::= word [pexp]

    stmt ::= DOC  {expr}

    % allow use of reserved words
    stmt ::= DOC  {expr} -- plaintext
```

## example:

### valid:

```diff
+//  ~@ NAME TsukiGva2
+//  ~@ NAME(AUTO)
+//  ~@ COPYRIGHT 2024 Rodrigo Monteiro Junior
+//  ~@ COPYRIGHT(NAME)
+//  ~@ COPYRIGHT(DATE NAME)
+//  ~@ COLOR(#FF00FF) TITLE Cool Title
+//  ~@ NAME -- NAME
```

### invalid:

```diff
-//  ~@ NAME TsukiGva2 
-//  ~@ NAME() TsukiGva2
```

## Configuration:

### format:

```
    path ::= [word "/"] word ["." ext]
    entry ::= word path
```

### example:

```
    NOT_IMPLEMENTED not_implemented.css
    TEST            test.css
```

$CONFIGDIR/not\_implemented.css
```css
    .NOT_IMPLEMENTED {
        color: #ff0000;
    }
```

$CONFIGDIR/test.css
```css
    .TEST {
        color: #f5f500;
    }
```

## Not planned (yet):

- [ ] Custom HTML support

## TODO

- [ ] Test each planned feature (by hand)
- [ ] Find bugs with the current plans
- [ ] Look for better implementation ideass
- [ ] Is this useful?
- [ ] Write code!

