(ns dice-engine.dice-roller-test
  (:require [clojure.test :refer [deftest is]]
            [dice-engine.dice :as dice]
            [dice-engine.dice-roller :as dice-roller]))

(def roll-1 {:expression    "2d6kh1"
             :roll          {:number-of-dice 2
                             :faces          6}
             :set-operation {:operator :drop
                             :selector :greater-than
                             :literal  2}})

(def roll-2 {:expression    "3d20rr1"
             :roll          {:number-of-dice 3
                             :faces          20}
             :set-operation {:operator :reroll
                             :selector :greater-than
                             :literal  7}})

(def roll-outcome-1 {:roll          {:number-of-dice 2
                                     :faces          6}
                     :set-operation {:operator :keep
                                     :selector :highest
                                     :literal  1}
                     :outcomes      [{:id              1
                                      :value           3
                                      :discarded       true
                                      :previous-values []}
                                     {:id              2
                                      :value           4
                                      :discarded       false
                                      :previous-values []}]})

(def roll-outcome-2 {:expression    "3d20rr7"
                     :roll          {:number-of-dice 3
                                     :faces          20}
                     :set-operation {:operator :reroll
                                     :selector :greater-than
                                     :literal  7}
                     :outcomes      [{:id              6
                                      :value           7
                                      :faces           20
                                      :previous-values []
                                      :discarded       false}
                                     {:id              7
                                      :value           5
                                      :faces           20
                                      :previous-values [10 20 8 20 20 19 20 9 18 19]
                                      :discarded       false}
                                     {:id              8
                                      :value           6
                                      :faces           20
                                      :previous-values [14 8 12 14 13]
                                      :discarded       false}]})

(deftest evaluate-roll-test
  (is (= (dice-roller/evaluate-roll roll-1)
         (:outcomes roll-outcome-1))))

(def dice-roll {:r1 {:expression    "2d6kh1"
                     :roll          {:number-of-dice 2
                                     :faces          6}
                     :set-operation {:operator :drop
                                     :selector :greater-than
                                     :literal  2}}
                :r2 {:expression    "3d20rr1"
                     :roll          {:number-of-dice 3
                                     :faces          20}
                     :set-operation {:operator :reroll
                                     :selector :greater-than
                                     :literal  7}}})

(def dice-roll-1 {:roll          {:number-of-dice 2
                                  :faces          6}
                  :set-operation {:operator :keep
                                  :selector :highest
                                  :literal  1}})

(dice-roller/roll-value (dice-roller/evaluate-roll (:r2 dice-roll)))
;(dice-roller/evaluate-roll '({:id 4, :value 5, :faces 6, :previous-values [], :discarded true}
;              {:id 5, :value 1, :faces 6, :previous-values [], :discarded false}))


(dice/drop dice/match 4 [{:id 6, :value 7, :faces 20, :previous-values []}
                         {:id 7, :value 16, :faces 20, :previous-values []}
                         {:id 8, :value 4, :faces 20, :previous-values []}])

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



