(ns mars-rover-clj.movement-test
  (:require [clojure.test :refer [deftest is testing]]
            [mars-rover-clj.movement :as movement]))

(deftest move-test
  (testing "Move returns the new co-ordinates and heading"
    (is (=
          (movement/move {:x 1 :y 2 :heading "N"})
          {:heading "N"
           :x       1
           :y       3}))
    (is (=
          (movement/move {:x 1 :y 2 :heading "S"})
          {:heading "S"
           :x       1
           :y       1}))
    (is (=
          (movement/move {:x 5 :y 2 :heading "E"})
          {:heading "E"
           :x       6
           :y       2}))))

(deftest check-bounds-test
  (is (= (movement/check-bounds {:heading "N"
                                 :x       3
                                 :y       4}
                                [5 5])
         true))
  (is (= (movement/check-bounds {:heading "N"
                                 :x       6
                                 :y       4}
                                [5 5])
         false))
  )

(deftest compass-angle-to-direction-test
  (testing "Tests for the correctness of compass-angle-to-direction"
    (is (=
          (movement/compass-angle-to-direction 0)
          "N"))
    (is (=
          (movement/compass-angle-to-direction 90)
          "E"))
    (is (=
          (movement/compass-angle-to-direction 180)
          "S"))
    (is (=
          (movement/compass-angle-to-direction 270)
          "W"))))

(deftest turn-rover-heading-test
  (testing "Tests for the correctness of turn-rover-heading"
    (is (=
          (movement/turn-rover-heading "N" + 90)
          "E"))
    (is (=
          (movement/turn-rover-heading "N" - 90)
          "W"))
    (is (=
          (movement/turn-rover-heading "N" + 180)
          "S"))
    (is (=
          (movement/turn-rover-heading "N" - 180)
          "S"))
    (is (=
          (movement/turn-rover-heading "S" - 90)
          "E"))))

(deftest move-rover-test
  (testing "Correctness of move-rover"
    (is (=
          (movement/move-rover
            {:x 1, :y 3, :heading "N"}
            ["L" "M" "L" "M" "L" "M" "L" "M" "M"]
            [5 5])
          {:x       1
           :y       4
           :heading "N"}))
    (is (=
          (movement/move-rover
            {:x 3, :y 3, :heading "E"}
            ["M" "M" "R" "M" "M" "R" "M" "R" "R" "M"]
            [5 5])
          {:x       5
           :y       1
           :heading "E"}))
    (is (=
          (movement/move-rover
            {:x 3, :y 3, :heading "E"}
            ["M" "M" "M" "L" "M" "R" "M" "R" "R" "M"]
            [5 5])
          movement/rover-out-of-bounds-error)))
  )

(deftest move-rovers-sequentially-test
  (testing "move-rovers"
    (is (= (->> (movement/init-rovers [5 5]
                                      [1 2 "N"]
                                      "LMLMLMLMM"
                                      [3 3 "E"]
                                      "MMRMMRMRRM")
             (movement/move-rovers-sequentially))
           '({:x 1, :y 3, :heading "N"} {:x 5, :y 1, :heading "E"})))))

(deftest init-rovers-test
  (testing "init-rovers"
    (is (= (movement/init-rovers [5 5]
                                 [1 2 "N"]
                                 "LMLMLMLMM"
                                 [3 3 "E"]
                                 "MMRMMRMRRM")
           '{:plateau-bounds [5 5],
             :rovers         (
                              {:rover        {:x 1, :y 2, :heading "N"}
                               :instructions ["L" "M" "L" "M" "L" "M" "L" "M" "M"]}
                              {
                               :rover        {:x 3, :y 3, :heading "E"}
                               :instructions ["M" "M" "R" "M" "M" "R" "M" "R" "R" "M"]})}))))