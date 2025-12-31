package main

import (
  "path/filepath"
  "regexp"
  "strconv"
  tk "modernc.org/tk9.0"
  "github.com/pdfcpu/pdfcpu/pkg/api"
)

// SplitPdf draws split menu
func SplitPdf() {
  var buttonsArr ButtonDefs
  var radioArr RadioDefs
  var buttons []*tk.TButtonWidget
  var fileToSplit []string
  var err error
  var pageCount, span int
  var pagesArr []int
  var dir string
  var chooseFile, splitFile, home, splitFile3 func()
  var msgLabel *tk.TLabelWidget
  var fileRow *tk.FrameWidget

  radioChoice := tk.Variable("")
  split := tk.App.Frame()

  t := Title{
    wmTitle: "PDF Tool -- split",
    title: "Split PDF",
    tipString: "Choose a file to split and pages, delimeted by commas (excluding).\nExample: if you specify page 2, the first page will be in the first file and all the other pages in the second.",
    isMainMenu: false,
    win: split,
    msgLabel: &msgLabel,
  }

  inputRow, msgRow, btnRow := MakeTitle(t)
  entryRow := inputRow.Frame() // a field to enter page numbers
  radioRow := inputRow.Frame() // a radio button to choose mode
  entryLine := entryRow.TEntry(tk.Placeholder("Example: 2, 4, 6"), tk.Textvariable(""))

  chooseFile = func() {
    if !ChooseOneFile(&fileToSplit, inputRow, &fileRow) { return }
    buttons[1].Configure(tk.State("normal"))
  }
  splitFile = func() {
    if fileToSplit == nil {
      PackMsg(msgRow, &msgLabel, "Don't press <Alt-S> unless you select a file to split")
      return
    }
    dir = filepath.Dir(fileToSplit[0])
    radioChoiceVal := radioChoice.Get()
    switch radioChoiceVal {
    case "0", "1", "2":
      // I don't check err here on purpose
      span, err = strconv.Atoi(radioChoiceVal)
      go func() {
        err = api.SplitFile(fileToSplit[0], dir, span, nil)
        r := Report{
          msgRow: msgRow,
          msgLabel: &msgLabel,
          msgSuccess: "Split files can be found in: ",
          msgFail: "Can't split the file: ",
          result: dir,
          err: err,
        }
        tk.PostEvent(func() { reportRes(r) }, false)
      }()
    case "3":
      splitFile3()
    default:
      PackMsg(msgRow, &msgLabel, "Choose one of the options with the radio buttons.")
    }
  }

  splitFile3 = func() {
    // must start with one digit and not 0, spaces, commas and digits after that
    validInput := regexp.MustCompile(`^[1-9](?:[,\s]*\d)*$`)
    if !validInput.MatchString(entryLine.Textvariable()) {
      PackMsg(msgRow, &msgLabel, "Invalid input: only single digits, commas, and spaces allowed, but at least one digit should be specified")
      return
    }
    pageCount, err = api.PageCountFile(fileToSplit[0])
    if err != nil {
      PackMsg(msgRow, &msgLabel, "Can't count pages of input file: " + err.Error())
      return
    }
    pagesArr, err = ParsePages(entryLine.Textvariable(), pageCount)
    if err != nil {
      PackMsg(msgRow, &msgLabel, err.Error())
      return
    }
    go func() {
      err = api.SplitByPageNrFile(fileToSplit[0], dir, pagesArr, nil);
      r := Report{
        msgRow: msgRow,
        msgLabel: &msgLabel,
        msgSuccess: "Split files can be found in: ",
        msgFail: "Can't split the file: ",
        result: dir,
        err: err,
      }
      tk.PostEvent(func() { reportRes(r) }, false)
    }()
  }
  
  home = func() {
    GoHome(buttonsArr, &split)
  }

  buttonsArr = ButtonDefs{
    {"Choose a file", "icons/icons8-add-file-24.png", "left", 0, chooseFile, "<Alt-c>"},
    {"Split PDF", "icons/icons8-split-24.png", "left", 0, splitFile, "<Alt-s>"},
    {"Home", "icons/icons8-home-24.png", "left", 0, home, "<Alt-h>"},
  }
  radioArr = RadioDefs{
    // numbers 0-2 are defined in pdfcpu docs
    {"Split by single page (burst)", "1"},
    {"Split by two pages", "2"},
    {"Split by bookmarks", "0"},
    {"Split by numbers specified below", "3"},
  }

  buttons = buttonsArr.CreateButtons(btnRow)
  buttons[1].Configure(tk.State("disabled"))
  CreateRadio(radioArr, radioChoice, radioRow)
  CreateEntry(entryRow, entryLine, "Pages to use for splitting: ")
  PackBottomBtns(btnRow)
  buttonsArr.SetHotkeys()
}
