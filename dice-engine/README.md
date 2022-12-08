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

**Die map**

```clojure
(def die {:id              1
          :value           1
          :faces           3
          :previous-values [2 2]})
```

**Dice-expression map**

`3d4kh2`

```clojure
{:expression    "3d4kh2"
 :roll          {:number-of-dice 3
                 :faces          4}
 :set-operation {:operator :keep
                 :selector :highest
                 :literal  2}}
```

**Result map -**

`3d4`

```clojure
 {:expression    "3d4kh2"
  :roll          {:number-of-dice 3
                  :faces          4}
  :set-operation {:operator :keep
                  :selector :highest
                  :literal  2}
  :outcomes      [{:id              1
                   :value           3
                   :faces           6
                   :discarded       true
                   :previous-values []}
                  {:id              2
                   :value           4
                   :faces           6
                   :discarded       false
                   :previous-values []}
                  {:id              1
                   :value           6
                   :faces           6
                   :discarded       false
                   :previous-values []}]}
``` 

- set operator functions
- set selector functions
- imperative shell
  _Parse dice expression_

- `parse-dice-expression`
    - Returns a map of dice-expr

**Pure functions -**

_Dice roll_

- `get-id`
    - Returns an auto incremented integer starting from 1
- `rand-int-natural [n]`
- `create-die [die-value faces]`
    - Returns a die map of the form:
      `{:id              a
      :value           n
      :faces           x
      :previous-values []}`
- `roll [n, faces]`
    - Returns `n` die maps. Calls `rand-int-natual` and `create-die`
- `reroll [die]`
    - Rerolls a single die. Updates its value and history.

_Set operations_

- `keep [selector literal dice]`
    - Applies selector on dice. Returns dice selected by selector.
- `drop [selector literal dice]`
    - Applies selector on dice. Returns dice un-selected by selector.
- `reroll [selector literal dice]`
    - Applies selector on dice. If no dice are selected, returns dice
      else, rerolls selected dice and appends the result to the un-selected dice.
      Recurses until no dice qualifies selector.

_Set selectors_

The following return a partition of the dice, based on their respective filtering conditions
`[[selected-dice] [unselected-dice]]`

- `highest [literal, dice]`
- `lowest [literal, dice]`
- `greater-than [literal, dice]`
- `lesser-than [literal, dice]`
- `match [literal, dice]`

**Imperative shell and parsing functions-**

The solution is expected to do the following

1. Operations on dice:
    - Roll dice
    - Apply set operators and selectors.
    - Track state of the roll.
    - Compute the final effective value of the roll
2. Compute result of the numeric operations on all the dice

- `roll-value`
    - Returns the sum of all valid dice in `:outcomes` of a roll. Valid dice being `:discarded false`

- `evaluate-roll`
    - The "imperative shell of the program". Takes a dice-expression map and returns the result map for that roll.
        - Rolls n dice of x faces
        - Applies set operation on the dice
        - Marks dice missing from the result as discarded
        - Associates the result to `:outcomes` and returns the result map

- `parse-roll-result`
    - Parses a single roll value in the form - `val`, `~discarded_val2~`
      or `val (~previousrollval~, ~previousrollval'~)`

- `parse-output`
    - Parses output of a single roll in the format \
      `2d6kh1: (~2~, 5) => 5`

### Design choices

**#1**

- Dice operators (drop, keep, reroll) return all dice
- They are responsible for updating dice state as follows
    - keep: all dice that are not selected by selector are marked as discarded
    - drop: all dice that are selected by selector are marked as discarded
    - reroll: all dice selected by selector are rerolled (update previous-values) \
      and recursed until no dice is selected by selector

**#2**

- Each die has an id
- Dice operators only return the result of the set operation and not all dice
- Selector functions return a subset of dice filtered based on the condition
- Operator functions iterate over the original dice list and return the resulting dice

> **Both the above approaches resulted in dice operators having to iterate over the original list of dice.**
**To avoid the redundant iteration, the below were considered.**

**#3**

- To avoid the extra iteration, we can make selectors only return predicates
- These predicates would be used by operators to filter dice
- Dice operators now only return predicates and not a collection of dice
- This wasn't extensible to range based selectors like `highest` and `lowest`. \
  These operations cannot be represented as predicates.
- This approach was later discarded

**#4**

- Operators return a partition of dice, as opposed to a subset of the dice
- The first sequence of the partition contains selected dice and the second the unselected
- Operators can now use either of these values without having to do another iteration on the original dice collection
