(ns solutions-4clojure.easy)

;Problem 19, Last Element
;Write a function which returns the last element in a sequence.
;
;(= (__ [1 2 3 4 5]) 5)
;(= (__ '(5 4 3)) 3)
;(= (__ ["b" "c" "d"]) "d")
(defn get-last [lst]
  (if (empty? (rest lst))
    (first lst)
    (get-last (rest lst))))

(= (get-last [1 2 3 4 5]) 5)
(= (get-last '(5 4 3)) 3)
(= (get-last ["b" "c" "d"]) "d")

(defn get-last-simple [l]
  (first (reverse l)))
(= (get-last-simple [1 2 3 4 5]) 5)
(= (get-last-simple '(5 4 3)) 3)
(= (get-last-simple ["b" "c" "d"]) "d")

;Problem 20, Penultimate Element
;Write a function which returns the second to last element from a sequence.
;(= (__ (list 1 2 3 4 5)) 4)
;(= (__ ["a" "b" "c"]) "b")
;(= (__ [[1 2] [3 4]]) [1 2])

(defn penultimate [lst]
  (if (empty? (rest (rest lst)))
    (first lst)
    (penultimate (rest lst))))

(= (penultimate (list 1 2 3 4 5)) 4)
(= (penultimate ["a" "b" "c"]) "b")
(= (penultimate [[1 2] [3 4]]) [1 2])
(= (penultimate [1]) nil)                                   ;fails, penultimate doesn't check for collections of length 1

(defn penultimate-v2 [lst]
  (let [c (count lst)]
    (cond
      (= c 2) (first lst)
      (> c 2) (penultimate-v2 (rest lst))
      :else nil)))

(= (penultimate-v2 (list 1 2 3 4 5)) 4)
(= (penultimate-v2 ["a" "b" "c"]) "b")
(= (penultimate-v2 [[1 2] [3 4]]) [1 2])
(= (penultimate-v2 [1]) nil)

(defn penultimate-simple [lst]
  (second (reverse lst)))

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

(defn get-nth [lst n]
  (reduce (fn [index element]
            (if (= n index)
              (reduced element)
              (+ 1 index)))
          0
          lst))

(= (get-nth '(4 5 6 7) 2) 6)
(= (get-nth [:a :b :c] 0) :a)
(= (get-nth [1 2 3 4] 1) 2)
(= (get-nth '([1 2] [3 4] [5 6]) 2) [5 6])

(defn get-nth-simple [lst index]
  (last (take (+ index 1) lst)))

(= (get-nth-simple '(4 5 6 7) 2) 6)
(= (get-nth-simple [:a :b :c] 0) :a)
(= (get-nth-simple [1 2 3 4] 1) 2)
(= (get-nth-simple '([1 2] [3 4] [5 6]) 2) [5 6])

(defn get-nth-simple-v2 [collection index]
  (first (drop
           index
           collection)))
(= (get-nth-simple-v2 '(4 5 6 7) 2) 6)
(= (get-nth-simple-v2 [:a :b :c] 0) :a)
(= (get-nth-simple-v2 [1 2 3 4] 1) 2)
(= (get-nth-simple-v2 '([1 2] [3 4] [5 6]) 2) [5 6])

;Problem 22, Count a Sequence
;Write a function which returns the total number of elements in a sequence.
;(= (__ '(1 2 3 3 1)) 5)
;(= (__ "Hello World") 11)
;(= (__ [[1 2] [3 4] [5 6]]) 3)
;(= (__ '(13)) 1)
;(= (__ '(:a :b :c)) 3)
(defn count-seq [collection]
  (reduce (fn [c _] (+ c 1))
          0
          collection))

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
(defn reverse-seq [collection]
  (reduce (fn [rev elem]
            (cons elem rev))
          '()
          collection))
(= (reverse-seq [1 2 3 4 5]) [5 4 3 2 1])
(= (reverse-seq (sorted-set 5 7 2 7)) '(7 5 2))
(= (reverse-seq [[1 2] [3 4] [5 6]]) [[5 6] [3 4] [1 2]])

;Problem 24, Sum It All Up
;Write a function which returns the sum of a sequence of numbers.
;(= (__ [1 2 3]) 6)
;(= (__ (list 0 -2 5 5)) 8)
;(= (__ #{4 2 1}) 7)
;(= (__ '(0 0 -1)) -1)
;(= (__ '(1 10 3)) 14)
(defn sum [lst]
  (reduce + lst))

(= (sum [1 2 3]) 6)
(= (sum (list 0 -2 5 5)) 8)
(= (sum #{4 2 1}) 7)
(= (sum '(0 0 -1)) -1)
(= (sum '(1 10 3)) 14)

(defn sum-v2 [collection]
  (apply + collection))

(= (sum-v2 [1 2 3]) 6)
(= (sum-v2 (list 0 -2 5 5)) 8)
(= (sum-v2 #{4 2 1}) 7)
(= (sum-v2 '(0 0 -1)) -1)
(= (sum-v2 '(1 10 3)) 14)


;Problem 25, Find the odd numbers
;Write a function which returns only the odd numbers from a sequence.
;(= (__ #{1 2 3 4 5}) '(1 3 5))
;(= (__ [4 2 1 6]) '(1))
;(= (__ [2 2 4 6]) '())
;(= (__ [1 1 1 3]) '(1 1 1 3))
(defn odd [lst]
  (filter odd? lst))

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
                      (+
                        (last out)
                        (nth out
                             (-
                               (count out)
                               2)))))
                  [1 1]
                  (range 0 (- n 2)))))

(= (fibonacci 3) '(1 1 2))
(= (fibonacci 6) '(1 1 2 3 5 8))
(= (fibonacci 8) '(1 1 2 3 5 8 13 21))

(defn fibonacci-simple [n]
  (last
    (take
      (dec n)
      (iterate
        (fn [fibonacci-sequence]
          (conj fibonacci-sequence
                (+ (last fibonacci-sequence)
                   (penultimate fibonacci-sequence))))
        [1 1]))))

(= (fibonacci-simple 3) '(1 1 2))
(= (fibonacci-simple 6) '(1 1 2 3 5 8))
(= (fibonacci-simple 8) '(1 1 2 3 5 8 13 21))

;Problem 45, Intro to Iterate
;The iterate function can be used to produce an infinite lazy sequence.
;(= __ (take 5 (iterate #(+ 3 %) 1)))
(= '(1 4 7 10 13) (take 5 (iterate #(+ 3 %) 1)))

;Problem 47, Contain Yourself
;The contains? function checks if a KEY is present in a given collection.
;(contains? #{4 5 6} __)
;(contains? [1 1 1 1 1] __)
;(contains? {4 :a 2 :b} __)
;(not (contains? [1 2 4] __))
(contains? #{4 5 6} 4)
(contains? [1 1 1 1 1] 4)
(contains? {4 :a 2 :b} 4)
(not (contains? [1 2 4] 4))

;Problem 48, Intro to some
;The some function takes a predicate function and a collection. It returns the first logical true value of (predicate x) where x is an item in the collection.
;
;(= __ (some #{2 7 6} [5 6 7 8]))
;(= __ (some #(when (even? %) %) [5 6 7 8]))
(= 6 (some #{2 7 6} [5 6 7 8]))
(= 6 (some #(when (even? %) %) [5 6 7 8]))

;Problem 49, Split a sequence
;Difficulty: easy
;Write a function which will split a sequence into two parts.
;
;(= (__ 3 [1 2 3 4 5 6]) [[1 2 3] [4 5 6]])
;(= (__ 1 [:a :b :c :d]) [[:a] [:b :c :d]])
;(= (__ 2 [[1 2] [3 4] [5 6]]) [[[1 2] [3 4]] [[5 6]]])

(defn custom-split-at [n collection]
  (vector
    (take n collection)
    (drop n collection)))

(= (custom-split-at 3 [1 2 3 4 5 6]) [[1 2 3] [4 5 6]])
(= (custom-split-at 1 [:a :b :c :d]) [[:a] [:b :c :d]])
(= (custom-split-at 2 [[1 2] [3 4] [5 6]]) [[[1 2] [3 4]] [[5 6]]])

(defn custom-split-at-v2 [n collection]
  ((juxt take drop)
   n
   collection))

(= (custom-split-at-v2 3 [1 2 3 4 5 6]) [[1 2 3] [4 5 6]])
(= (custom-split-at-v2 1 [:a :b :c :d]) [[:a] [:b :c :d]])
(= (custom-split-at-v2 2 [[1 2] [3 4] [5 6]]) [[[1 2] [3 4]] [[5 6]]])

;Problem 51, Advanced Destructuring
;Here is an example of some more sophisticated destructuring.
;(= [1 2 [3 4 5] [1 2 3 4 5]] (let [[a b & c :as d] __] [a b c d]))

(= [1 2 [3 4 5] [1 2 3 4 5]] (let [[a b & c :as d] (range 1 6)] [a b c d]))


