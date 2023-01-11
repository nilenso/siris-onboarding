(ns dice-engine.factories
  (:import (java.util UUID)))

(defn get-id
  "Returns an auto incremented id"
  []
  (UUID/randomUUID))

(defn create-die
  "Returns a die map of the form
  {:id              (get-id)
   :value           n
   :faces           x
   :previous-values []}"
  [die-value faces]
  {:id              (get-id)
   :value           die-value
   :faces           faces
   :previous-values []})