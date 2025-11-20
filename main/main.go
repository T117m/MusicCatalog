package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"

    "github.com/T117m/MusicCatalog/music"
    "github.com/T117m/MusicCatalog/storage"
    "github.com/T117m/MusicCatalog/player"
)

func main() {
    store, err := storage.New()
    if err != nil {
        fmt.Printf("Ошибка инициализации хранилища: %v\n", err)
        return
    }
    defer store.Close()

    p := player.New()
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Println("\n=== МУЗЫКАЛЬНЫЙ КАТАЛОГ ===")
        fmt.Println("1. Вывести все треки")
        fmt.Println("2. Добавить трек")
        fmt.Println("3. Удалить трек")
        fmt.Println("4. Найти треки по автору")
        fmt.Println("5. Воспроизвести трек")
        fmt.Println("6. Остановить воспроизведение")
        fmt.Println("0. Выход")
        fmt.Print("Выберите опцию: ")

        scanner.Scan()
        input := strings.TrimSpace(scanner.Text())

        switch input {
        case "1":
            displayAllTracks(store)
        case "2":
            addTrack(store, scanner)
        case "3":
            removeTrack(store, scanner)
        case "4":
            searchByArtist(store, scanner)
        case "5":
            playTrack(store, p, scanner)
        case "6":
            stopPlayback(p)
        case "0":
            fmt.Println("Выход...")
            return
        default:
            fmt.Println("Неверная опция!")
        }
    }
}

func displayAllTracks(store *storage.Storage) {
    tracks, err := store.GetAllTracks()
    if err != nil {
        fmt.Printf("Ошибка получения треков: %v\n", err)
        return
    }

    if len(tracks) == 0 {
        fmt.Println("Каталог пуст")
        return
    }

    fmt.Println("\n--- Все треки ---")
    for _, track := range tracks {
        fmt.Printf("ID: %d | %s - %s (%s) [%s]\n", 
            track.ID, track.Artist, track.Title, track.Genre, track.FileType)
    }
}

func addTrack(store *storage.Storage, scanner *bufio.Scanner) {
    fmt.Print("Название трека: ")
    scanner.Scan()
    title := strings.TrimSpace(scanner.Text())

    fmt.Print("Исполнитель: ")
    scanner.Scan()
    artist := strings.TrimSpace(scanner.Text())

    fmt.Print("Жанр: ")
    scanner.Scan()
    genre := strings.TrimSpace(scanner.Text())

    fmt.Print("Тип файла (mp3, wav, etc): ")
    scanner.Scan()
    fileType := strings.TrimSpace(scanner.Text())

    fmt.Print("Путь к файлу: ")
    scanner.Scan()
    filePath := strings.TrimSpace(scanner.Text())

    track := music.Track{
        Title:    title,
        Artist:   artist,
        Genre:    genre,
        FileType: fileType,
        FilePath: filePath,
    }

    if err := store.AddTrack(&track); err != nil {
        fmt.Printf("Ошибка добавления трека: %v\n", err)
    } else {
        fmt.Printf("Трек добавлен с ID: %d\n", track.ID)
    }
}

func removeTrack(store *storage.Storage, scanner *bufio.Scanner) {
    displayAllTracks(store)
    
    fmt.Print("Введите ID трека для удаления: ")
    scanner.Scan()
    idStr := strings.TrimSpace(scanner.Text())

    id, err := strconv.Atoi(idStr)
    if err != nil {
        fmt.Println("Неверный ID")
        return
    }

    if err := store.RemoveTrackByID(id); err != nil {
        fmt.Printf("Ошибка удаления: %v\n", err)
    } else {
        fmt.Println("Трек удален")
    }
}

func searchByArtist(store *storage.Storage, scanner *bufio.Scanner) {
    fmt.Print("Введите имя исполнителя: ")
    scanner.Scan()
    artist := strings.TrimSpace(scanner.Text())

    tracks, err := store.GetTracksByArtist(artist)
    if err != nil {
        fmt.Printf("Ошибка поиска: %v\n", err)
        return
    }

    if len(tracks) == 0 {
        fmt.Println("Треки не найдены")
        return
    }

    fmt.Printf("\n--- Треки исполнителя %s ---\n", artist)
    for _, track := range tracks {
        fmt.Printf("ID: %d | %s - %s (%s) [%s]\n", 
            track.ID, track.Artist, track.Title, track.Genre, track.FileType)
    }
}

func playTrack(store *storage.Storage, p *player.Player, scanner *bufio.Scanner) {
    displayAllTracks(store)
    
    fmt.Print("Введите ID трека для воспроизведения: ")
    scanner.Scan()
    idStr := strings.TrimSpace(scanner.Text())

    id, err := strconv.Atoi(idStr)
    if err != nil {
        fmt.Println("Неверный ID")
        return
    }

    track, err := store.GetTrackByID(id)
    if err != nil {
        fmt.Printf("Трек не найден: %v\n", err)
        return
    }

    if err := p.Play(&track); err != nil {
        fmt.Printf("Ошибка воспроизведения: %v\n", err)
    } else {
        fmt.Printf("Воспроизводится: %s - %s\n", track.Artist, track.Title)
        
        go func() {
            p.Wait()
        }()
    }
}

func stopPlayback(p *player.Player) {
    if p.IsPlaying() {
        p.Stop()
        fmt.Println("Воспроизведение остановлено")
    } else {
        fmt.Println("Нет активного воспроизведения")
    }
}
