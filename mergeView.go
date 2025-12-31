package main

import (
  "path/filepath"
  "fmt"
  tk "modernc.org/tk9.0"
  "github.com/pdfcpu/pdfcpu/pkg/api"
)

// display list of files in a text box
func (fileTags *FileTagArr) displayFiles(display *tk.TextWidget, mergeButton *tk.TButtonWidget) {
  display.Configure(tk.State("normal"))
  display.Clear()

  moveFile := func (i int, op string) {
    switch op {
    case "up":
      (*fileTags)[i], (*fileTags)[i-1] = (*fileTags)[i-1], (*fileTags)[i]
    case "down":
      (*fileTags)[i], (*fileTags)[i+1] = (*fileTags)[i+1], (*fileTags)[i]
    case "rm":
      (*fileTags) = append((*fileTags)[:i], (*fileTags)[i+1:]...)
    }
    fileTags.displayFiles(display, mergeButton)
  }

  if len(*fileTags) < 2 {
    mergeButton.Configure(tk.State("disabled"))
  } else {
    mergeButton.Configure(tk.State("normal"))
  }

  boldFont := tk.NewFont(tk.Weight("bold"))
  for i, f := range *fileTags {
    localI := i
    localF := f
    for _, tag := range []string{localF.upTagName, localF.downTagName, localF.removeTagName} {
      localTag := tag
      display.TagConfigure(
        localTag,
        tk.Foreground("blue"),
        tk.Font(boldFont),
      )
      display.TagBind(
        localTag,
        "<Button-1>",
        func(e *tk.Event) { 
          switch localTag {
          case localF.upTagName:
            moveFile(localI, "up")
          case localF.downTagName:
            moveFile(localI, "down")
          case localF.removeTagName:
            moveFile(localI, "rm")
          }
        },
      )
      display.TagBind(localTag, "<Enter>", func(e *tk.Event) {
        display.Configure(tk.Cursor("hand2"))
      })
      display.TagBind(localTag, "<Leave>", func(e *tk.Event) {
        display.Configure(tk.Cursor(""))
      })
    }
    switch {
    case len(*fileTags) == 1:
      display.InsertML(fmt.Sprintf("%d. %s %s", i+1, f.filename, f.removeTag), "<br>")
    case localI == 0:
      display.InsertML(fmt.Sprintf("%d. %s %s %s", i+1, f.filename, f.downTag, f.removeTag), "<br>")
    case localI == len(*fileTags) - 1:
      display.InsertML(fmt.Sprintf("%d. %s %s %s", i+1, f.filename, f.upTag, f.removeTag), "<br>")
    default: 
      display.InsertML(fmt.Sprintf("%d. %s %s %s %s", i+1, f.filename, f.upTag, f.downTag, f.removeTag), "<br>")
    }
  }
  display.Configure(tk.State("disabled"))
}

// MergePdf renders title, instruction, input and message rows
func MergePdf() {
  var fileTags FileTagArr
  var scroll *tk.TScrollbarWidget
  var buttonsArr ButtonDefs
  var display *tk.TextWidget
  var buttons []*tk.TButtonWidget
  var addAndUpdate, clearAndUpdate, home, actMerge func()
  var nextID int
  var msgLabel *tk.TLabelWidget

  merge := tk.App.Frame()
  t := Title{
    wmTitle: "PDF Tool -- merge",
    title: "Merge PDFs",
    tipString: "Choose files to merge.\nTip: hold Ctrl or Shift to select several files.",
    isMainMenu: false,
    win: merge,
    msgLabel: &msgLabel,
  }

  inputRow, msgRow, btnRow := MakeTitle(t)

  addAndUpdate = func() {
    filesToMerge := AddFiles(true)
    if len(filesToMerge) == 1 && filesToMerge[0] == "" { return }
    fileTags.CreateTags(filesToMerge, &nextID)
    fileTags.displayFiles(display, buttons[1])
  }
  clearAndUpdate = func() {
    fileTags = nil
    fileTags.displayFiles(display, buttons[1])
  }
  home = func() {
    GoHome(buttonsArr, &merge)
  }
  actMerge = func() {
    if len(fileTags) < 2 {
      PackMsg(msgRow, &msgLabel, "Don't press <Alt-M> unless you select at least 2 files to merge")
      return
    }
    dir := filepath.Dir(fileTags[0].fileWithPath)
    output := filepath.Join(dir, "merged.pdf")

    go func() {
      var inFiles []string

      for _, f := range fileTags {
        inFiles = append(inFiles, f.fileWithPath)
      }
      err := api.MergeCreateFile(inFiles, output, false, nil)
      r := Report{
        msgRow: msgRow,
        msgLabel: &msgLabel,
        msgSuccess: "Files merged to: ",
        msgFail: "Merge failed: ",
        result: output,
        err: err,
      }
      tk.PostEvent(func() { reportRes(r) }, false)
    }()
  }

  display = inputRow.Text(tk.Height(5), tk.State("disabled"), tk.Wrap("none"), tk.Setgrid(true), tk.Yscrollcommand(func(e *tk.Event) { e.ScrollSet(scroll) }))
  scroll = inputRow.TScrollbar(tk.Command(func(e *tk.Event) { e.Yview(display) }))
  buttonsArr = ButtonDefs{
    {"Add Files", "icons/icons8-add-file-24.png", "left", 0, addAndUpdate, "<Alt-a>"},
    {"Merge PDFs", "icons/icons8-merge-24.png", "left", 0, actMerge, "<Alt-m>"},
    {"Clear", "icons/icons8-clear-24.png", "left", 0, clearAndUpdate, "<Alt-c>"},
    {"Home", "icons/icons8-home-24.png", "left", 0, home, "<Alt-h>"},
  }
  buttons = buttonsArr.CreateButtons(btnRow)
  buttons[1].Configure(tk.State("disabled"))
  tk.Pack(display, tk.Side("left"), tk.Fill("both"), tk.Expand(true), tk.Padx("0"), tk.Pady("0"))
  tk.Pack(scroll, tk.Side("right"), tk.Fill("y"))
  PackBottomBtns(btnRow)
  buttonsArr.SetHotkeys()
}
