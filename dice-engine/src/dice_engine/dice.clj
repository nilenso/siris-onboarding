(ns dice-engine.dice)

(defn drop
  "Returns a new map with :numeric-value updated after discarding the dice that match n."
  [dice-roll & n]
  ())

(defn keep
  "Returns a new map with :numeric-value updated after discarding the dice that do not match n."
  [dice-roll & n]
  )

(defn reroll
  "Returns a new dice-roll map if none of the rerolled dice match n,
  otherwise recurses until none match. History of each die value is appended to
  :previous-values of the die map"
  [dice-roll & n]
  )

(defn highest
  "Returns a vector of integer values of the highest n dice.
  If the size of input dice vector(k) is less than n, returns the highest k dice"
  [dice n]
  )

(defn lowest
  "Returns a vector of integer values of the lowest n dice.
  If the size of input dice vector(k) is less than n, returns the lowest k dice"
  [dice n]
  )

(defn greater-than
  "Returns a vector of integers of value > n. Empty vector if none qualify."
  [dice n]
  )

(defn less-than
  "Returns a vector of integers of value > n. Empty vector if none qualify"
  [dice n]
  )