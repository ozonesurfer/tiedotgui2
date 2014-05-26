// TiedotGUI1
package main

import (
	"encoding/json"
	"fmt"
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"math/rand"
	"os"
	//	"strings"
	"time"
)

//var firstLE, lastLE, messageLE *walk.LineEdit
var mw *MyMainWindow

type Name struct {
	FirstName string `json: "firstname"`
	LastName  string `json: "lastname"`
}

const COLLECTION = "names"

type MyMainWindow struct {
	*walk.MainWindow
	model                      *EnvModel
	lb                         *walk.ListBox
	firstLE, lastLE, messageLE *walk.LineEdit
	//	te    *walk.TextEdit
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	//	InitDatabase()
	//	var mainWindow *walk.MainWindow
	mw = &MyMainWindow{model: NewEnvModel()}
	//	fmt.Println("Hello World!")
	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		MinSize:  Size{600, 400},
		Size:     Size{800, 600},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{Margins: Margins{Top: 30, Bottom: 30}},
				Children: []Widget{
					Label{
						Text: "First name:",
						Row:  0,
					},
					LineEdit{
						AssignTo: &mw.firstLE,
						ReadOnly: false,
						Text:     "",
						Row:      0,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "Last name:",
						Row:  0,
					},
					LineEdit{
						AssignTo: &mw.lastLE,
						ReadOnly: false,
						Text:     "",
						Row:      0,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text:      "Add Name",
						OnClicked: BtnAddName,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					LineEdit{
						AssignTo: &mw.messageLE,
						ReadOnly: true,
						Text:     "no errors",
					},
				},
			},
			VSplitter{
				Children: []Widget{
					ListBox{
						AssignTo: &mw.lb,
						Model:    mw.model,
						//		OnCurrentIndexChanged: mw.lb_CurrentIndexChanged,
						//		OnItemActivated:       mw.lb_ItemActivated,
					},
				},
			},
		},
	}.Create()); err != nil {
		fmt.Println("Window create error:", err)
		os.Exit(1)
	}
	mw.MainWindow.Run()
}

func GetDb() *db.DB {
	myDb, err := db.OpenDB("C:\\tmp\\tiedotgui2")
	if err != nil {
		panic(err)
	}
	return myDb
}

func InitDatabase() {
	database := GetDb()
	defer database.Close()
	database.Drop(COLLECTION)
	database.Create(COLLECTION, 1)
	collection := database.Use(COLLECTION)
	collection.Index([]string{"firstname"})
	collection.Index([]string{"lastname"})
}
func BtnAddName() {
	first := mw.firstLE.Text()
	last := mw.lastLE.Text()
	name := Name{first, last}
	fmt.Println("name =", name)
	database := GetDb()
	defer database.Close()
	collection := database.Use(COLLECTION)
	message := "no errors"
	// check to see if the name is already in the database
	var query interface{}
	result := make(map[uint64]struct{})
	queryString := `{"n":
	[{"eq": "` + first + `", "in": ["firstname"]},
	 {"eq": "` + last + `", "in": ["lastname"]}]
	}`
	json.Unmarshal([]byte(queryString), &query)
	db.EvalQuery(query, collection, &result)
	if len(result) > 0 {
		message = "That name is already on record"
	} else {
		// add the name if it is not already stored
		nameForDb := map[string]interface{}{"firstname": first, "lastname": last}
		key, _ := collection.Insert(nameForDb)
		message = "The name was safely added"

		str := name.FirstName + " " + name.LastName
		item := EnvItem{name: key, value: str}
		mw.model.items = append(mw.model.items, item)
		mw.model.PublishItemChanged(len(mw.model.items) - 1)
	}
	mw.messageLE.SetText(message)
	mw.firstLE.SetText("")
	mw.lastLE.SetText("")

}

type EnvItem struct {
	//	name  string
	name  uint64
	value string
}

type EnvModel struct {
	walk.ListModelBase
	items []EnvItem
}

func NewEnvModel() *EnvModel {
	/*	env := os.Environ()

		m := &EnvModel{items: make([]EnvItem, len(env))}

		for i, e := range env {
			j := strings.Index(e, "=")
			if j == 0 {
				continue
			}

			name := e[0:j]
			value := strings.Replace(e[j+1:], ";", "\r\n", -1)

			m.items[i] = EnvItem{name, value}
		}
	*/
	//	var m *EnvModel
	database := GetDb()
	defer database.Close()
	collection := database.Use(COLLECTION)
	var query interface{}
	result := make(map[uint64]struct{})
	e := json.Unmarshal([]byte(`"all"`), &query)
	if e != nil {
		panic(e)
	}
	err := db.EvalQuery(query, collection, &result)
	if err != nil {
		panic(err)
	}
	//	m := &EnvModel{items: make([]EnvItem, len(result))}
	m := &EnvModel{items: []EnvItem{}}
	fmt.Println("Result contains", len(result), "items")
	for key := range result {
		var value Name
		collection.Read(key, &value)
		str := value.FirstName + " " + value.LastName
		fmt.Println("found", str)
		myitem := EnvItem{name: key, value: str}
		m.items = append(m.items, myitem)
	}
	return m
}

func (m *EnvModel) ItemCount() int {
	return len(m.items)
}

func (m *EnvModel) Value(index int) interface{} {
	return m.items[index].value
}
