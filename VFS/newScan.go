package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
)

// ParticularFile тип структуры созданный для удобства отображения конкретного файла
// Под которым собрана информация: о имени, размере, и о том является ли экземпляр
// файлом или папкой
type ParticularFile struct {
	Name  string //Имя файла
	Size  int64  //Размер файла
	IsDir bool   //Папка ли конкретная сущность
}

// listRes карта содержащая все экземпляры ParticularFile найденные во время работы, для
// дальнейшего вывода в терминал
var listRes map[string]ParticularFile

// folderSize() рекурсивно проходит все вложенные папки и файлы и подсчитывает их вес
func folderSize(rootPath string) (int64, error) {
	var ValSize int64
	files, err := OpenPath(rootPath)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	for _, file := range files {
		//Если файл то пробавляем вес файла к весу папки
		if !file.IsDir() {
			ValSize += file.Size()
		} else {
			//Если фиректория то делаем рекурсивный вход
			folder_size, err := folderSize(path.Join(rootPath, file.Name()))
			if err != nil {
				return 0, err
			} else {
				//Если рекурсия закончена прибавляем размер к весу папки
				ValSize += folder_size
			}
		}

	}
	return ValSize, nil
}

// writeRes() выводит результаты работы программы и представляет их в удобном виде для понимания
func writeRes(d map[string]ParticularFile) error {
	for _, fileI := range d {
		//Определяем тип
		directORfile := "File"
		if fileI.IsDir {
			directORfile = "Directory"
		}
		//Определяем форматирование относительно веса файла
		if fileI.Size/1024/1024 != 0 {
			fmt.Println(directORfile, " - ", fileI.Name, " ", fileI.Size/1024/1024, " Mb")
		} else if fileI.Size/1024 != 0 {
			fmt.Println(directORfile, " - ", fileI.Name, " ", fileI.Size/1024, " Kb")
		} else {
			fmt.Println(directORfile, " - ", fileI.Name, " ", fileI.Size, " ba")
		}
	}
	return nil
}

// OpenPath() открывает и собирает данные в указанной директории
// после возвращает полученный []fs.FileInfo.
func OpenPath(rootPath string) ([]fs.FileInfo, error) {
	dir, err := os.Open(path.Join(rootPath))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer dir.Close()

	// Получаем список файлов и папок
	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return files, err
}

// ScanPath() получет список файлов и записывает их в зависимости от типа (файл или директория)
func ScanPath(rootPath string) error {
	//Получаем список []fs.FileInfo с файлами в указанной директории
	files, err := OpenPath(rootPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//Для каждого найденного файла
	for _, file := range files {
		//Проверяем папка ли
		if !file.IsDir() {
			//Если файл просто записываем в listRes
			listRes[file.Name()] = ParticularFile{file.Name(), file.Size(), false}
		} else {
			//Если папка находим размер
			folder_size, err := folderSize(path.Join(rootPath, file.Name()))
			if err != nil {
				return err
			} else {
				//Записываем в listRes с найденным размером
				listRes[file.Name()] = ParticularFile{file.Name(), folder_size, true}
			}
		}
	}
	return nil

}

// StartScan()  Эта функция начинает процесс
// сбора данных далее выводит результаты
func StartScan(pathScan string) error {

	listRes = make(map[string]ParticularFile)

	err := ScanPath(pathScan)
	if err != nil {
		return err
	}

	err = writeRes(listRes)
	if err != nil {
		return err
	}

	return nil

}
