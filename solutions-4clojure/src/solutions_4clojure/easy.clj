(ns solutions-4clojure.easy)

;Problem 19, Last Element
;Difficulty: easy
;Write a function which returns the last element in a sequence.
;
;(= (__ [1 2 3 4 5]) 5)
;(= (__ '(5 4 3)) 3)
;(= (__ ["b" "c" "d"]) "d")
(defn get-last [x]
  (if (empty? (rest x))
    (first x)
    (get-last (rest x))))

(= (get-last [1 2 3 4 5]) 5)
(= (get-last '(5 4 3)) 3)
(= (get-last ["b" "c" "d"]) "d")

;Problem 20, Penultimate Element
;Difficulty: easy
;Write a function which returns the second to last element from a sequence.
;
;(= (__ (list 1 2 3 4 5)) 4)
;(= (__ ["a" "b" "c"]) "b")
;(= (__ [[1 2] [3 4]]) [1 2])

;(fn my-pen
;  [lst]
;  (if (empty? (rest (rest lst)))
;    (first lst)
;    (my-pen (rest lst))))
(defn penultimate [x]
  (if (empty? (rest (rest x)))
    (first x)
    (penultimate (rest x))))

(= (penultimate (list 1 2 3 4 5)) 4)
(= (penultimate ["a" "b" "c"]) "b")
(= (penultimate [[1 2] [3 4]]) [1 2])
(= (penultimate [1]) nil)                                   ;fails, penultimate doesn't check for collections of length 1

(defn penultimate-v2 [x]
  (let [c (count x)]
    (cond
      (= c 2) (first x)
      (> c 2) (penultimate-v2 (rest x))
      :else nil)))

(= (penultimate-v2 (list 1 2 3 4 5)) 4)
(= (penultimate-v2 ["a" "b" "c"]) "b")
(= (penultimate-v2 [[1 2] [3 4]]) [1 2])
(= (penultimate-v2 [1]) nil)

(defn penultimate-simple [x]
  (second (reverse x)))

(= (penultimate-simple (list 1 2 3 4 5)) 4)
(= (penultimate-simple ["a" "b" "c"]) "b")
(= (penultimate-simple [[1 2] [3 4]]) [1 2])
(= (penultimate-simple [1]) nil)

;Problem 21, Nth Element
;Write a function which returns the Nth element from a sequence.
;(= (__ '(4 5 6 7) 2) 6)
;(= (__ [:a :b :c] 0) :a)
;(= (__ [1 2 3 4] 1) 2)
;(= (__ '([1 2] [3 4] [5 6]) 2) [5 6])
;(defn get-nth [x n]
;  reduce (fn [x index n] (if (= index n)
;                           (first x)
;                           (get-nth (rest x) n))) x)

(defn get-nth [x n]
  (reduce (fn [index elem]
            (if (= n index)
              (reduced elem)
              (+ 1 index))) 0 x))

(= (get-nth '(4 5 6 7) 2) 6)
(= (get-nth [:a :b :c] 0) :a)
(= (get-nth [1 2 3 4] 1) 2)
(= (get-nth '([1 2] [3 4] [5 6]) 2) [5 6])

(defn get-nth-simple [x n]
  (last (take (+ n 1) x)))

(= (get-nth-simple '(4 5 6 7) 2) 6)
(= (get-nth-simple [:a :b :c] 0) :a)
(= (get-nth-simple [1 2 3 4] 1) 2)
(= (get-nth-simple '([1 2] [3 4] [5 6]) 2) [5 6])

;Problem 22, Count a Sequence
;Write a function which returns the total number of elements in a sequence.
;(= (__ '(1 2 3 3 1)) 5)
;(= (__ "Hello World") 11)
;(= (__ [[1 2] [3 4] [5 6]]) 3)
;(= (__ '(13)) 1)
;(= (__ '(:a :b :c)) 3)
(defn count-seq [x]
  (reduce (fn [c _] (+ c 1)) 0 x))

(= (count-seq '(1 2 3 3 1)) 5)
(= (count-seq "Hello World") 11)
(= (count-seq [[1 2] [3 4] [5 6]]) 3)
(= (count-seq '(13)) 1)
(= (count-seq '(:a :b :c)) 3)

;Problem 23, Reverse a Sequence
;Write a function which reverses a sequence.
;(= (__ [1 2 3 4 5]) [5 4 3 2 1])
;(= (__ (sorted-set 5 7 2 7)) '(7 5 2))
;(= (__ [[1 2][3 4][5 6]]) [[5 6][3 4][1 2]])
(defn reverse-seq [x]
  (reduce (fn [rev elem] (cons elem rev)) nil x))
(= (reverse-seq [1 2 3 4 5]) [5 4 3 2 1])
(= (reverse-seq (sorted-set 5 7 2 7)) '(7 5 2))
(= (reverse-seq [[1 2][3 4][5 6]]) [[5 6][3 4][1 2]])

;Problem 24, Sum It All Up
;Write a function which returns the sum of a sequence of numbers.
;(= (__ [1 2 3]) 6)
;(= (__ (list 0 -2 5 5)) 8)
;(= (__ #{4 2 1}) 7)
;(= (__ '(0 0 -1)) -1)
;(= (__ '(1 10 3)) 14)
(defn sum [x]
  (reduce (fn [s elem] (+ s elem)) 0 x))
(= (sum [1 2 3]) 6)
(= (sum (list 0 -2 5 5)) 8)
(= (sum #{4 2 1}) 7)
(= (sum '(0 0 -1)) -1)
(= (sum '(1 10 3)) 14)

;Problem 25, Find the odd numbers
;Write a function which returns only the odd numbers from a sequence.
;(= (__ #{1 2 3 4 5}) '(1 3 5))
;(= (__ [4 2 1 6]) '(1))
;(= (__ [2 2 4 6]) '())
;(= (__ [1 1 1 3]) '(1 1 1 3))
(defn odd [x]
  (filter odd? x))
(= (odd #{1 2 3 4 5}) '(1 3 5))
(= (odd [4 2 1 6]) '(1))
(= (odd [2 2 4 6]) '())
(= (odd [1 1 1 3]) '(1 1 1 3))

;Problem 26, Fibonacci Sequence
;Write a function which returns the first X fibonacci numbers.
;(= (__ 3) '(1 1 2))
;(= (__ 6) '(1 1 2 3 5 8))
;(= (__ 8) '(1 1 2 3 5 8 13 21))
(defn fibonacci [n]
  (cond
  (= n 1) '(0)
  (= n 2) '(0 1)
  :else (reduce (fn [out _]
                  (conj
                    out
                    (+ (last out) (nth out (- (count out) 2)))))
                [1 1]
                (range 0 (- n 2)))))

(= (fibonacci 3) '(1 1 2))
(= (fibonacci 6) '(1 1 2 3 5 8))
(= (fibonacci 8) '(1 1 2 3 5 8 13 21))

