(ns dice-engine.dice-roller
  (:require [clojure.string :as string]
            [dice-engine.dice :as dice]))

(def dice-operators {:keep   dice/keep
                     :drop   dice/drop
                     :reroll dice/reroll-matched})

(def dice-selectors {:highest      dice/highest
                     :lowest       dice/lowest
                     :greater-than dice/greater-than
                     :less-than    dice/lesser-than
                     :match        dice/match})

(defn roll-value
  "Returns sum of all valid dice values"
  [dice]
  (->
    (filter #(not (:discarded %)) dice)
    (dice/sum)))

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
  "TODO: prefix and suffix 2d6kh1: (~2~, 5) => 5"
  [{:keys [expression outcomes]}]
  (str expression ": "
       "("
       (->> (map #(parse-roll-result %) outcomes)
         (string/join ", "))
       ")"
       " => "
       (roll-value outcomes)))

(defn evaluate-roll
  "Evaluates a dice roll, applies set operations and returns a dice collection"
  [{:keys                               [roll]
    {:keys [operator selector literal]} :set-operation}]
  (let [operator-fn (operator dice-operators)
        selector-fn (selector dice-selectors)
        roll-values (dice/roll
                      (:number-of-dice roll)
                      (:faces roll))
        result (operator-fn selector-fn literal roll-values)]
    (map
      (fn [{:keys [id] :as die}]
        (if-let [result-die
                 (->> result
                   (filter #(= (:id %) id))
                   (first))]
          (assoc result-die :discarded false)
          (assoc die :discarded true)))
      roll-values)))

(defn print-output
  "TODO:
  Takes a map of dice rolls and prints the result of each dice set
  operation and the operations"
  [dice-roll operations]
  (let [rolls (map (fn [[_ value :as roll]]
                     (assoc roll :outcomes (evaluate-roll value)))
                   dice-roll)]
    rolls))