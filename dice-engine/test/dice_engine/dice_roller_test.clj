(ns dice-engine.dice-roller-test
  (:require [clojure.test :refer [deftest is testing]]
            [dice-engine.dice :as dice]
            [dice-engine.dice-roller :as dice-roller]))

(def roll-1 {:expression    "2d6kh1"
             :roll          {:number-of-dice 2
                             :faces          6}
             :set-operation {:operator :keep
                             :selector :highest
                             :literal  1}})

(def roll-2 {:expression    "3d20rr>15"
             :roll          {:number-of-dice 3
                             :faces          20}
             :set-operation {:operator :reroll
                             :selector :greater-than
                             :literal  15}})

(def roll-outcome-1 {:expression    "2d6kh1"
                     :roll          {:number-of-dice 2
                                     :faces          6}
                     :set-operation {:operator :keep
                                     :selector :highest
                                     :literal  1}
                     :outcomes      [{:id              1
                                      :value           3
                                      :faces           6
                                      :discarded       true
                                      :previous-values []}
                                     {:id              2
                                      :value           4
                                      :faces           6
                                      :discarded       false
                                      :previous-values []}]})

(def roll-outcome-2 {:expression    "3d20rr>15"
                     :roll          {:number-of-dice 3
                                     :faces          20}
                     :set-operation {:operator :reroll
                                     :selector :greater-than
                                     :literal  15}
                     :outcomes      [{:id              1
                                      :value           4
                                      :faces           20
                                      :previous-values []
                                      :discarded       false}
                                     {:id              2
                                      :value           6
                                      :faces           20
                                      :previous-values []
                                      :discarded       false}
                                     {:id              3
                                      :value           11
                                      :faces           20
                                      :previous-values [20]
                                      :discarded       false}]})

(def roll-output "3d20rr>15: (4, 6, 11 (~20~)) => 21")

(deftest roll-value-test
  (is (=
        (dice-roller/roll-value roll-outcome-1)
        4))
  (is (=
        (dice-roller/roll-value roll-outcome-2)
        21)))

(deftest parse-roll-result-test
  (testing "value"
    (is (=
          (dice-roller/parse-roll-result {:id              2
                                          :value           4
                                          :discarded       false
                                          :previous-values []})
          "4")))
  (testing "previous values"
    (is (=
          (dice-roller/parse-roll-result {:id              8
                                          :value           6
                                          :faces           20
                                          :previous-values [14 8 12 14 13]
                                          :discarded       false})
          "6 (~14~, ~8~, ~12~, ~14~, ~13~)")))
  (testing "discarded value"
    (is (=
          (dice-roller/parse-roll-result {:id              1
                                          :value           3
                                          :discarded       true
                                          :previous-values []})
          "~3~"))))

(deftest parse-output-test
  (is (=
        (dice-roller/parse-output roll-outcome-1)
        "2d6kh1: (~3~, 4) => 4")))


(deftest evaluate-roll-test
  (let [rand-ints [3 4]]
    (with-local-vars [counter -1
                      id-counter 0]
      (with-redefs-fn {#'dice/rand-int-natural (fn [_]
                                                 (var-set counter (inc @counter))
                                                 (nth rand-ints (var-get counter)))
                       #'dice/get-id           (fn []
                                                 (var-set id-counter (inc @id-counter))
                                                 @id-counter)}
        #(is (=
               (dice-roller/evaluate-roll roll-1)
               roll-outcome-1))))))

(deftest dice-engine-test
  (let [rand-ints [4 6 20 11]]
    (with-local-vars [counter -1
                      id-counter 0]
      (with-redefs-fn {#'dice/rand-int-natural (fn [_]
                                                 (var-set counter (inc @counter))
                                                 (nth rand-ints (var-get counter)))
                       #'dice/get-id           (fn []
                                                 (var-set id-counter (inc @id-counter))
                                                 @id-counter)}
        #(is (=
               (dice-roller/parse-output (dice-roller/evaluate-roll roll-2))
               roll-output))))))

(deftest numeric-operation-test
  (is (= (+
           (dice-roller/roll-value roll-outcome-2)
           (- (dice-roller/roll-value roll-outcome-1) (* 3 4)))
         13)))


