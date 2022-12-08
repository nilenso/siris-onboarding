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
  [{:keys [outcomes]}]
  (->> outcomes
    (filter #(not (:discarded %)))
    (dice/sum)))

(defn evaluate-roll
  "Evaluates a dice roll, applies set operations and returns a dice collection"
  [{:keys                               [roll]
    {:keys [operator selector literal]} :set-operation :as dice}]
  (let [operator-fn (operator dice-operators)
        selector-fn (selector dice-selectors)
        roll-values (dice/roll
                      (:number-of-dice roll)
                      (:faces roll))
        result (operator-fn selector-fn literal roll-values)]
    (->> (map
           (fn [{:keys [id] :as die}]
             (if-let [result-die
                      (->> result
                        (filter #(= (:id %) id))
                        (first))]
               (assoc result-die :discarded false)
               (assoc die :discarded true)))
           roll-values)
      (assoc dice :outcomes))))

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
  "Parses a die roll in the output format.
  For instance, 2d6kh1: (~2~, 5) => 5"
  [{:keys [expression outcomes] :as roll}]
  (str expression ": "
       "("
       (->> (map #(parse-roll-result %) outcomes)
         (string/join ", "))
       ")"
       " => "
       (roll-value roll)))