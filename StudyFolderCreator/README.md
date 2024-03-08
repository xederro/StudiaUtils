Create study structure:

CSV file structure:

Nazwa;godz.Wykładów;godz.Ćwiczeń;godz.Labolatoriów;godz.Projektu;godz.Seminarium;ECTS

-in

path to csv file

-out

path where to create structure

-pre

insert before paths

-n

number of semester

```shell
.\SSC.exe -n=5 -out="..\Notatki" -in=".\Semestr 5.csv" -pre="Notatki/"
```
