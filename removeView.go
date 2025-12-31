package main

import (
  "path/filepath"
  tk "modernc.org/tk9.0"
  "github.com/pdfcpu/pdfcpu/pkg/api"
)

func RemovePdfPages() {
  var buttonsArr ButtonDefs
  var buttons []*tk.TButtonWidget
  var fileToCut, pagesArr []string
  var err error
  var chooseFile, home, actRemove func()
  var output string
  var msgLabel *tk.TLabelWidget
  var fileRow *tk.FrameWidget

  remove := tk.App.Frame()
  t := Title{
    wmTitle: "PDF Tool -- remove pages",
    title: "Remove pages from PDF",
    tipString: "Choose a file to remove pages",
    isMainMenu: false,
    win: remove,
    msgLabel: &msgLabel,
  }

  inputRow, msgRow, btnRow := MakeTitle(t)
  entryRow := inputRow.Frame() // a field to enter page numbers
  entryLine := entryRow.TEntry(tk.Placeholder("Example: -2, 3, 4-6, 7-"), tk.Textvariable(""))

  chooseFile = func() {
    if !ChooseOneFile(&fileToCut, inputRow, &fileRow) { return }
    buttons[1].Configure(tk.State("normal"))
  }

  home = func() {
    GoHome(buttonsArr, &remove)
  }

  actRemove = func() {
    if fileToCut == nil {
      PackMsg(msgRow, &msgLabel, "Don't press <Alt-O> unless you select a file")
      return
    }
    pagesArr, err = CreatePagesArr(entryLine.Textvariable())
    if err != nil {
      PackMsg(msgRow, &msgLabel, "Error parsing pages line: " + err.Error())
      return
    }
    go func() {
      output = fileToCut[0][:len(fileToCut[0])-len(filepath.Ext(fileToCut[0]))] + "_" + "cut" + ".pdf"
      err = api.RemovePagesFile(fileToCut[0], output, pagesArr, nil)
      r := Report{
        msgRow: msgRow,
        msgLabel: &msgLabel,
        msgSuccess: "Cut file: ",
        msgFail: "Can't remove the pages: ",
        result: output,
        err: err,
      }
      tk.PostEvent(func() { reportRes(r) }, false)
    }()
  }

  buttonsArr = ButtonDefs{
    {"Choose a file", "icons/icons8-add-file-24.png", "left", 0, chooseFile, "<Alt-c>"},
    {"Remove pages", "icons/icons8-remove-24.png", "left", 3, actRemove, "<Alt-o>"},
    {"Home", "icons/icons8-home-24.png", "left", 0, home, "<Alt-h>"},
  }

  buttons = buttonsArr.CreateButtons(btnRow)
  buttons[1].Configure(tk.State("disabled"))
  CreateEntry(entryRow, entryLine, "Pages to remove: ")
  PackBottomBtns(btnRow)
  buttonsArr.SetHotkeys()
}
