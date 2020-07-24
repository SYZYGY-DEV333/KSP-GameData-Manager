package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"encoding/json"
	"path/filepath"
	
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/plus3it/gorecurcopy"
	"github.com/gen2brain/dlgs"
)

type GameDataInfo struct {
	Name string
	Description	string
	KSPVersion string
}

func getGameDataInfo(rpath string) GameDataInfo {
	var g GameDataInfo
	
	content, err := ioutil.ReadFile(rpath + "/GameDataInfo.json")
	if err != nil {
        log.Fatal(err)
    }
	
	err = json.Unmarshal(content, &g)
	if err != nil {
        log.Fatal(err)
    }
    
    return g
}

func putGameDataInfo(rpath string, info GameDataInfo) []byte {
	jsondata, err := json.Marshal(info)
	if err != nil {
        log.Fatal(err)
    }
    
    err = ioutil.WriteFile(rpath + "/GameDataInfo.json", jsondata, 0666)
    if err != nil {
        log.Fatal(err)
    }
    
    return jsondata
}

func copyToArchive(gamedataPath string, info GameDataInfo) []byte {
    log.Println("GameData Path: " + gamedataPath)
    
    destinationPath, _ := filepath.Abs("GameArchive/" + info.Name)
    log.Println("Archive Path: " + destinationPath)
    
    if _, err1 := os.Stat(destinationPath); !os.IsNotExist(err1) {
		yes, _ := dlgs.Question("Copy To Archive", "A file with this name already exists. Do you wish to overwrite it?", true)
		if yes == true {
			err2 := os.RemoveAll(destinationPath)
			if err2 != nil {
				log.Fatal(err2)
			}
		}
	}
	
	err3 := gorecurcopy.CopyDirectory(gamedataPath, destinationPath + "/GameData")
	if err3 != nil {
		log.Fatal(err3)
	}
	
	data := putGameDataInfo(destinationPath, info)
	return data
}

func copyToArchiveDialog() {
	var (
		name string
		description string
		kspversion string
		path string
	)
	
	w2 := fyne.CurrentApp().NewWindow("Copy GameData Folder to Archive")
	
	label := widget.NewLabel("Put desired info below, and pick the folder to archive:")
	
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter name for new GameData archive:")
	
	descriptionEntry := widget.NewEntry()
	descriptionEntry.SetPlaceHolder("Enter a useful description:")
	
	kspversionEntry := widget.NewEntry()
	kspversionEntry.SetPlaceHolder("Enter KSP Version:")
	
	selectedFolder := widget.NewLabel("Selected Folder: ")
	
	fileSelectionButton := widget.NewButton("Select GameData Folder", func() {
		p, _, err := dlgs.File("Select GameData Folder", "", true)
		if err != nil {
			log.Fatal(err)
		}
		path = p
		selectedFolder.SetText("Selected Folder: " + path)
	})
	
	okbutton := widget.NewButton("Ok", func() {
		name = nameEntry.Text
		description = descriptionEntry.Text
		kspversion = kspversionEntry.Text
		
		info := GameDataInfo{
			Name: name,
			Description: description,
			KSPVersion: kspversion}
		
		log.Println(path)
		
		copyToArchive(path, info)
		dlgs.Info("GameData Manager", "Copy Complete")
		
		w2.Close()
	})
	
	cancelbutton := widget.NewButton("Cancel", func() {
		yes, _ := dlgs.Question("Cancel", "Are you sure you want to cancel?", true)
		
		if yes == true {
			w2.Close()
		}
	})
	
	w2.SetContent(widget.NewVBox(
		label, nameEntry, descriptionEntry, kspversionEntry, selectedFolder,
		fileSelectionButton, cancelbutton, okbutton))
	
	w2.Show()
}

func useArchivedGameDataDialog() {
	w2 := fyne.CurrentApp().NewWindow("Use Archived GameData")
	label := widget.NewLabel("Select an existing GameData archive to use, \nreplacing the one currently in use.")
 
	var gameDataList []string
	files, err := ioutil.ReadDir("GameArchive/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		gameDataList = append(gameDataList, f.Name())
	}
	
	var options []string
	for _, d := range gameDataList {
		info := getGameDataInfo("GameArchive/" + d)
		newOption := info.Name + " ― (" + info.Description + ") ― KSP Version: " + info.KSPVersion
		options = append(options, newOption)
	}
	
	selection := widget.NewRadio(options, func(string) { })
	
	okbutton := widget.NewButton("Ok", func() {
		yes, _ := dlgs.Question("Ok", "Are you sure you want to continue?", true)
		if yes == true {
			log.Println(strings.Split(selection.Selected, " ― (")[0])
			
			selected := strings.Split(selection.Selected, " ― (")[0]
			os.RemoveAll("GameData")
			
			selectedPath := "GameArchive/" + selected + "/GameData"
			err := gorecurcopy.CopyDirectory(selectedPath, "GameData")
			if err != nil {
				log.Fatal(err)
			}
			
			dlgs.Info("GameData Manager", "Copy Complete")
			w2.Close()
		}
	})
	
	cancelbutton := widget.NewButton("Cancel", func() {
		yes, _ := dlgs.Question("Cancel", "Are you sure you want to cancel?", true)
		if yes == true {
			w2.Close()
		}
	})
	
	w2.SetContent(widget.NewVBox(
		label, selection, cancelbutton, okbutton))
	w2.Show()
}

func deleteArchivedGameDataDialog() {
	w2 := fyne.CurrentApp().NewWindow("Delete Archived GameData")
	label := widget.NewLabel("Select an existing GameData archive to delete.")
 
	var gameDataList []string
	files, err := ioutil.ReadDir("GameArchive/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		gameDataList = append(gameDataList, f.Name())
	}
	
	var options []string
	for _, d := range gameDataList {
		info := getGameDataInfo("GameArchive/" + d)
		newOption := info.Name + " ― (" + info.Description + ") ― KSP Version: " + info.KSPVersion
		options = append(options, newOption)
	}
	
	selection := widget.NewRadio(options, func(string) { })
	
	okbutton := widget.NewButton("Ok", func() {
		yes, _ := dlgs.Question("Ok", "Are you sure you want to delete this archive?", true)
		if yes == true {
			log.Println(strings.Split(selection.Selected, " ― (")[0])
			
			selected := strings.Split(selection.Selected, " ― (")[0]
			
			selectedPath := "GameArchive/" + selected
			os.RemoveAll(selectedPath)
			
			dlgs.Info("GameData Manager", "Deletion Complete")
			w2.Close()
		}
	})
	
	cancelbutton := widget.NewButton("Cancel", func() {
		w2.Close()
	})
	
	w2.SetContent(widget.NewVBox(
		label, selection, cancelbutton, okbutton))
	w2.Show()
}

func cloneGameDataDialog() {
	w2 := fyne.CurrentApp().NewWindow("Clone Archived GameData")
	label := widget.NewLabel("Select an existing GameData archive to clone.")
	
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter name for new GameData archive:")
	
	descriptionEntry := widget.NewEntry()
	descriptionEntry.SetPlaceHolder("Enter a useful description:")
	
	kspversionEntry := widget.NewEntry()
	kspversionEntry.SetPlaceHolder("Enter KSP Version:")
 
	var gameDataList []string
	files, err := ioutil.ReadDir("GameArchive/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		gameDataList = append(gameDataList, f.Name())
	}
	
	var options []string
	for _, d := range gameDataList {
		info := getGameDataInfo("GameArchive/" + d)
		newOption := info.Name + " ― (" + info.Description + ") ― KSP Version: " + info.KSPVersion
		options = append(options, newOption)
	}
	
	selection := widget.NewRadio(options, func(string) { })
	
	okbutton := widget.NewButton("Ok", func() {
		yes, _ := dlgs.Question("Ok", "Are you sure you want to clone this archive?", true)
		if yes == true {
			selected := strings.Split(selection.Selected, " ― (")[0]
			selectedPath := "GameArchive/" + selected
			
			info := GameDataInfo{
				Name: nameEntry.Text,
				Description: descriptionEntry.Text,
				KSPVersion: kspversionEntry.Text}
			
			copyToArchive(selectedPath, info)
			
			dlgs.Info("GameData Manager", "Cloning Complete")
			w2.Close()
		}
	})
	
	cancelbutton := widget.NewButton("Cancel", func() {
		yes, _ := dlgs.Question("Cancel", "Are you sure you want to cancel?", true)
		if yes == true {
			w2.Close()
		}
	})
	
	w2.SetContent(widget.NewVBox(
		label, selection, nameEntry, descriptionEntry, kspversionEntry, cancelbutton, okbutton))
	w2.Show()
}

func editArchivedGameDataInfoDialog2(path string) {
	currentinfo := getGameDataInfo(path)
	
	w2 := fyne.CurrentApp().NewWindow("Edit Archived GameData Info")
	
	label := widget.NewLabel("Editing archive data for: " + currentinfo.Name)
	
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Name: (" + currentinfo.Name + ")")
	
	descriptionEntry := widget.NewEntry()
	descriptionEntry.SetPlaceHolder("Description: (" + currentinfo.Description + ")")
	
	kspversionEntry := widget.NewEntry()
	kspversionEntry.SetPlaceHolder("KSP Version: (" + currentinfo.KSPVersion + ")")
	
	okbutton := widget.NewButton("Ok", func() {
		info := GameDataInfo{
			Name: nameEntry.Text,
			Description: descriptionEntry.Text,
			KSPVersion: kspversionEntry.Text}
		
		putGameDataInfo(path, info)
		
		w2.Close()
	})
	
	cancelbutton := widget.NewButton("Cancel", func() {
		yes, _ := dlgs.Question("Cancel", "Are you sure you want to cancel?", true)
		
		if yes == true {
			w2.Close()
		}
	})
	
	w2.SetContent(widget.NewVBox(
		label, nameEntry, descriptionEntry, kspversionEntry, cancelbutton, okbutton))
	w2.Show()
}

func editArchivedGameDataInfoDialog1() {
	w2 := fyne.CurrentApp().NewWindow("Edit Archived GameData")
	label := widget.NewLabel("Select an existing GameData archive to edit.")
 
	var gameDataList []string
	files, err := ioutil.ReadDir("GameArchive/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		gameDataList = append(gameDataList, f.Name())
	}
	
	var options []string
	for _, d := range gameDataList {
		info := getGameDataInfo("GameArchive/" + d)
		newOption := info.Name + " ― (" + info.Description + ") ― KSP Version: " + info.KSPVersion
		options = append(options, newOption)
	}
	
	selection := widget.NewRadio(options, func(string) { })
	
	okbutton := widget.NewButton("Ok", func() {		
		selected := strings.Split(selection.Selected, " ― (")[0]
		
		selectedPath := "GameArchive/" + selected
		w2.Close()
		editArchivedGameDataInfoDialog2(selectedPath)
	})
	
	cancelbutton := widget.NewButton("Cancel", func() {
		w2.Close()
	})
	
	w2.SetContent(widget.NewVBox(
		label, selection, cancelbutton, okbutton))
	w2.Show()
}

func main() {
	a := app.New()
	w := a.NewWindow("KSP GameData Manager")
	
	w.SetContent(widget.NewVBox(
		widget.NewButton("Clone GameData", func() {
			cloneGameDataDialog()}),
		widget.NewButton("Edit Archived GameData Info", func() {
			editArchivedGameDataInfoDialog1()}),
		widget.NewButton("Use Archived GameData", func() {
			useArchivedGameDataDialog()}),
		widget.NewButton("Copy GameData Folder to Archive", func() {
			copyToArchiveDialog()}),
		widget.NewButton("Delete Archived GameData", func() {
			deleteArchivedGameDataDialog()}),
		widget.NewButton("Quit", func() {
			a.Quit()})))
	
	w.ShowAndRun()
}
