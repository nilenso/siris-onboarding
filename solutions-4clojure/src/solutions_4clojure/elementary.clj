(ns solutions-4clojure.elementary)
(require '(clojure.set))

;Problem 5, conj on lists
;(= __ (conj '(2 3 4) 1))
;(= __ (conj '(3 4) 2 1))
(= '(1 2 3 4) (conj '(2 3 4) 1))
(= '(1 2 3 4) (conj '(3 4) 2 1))

;Problem 6, Vectors
;(= [__] (list :a :b :c) (vec '(:a :b :c)) (vector :a :b :c))
(= [:a :b :c]
   (list :a :b :c)
   (vec '(:a :b :c))
   (vector :a :b :c))

;Problem 7, conj on vectors
;(= __ (conj [1 2 3] 4))
;(= __ (conj [1 2] 3 4))
(= [1 2 3 4] (conj [1 2 3] 4))
(= [1 2 3 4] (conj [ 1 2] 3 4))

;Problem 8, Sets
;(= __ (set '(:a :a :b :c :c :c :c :d :d)))
;(= __ (clojure.set/union #{:a :b :c} #{:b :c :d}))
(= #{:a :b :c :d} (set '(:a :a :b :c :c :c :c :d :d)))
(= #{:a :b :c :d} (clojure.set/union #{:a :b :c} #{:b :c :d}))

;Problem 9, conj on sets
;(= #{1 2 3 4} (conj #{1 4 3} __))
(= #{1 2 3 4} (conj #{1 4 3} 2))

;Problem 10, Maps
;(= __ ((hash-map :a 10, :b 20, :c 30) :b))
;(= __ (:b {:a 10, :b 20, :c 30}))
(= 20 ((hash-map :a 10, :b 20, :c 30) :b))
(= 20 (:b { :a 10, :b 20, :c 30}))

;Problem 11, conj on maps
;(= {:a 1, :b 2, :c 3} (conj {:a 1} __ [:c 3]))
(= {:a 1, :b 2, :c 3} (conj {:a 1} [:b 2] [:c 3]))

;Problem 12, Sequences
;(= __ (first '(3 2 1)))
;(= __ (second [2 3 4]))
;(= __ (last (list 1 2 3)))
(= 3 (first '(3 2 1)))
(= 3 (second [2 3 4]))
(= 3 (last (list 1 2 3)))

;Problem 13, rest
;(= __ (rest [10 20 30 40]))
(= [20 30 40] (rest [10 20 30 40]))

;Problem 14, Functions
;(= __ ((fn add-five [x] (+ x 5)) 3))
;(= __ ((fn [x] (+ x 5)) 3))
;(= __ (#(+ % 5) 3))
;(= __ ((partial + 5) 3))
(= 8 ((fn add-five [x] (+ x 5)) 3))
(= 8 ((fn [x] (+ x 5)) 3))
(= 8 (#(+ % 5) 3))
(= 8 ((partial + 5) 3))

;Problem 15, Double Down
;Write a function which doubles a number.
;(= (__ 2) 4)
;(= (__ 3) 6)
;(= (__ 11) 22)
;(= (__ 7) 14)
(= ((fn twice [x] (* x 2)) 2) 4)
(= ((fn twice [x] (* x 2)) 3) 6)
(= ((fn twice [x] (* x 2)) 11) 22)
(= ((fn twice [x] (* x 2)) 7) 14)

;Problem 16, Hello World
;Write a function which returns a personalized greeting.
;(= (__ "Dave") "Hello, Dave!")
;(= (__ "Jenn") "Hello, Jenn!")
;(= (__ "Rhea") "Hello, Rhea!")
(defn hello [name]
  (format "Hello, %s!" name))
(= (hello "Dave") "Hello, Dave!")
(= (hello "Jenn") "Hello, Jenn!")
(= (hello "Rhea") "Hello, Rhea!")

;Problem 17, map
;(= __ (map #(+ % 5) '(1 2 3)))
(= '(6 7 8) (map #(+ % 5) '(1 2 3)))

;Problem 18, filter
;(= __ (filter #(> % 5) '(3 4 5 6 7)))
(= '(6 7) (filter
            #(> % 5)
            '(3 4 5 6 7)))

;Problem 35, Local bindings
;(= __ (let [x 5] (+ 2 x)))
;(= __ (let [x 3, y 10] (- y x)))
;(= __ (let [x 21] (let [y 3] (/ x y))))
(= 7 (let [x 5] (+ 2 x)))
(= 7 (let [x 3, y 10] (- y x)))
(= 7 (let [x 21] (let [y 3] (/ x y))))

;Problem 36, Let it Be
;(= 10 (let __ (+ x y)))
;(= 4 (let __ (+ y z)))
;(= 1 (let __ z))
(= 10 (let [x 7 y 3] (+ x y)))
(= 4 (let [y 3 z 1] (+ y z)))
(= 1 (let [z 1] z))

;Problem 37, Regular Expressions
;(= __ (apply str (re-seq #"[A-Z]+" "bA1B3Ce ")))
(= "ABC" (apply str
                (re-seq #"[A-Z]+" "bA1B3Ce ")))

;Problem 52, Intro to Destructuring
;Let bindings and function parameter lists support destructuring.
;(= [2 4] (let [[a b c d e f g] (range)] __))

(= [2 4] (let [[a b c d e f g] (range)] [c e]))



