(ns dice-engine.dice-test
  (:require [clojure.test :refer [deftest is testing]]
            [dice-engine.dice :as dice]))

(deftest drop-test
  (testing "should mark the dice that match the exact value as discarded and update :numeric-value"
    (is (= (dice/drop {:numeric-value 8
                       :dice          [
                                       {:value           1
                                        :discarded       false
                                        :faces           4
                                        :previous-values []}
                                       {:value           3
                                        :discarded       false
                                        :faces           4
                                        :previous-values []}
                                       {:value           4
                                        :discarded       false
                                        :faces           4
                                        :previous-values []}]}
                      3)
           '{:numeric-value 5
             :dice          [[
                              {:value           1
                               :discarded       false
                               :faces           4
                               :previous-values []}
                              {:value           3
                               :discarded       true
                               :faces           4
                               :previous-values []}
                              {:value           4
                               :discarded       false
                               :faces           4
                               :previous-values []}]]})))

  (testing "should return the dice-roll map unchanged if none of the dice value match."
    (is (= (dice/drop {:numeric-value 2
                       :dice          [{:value           2
                                        :discarded       false
                                        :faces           2
                                        :previous-values []}]}
                      1)
           '{:numeric-value 2
             :dice          [{:value           2
                              :discarded       false
                              :faces           2
                              :previous-values []}]})))

  )

(deftest keep-test
  (testing "should mark the dice that do not match as discarded and update :numeric-value"
    (is (= (dice/keep {:numeric-value 8
                       :dice          [
                                       {:value           1
                                        :discarded       false
                                        :faces           4
                                        :previous-values []}
                                       {:value           3
                                        :discarded       false
                                        :faces           4
                                        :previous-values []}
                                       {:value           4
                                        :discarded       false
                                        :faces           4
                                        :previous-values []}]}
                      4)
           '{:numeric-value 4
             :dice          [
                             {:value           1
                              :discarded       true
                              :faces           4
                              :previous-values []}
                             {:value           3
                              :discarded       true
                              :faces           4
                              :previous-values []}
                             {:value           4
                              :discarded       false
                              :faces           4
                              :previous-values []}]})))
  (testing "should return mark all dice as discarded and update :numeric-value to 0 if none match"
    (is (= dice/keep {:numeric-value 8
                      :dice          [
                                      {:value           1
                                       :discarded       false
                                       :faces           4
                                       :previous-values []}
                                      {:value           3
                                       :discarded       false
                                       :faces           4
                                       :previous-values []}
                                      {:value           4
                                       :discarded       false
                                       :faces           4
                                       :previous-values []}]}
           2)
        '{:numeric-value 0
          :dice          [
                          {:value           1
                           :discarded       true
                           :faces           4
                           :previous-values []}
                          {:value           3
                           :discarded       true
                           :faces           4
                           :previous-values []}
                          {:value           4
                           :discarded       true
                           :faces           4
                           :previous-values []}]})))

(deftest highest-test
  (is (= (dice/highest [
                        {:value           1
                         :discarded       false
                         :faces           4
                         :previous-values []}
                        {:value           3
                         :discarded       false
                         :faces           4
                         :previous-values []}
                        {:value           4
                         :discarded       false
                         :faces           4
                         :previous-values []}]
                       2)
         [3 4])))

(deftest lowest-test
  (is (= (dice/lowest [
                       {:value           1
                        :discarded       false
                        :faces           4
                        :previous-values []}
                       {:value           3
                        :discarded       false
                        :faces           4
                        :previous-values []}
                       {:value           4
                        :discarded       false
                        :faces           4
                        :previous-values []}]
                      2)
         [1 3])))

(deftest greater-than-test
  (is (= (dice/greater-than [
                             {:value           1
                              :discarded       false
                              :faces           4
                              :previous-values []}
                             {:value           3
                              :discarded       false
                              :faces           4
                              :previous-values []}
                             {:value           4
                              :discarded       false
                              :faces           4
                              :previous-values []}]
                            3)
         [4]))
  )

(deftest less-than-test
  (is (= (dice/less-than [
                          {:value           1
                           :discarded       false
                           :faces           4
                           :previous-values []}
                          {:value           3
                           :discarded       false
                           :faces           4
                           :previous-values []}
                          {:value           4
                           :discarded       false
                           :faces           4
                           :previous-values []}]
                         2)
         [1]))
  )

