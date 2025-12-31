package main

import (
  "path/filepath"
  tk "modernc.org/tk9.0"
  "github.com/pdfcpu/pdfcpu/pkg/api"
)

// ExtractPdf extracts selected pages to a new file
func ExtractPdf() {
  var buttonsArr ButtonDefs
  var buttons []*tk.TButtonWidget
  var fileToExtract, pagesArr []string
  var err error
  var chooseFile, home, actExtract func()
  var dir string
  var msgLabel *tk.TLabelWidget
  var fileRow *tk.FrameWidget

  extract := tk.App.Frame()
  t := Title{
    wmTitle: "PDF Tool -- extract",
    title: "Extract PDF pages",
    tipString: "Choose a file to extract pages from, 1 file for each page will be created",
    isMainMenu: false,
    win: extract,
    msgLabel: &msgLabel,
  }

  inputRow, msgRow, btnRow := MakeTitle(t)
  entryRow := inputRow.Frame() // a field to enter page numbers
  entryLine := entryRow.TEntry(tk.Placeholder("Example: -2, 3, 4-6, 7-"), tk.Textvariable(""))

  chooseFile = func() {
    if !ChooseOneFile(&fileToExtract, inputRow, &fileRow) { return }
    buttons[1].Configure(tk.State("normal"))
  }
  home = func() {
    GoHome(buttonsArr, &extract)
  }

  actExtract = func() {
    if fileToExtract == nil {
      PackMsg(msgRow, &msgLabel, "Don't press <Alt-X> unless you select a file to extract pages from")
      return
    }
    pagesArr, err = CreatePagesArr(entryLine.Textvariable())
    if err != nil {
      PackMsg(msgRow, &msgLabel, "Error parsing pages line: " + err.Error())
      return
    }
    dir = filepath.Dir(fileToExtract[0])
    go func() {
      err = api.ExtractPagesFile(fileToExtract[0], dir, pagesArr, nil)
      r := Report{
        msgRow: msgRow,
        msgLabel: &msgLabel,
        msgSuccess: "The file with extracted pages: ",
        msgFail: "Can't extract pages from the file: ",
        result: dir,
        err: err,
      }
      tk.PostEvent(func() { reportRes(r) }, false)
    }()
  }

  buttonsArr = ButtonDefs{
    {"Choose a file", "icons/icons8-add-file-24.png", "left", 0, chooseFile, "<Alt-c>"},
    {"Extract pages", "icons/icons8-tweezers-24.png", "left", 1, actExtract, "<Alt-x>"},
    {"Home", "icons/icons8-home-24.png", "left", 0, home, "<Alt-h>"},
  }
  buttons = buttonsArr.CreateButtons(btnRow)
  buttons[1].Configure(tk.State("disabled"))
  CreateEntry(entryRow, entryLine, "Pages to extract: ")
  PackBottomBtns(btnRow)
  buttonsArr.SetHotkeys()
}
