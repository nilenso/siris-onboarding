(ns dice-engine.dice)

(defn rand-int-natural
  "Returns a random integer between 1 to n (both inclusive)."
  [n]
  (+ 1 (rand-int n)))

(defn create-die
  "Returns a die map of the form
  {:value           n
   :discarded       false
   :faces           x
   :previous-values []}"
  [die-value faces]
  (into {:value           die-value
         :discarded       false
         :faces           faces
         :previous-values []}))

(defn sum
  "Takes a list of die maps and returns the sum of the numeric values as an integer."
  [dice]
  (->>
    (map :value dice)
    (reduce +)))

(defn roll
  "Returns a dice-roll map after calling rand [faces] n times"
  [n faces]
  (->> (repeatedly
         n
         #(rand-int-natural faces))
    (map #(create-die % faces))))

(defn discard
  "Applies the selector on :dice with n. Returns a new map with :numeric-value updated
  after discarding the dice that match the result of the selector."
  [dice selector n]
  (let [filtered-dice (selector dice n)]
    (map #(if (some
                #{(:value %)}
                filtered-dice)
            (assoc % :discarded true)
            %)
         dice)))

(defn pick
  "Returns a new map with :numeric-value updated after discarding the dice that do not match n."
  [dice selector n]
  (let [filtered-dice (selector dice n)]
    (map #(if-not (some
                    #{(:value %)}
                    filtered-dice)
            (assoc % :discarded true)
            %)
         dice)))

(defn reroll
  "Returns a new dice-roll map if none of the rerolled dice match n,
  otherwise recurses until none match. History of each die value is appended to
  :previous-values of the die map."
  [dice-roll selector n]
  ;TODO
  )

(defn highest
  "Returns a vector of integer values of the highest n dice.
  If the size of input dice vector (k) is less than n, returns the highest k dice."
  [dice n]
  (->>
    (map :value dice)
    (sort >)
    (take n)))

(defn lowest
  "Returns a vector of integer values of the lowest n dice.
  If the size of input dice vector (k) is less than n, returns the lowest k dice."
  [dice n]
  (->>
    (map :value dice)
    (sort)
    (take n)))

(defn greater-than
  "Returns a vector of integers of value > n. Empty vector if none qualify."
  [dice n]
  (->>
    (map :value dice)
    (filter #(> % n))))

(defn less-than
  "Returns a vector of integers of value > n. Empty vector if none qualify."
  [dice n]
  (->>
    (map :value dice)
    (filter #(< % n))))

(defn match
  "Returns a vector of integers of value = n. Empty vector if none match."
  [dice n]
  (->>
    (map :value dice)
    (filter #(= n %))))
