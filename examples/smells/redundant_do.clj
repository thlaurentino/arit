(ns examples.smells.redundant-do)

;; O `do` é desnecessário dentro de `when`.
(defn print-message [message]
  (when message
    (do
      (println message)
      true)))
