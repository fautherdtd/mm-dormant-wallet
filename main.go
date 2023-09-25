package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"runtime"
	"sync"

	"golang.org/x/crypto/ripemd160"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // Используем все доступные ядра процессора

	var wg sync.WaitGroup

	// Создаем случайный seed (32 байта)
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		log.Fatal(err)
	}
	counter := 0
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				counter++
				fmt.Printf("Проверенно адресов: %d\n", counter)
				privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
				if err != nil {
					fmt.Println("Ошибка при генерации приватного ключа:", err)
					return
				}
				publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
				publicKeyHash := hash160(publicKey)

				// Создание адреса
				address := createAddress(publicKeyHash)
				if containsAddress(address) {
					// Запись адреса в файл wallet.txt
					err = appendToFile("wallet.txt", "Приватный ключ (hex): "+hex.EncodeToString(privateKey.D.Bytes()))
					if err != nil {
						fmt.Println("Ошибка при записи в файл:", err)
						return
					}
					err = appendToFile("wallet.txt", "Публичный ключ (hex): "+hex.EncodeToString(publicKey))
					if err != nil {
						fmt.Println("Ошибка при записи в файл:", err)
						return
					}
					err = appendToFile("wallet.txt", "Биткоин-адрес (Legacy P2PKH): "+address)
					if err != nil {
						fmt.Println("Ошибка при записи в файл:", err)
						return
					}

					// Вывод информации
					fmt.Println("Приватный ключ (hex):", hex.EncodeToString(privateKey.D.Bytes()))
					fmt.Println("Публичный ключ (hex):", hex.EncodeToString(publicKey))
					fmt.Println("Биткоин-адрес (Legacy P2PKH):", address)
					break
				}
			}
		}()
	}

	wg.Wait()
}

func containsAddress(address string) bool {
	targetAddress := [105]string{
		"1P1iThxBH542Gmk1kZNXyji4E4iwpvSbrt",
		"1FeexV6bAHb8ybZjqQMjJrcCrHGW9sb6uF",
		"1LdRcdxfbSnmCYYNdeYpUnztiYzVfBEQeC",
		"1AC4fMwgY8j9onSbXEWeH6Zan8QGMSdmtA",
		"12ib7dApVFvg82TXKycWBNpN8kFyiAN1dr",
		"17rm2dvb439dZqyMe2d4D6AQJSgg6yeNRn",
		"1GR9qNz7zgtaW5HwwVpEJWMnGWhsbsieCG",
		"1BZaYtmXka1y3Byi2yvXCDG92Tjz7ecwYj",
		"1F34duy2eeMz5mSrvFepVzy7Y1rBsnAyWC",
		"14YK4mzJGo5NKkNnmVJeuEAQftLt795Gec",
		"1Ki3WTEEqTLPNsN5cGTsMkL2sJ4m5mdCXT",
		"1KbrSKrT3GeEruTuuYYUSQ35JwKbrAWJYm",
		"1ucXXZQSEf4zny2HRwAQKtVpkLPTUKRtt",
		"1CPaziTqeEixPoSFtJxu74uDGbpEAotZom",
		"1P9fAFAsSLRmMu2P7wZ5CXDPRfLSWTy9N8",
		"12ib7dApVFvg82TXKycWBNpN8kFyiAN1dr",
		"12ib7dApVFvg82TXKycWBNpN8kFyiAN1dr",
		"12ib7dApVFvg82TXKycWBNpN8kFyiAN1dr",
		"12ib7dApVFvg82TXKycWBNpN8kFyiAN1dr",
		"12ib7dApVFvg82TXKycWBNpN8kFyiAN1dr",
		"12ib7dApVFvg82TXKycWBNpN8kFyiAN1dr",
		"12ib7dApVFvg82TXKycWBNpN8kFyiAN1dr",
		"12ib7dApVFvg82TXKycWBNpN8kFyiAN1dr",
		"12ib7dApVFvg82TXKycWBNpN8kFyiAN1dr",
		"12ib7dApVFvg82TXKycWBNpN8kFyiAN1dr",
		"1P9fAFAsSLRmMu2P7wZ5CXDPRfLSWTy9N8",
		"1HLvaTs3zR3oev9ya7Pzp3GB9Gqfg6XYJT",
		"167ZWTT8n6s4ya8cGjqNNQjDwDGY31vmHg",
		"18zuLTKQnLjp987LdxuYvjekYnNAvXif2b",
		"198aMn6ZYAczwrE5NvNTUMyJ5qkfy4g3Hi",
		"15Z5YJaaNSxeynvr6uW6jQZLwq3n1Hu6RX",
		"3HCdgNiAsjidcGx4eeMK9AXvXBfrnteunW",
		"1JfXLzQvYPZHNzX4vhH6aoetGDfcPD1YEX",
		"13eEt6myAo1zAC7o7RK5sVxxCNCAgd6ApH",
		"1DzjE3ANaKLasY2n6e5ToJ4CQCXrvDvwsf",
		"15YMdTNT83UJqfpaZfDccy9yBYQFxHxVFt",
		"1FJuzzQFVMbiMGw6JtcXefdD64amy7mSCF",
		"1Ac2JdpQ5c9NeSajdGx6dofxeXkn4S35ft",
		"1AYLzYN7SGu5FQLBTADBzqKm4b6Udt6Bw6",
		"1JxmKkNK1b3p7r8DDPtnNmGeLZDcgPadJb",
		"1LBBmkr9muf7RjjBbzQQvzNQpRRaVEnavs",
		"138EMxwMtKuvCEUtm4qUfT2x344TSReyiT",
		"1DR93bfKVCUJkDvPuxbUAEtzYRaJEnwjNt",
		"1CCqLR8YrUMPFgYZWwLW8FkezbFjfeXD8n",
		"1BVtDi7txPCG2TH5Crd2Rw5MtpivbmoKgB",
		"1BeouDc6jtHpitvPz3gR3LQnBGb7dKRrtC",
		"1ARWCREnmdKyHgNg2c9qih8UzRr4MMQEQS",
		"1DaCQDfStUgkPQXcf53Teeo6LPiKcVMBM9",
		"19z6WynrjHeD5MMv6919BuQRwybuen1sRv",
		"1NQEV6T4avmPqUVTvgsKkeB6yc8qnSWfhR",
		"1NJQZhzYac89fDhQCmb1khdjekKNVYLFMY",
		"12ytiN9oWQTRGb6JjZiaoWMAvF9nPWdGX1",
		"1Btud1pqADgGzgBCZzxzc2b1o1ytk1HYWC",
		"1BXZng4dcXDnYNRXRgHqWjzT5RwxHHBSHo",
		"1BvNwfxEQwZNRmYQ3eno6e976XyxhCsRXj",
		"17spLhCpZVdQXFz2ZL1aP5gRci6RFVNhrD",
		"1Miy5sJZSamDZN6xcJJidp9zYxhSrpDeJm",
		"1Kq6hXXiSpdp9bg9hDDyqm8ZfvgZmzchjn",
		"16aEn4p6hK4FMpLtJGpoQZMZ946sDg1Z6n",
		"1MLiPwYjNACQHREFKwGtkPpWgd8PqpbuQ4",
		"18Hp8j2JMvwtPs1eqNaYEEVvuFpjQJRFVY",
		"16eb495TbiCRbRbZv4WBdaUvNGxUYJ4jed",
		"3N5Nny9doXSfjYh5k9XNdG8baE1kiayx5o",
		"1JCrPqogEKEpM9fuFQV7LpF9e8cgf3YZ8m",
		"1GX7i8jG8DD1mG85BNnz7xybVhSmw84Uii",
		"1N5NqDWiLVqtU8mEzCNEeEbQVHwuGGChJs",
		"124YoiaSaUssbBeP5RukbSN9Evc3UJfwPj",
		"1VeMPNgEtQGurwjW2WsYXQaw4boAX5k6S",
		"18eY9oWL2mkXCL1VVwPme2NMmAVhX6EfyM",
		"1ALXLVNj7yKRU2Yki3K3yQGB5TBPof7jyo",
		"1LwBdypLh3WPawK1WUqGZXgs4V8neHHqb7",
		"15MZvKjqeNz4AVz2QrHumQcRJq2JVHjFUz",
		"35h3wuHCVzuULtL7nRtYq9bzWVRR554QbY",
		"1yAFNheT6MyMddhYXqjW9yYgNh6KiKTWb",
		"1GMFSWQQQhCQyRNQcac9tDKcvqYCuripVs",
		"14mPMrRm6TdjqHZhd7aBUbuWt5MYWReukR",
		"1FvUkW8thcqG6HP7gAvAjcR52fR7CYodBx",
		"1Gn1GzVa88T1X3fdhejyq6jrZs43T24xW6",
		"1PTYXwamXXgQoAhDbmUf98rY2Pg1pYXhin",
		"16oKJMcUZkDbq2tXDr9Fm2HwgBAkJPquyU",
		"193zge613zys7eSkQYQVd9xubovLN8Sr6j",
		"1L9ipUywwErf9EaKCgLqrkoSM5ab3wrjvh",
		"1EgH7EUfgjr8gAK9t1BeHLDC1ijrVvdec3",
		"18jGeboNHt1YpsDFcCeYPKm8qnAe9942BG",
		"19HhmfxGsznL8K7wXjZiFnhqddQucgfZzB",
		"1ArZGb5V24gAgN51FeQknobi6kNyGx739r",
		"1JjMoB212ctAiuDvURyWhs813yY4c75cap",
		"1FDVbVJYKkWPFcJEzCxi99vpKTYxEY3zdj",
		"1Hc4EvgZmWECnETeTL4w4ySz76JxubyRjw",
		"1AGAvShyB22eUxz1DKfBBgGENDSZP8dcq9",
		"1KpwMa1w9DTUCB5asCgUdLRA22hto1Qgqv",
		"17TZNT8CBPzUPDfKTXC25RQHrW6M2q6kRo",
		"1LQaq7LLoyjdfH3vczuusa17WsRokhsRvG",
		"15HiQkbvQMoAzXyKdQbuCKTGDxTswYBUf5",
		"1AenFm1zSRkhtPHwZmP2UuRQbWpakD8cVZ",
		"1NY5KheH3koPcuQrBLXVGq87YbijtXdZXD",
		"13KYdPnzGh5H8exFY3FhUo9Rvvs6kKAcL8",
		"1EUJKGm3FB65rr5W9anAEoWA3m71WpDayZ",
		"1LawddwNrqHySAMK528KLhPC2d9aWt9YMQ",
		"1PEUv3FjSWq88AgNYefeYaEhLWSiMW2vuy",
		"18cKGtwdQHmnDXD6w6AhBhHsmxnK8gsVHf",
		"19DdkMxutkLGY67REFPLu51imfxG9CUJLD",
		"1BrSzBwx2RLuppgGziqgF7oMuneHQVhsNc",
		"17KcBp8g76Ue8pywgjta4q8Ds6wK4bEKp7",
		"17j45BXWrjSDttuurcSQubYLdLescJ7eJH",
	}
	for _, n := range targetAddress {
		if address == n {
			return true
		}
	}
	return false
}

// hash160 вычисляет хэш160 от данных
func hash160(data []byte) []byte {
	sha256 := sha256.New()
	sha256.Write(data)
	hash := sha256.Sum(nil)

	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)
	return ripemd160.Sum(nil)
}

// createAddress создает биткоин-адрес из хэша от публичного ключа
func createAddress(publicKeyHash []byte) string {
	networkBytes := []byte{0x00} // 0x00 для основной сети (MainNet)
	payload := append(networkBytes, publicKeyHash...)

	// Вычисление контрольной суммы (SHA256(SHA256(payload)))
	checksum := doubleSHA256(payload)[:4]

	// Создание окончательного байтового массива
	addressBytes := append(payload, checksum...)

	// Кодирование в base58
	address := base58Encode(addressBytes)

	return address
}

// doubleSHA256 вычисляет SHA-256 хэш от SHA-256 хэша входных данных
func doubleSHA256(data []byte) []byte {
	sha256 := sha256.New()
	sha256.Write(data)
	hash := sha256.Sum(nil)

	sha256.Reset()
	sha256.Write(hash)
	return sha256.Sum(nil)
}

// base58Encode кодирует байты в строку base58
func base58Encode(input []byte) string {
	const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	var result []byte
	x := new(big.Int).SetBytes(input)

	for x.Cmp(big.NewInt(0)) > 0 {
		mod := new(big.Int)
		x.DivMod(x, big.NewInt(58), mod)
		result = append(result, base58Alphabet[mod.Int64()])
	}

	// Reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	// Add '1' characters for leading zeros
	for i := 0; i < len(input) && input[i] == 0; i++ {
		result = append([]byte{'1'}, result...)
	}

	return string(result)
}

func appendToFile(fileName string, data string) error {
	// Открываем файл для добавления данных (или создания, если файл не существует)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Добавляем данные в файл, разделяя новой строкой
	_, err = file.WriteString(data + "\n")
	if err != nil {
		return err
	}
	return nil
}
