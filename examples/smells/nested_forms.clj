(ns examples.smells.nested-forms)

;; Aninhamento excessivo de condicionais.
(defn classify [value]
  (if value
    (if (number? value)
      (if (pos? value)
        :positive
        :non-positive)
      :not-a-number)
    :empty))
