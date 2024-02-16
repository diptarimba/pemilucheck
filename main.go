package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

var baseUrl = "https://sirekap-obj-data.kpu.go.id/wilayah/pemilu/ppwp%s.json"
var baseUrlCekData = "https://sirekap-obj-data.kpu.go.id/pemilu/hhcw/ppwp%s.json"
var addPath = ""
var kotaNama = ""

type WilayahPemilu []struct {
	Nama    string `json:"nama"`
	ID      int    `json:"id"`
	Kode    string `json:"kode"`
	Tingkat int    `json:"tingkat"`
}

type DataPemilu struct {
	Chart *struct {
		Num100025 int         `json:"100025"`
		Num100026 int         `json:"100026"`
		Num100027 int         `json:"100027"`
		Null      interface{} `json:"null"`
	} `json:"chart"`
	Images       []string `json:"images"`
	Administrasi *struct {
		SuaraSah        int `json:"suara_sah"`
		SuaraTotal      int `json:"suara_total"`
		PemilihDptJ     int `json:"pemilih_dpt_j"`
		PemilihDptL     int `json:"pemilih_dpt_l"`
		PemilihDptP     int `json:"pemilih_dpt_p"`
		PenggunaDptJ    int `json:"pengguna_dpt_j"`
		PenggunaDptL    int `json:"pengguna_dpt_l"`
		PenggunaDptP    int `json:"pengguna_dpt_p"`
		PenggunaDptbJ   int `json:"pengguna_dptb_j"`
		PenggunaDptbL   int `json:"pengguna_dptb_l"`
		PenggunaDptbP   int `json:"pengguna_dptb_p"`
		SuaraTidakSah   int `json:"suara_tidak_sah"`
		PenggunaTotalJ  int `json:"pengguna_total_j"`
		PenggunaTotalL  int `json:"pengguna_total_l"`
		PenggunaTotalP  int `json:"pengguna_total_p"`
		PenggunaNonDptJ int `json:"pengguna_non_dpt_j"`
		PenggunaNonDptL int `json:"pengguna_non_dpt_l"`
		PenggunaNonDptP int `json:"pengguna_non_dpt_p"`
	} `json:"administrasi"`
	Psu         *interface{} `json:"psu"`
	Ts          string       `json:"ts"`
	StatusSuara bool         `json:"status_suara"`
	StatusAdm   bool         `json:"status_adm"`
}

var headerHttp = map[string]string{
	"Accept":          "application/json, text/plain, */*",
	"Accept-Language": "en-US,en;q=0.9",
	"Cache-Control":   "no-cache",
	"Connection":      "keep-alive",
	"Origin":          "https://pemilu2024.kpu.go.id",
	"Pragma":          "no-cache",
	"Referer":         "https://pemilu2024.kpu.go.id/",
	"Sec-Fetch-Dest":  "empty",
	"Sec-Fetch-Mode":  "cors",
	"Sec-Fetch-Site":  "same-site",
	"User-Agent":      "Mozilla/5.0 (iPad; CPU OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1 Edg/121.0.0.0",
}

func main() {
	fetchData("", nil)
}

func fetchData(id string, kota *string) {

	var dataRes WilayahPemilu
	if id == "" {
		id = "0"
	}
	if addPath == "/0" {
		addPath = ""
	}
	addPath += "/" + id
	url := fmt.Sprintf(baseUrl, addPath)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for key, value := range headerHttp {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &dataRes)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, each := range dataRes {
		kotaNama += "-" + each.Nama
		if each.Tingkat < 5 {
			fetchData(each.Kode, &kotaNama)
		} else {
			checkData(each.Kode, &kotaNama)
			addPath = addPath[:len(addPath)-len(each.Kode)-1]
		}
		kotaNama = kotaNama[:len(kotaNama)-len(each.Nama)-1]
		time.Sleep(500 * time.Millisecond)
	}
	addPath = addPath[:len(addPath)-len(id)-1]
}

func checkData(kode string, kota *string) {
	var dataRes DataPemilu

	addPath += "/" + kode
	url := fmt.Sprintf(baseUrlCekData, addPath)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for key, value := range headerHttp {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(body, &dataRes)
	if err != nil {
		log.Println(err)
		return
	}
	if dataRes.Administrasi == nil {
		log.Println("Data belum tersedia, link : " + url)
		logError("unavailable.txt", "Data belum tersedia", *kota, url)
		return
	}

	if dataRes.Chart == nil {
		log.Println("Data belum tersedia v2, link : " + url)
		logError("unavailablev2.txt", "Data Aneh", *kota, url)
		return
	}
	total := dataRes.Chart.Num100025 + dataRes.Chart.Num100026 + dataRes.Chart.Num100027
	rawSelisih := total - dataRes.Administrasi.SuaraSah
	selisih := int(math.Abs(float64(rawSelisih)))
	dataToSave := fmt.Sprintf("Anis: %d, Prabowo: %d, Ganjar: %d, Total 3 Paslon: %d, Total Sah: %d, Selisih: %d, Wilayah: %s, Keterangan: %s", dataRes.Chart.Num100025, dataRes.Chart.Num100026, dataRes.Chart.Num100027, total, dataRes.Administrasi.SuaraSah, selisih, *kota, "Data Tidak Sesuai")

	if total != dataRes.Administrasi.SuaraSah {
		log.Println("Perhitungan Tidak Sesuai, link : " + url)
		logError("invalid.txt", dataToSave, *kota, url)
		return
	}

	log.Println("Perhitungan Sesuai, link : " + url)
	logError("valid.txt", "Perhitungan Sesuai", *kota, url)
}

func logError(fileName, message, wilayah, url string) {
	url = strings.TrimPrefix(url, "https://sirekap-obj-data.kpu.go.id/pemilu/hhcw/ppwp")
	url = strings.TrimSuffix(url, ".json")
	url = "https://pemilu2024.kpu.go.id/pilpres/hitung-suara" + url
	fileContent := map[string]string{"message": message, "link": url, "wilayah": wilayah}
	fileBytes, err := json.Marshal(fileContent)
	if err != nil {
		log.Println("Error marshalling to JSON:", err)
		return
	}

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(string(fileBytes) + "\n"); err != nil {
		log.Println(err)
	}
}
