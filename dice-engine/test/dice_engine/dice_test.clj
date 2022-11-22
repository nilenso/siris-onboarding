(ns dice-engine.dice-test
  (:require [clojure.test :refer [deftest is testing]]
            [dice-engine.dice :as dice]))

(deftest drop-test
  (testing "should mark the dice that match the exact value as discarded"
    (is (= (dice/discard [{:value           1
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
                         dice/match
                         3)
           [{:value           1
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
             :previous-values []}])))

  (testing "should return the dice unchanged if none of the dice value match."
    (is (= (dice/discard [{:value           2
                           :discarded       false
                           :faces           2
                           :previous-values []}]
                         dice/match
                         1)
           [{:value           2
             :discarded       false
             :faces           2
             :previous-values []}])))

  )

(deftest keep-test
  (testing "should mark the dice that do not match as discarded and update :numeric-value"
    (is (= (dice/pick [{:value           1
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
                      dice/match
                      4)
           [{:value           1
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
             :previous-values []}])))

  (testing "should return mark all dice as discarded and update :numeric-value to 0 if none match"
    (is (= (dice/pick [{:value           1
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
                      dice/match
                      2)
           [{:value           1
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
             :previous-values []}]))))

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
         [4 3])))

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

(deftest equal-test
  (is (= (dice/match [
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
                     4)
         [4])))

(deftest sum-test
  (is (= (dice/sum [
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
                     :previous-values []}])
         8))
  )

