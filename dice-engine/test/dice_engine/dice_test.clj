(ns dice-engine.dice-test
  (:require [clojure.test :refer [deftest is testing]]
            [dice-engine.dice :as dice]
            [dice-engine.factories :as factories]))

(def dice-1 [{:id              1
              :value           3
              :faces           6
              :previous-values []}
             {:id              2
              :value           4
              :faces           6
              :previous-values []}])

(def dice-2 [{:id              6
              :value           7
              :faces           20
              :previous-values []}
             {:id              7
              :value           5
              :faces           20
              :previous-values []}
             {:id              8
              :value           6
              :faces           20
              :previous-values []}])

(deftest drop-test
  (testing "should mark the dice that match the exact value as discarded"
    (let [die-4 (factories/create-die 4 6)
          die-3 (factories/create-die 3 6)]
      (is (= (dice/drop dice/match 3 [die-4 die-3])
             [die-4]))))

  (testing "should return the dice unchanged if none of the dice value match."
    (let [die-4 (factories/create-die 4 6)
          die-3 (factories/create-die 3 6)]
      (is (=
            (dice/drop dice/greater-than 5 [die-4 die-3])
            [die-4 die-3])))))

(deftest keep-test
  (testing "should return dice that match the value"
    (let [die-4 (factories/create-die 4 6)
          die-3 (factories/create-die 3 6)
          all-dice [die-4 die-3]
          expected-dice [die-4]]
      (is (=
            (dice/keep dice/match 4 all-dice)
            expected-dice))))

  (testing "should return empty vector if none of the dice values match"
    (let [die-4 (factories/create-die 4 6)
          die-3 (factories/create-die 3 6)
          all-dice [die-4 die-3]
          expected-dice []]
      (is (=
            (dice/keep dice/match 2 all-dice)
            expected-dice)))))

(deftest reroll-matched-test
  (testing "should reroll the matched die"
    (with-redefs-fn {#'dice/rand-int-natural (fn [_] 2)}
      #(is (=
             (dice/reroll-matched dice/match 3 dice-1)
             [{:id              1
               :value           2
               :faces           6
               :previous-values [3]}
              {:id              2
               :value           4
               :faces           6
               :previous-values []}]))))
  (testing "should reroll until no die qualifies selector"
    (let [rand-ints [5 9 1]]
      (with-local-vars [counter -1]
        (with-redefs-fn {#'dice/rand-int-natural (fn [_]
                                                   (var-set counter (inc @counter))
                                                   (nth rand-ints (var-get counter)))}
          #(is (=
                 (dice/reroll-matched dice/greater-than 3 dice-1)
                 [{:id              2
                   :value           1
                   :faces           6
                   :previous-values '(9 5 4)}
                  {:id              1
                   :value           3
                   :faces           6
                   :previous-values []}])))))))

(deftest highest-test
  (is (=
        (dice/highest 2 dice-2)
        ['({:id              6
            :value           7
            :faces           20
            :previous-values []}
           {:id              8
            :value           6
            :faces           20
            :previous-values []})
         '({:id              7
            :value           5
            :faces           20
            :previous-values []})])))

(deftest lowest-test
  (is (=
        (dice/lowest 2 dice-2)
        ['({:id              7
            :value           5
            :faces           20
            :previous-values []}
           {:id              8
            :value           6
            :faces           20
            :previous-values []})
         '({:id              6
            :value           7
            :faces           20
            :previous-values []})])))

(deftest greater-than-test
  (is (= (dice/greater-than 5 dice-2)
         ['({:id              6
             :value           7
             :faces           20
             :previous-values []}
            {:id              8
             :value           6
             :faces           20
             :previous-values []})
          '({:id              7
             :value           5
             :faces           20
             :previous-values []})])))

(deftest less-than-test
  (is (= (dice/lesser-than 7 dice-2)
         ['({:id              7
             :value           5
             :faces           20
             :previous-values []}
            {:id              8
             :value           6
             :faces           20
             :previous-values []})
          '({:id              6
             :value           7
             :faces           20
             :previous-values []})])))

(deftest match-test
  (is (=
        (dice/match 4 dice-1)
        ['({:id              2
            :value           4
            :faces           6
            :previous-values []})
         '({:id              1
            :value           3
            :faces           6
            :previous-values []})])))

(deftest sum-test
  (is (=
        (dice/sum dice-1)
        7)))

(deftest partition-dice-by-value-test
  (is (= (dice/partition-dice-by-value (partial #(> % 5)) dice-2)
         ['({:id              6
             :value           7
             :faces           20
             :previous-values []}
            {:id              8
             :value           6
             :faces           20
             :previous-values []})
          '({:id              7
             :value           5
             :faces           20
             :previous-values []})])))

(deftest roll-test
  (let [rand-ints [4 6 8]]
    (with-local-vars [counter -1
                      id-counter 0]
      (with-redefs-fn {#'dice/rand-int-natural (fn [_]
                                                 (var-set counter (inc @counter))
                                                 (nth rand-ints (var-get counter)))
                       #'dice/get-id           (fn []
                                                 (var-set id-counter (inc @id-counter))
                                                 @id-counter)}
        #(is (= (dice/roll 3 8)
                [{:id              1
                  :value           4
                  :faces           8
                  :previous-values []}
                 {:id              2
                  :value           6
                  :faces           8
                  :previous-values []}
                 {:id              3
                  :value           8
                  :faces           8
                  :previous-values []}]))))))

(deftest create-dice-test
  (with-redefs-fn {#'dice/get-id (fn [] 1)}
    #(is (= (dice/create-die 4 5)
            {:id              1
             :value           4
             :faces           5
             :previous-values []}))))

(deftest rand-int-natural-test
  (with-redefs-fn {#'rand-int (fn [_] 0)}
    #(is (= (dice/rand-int-natural 5) 1))))
