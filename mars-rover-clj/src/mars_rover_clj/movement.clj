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

(def rover-out-of-bounds-error
  "Error string when a bad instruction is passed to rover making it go off-grid."
  "Bad instructions. Rover overshoots plateau bounds.")

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
  [rover plateau-bounds]
  (and
    (<= 0 (:x rover) (first plateau-bounds))
    (<= 0 (:y rover) (second plateau-bounds))))

(defn move
  "Returns a rover with updated position (x and y coordinates) and the heading.
  If rover position results beyond the plateau bounds, it returns the rover unchanged."
  [rover]
  ((:move-f (directions (:heading rover))) rover))

(def rover-movements
  "A map of instruction symbols and the movement functions to move rovers"
  {"L" (fn [rover]
         (update rover :heading #(turn-rover-heading % - 90)))
   "R" (fn [rover]
         (update rover :heading #(turn-rover-heading % + 90)))
   "M" move})

(defn move-rover
  "Moves a rover as specified by the instructions string"
  [rover instructions plateau-bounds]
  (if (not-empty instructions)
    (move-rover
      ((rover-movements (first instructions)) rover)
      (rest instructions)
      plateau-bounds)
    (if (check-bounds rover plateau-bounds)
      rover
      rover-out-of-bounds-error)))

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
                              {:rover        {:x       x
                                              :y       y
                                              :heading heading}
                               :instructions (->>
                                               (second rover-instructions)
                                               (map str)
                                               (vec))}))))})

(defn move-rovers-sequentially
  "Moves rovers based on the instruction string"
  [{:keys [plateau-bounds rovers]}]
  (map #(move-rover
          (:rover %)
          (:instructions %)
          plateau-bounds)
       rovers))