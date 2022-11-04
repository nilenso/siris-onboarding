(ns mars-rover-clj.movement)

; X-Y Plane

(def directions {"N"     {:compass-angle 0
                          :move          (fn [rover]
                                           [(:x rover)
                                            (inc (:y rover))
                                            (:heading rover)])}
                     "E" {:compass-angle 90
                          :move          (fn [rover]
                                           [(inc (:x rover))
                                            (:y rover)
                                            (:heading rover)])}
                     "S" {:compass-angle 180
                          :move          (fn [rover] [(:x rover)
                                                      (dec (:y rover))
                                                      (:heading rover)])}
                     "W" {:compass-angle 270
                          :move          (fn [rover] [(dec (:x rover))
                                                      (:y rover)
                                                      (:heading rover)])}})

(defn compass-angle-to-direction
  "Returns the direction (N, E, S or W) from the compass angle"
  [compass-angle]
  (key (first (filter
                #(= (:compass-angle (val %)) (mod compass-angle 360))
                directions))))

(defn turn-rover-heading
  "Returns the new heading of the rover after turning it by direction-fn"
  [direction-fn rover turn-by]
  (update rover :heading
          #(->> (direction-fn
                  (:compass-angle (directions %))
                  turn-by)
                (compass-angle-to-direction))))

; Rover movements
(defn move
  "Returns the new position (x and y coordinates) and the heading of the rover"
  [rover]
  ((:move (directions (:heading rover))) rover))

(def rover-movements
  "A map of instruction symbols and the movement functions to move rovers"
  {"L" (fn [rover]
         (turn-rover-heading + rover 90))
   "R" (fn [rover]
         (turn-rover-heading - rover 90))
   "M" (fn [rover]
         (move rover))})

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

(defn move-rovers
  "Moves rovers based on the instruction string"
  [rover-input]
  (->> (rover-input :rovers)
       (map
         (fn [[[x y heading] [movement-instructions]]]
           (let [rover {:x       x
                        :y       y
                        :heading heading}]
             (map
               #(->>
                  (str %)
                  (rover-movements rover))
               movement-instructions))))))


;Testing stuff
(move {:x       1
       :y       2
       :heading "N"})

(((directions "N") :move) {:x           1
                               :y       2
                               :heading "N"})
(->> (partition 2 [[1 2 "N"]
                   ["LMLMLMLMM"]
                   [3 3 "E"]
                   ["MMRMMRMRRM"]])
     (map vec))
(move-rovers (init-rovers [5 5]
                          [1 2 "N"]
                          ["LMLMLMLMM"]
                          [3 3 "E"]
                          ["MMRMMRMRRM"]))
(compass-angle-to-direction 180)
(map str "MMRMMRMRRM")
(turn-rover-heading + {:x       1
                       :y       2
                       :heading "N"} 90)
(rover-movements "R")

