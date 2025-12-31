package main

import (
  "path/filepath"
  "strconv"
  tk "modernc.org/tk9.0"
  "github.com/pdfcpu/pdfcpu/pkg/api"
)

// RotatePdf rotates all pages in the file for a given angle
func RotatePdf() {
  var buttonsArr ButtonDefs
  var radioArr1, radioArr2 RadioDefs
  var buttons []*tk.TButtonWidget
  var fileToRotate, pagesArr []string
  var err error
  var angle int
  var chooseFile, home, actRotate func()
  var output string
  var msgLabel *tk.TLabelWidget
  var fileRow *tk.FrameWidget

  degreeChoice := tk.Variable("")
  pagesChoice := tk.Variable("")

  rotate := tk.App.Frame()
  t := Title{
    wmTitle: "PDF Tool -- rotate",
    title: "Rotate PDF",
    tipString: "Choose a file to rotate and an angle",
    isMainMenu: false,
    win: rotate,
    msgLabel: &msgLabel,
  }

  inputRow, msgRow, btnRow := MakeTitle(t)
  radioRow1 := inputRow.Frame() // degrees to rotate
  radioRow2 := inputRow.Frame() // all pages or selected only
  entryRow := inputRow.Frame() // a field to enter page numbers
  entryLine := entryRow.TEntry(tk.Placeholder("Example: -2, 3, 4-6, 7-"), tk.Textvariable(""))

  chooseFile = func() {
    if !ChooseOneFile(&fileToRotate, inputRow, &fileRow) { return }
    buttons[1].Configure(tk.State("normal"))
  }
  home = func() {
    GoHome(buttonsArr, &rotate)
  }
  actRotate = func() {
    if fileToRotate == nil {
      PackMsg(msgRow, &msgLabel, "Don't press <Alt-R> unless you select a file to rotate")
      return
    }
    degreeChoiceVal := degreeChoice.Get()
    pagesChoiceVal := pagesChoice.Get()
    if pagesChoiceVal == "" {
      PackMsg(msgRow, &msgLabel, "Select mode: all pages or selected pages")
      return
    } else if pagesChoiceVal == "selected" {
      pagesArr, err = CreatePagesArr(entryLine.Textvariable())
      if err != nil {
        PackMsg(msgRow, &msgLabel, "Error parsing pages line: " + err.Error())
        return
      }
    } // there's one other case, pageChoice == 'all', then we need pagesToRotate = nil, which is already true

    switch degreeChoiceVal {
    case "90", "180", "270":
      // I don't check err here on purpose
      angle, err = strconv.Atoi(degreeChoiceVal)
      go func() {
        output = fileToRotate[0][:len(fileToRotate[0])-len(filepath.Ext(fileToRotate[0]))] + "_" + degreeChoiceVal + ".pdf"
        err = api.RotateFile(fileToRotate[0], output, angle, pagesArr, nil)
        r := Report{
          msgRow: msgRow,
          msgLabel: &msgLabel,
          msgSuccess: "Rotated file: ",
          msgFail: "Can't rotate the file: ",
          result: output,
          err: err,
        }
        tk.PostEvent(func() { reportRes(r) }, false)
      }()
    default:
      PackMsg(msgRow, &msgLabel, "Choose one of the options with the radio buttons.")
      return
    }
  }

  buttonsArr = ButtonDefs{
    {"Choose a file", "icons/icons8-add-file-24.png", "left", 0, chooseFile, "<Alt-c>"},
    {"Rotate PDF", "icons/icons8-rotate-24.png", "left", 0, actRotate, "<Alt-r>"},
    {"Home", "icons/icons8-home-24.png", "left", 0, home, "<Alt-h>"},
  }
  radioArr1 = RadioDefs{
    {"Rotate by 90 degrees clockwise", "90"},
    {"Rotate by 180 degrees clockwise", "180"},
    {"Rotate by 270 degrees clockwise", "270"},
  }
  radioArr2 = RadioDefs{
    // numbers 0-2 are defined in pdfcpu docs
    {"Rotate all the pages", "all"},
    {"Rotate only pages indicated in the field below", "selected"},
  }

  buttons = buttonsArr.CreateButtons(btnRow)
  buttons[1].Configure(tk.State("disabled"))

  CreateRadio(radioArr1, degreeChoice, radioRow1)
  separator := inputRow.TSeparator()
  tk.Pack(separator, tk.Fill("x"))
  CreateRadio(radioArr2, pagesChoice, radioRow2)
  CreateEntry(entryRow, entryLine, "Pages to rotate: ")
  PackBottomBtns(btnRow)
  buttonsArr.SetHotkeys()
}
