(ns dice-engine.dice-roller
  (:require [dice-engine.dice :as dice]))

(comment {:faces
          :number-of-dice
          :selector
          :operator})

(def dice-operators {:keep   dice/keep
                     :drop   dice/drop
                     :reroll dice/reroll})

(defn process-dice
  "docstring"
  [{:keys [faces number-of-dice selector operator] :as die-roll}])

;2d6kh1
(def dice-roll-1 {:roll          {:number-of-dice 2
                                  :faces          6}
                  :set-operation {:operator :keep
                                  :selector :highest
                                  :literal  1}})

(def dice-roll-1-intermediate {:roll          {:number-of-dice 2
                                               :faces          6}
                               :set-operation {:operator :keep
                                               :selector :highest
                                               :literal  1}}
                              :outcomes [{:value 3 :discared: false :rerolled true}])

(defn intermediate->value
  [intermediate-roll]
  10)

(defn intermediate->output-format
  [intermediate-roll]
  "(val1, ~discarded_val2~, val3 (~previousrollval3~, ~previousrollval3'~))")

(def dice-roll-2 {:roll          {:number-of-dice 3
                                  :faces          4}
                  :set-operation {:operator :drop
                                  :selector :lowest
                                  :literal  2}})

(def numerals-1 2)

; + - * /

(def die-1 {:id              3
            :value           die-value
            :faces           faces
            :previous-values []})

; (val1, ~discarded_val2~, val3 (~previousrollval3~, ~previousrollval3'~))

(defn parse []
  {:outcome [{:id              3
              :value           59
              :faces           60
              :previous-values []}
             {:id              4
              :value           49
              :faces           50
              :previous-values []}]
   :final-state [{:id              3
                  :value           59
                  :faces           60
                  :previous-values []}]
   :set-operation {:operator :drop
                   :selector :lowest
                   :literal  1}})

(* (parse dice-roll-1)
   (+ (parse dice-roll-2)
      (parse numerals-1)))

;2d3 3d4 1 +
;[2 3] [1 2 4] 1 +
(fn [operator & n]
  ; check for type
  ; dice -> dice/sum
  )
