# Instruksi Build

Untuk membangun file `main.go` dengan nama kustom, gunakan perintah `go build` dengan opsi `-o` seperti berikut:

Untuk sistem operasi Windows:
```bash
go build -o nama_kustom.exe .\main.go
```
Untuk sistem operasi lain seperti Linux atau macOS:
```bash
go build -o nama_kustom ./main.go
```