(ns mars-rover-clj.movement
  (:require [clojure.set :as set]))

; X-Y Plane
(def directions {"N" 0
                 "E" 90
                 "S" 180
                 "W" 270})
(def directions-map {"N" {:angle 0
                          :move  (fn [x y heading] [x (inc y) heading])}
                     "E" {:angle 90
                          :move  (fn [x y heading] [(inc x) y heading])}
                     "S" {:angle 180
                          :move  (fn [x y heading] [x (dec y) heading])}
                     "W" {:angle 270
                          :move  (fn [x y heading] [(dec x) y heading])}})

(defn turn-left
  "Returns an angle of the compass point after turning left by the given angle."
  [heading turn-by]
  (-
    ((directions-map heading) :angle)
    turn-by))

(defn turn-right
  "Returns an angle of the compass point after turning right by the given angle."
  [heading turn-by]
  (+
    ((directions-map heading) :angle)
    turn-by))

(def turn-directions {"L" (fn [heading turn-by]
                            (turn-left heading turn-by))
                      "R" (fn [heading turn-by]
                            (turn-right heading turn-by))})

(defn angle-to-direction
  "Returns a compass point given the angle. North 0, East 90, South 180, West 270."
  [angle]
  ((set/map-invert directions)
   (mod angle 360)))

; Rover movements
(defn move
  "Returns the new position (x and y coordinates) and the heading of the rover"
  [x y heading]
  ((:move (directions-map heading)) x y heading))

(defn turn-90
  "Turns the heading of the rover in the given direction"
  [heading turn-direction]
  (angle-to-direction
    ((turn-directions turn-direction)
     heading
     90)))

; I/O functions
(defn init-rovers
  "Takes an input for the Mars Rover problem
  as a vector of the plateau's upper right coordinates and an infinite sequence of rover's initial position
  and a string of a series of instructions.
  Returns a map of plateau bounds and rovers."
  [plateau-upper-right-coordinates & rover-details]
  {:plateau-bounds plateau-upper-right-coordinates
   :rovers         (->> (partition 2 rover-details)
                        (map vec))})


; Testing stuff
(move 1 2 "N")
(turn-90 "N" "R")
(angle-to-direction 0)
(((directions-map "N") :move) 1 2 "N")
(->> (partition 2 [[1 2 "N"]
                   ["LMLMLMLMM"]
                   [3 3 "E"]
                   ["MMRMMRMRRM"]])
     (map vec))
