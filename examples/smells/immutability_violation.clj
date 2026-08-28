(ns examples.smells.immutability-violation)

;; Alteração de estado global dentro de uma função.
(def counter (atom 0))

(defn increment-counter []
  (reset! counter (inc @counter)))
