(ns dice-engine.dice
  (:refer-clojure :exclude [drop keep]))

(def id-counter (atom -1))

(defn get-id
  "Returns an auto incremented id"
  []
  (swap! id-counter inc))

(defn rand-int-natural
  "Returns a random integer between 1 to n (both inclusive)."
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
  "Takes a list of die maps and returns the sum of the numeric values as an integer."
  [dice]
  (->>
    (map :value dice)
    (reduce +)))

(defn roll
  "Returns literal dice, each with a value <= faces"
  [literal faces]
  (->>
    (repeatedly literal #(rand-int-natural faces))
    (map #(create-die % faces))))

(defn reroll
  "Rerolls a die"
  [{:keys [value faces previous-values] :as die}]
  (assoc die :value (rand-int-natural faces)
             :previous-values (conj previous-values value)))

(defn drop
  "Applies the selector on dice with literal. Returns a new map
  after discarding the dice that match the result of the selector."
  [dice literal selector]
  (filter
    (complement
      (partial selector literal))
    dice))

(defn keep
  "Returns a new map after discarding the dice that do not match literal."
  [dice literal selector]
  (filter
    (partial selector literal)
    dice))

(defn reroll-matched
  "Rerolls dice filtered by selector and returns the rerolled dice.
   Rerolls until none of the dice can be filtered by selector."
  [dice literal selector]
  (prn dice)
  (let [grouped-dice (group-by (partial selector literal) dice)
        matched-dice (get grouped-dice true)
        rest-dice (get grouped-dice false)]
    (if (empty? matched-dice)
      dice
      (recur (concat (map reroll matched-dice)
                     rest-dice)
             literal
             selector))))

(defn highest
  "Returns a vector of integer values of the highest literal dice.
  If the size of input dice vector (k) is less than literal, returns the highest k dice."
  [literal dice]
  (->>
    dice
    (sort-by :value >)
    (take literal)))

(defn lowest
  "Returns a vector of integer values of the lowest literal dice.
  If the size of input dice vector (k) is less than literal, returns the lowest k dice."
  [literal dice]
  (->>
    dice
    (sort-by :value)
    (take literal)))

(defn greater-than
  "Returns a vector of integers of value > literal. Empty vector if none qualify."
  [literal die]
  (->
    (:value die)
    (> literal)))

(defn lesser-than
  "Returns a vector of integers of value > literal. Empty vector if none qualify."
  [literal die]
  (->
    (:value die)
    (< literal)))

(defn match
  "Returns a vector of integers of value = literal. Empty vector if none match."
  [literal die]
  (->
    (:value die)
    (= literal)))
