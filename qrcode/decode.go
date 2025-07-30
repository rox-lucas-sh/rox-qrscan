package qrcode

import (

	// Importa os decodificadores de imagem (PNG, JPEG, etc.)
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"gocv.io/x/gocv"
)

// None of these functions really work.

func Decode(imageBytes []byte) (url string, err error) {
	// 1. Decodifica o []byte para um objeto image.Image
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return "", fmt.Errorf("falha ao decodificar a imagem: %w", err)
	}

	// 2. Prepara a imagem para o gozxing
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("falha ao criar bitmap: %w", err)
	}

	// 3. Decodifica o QR Code
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", fmt.Errorf("falha ao decodificar o QR Code: %w", err)
	}

	return result.GetText(), nil
}

func DecodeQRCodeWithCV(imageBytes []byte) (string, error) {
	// Decodifica os bytes para a matriz de imagem da OpenCV
	img, err := gocv.IMDecode(imageBytes, gocv.IMReadUnchanged)
	if err != nil {
		return "", fmt.Errorf("falha ao decodificar bytes para imagem: %w", err)
	}
	if img.Empty() {
		return "", fmt.Errorf("imagem de entrada está vazia")
	}
	defer img.Close()

	// Cria o detector
	detector := gocv.NewQRCodeDetector()
	defer detector.Close()

	// --- CORREÇÃO AQUI ---
	// 1. Crie as matrizes para receber os dados de saída
	points := gocv.NewMat()
	straightQR := gocv.NewMat()

	// 2. Garanta que elas serão fechadas ao final da função
	defer points.Close()
	defer straightQR.Close()

	// 3. Passe os ponteiros para as matrizes na chamada da função
	result := detector.DetectAndDecode(img, &points, &straightQR)
	if result == "" {
		return "", fmt.Errorf("nenhum QR Code foi encontrado")
	}

	// --- Opcional: O que fazer com as informações ---

	// 1. Desenhar um contorno ao redor do QR Code detectado
	if !points.Empty() {
		// CORREÇÃO AQUI: Iteramos manualmente para extrair os pontos
		// CORREÇÃO FINAL AQUI: Usar o tipo image.Point da biblioteca padrão
		numPoints := points.Rows()
		pts := make([]image.Point, numPoints) // Usa image.Point
		for i := 0; i < numPoints; i++ {
			x := points.GetFloatAt(i, 0)
			y := points.GetFloatAt(i, 1)
			pts[i] = image.Pt(int(x), int(y)) // Usa image.Pt para criar o ponto
		}

		poly := gocv.NewPointsVectorFromPoints([][]image.Point{pts}) // O slice também usa image.Point
		gocv.Polylines(&img, poly, true, color.RGBA{R: 0, G: 255, B: 0, A: 255}, 2)

		// Salva a imagem com o contorno para verificação
		outputFile := "imagem_com_contorno.png"
		if ok := gocv.IMWrite(outputFile, img); !ok {
			fmt.Println("Falha ao salvar a imagem com contorno")
		} else {
			fmt.Printf("Imagem com contorno salva em: %s\n", outputFile)
		}
	}

	// 2. Salvar a imagem do QR Code "limpa" e retificada
	if !straightQR.Empty() {
		outputFile := "qrcode_retificado.png"
		if ok := gocv.IMWrite(outputFile, straightQR); !ok {
			fmt.Println("Falha ao salvar a imagem retificada")
		} else {
			fmt.Printf("QR Code retificado salvo em: %s\n", outputFile)
		}
	}

	return result, nil
}

func FindAndDrawQRCode_Safe(imageBytes []byte) (result string, err error) {
	// --- Bloco de Recuperação de Pânico ---
	// Se qualquer chamada da GoCV quebrar (panic), este bloco será executado.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("pânico recuperado na GoCV: %v. Isso indica um problema sério na instalação da OpenCV ou um bug na biblioteca. Verifique sua instalação (pkg-config --cflags --libs opencv4) e reinstale se necessário", r)
		}
	}()

	img, err := gocv.IMDecode(imageBytes, gocv.IMReadColor)
	if err != nil {
		return "", fmt.Errorf("falha ao decodificar bytes para imagem: %w", err)
	}
	if img.Empty() {
		return "", fmt.Errorf("imagem de entrada está vazia")
	}
	defer img.Close()

	// Simplificando o pré-processamento: Às vezes, menos é mais.
	// O thresholding pode ter sido agressivo demais.
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	detector := gocv.NewQRCodeDetector()
	defer detector.Close() // Defer ainda é seguro.

	points := gocv.NewMat()
	defer points.Close()

	// Tentaremos a detecção na imagem em escala de cinza.
	result = detector.DetectAndDecode(gray, &points, nil)
	if result == "" {
		return "", fmt.Errorf("nenhum QR Code foi encontrado na imagem")
	}

	// Se chegou aqui, sucesso! Vamos desenhar.
	if !points.Empty() {
		fmt.Println("QR Code encontrado! Desenhando contorno...")

		numPoints := points.Rows()
		pts := make([]image.Point, numPoints)
		for i := range numPoints {
			x := points.GetFloatAt(i, 0)
			y := points.GetFloatAt(i, 1)
			pts[i] = image.Pt(int(x), int(y))
		}

		poly := gocv.NewPointsVectorFromPoints([][]image.Point{pts})
		gocv.Polylines(&img, poly, true, color.RGBA{R: 0, G: 255, B: 0, A: 255}, 3)

		outputFile := "resultado_final_com_contorno.png"
		if ok := gocv.IMWrite(outputFile, img); !ok {
			fmt.Println("Falha ao salvar a imagem com contorno")
		} else {
			fmt.Printf("Imagem com contorno salva em: %s\n", outputFile)
		}
	}

	return result, nil
}
