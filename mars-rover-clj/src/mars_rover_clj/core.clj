(ns mars-rover-clj.core)
(require '(clojure.set))

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
  ((clojure.set/map-invert directions)
   (mod angle 360)))

(defn move
  "Returns the new position (x and y coordinates) and the heading of the rover"
  [x y heading]
  (cond (= heading "N")
        [x (inc y) heading]
        (= heading "S")
        [x (dec y) heading]
        (= heading "E")
        [(inc x) y heading]
        (= heading "W")
        [(dec x) y heading]))

(defn turn-90
  "Turns the heading of the rover in the given direction"
  [heading turn-direction]
  (angle-to-direction
    ((turn-directions turn-direction)
     heading
     90)))

(turn-90 "N" "R")
(angle-to-direction 0)
