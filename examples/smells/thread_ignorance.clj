(ns examples.smells.thread-ignorance)

;; O fluxo de dados pode ser escrito com a macro `->`.
(defn normalize-name [name]
  (clojure.string/trim (clojure.string/lower-case name)))
