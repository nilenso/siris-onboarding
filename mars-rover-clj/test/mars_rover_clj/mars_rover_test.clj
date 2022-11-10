(ns mars-rover-clj.mars-rover-test
  (:require [clojure.test :refer [deftest is testing]]
            [mars-rover-clj.mars-rover :as mars-rover]))

(deftest move-test
  (testing "Move returns the new co-ordinates and heading"
    (is (=
          (mars-rover/move {:x 1 :y 2 :heading "N"})
          {:heading "N"
           :x       1
           :y       3}))
    (is (=
          (mars-rover/move {:x 1 :y 2 :heading "S"})
          {:heading "S"
           :x       1
           :y       1}))
    (is (=
          (mars-rover/move {:x 5 :y 2 :heading "E"})
          {:heading "E"
           :x       6
           :y       2}))))

(deftest check-bounds-test
  (is (= (mars-rover/check-bounds {:heading "N"
                                   :x       3
                                   :y       4}
                                  [5 5])
         true))
  (is (= (mars-rover/check-bounds {:heading "N"
                                   :x       6
                                   :y       4}
                                  [5 5])
         false))
  )

(deftest compass-angle-to-direction-test
  (testing "Tests for the correctness of compass-angle-to-direction"
    (is (=
          (mars-rover/compass-angle-to-direction 0)
          "N"))
    (is (=
          (mars-rover/compass-angle-to-direction 90)
          "E"))
    (is (=
          (mars-rover/compass-angle-to-direction 180)
          "S"))
    (is (=
          (mars-rover/compass-angle-to-direction 270)
          "W"))))

(deftest turn-rover-heading-test
  (testing "Tests for the correctness of turn-rover-heading"
    (is (=
          (mars-rover/turn-heading "N" + 90)
          "E"))
    (is (=
          (mars-rover/turn-heading "N" - 90)
          "W"))
    (is (=
          (mars-rover/turn-heading "N" + 180)
          "S"))
    (is (=
          (mars-rover/turn-heading "N" - 180)
          "S"))
    (is (=
          (mars-rover/turn-heading "S" - 90)
          "E"))))

(deftest move-rover-test
  (testing "Correctness of move-rover"
    (is (=
          (mars-rover/move-rover
            {:x 1, :y 3, :heading "N"}
            ["L" "M" "L" "M" "L" "M" "L" "M" "M"]
            [5 5])
          {:x       1
           :y       4
           :heading "N"}))
    (is (=
          (mars-rover/move-rover
            {:x 3, :y 3, :heading "E"}
            ["M" "M" "R" "M" "M" "R" "M" "R" "R" "M"]
            [5 5])
          {:x       5
           :y       1
           :heading "E"}))
    (is (=
          (mars-rover/move-rover
            {:x 3, :y 3, :heading "E"}
            ["M" "M" "M" "L" "M" "R" "M" "R" "R" "M"]
            [5 5])
          mars-rover/out-of-bounds-error)))
  )

(deftest move-rovers-sequentially-test
  (testing "move-rovers"
    (is (= (->> (mars-rover/init-rovers [5 5]
                                        [1 2 "N"]
                                        "LMLMLMLMM"
                                        [3 3 "E"]
                                        "MMRMMRMRRM")
             (mars-rover/move-rovers-sequentially))
           '({:x 1, :y 3, :heading "N"} {:x 5, :y 1, :heading "E"})))))

(deftest init-rovers-test
  (testing "init-rovers"
    (is (= (mars-rover/init-rovers [5 5]
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
