(ns dice-engine.dice-roller
  (:require [clojure.string :as string]
            [dice-engine.dice :as dice]))

(comment {:faces
          :number-of-dice
          :selector
          :operator})

(def dice-operators {:keep   dice/keep
                     :drop   dice/drop
                     :reroll dice/reroll-matched})

(def dice-selectors {:highest      dice/highest
                     :lowest       dice/lowest
                     :greater-than dice/greater-than
                     :less-than    dice/lesser-than
                     :match        dice/match})

(defn roll-value
  [outcomes]
  (->>
    (filter #(not (:discarded %)) outcomes)
    (map :value)
    (reduce +)))

(defn parse-roll-result
  "Parses a dice roll to in the format:
  (val1, ~discarded_val2~, val3 (~previousrollval3~, ~previousrollval3'~))"
  [{:keys [value discarded previous-values]}]
  (cond
    discarded (str "~" value "~")
    (not-empty previous-values) (str
                                  value
                                  " "
                                  "(" (->>
                                        (map #(str "~" % "~") previous-values)
                                        (string/join ", "))
                                  ")")
    :else (str value)))

(defn parse-output
  [{:keys [outcomes]}]
  (str "("
       (->> (map #(parse-output %) outcomes)
         (string/join ", "))
       ")"))

(defn evaluate-roll
  "docstring"
  [{:keys                               [roll]
    {:keys [operator selector literal]} :set-operation}]
  (let [operator (operator dice-operators)
        selector (selector dice-selectors)
        roll-values (dice/roll
                      (:number-of-dice roll)
                      (:faces roll))
        result (operator roll-values literal selector)]
    (prn roll-values)
    (map
      (fn [{:keys [id] :as die}]
        (if-let [result-die
                 (->> result
                   (filter #(= (:id %) id))
                   (first))]
          (assoc result-die :discarded false)
          (assoc die :discarded true)))
      roll-values)))

(def dice-roll {:r1 {:expression    "2d6kh1"
                     :roll          {:number-of-dice 2
                                     :faces          6}
                     :set-operation {:operator :keep
                                     :selector :highest
                                     :literal  1}}
                :r2 {:expression    "3d20rr1"
                     :roll          {:number-of-dice 3
                                     :faces          20}
                     :set-operation {:operator :reroll
                                     :selector :match
                                     :literal  8}}})

(defn print-output
  "Takes a map of dice rolls and prints the result of each dice set operation and the operations"
  [dice-roll operations]
  (let [rolls (map (fn [[_ value :as roll]]
                     (assoc roll :outcomes (evaluate-roll value)))
                   dice-roll)]
    rolls))



; Tests
(def dice-roll-1 {:roll          {:number-of-dice 2
                                  :faces          6}
                  :set-operation {:operator :keep
                                  :selector :highest
                                  :literal  1}})

(def dice-roll-1-intermediate {:roll          {:number-of-dice 2
                                               :faces          6}
                               :set-operation {:operator :keep
                                               :selector :highest
                                               :literal  1}
                               :outcomes      [{:value           3
                                                :discarded       true
                                                :previous-values []}
                                               {:value           4
                                                :discarded       false
                                                :previous-values []}]})



(def dice-roll-2 {:roll          {:number-of-dice 3
                                  :faces          4}
                  :set-operation {:operator :drop
                                  :selector :lowest
                                  :literal  2}})

(def numerals-1 2)
