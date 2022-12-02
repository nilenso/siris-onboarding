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
  "Returns n dice, each with a value <= faces"
  [n faces]
  (->>
    (repeatedly n #(rand-int-natural faces))
    (map #(create-die % faces))))

(defn reroll
  "Rerolls a die"
  [{:keys [value faces previous-values] :as die}]
  (assoc die :value (rand-int-natural faces)
             :previous-values (conj previous-values value)))

(defn drop
  "Applies the selector on dice with n. Returns a new map
  after discarding the dice that match the result of the selector."
  [dice n selector]
  (let [filtered-ids (->>
                       (selector n dice)
                       (map :id))]
    (filter (fn [die] (->
                        (some
                          #(= (:id die) %)
                          filtered-ids)
                        (not)))
            dice)))

(defn keep
  "Returns a new map after discarding the dice that do not match n."
  [dice n selector]
  (selector n dice))

(defn reroll-matched
  "Rerolls dice filtered by selector and returns the rerolled dice.
   Rerolls until none of the dice can be filtered by selector."
  [dice n selector]
  (let [filtered-ids (->>
                       (selector n dice)
                       (map :id))]
    (if (empty? filtered-ids)
      dice
      (recur (map
               (fn [die]
                 (if (some
                       #(= (:id die) %)
                       filtered-ids)
                   (reroll die)
                   die))
               dice)
             n
             selector))))

(defn highest
  "Returns a vector of integer values of the highest n dice.
  If the size of input dice vector (k) is less than n, returns the highest k dice."
  [n dice]
  (->>
    dice
    (sort-by :value >)
    (take n)))

(defn lowest
  "Returns a vector of integer values of the lowest n dice.
  If the size of input dice vector (k) is less than n, returns the lowest k dice."
  [n dice]
  (->>
    dice
    (sort-by :value)
    (take n)))

(defn greater-than
  "Returns a vector of integers of value > n. Empty vector if none qualify."
  [n dice]
  (filter
    #(->
       (:value %)
       (> n))
    dice))

(defn lesser-than
  "Returns a vector of integers of value > n. Empty vector if none qualify."
  [n dice]
  (filter
    #(->
       (:value %)
       (< n))
    dice))

(defn match
  "Returns a vector of integers of value = n. Empty vector if none match."
  [n dice]
  (filter
    #(->
       (:value %)
       (= n))
    dice))

