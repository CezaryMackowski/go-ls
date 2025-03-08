package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetFiles(t *testing.T) {
	// Tworzymy tymczasowy katalog
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Nie udało się utworzyć katalogu: %v", err)
	}
	// Usuwamy katalog po teście
	defer os.RemoveAll(tmpDir)

	// Tworzymy przykładowy plik
	filePath := filepath.Join(tmpDir, "testfile.txt")
	f, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Nie udało się utworzyć pliku: %v", err)
	}
	f.Close()

	// Tworzymy przykładowy podkatalog
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Nie udało się utworzyć podkatalogu: %v", err)
	}

	// Przygotowujemy przykładową konfigurację (upewnij się, że pola odpowiadają Twojej strukturze Config)
	config := NewConfig()

	files, columnsWidth, err := GetFiles(tmpDir, config)
	if err != nil {
		t.Fatalf("Błąd przy wywołaniu GetFiles: %v", err)
	}

	// Przykładowe asercje
	if len(files) != 2 {
		t.Errorf("Oczekiwano 2 wpisów, otrzymano %d", len(files))
	}
	// Dodatkowe testy: sprawdzenie czy plik o nazwie "testfile.txt" istnieje, a także czy są ustawione poprawne kolory, uprawnienia, itd.
	_ = columnsWidth // Możesz tu dodać asercje dotyczące szerokości kolumn
}
