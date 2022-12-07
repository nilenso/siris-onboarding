(ns dice-engine.dice)

(defn drop
  "Takes a dice map, discards dice that match n and returns a new dice map."
  [dice-roll & n]
  ())

(defn keep
  "Takes a dice map, keeps dice that match n and returns a new dice map."
  [dice-roll & n]
  )

(defn reroll
  "Takes a dice map, re-rolls dice that match n.
  Returns a new dice map if none of the rerolled dice match n,re-rolls until none match.
  Updates history of each die on every roll."
  [dice-roll & n]
  )

(defn highest
  "Returns a vector of integer values of the highest n dice.
  If the size of input dice vector (k) is less than n, returns the highest k dice"
  [dice n]
  )

(defn lowest
  "Returns a vector of integer values of the lowest n dice.
  If the size of input dice vector (k) is less than n, returns the lowest k dice"
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