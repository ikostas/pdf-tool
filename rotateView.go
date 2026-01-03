package main

import (
  "path/filepath"
  "strconv"
  tk "modernc.org/tk9.0"
  "github.com/pdfcpu/pdfcpu/pkg/api"
)

// rotatePdf rotates all pages in the file for a given angle
func rotatePdf() {
  var buttonsArr buttonDefs
  var radioArr1, radioArr2 radioDefs
  var buttons []*tk.TButtonWidget
  var pagesArr []string
  var err error
  var angle int
  var chooseFile, home, actRotate func()
  var inputFile, output string
  var msgLabel *tk.TLabelWidget
  var fileRow *tk.FrameWidget

  degreeChoice := tk.Variable("")
  pagesChoice := tk.Variable("")

  rotate := tk.App.Frame()
  t := title{
    wmTitle: "PDF Tool -- rotate",
    title: "Rotate PDF",
    tipString: "Choose a file to rotate and an angle",
    isMainMenu: false,
    win: rotate,
    msgLabel: &msgLabel,
  }

  r := makeTitle(t)
  inputRow, msgRow, btnRow := r.ir, r.mr, r.br

  radioRow1 := inputRow.Frame() // degrees to rotate
  radioRow2 := inputRow.Frame() // all pages or selected only
  entryRow := inputRow.Frame() // a field to enter page numbers
  entryLine := entryRow.TEntry(tk.Placeholder("Example: -2, 3, 4-6, 7-"), tk.Textvariable(""))

  chooseFile = func() {
    if !chooseOneFile(&inputFile, inputRow, &fileRow) { return }
    buttons[1].Configure(tk.State("normal"))
  }
  home = func() {
    goHome(buttonsArr, &rotate)
  }
  actRotate = func() {
    if inputFile == "" {
      packMsg(msgRow, &msgLabel, "Don't press <Alt-R> unless you select a file to rotate")
      return
    }
    degreeChoiceVal := degreeChoice.Get()
    pagesChoiceVal := pagesChoice.Get()
    if pagesChoiceVal == "" {
      packMsg(msgRow, &msgLabel, "Select mode: all pages or selected pages")
      return
    } else if pagesChoiceVal == "selected" {
      pagesArr, err = createPagesArr(entryLine.Textvariable())
      if err != nil {
        packMsg(msgRow, &msgLabel, "Error parsing pages line: " + err.Error())
        return
      }
    } // there's one other case, pageChoice == 'all', then we need pagesToRotate = nil, which is already true

    switch degreeChoiceVal {
    case "90", "180", "270":
      // I don't check err here on purpose
      angle, _ = strconv.Atoi(degreeChoiceVal)
      go func() {
        output = inputFile[:len(inputFile)-len(filepath.Ext(inputFile))] + "_" + degreeChoiceVal + ".pdf"
        err = api.RotateFile(inputFile, output, angle, pagesArr, nil)
        r := report{
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
      packMsg(msgRow, &msgLabel, "Choose one of the options with the radio buttons.")
      return
    }
  }

  buttonsArr = buttonDefs{
    {"Choose a file", "icons/icons8-add-file-24.png", "left", 0, chooseFile, "<Alt-c>"},
    {"Rotate PDF", "icons/icons8-rotate-24.png", "left", 0, actRotate, "<Alt-r>"},
    {"Home", "icons/icons8-home-24.png", "left", 0, home, "<Alt-h>"},
  }
  radioArr1 = radioDefs{
    {"Rotate by 90 degrees clockwise", "90"},
    {"Rotate by 180 degrees clockwise", "180"},
    {"Rotate by 270 degrees clockwise", "270"},
  }
  radioArr2 = radioDefs{
    // numbers 0-2 are defined in pdfcpu docs
    {"Rotate all the pages", "all"},
    {"Rotate only pages indicated in the field below", "selected"},
  }

  buttons = buttonsArr.createButtons(btnRow)
  buttons[1].Configure(tk.State("disabled"))

  createRadio(radioArr1, degreeChoice, radioRow1)
  separator := inputRow.TSeparator()
  tk.Pack(separator, tk.Fill("x"))
  createRadio(radioArr2, pagesChoice, radioRow2)
  createEntry(entryRow, entryLine, "Pages to rotate: ")
  packBottomBtns(btnRow)
  buttonsArr.setHotkeys()
}
