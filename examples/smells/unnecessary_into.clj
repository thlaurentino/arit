(ns examples.smells.unnecessary-into)

;; `into` é usado sem necessidade para transformar o resultado de `map`.
(defn double-values [values]
  (into [] (map #(* 2 %) values)))
