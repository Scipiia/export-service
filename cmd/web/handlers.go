package main

import (
	"bufio"
	"dem3_demo_v2/pkg/models"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//func (app *application) home(w http.ResponseWriter, r *http.Request) {
//	if r.URL.Path != "/" {
//		app.notFound(w)
//		return
//	}
//
//	s, err := app.snippets.Latest()
//	if err != nil {
//		app.serverError(w, err)
//		return
//	}
//
//	//files := []string{
//	//	"./ui/html/home.page.tmpl",
//	//	"./ui/html/base.layout.tmpl",
//	//	"./ui/html/footer.partial.tmpl",
//	//}
//	//
//	//parseFiles, err := template.ParseFiles(files...)
//	//if err != nil {
//	//	log.Println(err.Error())
//	//	http.Error(w, "Internal Server Error", 500)
//	//	return
//	//}
//	//
//	//err = parseFiles.Execute(w, s)
//	//if err != nil {
//	//	log.Println(err.Error())
//	//	http.Error(w, "Internal Server Error", 500)
//	//}
//
//	//Используем помощника render() для отображения шаблона.
//	app.render(w, r, "home.page.tmpl", &templateData{
//		Snippets: s,
//	})
//}
//
//func (app *application) show(w http.ResponseWriter, r *http.Request) {
//	id, err := strconv.Atoi(r.URL.Query().Get("id"))
//	if err != nil || id < 1 {
//		app.notFound(w)
//		return
//	}
//
//	s, err := app.snippets.Get(id)
//	if err != nil {
//		if errors.Is(err, models.ErrNoRecord) {
//			app.notFound(w)
//		} else {
//			app.serverError(w, err)
//		}
//		return
//	}
//
//	// Используем помощника render() для отображения шаблона.
//	app.render(w, r, "show.page.tmpl", &templateData{
//		ProfData: s,
//	})
//}
//
//func (app *application) create(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		w.Header().Set("Allow", http.MethodPost)
//		app.clientError(w, http.StatusMethodNotAllowed)
//		return
//	}
//
//	// Создаем несколько переменных, содержащих тестовые данные. Мы удалим их позже.
//	title := "История про улитку"
//	content := "Улитка выползла из раковины,\nвытянула рожки,\nи опять подобрала их."
//	expires := "7"
//
//	// Передаем данные в метод SnippetModel.Insert(), получая обратно
//	// ID только что созданной записи в базу данных.
//	id, err := app.snippets.Insert(title, content, expires)
//	if err != nil {
//		app.serverError(w, err)
//		return
//	}
//
//	// Перенаправляем пользователя на соответствующую страницу заметки.
//	http.Redirect(w, r, fmt.Sprintf("/dem/show?id=%d", id), http.StatusSeeOther)
//}

// my logics
func (app *application) ExportFromProfstroi(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	dir := "temp"
	files, err := readDir(dir)
	if err != nil {
		app.logger.Error(err)
	}

	fileNames := ""
	for _, file := range files {
		if strings.Contains(file, ".txt") {
			fileNames = file
		}
	}

	data, err := os.OpenFile(filepath.Join(dir, fileNames), os.O_RDONLY, 0644)
	if err != nil {
		app.logger.Error(err)
	}

	defer data.Close()

	decoder := charmap.Windows1251.NewDecoder()
	reader := decoder.Reader(data)
	scanner := bufio.NewScanner(reader)

	dataStr := &models.ProfData{}

	mapDet := make(map[int][]string)

	n := 0
	for scanner.Scan() {
		n++
		a := scanner.Text()
		dataStr.D2Number += searchPos(a, "Номер_заказа::::")
		dataStr.D2Profstroi += searchPos(a, "Профстрой::::")
		dataStr.D2Object += searchPos(a, "Объект::::")
		dataStr.D2Manager += searchPos(a, "Менеджер::::")
		dataStr.D2Kontragent += searchPos(a, "Контрагент.Мен::::")
		dataStr.D2ID += searchPos(a, "Идентификатор_ID::::")
		dataStr.D2Diler += searchPos(a, "Контрагент::::")
		dataStr.D2City += searchPos(a, "Город::::")
		dataStr.D2Napr += searchPos(a, "Направление_источник::::")
		dataStr.D2SumProjToSk += searchPos(a, "Стоимость_проекта_без_скидок::::")
		dataStr.D2SumSkidka += searchPos(a, "Скидка_общая::::")
		dataStr.D2SumProjWithSkidka += searchPos(a, "Стоимость_проекта_со_скидками::::")
		dataStr.D2SumConstrWithSkidka += searchPos(a, "Стоимость_конструкций_со_скидками::::")
		dataStr.D2SumRabWithSkidka += searchPos(a, "Стоимость_работ_со_скидками::::")
		dataStr.D2Status += searchPos(a, "Статус_проекта::::")
		dataStr.NoteOrder += searchPos(a, "Примечание::::")

		material := searchPos(a, "Материал::::")
		split := strings.Split(material, ";;;;")
		for _, s := range split {
			if s != "" {
				split2 := strings.Split(s, ";;")
				for _, s2 := range split2 {
					mapDet[n] = append(mapDet[n], s2)
				}
			}
		}
	}

	id, err := app.profData.InsertProfData(dataStr.D2Number, dataStr.D2Profstroi, dataStr.D2Object, dataStr.D2Manager, dataStr.D2Kontragent, dataStr.D2ID,
		dataStr.D2Diler, dataStr.D2City, dataStr.D2Napr, dataStr.D2SumProjToSk, dataStr.D2SumSkidka, dataStr.D2SumProjWithSkidka,
		dataStr.D2SumConstrWithSkidka, dataStr.D2SumRabWithSkidka, dataStr.D2Status, dataStr.NoteOrder)
	if err != nil {
		app.logger.Error(err)
	}

	for _, i2 := range mapDet {
		dataStr.Details.Size = i2[0]
		dataStr.Details.Name = i2[1]
		dataStr.Details.Count = i2[2]
		dataStr.Details.Allowances = i2[4]
		dataStr.Details.Color = i2[5]
		dataStr.Details.Height = i2[6]

		app.profData.InsertDemMaterial(id, dataStr.Details.Size, dataStr.Details.Name, dataStr.Details.Count, dataStr.Details.Allowances,
			dataStr.Details.Color, dataStr.Details.Height)
		if err != nil {
			app.logger.Error(err)
		}
	}

	app.render(w, r, "export.page.tmpl", &templateData{
		ProfData: dataStr,
	})

	if err := scanner.Err(); err != nil {
		app.logger.Error(err)
	}

	//Таймер
	duration := time.Since(start)
	fmt.Println(duration.Seconds())

	//app.render(w, r, "export.page.tmpl", nil)
	//http.Redirect(w, r, fmt.Sprintf("/dem/show?id=%d", id), http.StatusSeeOther)
}

func readDir(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	os.IsNotExist(err)
	if err != nil {
		log.Fatal("Readdir error: ", err)
	}

	var fileSlice []string

	for _, file := range files {
		if file.IsDir() {
			p, err := readDir(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("dirwalk %s: %w", filepath.Join(dir, file.Name()), err)
			}
			filePiz = append(fileSlice, p...)
			continue
		}
		filePiz = append(fileSlice, file.Name())
	}
	return fileSlice, nil
}

func searchPos(str, str1 string) string {
	ret := ""
	if strings.Contains(str, str1) {
		_, after, _ := strings.Cut(str, str1)
		ret += after
	}
	return ret
}
