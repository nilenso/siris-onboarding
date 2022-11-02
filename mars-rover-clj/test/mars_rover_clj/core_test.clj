(ns mars-rover-clj.core-test
  (:require [clojure.set :refer :all]
            [clojure.test :refer [deftest is testing]]
            [mars-rover-clj.core :refer [move turn-90]]))

(deftest move-test
  (testing "Move returns the new co-ordinates and heading"
    (is (= (move 1 2 "N") [1 3 "N"]))))

(deftest turn-90-test
  (testing "Tests for the correctness of turn-90"
    (is (= (turn-90 "N" "R") "E"))))

;(deftest turn-left
;  (testing "Tests for the correctness of turn-left"
;    (is (= (turn-left )))))