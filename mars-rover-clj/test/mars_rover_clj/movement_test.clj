(ns mars-rover-clj.movement-test
  (:require [clojure.test :refer [deftest is testing]]
            [mars-rover-clj.movement :as movement]))

(deftest move-test
  (testing "Move returns the new co-ordinates and heading"
    (is (= (movement/move {:x 1 :y 2 :heading "N"}) [1 3 "N"]))
    (is (= (movement/move {:x 1 :y 2 :heading "N"}) [1 3 "N"]))))

(deftest compass-angle-to-direction-test
  (testing "Tests for the correctness of compass-angle-to-direction"
    (is (= (movement/compass-angle-to-direction -90) "W"))
    (is (= (movement/compass-angle-to-direction 90) "E"))
    (is (= (movement/compass-angle-to-direction -180) "S"))
    (is (= (movement/compass-angle-to-direction 360) "N"))
    (is (= (movement/compass-angle-to-direction 270) "W"))
    (is (= (movement/compass-angle-to-direction 270) "W"))))

(deftest turn-rover-heading-test
  (testing "Tests for the correctness of turn-rover-heading"
    (is (= (movement/turn-rover-heading + {:x 1 :y 2 :heading "N"} 90) {:x 1 :y 2 :heading "E"}))
    (is (= (movement/turn-rover-heading - {:x 1 :y 2 :heading "N"} 90) {:x 1 :y 2 :heading "W"}))
    (is (= (movement/turn-rover-heading + {:x 1 :y 2 :heading "N"} 180) {:x 1 :y 2 :heading "S"}))
    (is (= (movement/turn-rover-heading - {:x 1 :y 2 :heading "N"} 180) {:x 1 :y 2 :heading "S"}))
    (is (= (movement/turn-rover-heading - {:x 1 :y 2 :heading "S"} 90) {:x 1 :y 2 :heading "E"}))))

(deftest move-rovers-test
  (testing "move-rovers"
    (is (=
          (movement/move-rovers (movement/init-rovers [5 5] [1 2 "N"]
                                              ["LMLMLMLMM"]
                                              [3 3 "E"]
                                              ["MMRMMRMRRM"]))
           ))))


