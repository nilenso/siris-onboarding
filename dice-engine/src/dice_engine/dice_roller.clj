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
  [{:keys [faces number-of-dice selector operator] :as die-roll}]
  )

;2d6kh1
(def dice-roll-1 {:roll          {:number-of-dice 2
                                  :faces          6}
                  :set-operation {:operator :keep
                                  :selector :highest
                                  :literal  1}})



;2d3 3d4 1 +
;[2 3] [1 2 4] 1 +
(fn [operator & n]
  ; check for type
  ; dice -> dice/sum

  )

