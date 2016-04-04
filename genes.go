package main

type fenotype struct {
  red   uint8
  green uint8
  blue  uint8
  speed float64
}

type genotype struct {
  frontLegs bool
}

type gen struct {
  name  string
  apply func(genotype, fenotype) fenotype
}

var genes []gen

func init() {
}
