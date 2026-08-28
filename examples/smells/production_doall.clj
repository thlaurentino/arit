(ns examples.smells.production-doall)

;; `doall` força a realização de uma sequência lazy.
(defn print-users [users]
  (doall (map println users)))
