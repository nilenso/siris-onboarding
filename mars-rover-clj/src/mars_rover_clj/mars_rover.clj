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
  "A map of command symbols and functions to move rovers."
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
  "Moves a rover as specified by the commands vector."
  [rover
   [current-command & subsequent-commands]
   plateau-bounds]
  (cond
    (nil? current-command)
    (if (check-bounds rover plateau-bounds)
      rover
      out-of-bounds-error)
    :else (move-rover
            ((commands current-command) rover)
            subsequent-commands
            plateau-bounds)))

(defn parse-rover-input
  "Takes a collection of vectors with alternate rover positions and instruction strings.
  Returns a map of rovers' initial positions and commands."
  [rover-input]
  (->> (partition 2 rover-input)
    (map (fn [[rover instructions]]
           (let [[x y heading] rover]
             {:rover    {:x       x
                             :y       y
                             :heading heading}
              :commands (->>
                              instructions
                              (map str)
                              (vec))})))))

(defn init-rovers
  "Takes an input for the Mars Rover problem
  as a vector of the plateau's upper right coordinates and an
  infinite sequence of rover's initial position
  and a string of a series of instructions.
  Returns a map of plateau bounds and rovers."
  [plateau-upper-right-coordinates & rover-input]
  {:plateau-bounds plateau-upper-right-coordinates
   :rovers         (parse-rover-input rover-input)})

(comment
  (= (init-rovers [5 5]
                  "LMLMLMLMM"
                  [1 2 "N"]
                  [3 3 "E"]
                  "MMRMMRMRRM")
     '{:plateau-bounds [5 5],
       :rovers         (
                        {:rover    {:x 1, :y 2, :heading "N"}
                         :commands ["L" "M" "L" "M" "L" "M" "L" "M" "M"]}
                        {
                         :rover    {:x 3, :y 3, :heading "E"}
                         :commands ["M" "M" "R" "M" "M" "R" "M" "R" "R" "M"]})}))

(defn move-rovers-sequentially
  "Returns positions of rovers after moving them sequentially."
  [{:keys [plateau-bounds rovers]}]
  (map #(move-rover
          (:rover %)
          (:commands %)
          plateau-bounds)
       rovers))

