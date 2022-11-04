(ns mars-rover-clj.movement-test
  (:require [clojure.test :refer [deftest is testing]]
            [mars-rover-clj.movement :as movement]))

(deftest move-test
  (testing "Move returns the new co-ordinates and heading"
    (is (= (movement/move 1 2 "N") [1 3 "N"]))
    (is (= (movement/move 1 2 "N") [1 3 "N"]))))

(deftest turn-90-test
  (testing "Tests for the correctness of turn-90"
    (is (= (movement/turn-90 "N" "R") "E"))))

(deftest turn-left-test
  (testing "Tests for the correctness of turn-left"
    (is (= (movement/turn-left "N" 90) -90))
    (is (= (movement/turn-left "S" 90) 90))
    (is (= (movement/turn-left "E" 90) 0))
    (is (= (movement/turn-left "W" 90) 180))))


(deftest turn-right-test
  (testing "Tests for the correctness of turn-left"
    (is (= (movement/turn-right "N" 90) 90))
    (is (= (movement/turn-right "S" 180) 360))
    (is (= (movement/turn-right "E" 90) 180))
    (is (= (movement/turn-right "W" 90) 360))))

(deftest angle-to-direction-test
  (testing "Tests the direction returned for an angle which is a multiple of 90."
    (is (= (movement/angle-to-direction -90) "W"))
    (is (= (movement/angle-to-direction 90) "E"))
    (is (= (movement/angle-to-direction -180) "S"))
    (is (= (movement/angle-to-direction 360) "N"))
    (is (= (movement/angle-to-direction 270) "W"))
    (is (= (movement/angle-to-direction 270) "W"))))