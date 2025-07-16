package main

import (
	"fmt"

	"gocv.io/x/gocv"
)

func MinimalDetectorTest() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("pânico no teste mínimo: %v", r)
		}
	}()

	fmt.Println("Criando QRCodeDetector...")
	detector := gocv.NewQRCodeDetector()

	// A chamada .Close() também pode causar o pânico se a criação falhou
	// e deixou o ponteiro interno C inválido.
	fmt.Println("Fechando QRCodeDetector...")
	detector.Close()

	fmt.Println("Teste mínimo concluído sem pânico.")
	return nil
}

func main() {
	fmt.Println("Iniciando teste de isolamento...")
	err := MinimalDetectorTest()
	if err != nil {
		fmt.Printf("\n--- ERRO CAPTURADO ---\n%v\n", err)
	}
}
