package main

import (
	"flag"
	"log"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	_ = iota
	W
	C
	L
	P
	S
)

var (
	num int
	in  string
	out string
	pre string
)

func main() {
	inFlag := flag.String("in", ".", "path to csv file")
	outFlag := flag.String("out", ".", "path where to create structure")
	preFlag := flag.String("pre", "", "insert before paths")
	nFlag := flag.Int("n", 1, "number of semester")
	flag.Parse()

	num = *nFlag
	in = *inFlag
	out = *outFlag
	pre = *preFlag

	pd := ParseCSV()
	CreateFolderStructure(pd)
}

func ParseCSV() *[][]string {
	data, err := os.ReadFile(in)
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}

	var parsedData [][]string

	s := strings.Split(string(data), "\n")
	for _, l := range s {
		p := strings.Split(l, ";")
		parsedData = append(parsedData, p)
	}

	return &parsedData
}

func CreateFolderStructure(parsedData *[][]string) {
	name := "Semestr " + strconv.Itoa(num)
	err := os.Mkdir(out+"/"+name, 0750)
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}

	created, err := os.Create(out + "/" + name + "/" + name + ".md")
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}

	_, err = created.Write([]byte("# " + name + "\n"))
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}

	for k, p := range *parsedData {
		n := strings.Trim(p[0], "\"")
		sub := name + "/" + n

		err = os.Mkdir(out+"/"+sub, 0750)
		if err != nil {
			log.Fatalln("An Error Occurred:", err)
		}

		CreateSubDirectory(sub, p)
		escape, _ := url.JoinPath("/" + pre + "/" + sub + "/" + n + ".md")
		_, err = created.Write([]byte(
			strconv.Itoa(k+1) + ". [" + n + "](" + escape + ")\n",
		))
		if err != nil {
			log.Fatalln("An Error Occurred:", err)
		}
	}
}

func CreateSubDirectory(where string, what []string) {
	n := strings.Trim(what[0], "\"")
	created, err := os.Create(out + "/" + where + "/" + n + ".md")
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}

	_, err = created.Write([]byte("---\nsemestr: " + strconv.Itoa(num) + "\nocena: \nects: " + what[5] + "\ntyp: 'GK'\n---\n\n# Kurs:\n"))
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}

	for k, v := range what {
		if k == 0 {
			continue
		} else if k == len(what)-1 {
			break
		}

		if v != "" {
			atoi, err := strconv.Atoi(v)
			if err != nil {
				log.Fatalln("An Error Occurred:", err)
			}

			p := CreateLectures(where, int(math.Floor(7.5*float64(atoi))), k)

			_, err = created.Write([]byte(p))
			if err != nil {
				log.Fatalln("An Error Occurred:", err)
			}
		}
	}
}

func CreateLectures(where string, n int, kind int) string {
	var (
		name   string
		single string
		short  string
	)
	switch kind {
	case W:
		name = "Wykłady"
		single = "Wykład"
		short = "W"
		break
	case C:
		name = "Ćwiczenia"
		single = "Ćwiczenie"
		short = "C"
		break
	case L:
		name = "Labolatoria"
		single = "Labolatorium"
		short = "L"
		break
	case P:
		name = "Projekt"
		single = "Projekt"
		short = "P"
		break
	case S:
		name = "Seminaria"
		single = "Seminarium"
		short = "S"
		break
	default:
		log.Fatalln("An Error Occurred: bad type of lecture")
	}

	folder := where + "/" + name

	err := os.Mkdir(out+"/"+folder, 0750)
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}

	created, err := os.Create(out + "/" + where + "/" + name + "/" + name + ".md")
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}

	_, err = created.Write([]byte("---\nsemestr: " + strconv.Itoa(num) + "\nocena: \ntyp: '" + short + "'\n---\n\n# " + name + ":\n"))
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}

	for i := 1; i <= n; i++ {
		err := os.Mkdir(out+"/"+folder+"/"+single+" "+strconv.Itoa(i), 0750)
		if err != nil {
			log.Fatalln("An Error Occurred:", err)
		}

		path := folder + "/" + single + " " + strconv.Itoa(i) + "/" + single + " " + strconv.Itoa(i) + ".md"

		_, err = os.Create(out + "/" + path)
		if err != nil {
			log.Fatalln("An Error Occurred:", err)
		}

		escape, _ := url.JoinPath("/" + pre + "/" + path)
		_, err = created.Write([]byte(
			strconv.Itoa(i) + ". [" + single + " " + strconv.Itoa(i) + "](" + escape + ")\n",
		))
		if err != nil {
			log.Fatalln("An Error Occurred:", err)
		}
	}

	escape, _ := url.JoinPath("/" + pre + "/" + where + "/" + name + "/" + name + ".md")
	return "# [" + name + " ](" + escape + ")\n"
}
