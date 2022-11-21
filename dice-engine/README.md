# dice-engine

Solution to the [Dice Engine](https://github.com/nilenso/winter-onboarding-2021/blob/main/doc/dice-engine-problem.md)
problem in Clojure.

### Scope of the problem statement ###

1. Numeric Operations
    1. Unary
    2. Binary

2. Brackets to override natural order od precedence
3. Set & literal set representations
4. Set Operations

| Syntax | Name           | Description                                                                  |
|--------|----------------|------------------------------------------------------------------------------|
| k      | keep           | keep all matched values                                                      |
| d      | drop           | Drop all matched values                                                      |
| rr     | reroll         | Rerolls all matched values until none match.                                 |

5. Set selectors

| Syntax | Name              | Description                                   |
|--------|-------------------|-----------------------------------------------|
| X      | literal           | All values that are literally this value      |
| hX     | highest X         | The highest X values in the set               |
| lX     | lowest X          | The lowest X values in the set                |
| \>X    | greater than X    | Values in the set that are greater than X     |
| \<X    | lesser than X     | Values in the set that are lesser than X      |

### Solution ###

**Dice-roll notation:** `NdX`

**Die map representation**

```clojure
(def die {:value           1
          :discarded       false
          :faces           3
          :previous-values [2 2]})
```

**Dice-roll map representation -**

`3d4`

```clojure
(def dice-roll {:numeric-value 8
                :dice          [d1 d2 d3]})  ; d1, d2 and d3 are die maps of face value 4, 3 & 1
``` 

**Dice-expression map representation**

`3d4kh2`

```clojure
(def dice-expr {:expr         "3d4kh2"
                :set-selector highest       ; highest and keep are functions
                :set-operator keep})
```

_Parse dice expression_

- `parse-dice-expression`
    - Returns a map of dice-expr

_Dice roll_

- `rand [n]`
- `roll [n, faces]`
    - Returns a dice-roll map after calling `rand [faces]` n times

**Pure functions -**

_Set operations_

- `keep [dice-roll n]`
    - Returns a new dice-roll map with `:numeric-value` updated \
      and `discarded` updated to true where die value does not match `n`.
- `drop [dice-roll n]`
    - Returns a new dice-roll map with `:numeric-value` updated \
      and `discarded` updated to true where die value matches `n`.
- `reroll [dice-roll n]`
    - Returns a new dice-roll map if none of the rerolled dice match `n`,\
      otherwise recurses until none match. `previous-values` are updated for each dice \
      on every roll.

_Set selectors_

(Return a new dice-roll map with `:numeric-value` updated \
and `discarded` updated to true where die value matches `n`)

- `highest [x, dice]`
- `lowest [x, dice]`
- `greater-than [x, dice]`
- `less-than [x, dice]`

**Imperative shell -**

1. `parse-input`
    - represent the input expression as a tree
        - numeric operators (`+`/`-`/`*`/`/`) are parent nodes
        - child nodes can either be a `dice-roll` map or a number
2. Solve the expressions bottom-up
    - Order of execution
        1. parse dice expression - `parse-dice-expression`
        2. roll the dice
        3. if `:set-selector` is not nil, call set-selector function
        4. if `:set-operator` is not nil, pass the result of set-selector to set-operator
