(ns solutions-4clojure.medium)

;Problem 43, Reverse Interleave
;Write a function which reverses the interleave process into x number of subsequences.
;(= (rve [1 2 3 4 5 6] 2) '((1 3 5) (2 4 6)))
;(= (__ (range 9) 3) '((0 3 6) (1 4 7) (2 5 8)))
;(= (__ (range 10) 5) '((0 5) (1 6) (2 7) (3 8) (4 9)))

(defn custom-split [sequence index]
  (for [i (range index)]
    (take-nth
      index
      (drop i sequence))))

(= (custom-split [1 2 3 4 5 6 7] 2) '((1 3 5 7) (2 4 6)))
(= (custom-split (range 9) 3) '((0 3 6) (1 4 7) (2 5 8)))
(= (custom-split (range 10) 5) '((0 5) (1 6) (2 7) (3 8) (4 9)))

(defn custom-split-simple [collection index]
  (apply
    map
    vector
    (partition
      index
      collection)))

(= (custom-split-simple [1 2 3 4 5 6] 2) '((1 3 5) (2 4 6)))
(= (custom-split-simple (range 9) 3) '((0 3 6) (1 4 7) (2 5 8)))
(= (custom-split-simple (range 10) 5) '((0 5) (1 6) (2 7) (3 8) (4 9)))

;Problem 44, Rotate Sequence
;Write a function which can rotate a sequence in either direction.
;(= (__ 2 [1 2 3 4 5]) '(3 4 5 1 2))
;(= (__ -2 [1 2 3 4 5]) '(4 5 1 2 3))
;(= (__ 6 [1 2 3 4 5]) '(2 3 4 5 1))
;(= (__ 1 '(:a :b :c)) '(:b :c :a))
;(= (__ -4 '(:a :b :c)) '(:c :a :b))

(defn rotate-sequence [n collection]
  (flatten (reverse
             (split-at
               (mod n (count collection))
               collection))))

(= (rotate-sequence 2 [1 2 3 4 5]) '(3 4 5 1 2))
(= (rotate-sequence -2 [1 2 3 4 5]) '(4 5 1 2 3))
(= (rotate-sequence 6 [1 2 3 4 5]) '(2 3 4 5 1))
(= (rotate-sequence 1 '(:a :b :c)) '(:b :c :a))
(= (rotate-sequence -4 '(:a :b :c)) '(:c :a :b))


;Problem 46, Flipping out
;Write a higher-order function which flips the order of the arguments of an input function.
;(= true ((() >) 7 8))
;(= 4 ((() quot) 2 8))
;(= [1 2 3] ((() take) [1 2 3 4 5] 3))

(defn flip-args [function]
  (fn [& args]
    (apply function (reverse args))))

(= 3 ((flip-args nth) 2 [1 2 3 4 5]))
(= true ((flip-args >) 7 8))
(= 4 ((flip-args quot) 2 8))
(= [1 2 3] ((flip-args take) [1 2 3 4 5] 3))



