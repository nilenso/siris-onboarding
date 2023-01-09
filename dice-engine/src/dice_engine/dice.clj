(ns dice-engine.dice
  (:refer-clojure :exclude [drop keep]))

(def id-counter (atom 0))

(defn get-id
  "Returns an auto incremented id"
  []
  (swap! id-counter inc))

(defn rand-int-natural
  "Returns a random integer between 1 to n (both inclusive)"
  [n]
  (+ 1 (rand-int n)))

(defn create-die
  "Returns a die map of the form
  {:id              (get-id)
   :value           n
   :faces           x
   :previous-values []}"
  [die-value faces]
  {:id              (get-id)
   :value           die-value
   :faces           faces
   :previous-values []})

(defn sum
  "Takes a list of die maps and returns the sum of the values"
  [dice]
  (->>
    dice
    (map :value)
    (reduce +)))

(defn roll
  "Returns n dice, each with a value <= faces"
  [n faces]
  (->>
    (repeatedly n #(rand-int-natural faces))
    (map #(create-die % faces))))

(defn reroll
  "Rerolls a die"
  [{:keys [value faces previous-values] :as die}]
  (assoc die :value (rand-int-natural faces)
             :previous-values (cons value previous-values)))

(defn drop
  "Applies the selector on dice with literal.
  Returns dice unselected by selector."
  [selector literal dice]
  (->>
    dice
    (selector literal)
    (second)))

(defn keep
  "Applies the selector on dice with literal.
  Returns dice selected by selector."
  [selector literal dice]
  (->>
    dice
    (selector literal)
    (first)))

(defn reroll-matched
  "Rerolls dice selected by selector and returns the rerolled dice.
   Rerolls until none of the dice qualify the selector."
  [selector literal dice]
  (let [[matched-dice rest-dice] (selector literal dice)]
    (if (empty? matched-dice)
      dice
      (recur
        selector
        literal
        (concat (map reroll matched-dice) rest-dice)))))

(defn highest
  "Splits dice into two collections, by the first literal highest values and the rest"
  [literal dice]
  (->>
    dice
    (sort-by :value >)
    (split-at literal)))

(defn lowest
  "Splits dice into two collections, by the first literal lowest values and the rest"
  [literal dice]
  (->>
    dice
    (sort-by :value)
    (split-at literal)))

(defn partition-dice-by-value
  "Splits dice into two collections, grouping them by applying pred on a die's value"
  [pred dice]
  (let [grouped-dice (group-by
                       #(->
                          (:value %)
                          (pred))
                       dice)]
    (vector
      (get grouped-dice true [])
      (get grouped-dice false []))))

(defn greater-than
  "Splits dice into two collections, the first with values greater than literal and the second with the rest"
  [literal dice]
  (partition-dice-by-value
    #(> % literal)
    dice))

(defn lesser-than
  "Splits dice into two collections, the first with values lesser than literal and the second with the rest"
  [literal dice]
  (partition-dice-by-value
    #(< % literal)
    dice))

(defn match
  "Splits dice into two collections, the first with values that match literal and the second with the rest"
  [literal dice]
  (partition-dice-by-value
    #(= literal %)
    dice))

