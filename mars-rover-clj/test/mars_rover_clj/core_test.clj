(ns mars-rover-clj.core-test
  (:require [clojure.test :refer :all]
            [mars-rover-clj.core :refer :all]))

(deftest move-test
  (testing "Move returns the new co-ordinates and heading"
    (is (= (move 1 2 "N") '(1 3 "N")))))