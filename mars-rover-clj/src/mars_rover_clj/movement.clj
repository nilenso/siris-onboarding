(ns mars-rover-clj.movement)

; X-Y Plane
(def directions
  {"N" {:compass-angle 0
        :move-f        (fn [rover] (update rover :y inc))}
   "E" {:compass-angle 90
        :move-f        (fn [rover] (update rover :x inc))}
   "S" {:compass-angle 180
        :move-f        (fn [rover] (update rover :y dec))}
   "W" {:compass-angle 270
        :move-f        (fn [rover] (update rover :x dec))}})

(def
  ^{:doc "Returns the direction (N, E, S or W) from the compass angle"}
  compass-angle-to-direction
  (->> directions
    (map (fn [[direction {:keys [compass-angle]}]] [compass-angle direction]))
    (into {})))

(defn turn-rover-heading
  "Returns the new heading of the rover after turning it by direction-fn"
  [heading direction-fn turn-by]
  (-> (:compass-angle (directions heading))
      (#(direction-fn % turn-by))
      (mod 360)
      (compass-angle-to-direction)))

; Rover movements
(defn check-bounds
  "Returns true if the rover's position is within the plateau bounds, false otherwise"
  [plateau-bounds rover]
  (and
    (<= 0 (:x rover) (first plateau-bounds))
    (<= 0 (:y rover) (second plateau-bounds))))

(defn move
  "Returns a rover with updated position (x and y coordinates) and the heading.
  If rover position results beyond the plateau bounds, it returns the rover unchanged."
  [plateau-bounds rover]
  (let [moved-rover ((:move-f (directions (:heading rover))) rover)]
    (if (check-bounds plateau-bounds moved-rover)
      moved-rover
      rover)))

(def rover-movements
  "A map of instruction symbols and the movement functions to move rovers"
  {"L" (fn [_ rover]
         (update rover :heading #(turn-rover-heading % - 90)))
   "R" (fn [_ rover]
         (update rover :heading #(turn-rover-heading % + 90)))
   "M" move})

; I/O functions
(defn init-rovers
  "Takes an input for the Mars Rover problem
  as a vector of the plateau's upper right coordinates and an infinite sequence of rover's initial position
  and a string of a series of instructions.
  Returns a map of plateau bounds and rovers."
  [plateau-upper-right-coordinates & rover-details]
  {:plateau-bounds plateau-upper-right-coordinates
   :rovers         (->> (partition 2 rover-details)
                     (map (fn [rover-instructions]
                            (let [[x y heading] (first rover-instructions)]
                              [{:x       x
                                :y       y
                                :heading heading}
                               (second rover-instructions)]))))})


(defn move-rover
  "Moves a rover as specified by the instructions string"
  [plateau-bounds [rover instructions]]
  (reduce (fn [acc instruction]
            ((->>
               (str instruction)
               (rover-movements)) plateau-bounds acc))
          rover instructions))

(defn move-rovers
  "Moves rovers based on the instruction string"
  [rover-input]
  (->> (rover-input :rovers)
    (map move-rover (:plateau-bounds rover-input))))

;Testing stuff

(move [5 5] {:x       1
             :y       2
             :heading "N"})

(((directions "N") :move-f) {:x       1
                             :y       2
                             :heading "N"})
(->> (partition 2 [[1 2 "N"]
                   "LMLMLMLMM"
                   [3 3 "E"]
                   "MMRMMRMRRM"])
  (map vec))
(move-rovers (init-rovers [5 5]
                          [1 2 "N"]
                          "LMLMLMLMM"
                          [3 3 "E"]
                          "MMRMMRMRRM"))

(move-rover [5 5] [{:x       1
                    :y       2
                    :heading "N"} "LMLMLMLMM"])
(move-rover [5 5] [{:x       3
                    :y       3
                    :heading "E"} "MMRMMRMRRM"])

(compass-angle-to-direction 180)
(map str "MMRMMRMRRM")
(turn-rover-heading "N" - 90)


(map
  (fn [[[x y heading] [movement-instructions]]]
    (let [rover {:x       x
                 :y       y
                 :heading heading}]
      (map
        #(->>
           (str %)
           ((rover-movements) rover))
        movement-instructions))))
(move [5 5] {:x 5, :y 3, :heading "E"})

((->>
   (str "M")
   (rover-movements))
 [5 5] {:x 4, :y 3, :heading "E"})

(compass-angle-to-direction 90)