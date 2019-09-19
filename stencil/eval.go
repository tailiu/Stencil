/*
 * Physical Migration Handler
 */

 package main

 import (
	 "stencil/evaluation"
 )
 
 func main() {
	 evalConfig := evaluation.InitializeEvalConfig()
 
	 evaluation.GetDataBagOfUser("2038186478", "diaspora", "gnusocial", evalConfig)
		 
 }
 