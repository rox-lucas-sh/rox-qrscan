package bucket

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

// bucketKeyPath := "./gcp-bucket-key.json"
// os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", bucketKeyPath)

// file, err := os.Open("qrcodes.csv")
// if err != nil {
// 	panic(fmt.Sprintf("falha ao abrir o arquivo qrcodes.csv: %v", err))
// }
// defer file.Close()
// scrapping.ScrapFromFile(file)

var BucketName = "bucket-test-rox"

func DownloadFile(bucket, object string) ([]byte, error) {
	// `context.Background()` é um bom ponto de partida.
	ctx := context.Background()

	// `storage.NewClient` cria um novo cliente.
	// O mais legal é que ele encontra automaticamente suas credenciais
	// (aquelas que configuramos com `gcloud auth application-default login`).
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar cliente do storage: %w", err)
	}
	defer client.Close()

	// Usamos um timeout para a operação de download para não ficar esperando para sempre.
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Criamos um "leitor" (Reader) para o objeto no bucket.
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter o reader do objeto %q: %w", object, err)
	}
	defer rc.Close()

	// `io.ReadAll` lê todo o conteúdo do arquivo.
	body, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler o conteúdo do arquivo: %w", err)
	}

	fmt.Printf("Arquivo baixado com sucesso! Tamanho: %d mega bytes.\n", len(body)/(1024*1024))

	return body, nil
}

func UploadFile(bucket, object, filePath string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("falha ao criar cliente do storage: %w", err)
	}
	defer client.Close()

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	// O `Close` faz com que o buffer seja enviado para o GCS.
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	fmt.Printf("Arquivo %q enviado com sucesso para o bucket %q como %q.\n", filePath, bucket, object)
	return nil
}

func UploadImage(bucket, objectName string, image io.Reader) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("falha ao criar cliente do storage: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	wc := client.Bucket(bucket).Object(objectName).NewWriter(ctx)
	if _, err = io.Copy(wc, image); err != nil {
		return fmt.Errorf("falha ao copiar imagem: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("falha ao fechar writer: %w", err)
	}

	fmt.Printf("Imagem %q enviada com sucesso para o bucket %q.\n", objectName, bucket)
	return nil
}
