package main

import (
  "path/filepath"
  "regexp"
  "strconv"
  tk "modernc.org/tk9.0"
  "github.com/pdfcpu/pdfcpu/pkg/api"
)

// splitPdf draws split menu
func splitPdf() {
  var buttonsArr buttonDefs
  var radioArr radioDefs
  var buttons []*tk.TButtonWidget
  var inputFile string
  var err error
  var pageCount, span int
  var pagesArr []int
  var dir string
  var chooseFile, splitFile, home, splitFile3 func()
  var msgLabel *tk.TLabelWidget
  var fileRow *tk.FrameWidget

  radioChoice := tk.Variable("")
  split := tk.App.Frame()

  t := title{
    wmTitle: "PDF Tool -- split",
    title: "Split PDF",
    tipString: "Choose a file to split and pages, delimeted by commas (excluding).\nExample: if you specify page 2, the first page will be in the first file and all the other pages in the second.",
    isMainMenu: false,
    win: split,
    msgLabel: &msgLabel,
  }

  r := makeTitle(t)
  inputRow, msgRow, btnRow := r.ir, r.mr, r.br

  entryRow := inputRow.Frame() // a field to enter page numbers
  radioRow := inputRow.Frame() // a radio button to choose mode
  entryLine := entryRow.TEntry(tk.Placeholder("Example: 2, 4, 6"), tk.Textvariable(""))

  chooseFile = func() {
    if !chooseOneFile(&inputFile, inputRow, &fileRow) { return }
    buttons[1].Configure(tk.State("normal"))
  }
  splitFile = func() {
    if inputFile == "" {
      packMsg(msgRow, &msgLabel, "Don't press <Alt-S> unless you select a file to split")
      return
    }
    dir = filepath.Dir(inputFile)
    radioChoiceVal := radioChoice.Get()
    switch radioChoiceVal {
    case "0", "1", "2":
      // I don't check err here on purpose
      span, _ = strconv.Atoi(radioChoiceVal)
      go func() {
        err = api.SplitFile(inputFile, dir, span, nil)
        r := report{
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
      packMsg(msgRow, &msgLabel, "Choose one of the options with the radio buttons.")
    }
  }

  splitFile3 = func() {
    // must start with one digit and not 0, spaces, commas and digits after that
    validInput := regexp.MustCompile(`^[1-9](?:[,\s]*\d)*$`)
    if !validInput.MatchString(entryLine.Textvariable()) {
      packMsg(msgRow, &msgLabel, "Invalid input: only single digits, commas, and spaces allowed, but at least one digit should be specified")
      return
    }
    pageCount, err = api.PageCountFile(inputFile)
    if err != nil {
      packMsg(msgRow, &msgLabel, "Can't count pages of input file: " + err.Error())
      return
    }
    pagesArr, err = parsePages(entryLine.Textvariable(), pageCount)
    if err != nil {
      packMsg(msgRow, &msgLabel, err.Error())
      return
    }
    go func() {
      err = api.SplitByPageNrFile(inputFile, dir, pagesArr, nil);
      r := report{
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
    goHome(buttonsArr, &split)
  }

  buttonsArr = buttonDefs{
    {"Choose a file", "icons/icons8-add-file-24.png", "left", 0, chooseFile, "<Alt-c>"},
    {"Split PDF", "icons/icons8-split-24.png", "left", 0, splitFile, "<Alt-s>"},
    {"Home", "icons/icons8-home-24.png", "left", 0, home, "<Alt-h>"},
  }
  radioArr = radioDefs{
    // numbers 0-2 are defined in pdfcpu docs
    {"Split by single page (burst)", "1"},
    {"Split by two pages", "2"},
    {"Split by bookmarks", "0"},
    {"Split by numbers specified below", "3"},
  }

  buttons = buttonsArr.createButtons(btnRow)
  buttons[1].Configure(tk.State("disabled"))
  createRadio(radioArr, radioChoice, radioRow)
  createEntry(entryRow, entryLine, "Pages to use for splitting: ")
  packBottomBtns(btnRow)
  buttonsArr.setHotkeys()
}
