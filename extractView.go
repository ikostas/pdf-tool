package main

import (
  "fmt"
  "path/filepath"
  tk "modernc.org/tk9.0"
  "github.com/pdfcpu/pdfcpu/pkg/api"
)

// extractPdf extracts selected pages to a new file
func extractPdf() {
  var buttonsArr buttonDefs
  var radioArr radioDefs
  var buttons []*tk.TButtonWidget
  var chooseFile, home, actExtract func()
  var getPagesFromUI func() ([]string, error)
  var runExtract func(string, []string, report)
  var inputFile, dir string
  var msgLabel *tk.TLabelWidget
  var fileRow *tk.FrameWidget

  radioChoice := tk.Variable("")
  extract := tk.App.Frame()
  t := title{
    wmTitle: "PDF Tool -- extract",
    title: "Extract PDF pages",
    tipString: "Choose a file to extract pages from, 1 file for each page will be created",
    isMainMenu: false,
    win: extract,
    msgLabel: &msgLabel,
  }

  r := makeTitle(t)
  inputRow, msgRow, btnRow := r.ir, r.mr, r.br

  entryRow := inputRow.Frame() // a field to enter page numbers
  radioRow := inputRow.Frame() // a radio button to choose mode
  entryLine := entryRow.TEntry(tk.Placeholder("Example: -2, 3, 4-6, 7-"), tk.Textvariable(""))

  chooseFile = func() {
    if !chooseOneFile(&inputFile, inputRow, &fileRow) { return }
    buttons[1].Configure(tk.State("normal"))
  }
  home = func() {
    goHome(buttonsArr, &extract)
  }

  actExtract = func() {
    if inputFile == "" {
      packMsg(msgRow, &msgLabel, "Don't press <Alt-X> unless you select a file to extract pages from")
      return
    }
    dir = filepath.Dir(inputFile)
    pages, errUI := getPagesFromUI()
    if errUI != nil {
      packMsg(msgRow, &msgLabel, "Input error: " + errUI.Error())
      return
    }
    r := report{
      msgRow: msgRow,
      msgLabel: &msgLabel,
      msgSuccess: "The folder with extracted pages: ",
      msgFail: "Can't extract pages from the file: ",
      result: dir,
    }
    runExtract(inputFile, pages, r)
  }

  getPagesFromUI = func() ([]string, error) {
    radioChoiceVal := radioChoice.Get()
    switch radioChoiceVal  {
    case "odd", "even":
      return []string{radioChoiceVal}, nil
    case "spec":
      return createPagesArr(entryLine.Textvariable())
    default:
      return nil, fmt.Errorf("choose one of the options with the radio buttons")
    }
  }

  runExtract = func(f string, pages []string, r report) {
    go func() {
      r.err = api.ExtractPagesFile(f, r.result, pages, nil)
      tk.PostEvent(func() { reportRes(r) }, false)
    }()
  }

  buttonsArr = buttonDefs{
    {"Choose a file", "icons/icons8-add-file-24.png", "left", 0, chooseFile, "<Alt-c>"},
    {"Extract pages", "icons/icons8-tweezers-24.png", "left", 1, actExtract, "<Alt-x>"},
    {"Home", "icons/icons8-home-24.png", "left", 0, home, "<Alt-h>"},
  }
  radioArr = radioDefs{
    {"Extract odd pages", "odd"},
    {"Extract even pages", "even"},
    {"Extract pages specified below", "spec"},
  }
  buttons = buttonsArr.createButtons(btnRow)
  buttons[1].Configure(tk.State("disabled"))
  createRadio(radioArr, radioChoice, radioRow)
  createEntry(entryRow, entryLine, "Pages to extract: ")
  packBottomBtns(btnRow)
  buttonsArr.setHotkeys()
}
