(ns mars-rover-clj.mars-rover)

(def directions
  {"N" {:compass-angle 0
        :move-f        (fn [rover] (update rover :y inc))}
   "E" {:compass-angle 90
        :move-f        (fn [rover] (update rover :x inc))}
   "S" {:compass-angle 180
        :move-f        (fn [rover] (update rover :y dec))}
   "W" {:compass-angle 270
        :move-f        (fn [rover] (update rover :x dec))}})

(def out-of-bounds-error
  "Error string when a bad instruction is passed making the rover go off-grid."
  "Bad instructions. Rover overshoots plateau bounds.")

(def
  ^{:doc "Returns the direction (N, E, S or W) from the compass angle."}
  compass-angle-to-direction
  (->> directions
    (map (fn [[direction {:keys [compass-angle]}]]
           [compass-angle direction]))
    (into {})))

(defn turn-heading
  "Returns the new heading of the rover after turning it by direction-fn."
  [heading direction-fn turn-by]
  (-> (:compass-angle (directions heading))
      (#(direction-fn % turn-by))
      (mod 360)
      (compass-angle-to-direction)))

(defn move
  "Returns a rover with updated position (x and y coordinates) and the heading."
  [rover]
  ((:move-f (->
              (:heading rover)
              (directions)))
   rover))

(def commands
  "A map of instruction symbols and functions to move rovers."
  {"L" (fn [rover]
         (update rover :heading #(turn-heading % - 90)))
   "R" (fn [rover]
         (update rover :heading #(turn-heading % + 90)))
   "M" move})

(defn check-bounds
  "Returns true if the rover's position is within the plateau bounds, false otherwise."
  [rover plateau-bounds]
  (and
    (<= 0 (:x rover) (first plateau-bounds))
    (<= 0 (:y rover) (second plateau-bounds))))

(defn move-rover
  "Moves a rover as specified by the instructions vector."
  [rover
   [current-instruction & next-instructions]
   plateau-bounds]
  (cond
    (nil? current-instruction)
    (if (check-bounds rover plateau-bounds)
      rover
      out-of-bounds-error)
    :else (move-rover
            ((commands current-instruction) rover)
            next-instructions
            plateau-bounds)))

(defn init-rovers
  "Takes an input for the Mars Rover problem
  as a vector of the plateau's upper right coordinates and an
  infinite sequence of rover's initial position
  and a string of a series of instructions.
  Returns a map of plateau bounds and rovers."
  [plateau-upper-right-coordinates & rover-details]
  {:plateau-bounds plateau-upper-right-coordinates
   :rovers         (->> (partition 2 rover-details)
                     (map (fn [[rover instructions]]
                            (let [[x y heading] rover]
                              {:rover        {:x       x
                                              :y       y
                                              :heading heading}
                               :instructions (->>
                                               instructions
                                               (map str)
                                               (vec))}))))})
(comment
  (= (init-rovers [5 5]
                  "LMLMLMLMM"
                  [1 2 "N"]
                  [3 3 "E"]
                  "MMRMMRMRRM")
     '{:plateau-bounds [5 5],
       :rovers         (
                        {:rover        {:x 1, :y 2, :heading "N"}
                         :instructions ["L" "M" "L" "M" "L" "M" "L" "M" "M"]}
                        {
                         :rover        {:x 3, :y 3, :heading "E"}
                         :instructions ["M" "M" "R" "M" "M" "R" "M" "R" "R" "M"]})}))

(defn move-rovers-sequentially
  "Moves rovers based on the instruction string."
  [{:keys [plateau-bounds rovers]}]
  (map #(move-rover
          (:rover %)
          (:instructions %)
          plateau-bounds)
       rovers))

