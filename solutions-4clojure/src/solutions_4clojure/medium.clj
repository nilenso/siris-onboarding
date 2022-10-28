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

;Problem 50, Split by Type
;Write a function which takes a sequence consisting of items with different
; types and splits them up into a set of homogeneous sub-sequences.
; The internal order of each sub-sequence should be maintained, but the
; sub-sequences themselves can be returned in any order (this is why 'set'
; is used in the test cases).
;(= (set (__ [1 :a 2 :b 3 :c])) #{[1 2 3] [:a :b :c]})
;(= (set (__ [:a "foo"  "bar" :b])) #{[:a :b] ["foo" "bar"]})
;(= (set (__ [[1 2] :a [3 4] 5 6 :b])) #{[[1 2] [3 4]] [:a :b] [5 6]})

(defn split-by-type [collection]
  (vals
    (group-by type collection)))

(= (set (split-by-type [1 :a 2 :b 3 :c])) #{[1 2 3] [:a :b :c]})
(= (set (split-by-type [:a "foo" "bar" :b])) #{[:a :b] ["foo" "bar"]})
(= (set (split-by-type [[1 2] :a [3 4] 5 6 :b])) #{[[1 2] [3 4]] [:a :b] [5 6]})

;Problem 54, Partition a Sequence
;Write a function which returns a sequence of lists of x items each. Lists of less than x items should not be returned.
;(= (__ 3 (range 9)) '((0 1 2) (3 4 5) (6 7 8)))
;(= (__ 2 (range 8)) '((0 1) (2 3) (4 5) (6 7)))
;(= (__ 3 (range 8)) '((0 1 2) (3 4 5)))
(defn custom-partition [n collection]
  (if (>= (count collection) n)
    (conj
      (custom-partition n (drop n collection))
      (take n collection))
    '()))

(= (custom-partition 3 (range 9)) '((0 1 2) (3 4 5) (6 7 8)))
(= (custom-partition 2 (range 8)) '((0 1) (2 3) (4 5) (6 7)))
(= (custom-partition 3 (range 8)) '((0 1 2) (3 4 5)))

(custom-partition 3 (range 8))

;Problem 55, Count Occurrences
;Write a function which returns a map containing the number of occurrences of each distinct item in a sequence.
;(= (__ [:b :a :b :a :b]) {:a 2, :b 3})
;(= (__ '([1 2] [1 3] [1 3])) {[1 2] 1, [1 3] 2})
;(= (__ [1 1 2 3 2 1 1]) {1 4, 2 2, 3 1})
(defn custom-frequencies [collection]
  (reduce
    (fn [frequency-map element]
      (if-let [cnt (get frequency-map element)]
        (assoc frequency-map element (+ cnt 1))
        (assoc frequency-map element 1)))
    {}
    collection))

(= (custom-frequencies [:b :a :b :a :b]) {:a 2, :b 3})
(= (custom-frequencies '([1 2] [1 3] [1 3])) {[1 2] 1, [1 3] 2})
(= (custom-frequencies [1 1 2 3 2 1 1]) {1 4, 2 2, 3 1})

(defn custom-frequencies-v2 [collection]
  (reduce
    (fn [frequency-map element]
      (update
        frequency-map
        element
        #(inc (or % 0))))
    {}
    collection))

(= (custom-frequencies-v2 [:b :a :b :a :b]) {:a 2, :b 3})
(= (custom-frequencies-v2 '([1 2] [1 3] [1 3])) {[1 2] 1, [1 3] 2})
(= (custom-frequencies-v2 [1 1 2 3 2 1 1]) {1 4, 2 2, 3 1})

(defn custom-frequencies-simple [collection]
  (reduce-kv
    (fn [m k v]
      (assoc m k (count v)))
    {}
    (group-by identity collection)))

(= (custom-frequencies-simple [:b :a :b :a :b]) {:a 2, :b 3})
(= (custom-frequencies-simple '([1 2] [1 3] [1 3])) {[1 2] 1, [1 3] 2})
(= (custom-frequencies-simple [1 1 2 3 2 1 1]) {1 4, 2 2, 3 1})
