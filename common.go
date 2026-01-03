package main

import (
  "embed"
  "fmt"
  "path/filepath"
  "regexp"
  "sort"
  "strconv"
  "strings"
  tk "modernc.org/tk9.0"
)

//go:embed icons/*.png
var iconFS embed.FS

// buttonDef is a struct for tk button
type buttonDef struct {
  txt string
  image string
  compound string
  underline int
  command func()
  key string
}

// buttonDefs is an array of tk buttons
type buttonDefs []buttonDef

// radioDef is tk radio switch
type radioDef struct {
  txt string
  value string
}

// radioDefs is an array of radio switches
type radioDefs []radioDef

// fileTag contains data for the files to display
type fileTag struct {
  fileWithPath string
  filename string
  upTagName string
  downTagName string
  removeTagName string
  upTag string
  downTag string
  removeTag string
}

// report is used to report result of func call
type report struct {
  msgRow *tk.TLabelframeWidget
  msgLabel **tk.TLabelWidget
  msgSuccess string
  msgFail string
  result string
  err error
}

// title is used to create window title and rows
type title struct {
  wmTitle string
  title string
  tipString string
  isMainMenu bool
  win *tk.FrameWidget
  msgLabel **tk.TLabelWidget
}

// rows are used to return 3 rows for UI
type rows struct {
  ir *tk.TLabelframeWidget
  mr *tk.TLabelframeWidget
  br *tk.FrameWidget
}

// fileTagArr is an array of files to display
type fileTagArr []fileTag

// createButtons creates tk buttons
func (defs buttonDefs) createButtons(parent *tk.FrameWidget) []*tk.TButtonWidget{
  var buttons []*tk.TButtonWidget
  for _, b := range defs {
    data, err := iconFS.ReadFile(b.image)
		if err != nil {
			panic(err)
		}

    btn := parent.TButton(
      tk.Txt(b.txt),
      tk.Image(tk.NewPhoto(tk.Data(data))),
      tk.Compound(b.compound),
      tk.Underline(b.underline),
      tk.Command(b.command),
    )
    tk.Pack(btn, tk.Side("left"), tk.Padx("2m"), tk.Pady("1m"))
    buttons = append(buttons,  btn)
  }
  return buttons
}

// setHotkeys sets hotkeys for buttons
func (defs buttonDefs) setHotkeys() {
  for _, i := range defs {
    tk.Bind(tk.App, i.key, tk.Command(i.command))
  }
}

// unbind unbinds hotkeys for buttons
func (defs buttonDefs) unbind() {
  for _, i := range defs {
    tk.Bind(tk.App, i.key, "")
  }
}

// addFiles opens file dialog
func addFiles(multiple bool) []string{
  file := tk.GetOpenFile(
    tk.Multiple(multiple),
    tk.Filetypes([]tk.FileType{
      {TypeName: "PDF files", Extensions: []string{".pdf"}},
      {TypeName: "All files", Extensions: []string{"*"}},
    }),
  )
  if len(file) == 0 {
    return nil
  }
  return file
}

// createTags creates tk tags to move files up/down or remove
func (fileTags *fileTagArr) createTags(fileEntries []string, nextID *int) {

  for _, f := range fileEntries {
    i := *nextID
    *nextID++
    upTagName := fmt.Sprintf("up_tag_%d", i)
    downTagName := fmt.Sprintf("down_tag_%d", i)
    removeTagName := fmt.Sprintf("rm_tag_%d", i)

    ft := fileTag{
      fileWithPath: f,
      filename: filepath.Base(f),
      upTagName: upTagName,
      downTagName: downTagName,
      removeTagName: removeTagName,
      upTag: fmt.Sprintf("<%s>[↑]</%s>", upTagName, upTagName),
      downTag: fmt.Sprintf("<%s>[↓]</%s>", downTagName, downTagName),
      removeTag: fmt.Sprintf("<%s>[✖]</%s>", removeTagName, removeTagName),
    }
    *fileTags = append(*fileTags, ft)
  }
}

// makeTitle creates window title, displays tips an returns 3 rows
func makeTitle(t title) rows {
  var inputRow *tk.TLabelframeWidget
  var msgRow *tk.TLabelframeWidget

  tk.App.WmTitle(t.wmTitle)
  titleFont := tk.NewFont(tk.Family("TkHeadingFont"), tk.Size(22), tk.Weight("bold"))
  titleRow := t.win.Frame()
  instRow := t.win.TLabelframe(tk.Txt("Instructions"), tk.Padding("10 10"))
  btnRow := t.win.Frame()
  windowTitle := titleRow.TLabel(tk.Txt(t.title), tk.Font(titleFont)) 
  tips := instRow.TLabel(tk.Txt(t.tipString))

  // always pack title & instructions row
  tk.Pack(windowTitle, tips, tk.Side("left"))
  tk.Pack(titleRow, instRow, tk.Pady("5"), tk.Fill("x"))
  tk.Pack(t.win, tk.Fill("both"), tk.Expand(true))
  
  // pack Input and Message rows only if required
  if !t.isMainMenu {
    inputRow = t.win.TLabelframe(tk.Txt("User input"), tk.Padding("10 10"))
    tk.Pack(inputRow, tk.Pady("5"), tk.Fill("x"))
    msgRow = t.win.TLabelframe(tk.Txt("Messages"), tk.Padding("10 10"))
    tk.Pack(msgRow, tk.Pady("5"), tk.Fill("both"), tk.Expand(true))
    // without this Message window is not visible
    *t.msgLabel = msgRow.TLabel(tk.Txt(""))
    tk.Pack(*t.msgLabel)
  }

  // inputRow and msgRow can be returned as nil
  r := rows{inputRow, msgRow, btnRow}
  return r
}

// packBottomBtns makes the end row centered
func packBottomBtns(btnRow *tk.FrameWidget) {
  leftSpacer := btnRow.Frame()
  rightSpacer := btnRow.Frame()
  centeredButtons := btnRow.Frame()
  tk.Pack(leftSpacer, tk.Side("left"), tk.Expand(true))
  tk.Pack(centeredButtons, tk.Side("left"))
  tk.Pack(rightSpacer, tk.Side("left"), tk.Expand(true))
  tk.Pack(btnRow)
}

// packMsg displays messages
func packMsg(msgRow *tk.TLabelframeWidget, msgLabel **tk.TLabelWidget, msg string) {
  // msgLabel always exists as it's initalized with ""
  tk.Destroy(*msgLabel)
  *msgLabel = msgRow.TLabel(tk.Txt(msg))
  tk.Pack(*msgLabel, tk.Side("left"), tk.Padx("5"), tk.Pady("5"))
}

// parsePages parses user input, the pages to work on
func parsePages(pages string, pageCount int) ([]int, error) {
  /* cases to consider:
1. duplicated numbers - filter
2. wrong order - sort ascending
3. any number more than page count in file - return err
4. spaces - remove all first
*/
  var pagesArr []int
  // some 'production code' here he he :)
  pages = strings.ReplaceAll(pages, " ", "")
  pagesStr := strings.Split(pages, ",")
  seen := make(map[int]bool)
  for _, p := range pagesStr {
    if p == "" { continue }
    n, err := strconv.Atoi(p)
    if err != nil {
      return nil, fmt.Errorf("error converting to number: %s", p)
    }
    if n > pageCount {
      return nil, fmt.Errorf("page you entered exceeds number of pages in file: %d", n)
    }
    if n < 2 {
      return nil, fmt.Errorf("numbers less than 2 don't work")
    }
    if !seen[n] {
      seen[n] = true
      pagesArr = append(pagesArr, n)
    }
  }
  sort.Ints(pagesArr)
  return pagesArr, nil
}

// chooseOneFile opens file open dialog, adds a file to array and outputs the filename in the message area
func chooseOneFile(file *string, inputRow *tk.TLabelframeWidget, fileRow **tk.FrameWidget) bool {
  tmpFile := addFiles(false)
  // don't change *file
  if len(tmpFile) == 1 && (tmpFile)[0] == "" { 
    return false
  }
  *file = tmpFile[0]
  if *fileRow != nil {
    tk.Destroy(*fileRow)
  }
  *fileRow = inputRow.Frame() // output the chosen file
  fileLabel := (*fileRow).TLabel(tk.Txt("Chosen file: " + filepath.Base(*file)))
  tk.Pack(fileLabel, tk.Side("left"), tk.Pady(5))
  tk.Pack(*fileRow, tk.Fill("x"), tk.Expand(true))
  return true
}

// createRadio created radio buttons and places them in Grid
func createRadio(radioArr radioDefs, radioVar *tk.VariableOpt, radioRow *tk.FrameWidget) {
  for i, r := range radioArr {
    rb := radioRow.TRadiobutton(tk.Txt(r.txt), radioVar, tk.Value(r.value))
    tk.Grid(rb, tk.Row(i), tk.Column(0), tk.Padx("10"), tk.Pady("10"), tk.Sticky("nsew"))
  }
  tk.Pack(radioRow)
}

// createEntry creates a field and label for text entry
func createEntry(entryRow *tk.FrameWidget, entryLine *tk.TEntryWidget, desc string) {
  tk.Pack(entryRow.TLabel(tk.Txt(desc)), tk.Side("left"), tk.Pady("5"))
  tk.Pack(entryLine, tk.Side("left"))
  tk.Pack(entryRow, tk.Fill("x"), tk.Expand(true))
}

// createPagesArr creates an array of strings from user input, one string
func createPagesArr(input string) ([]string, error){
  var pagesArr []string

  // must be at least one digit, can contain commas and spaces
  validInput := regexp.MustCompile(`^(?:[,\s]*-?\d+-?)+$`)
  if !validInput.MatchString(input) {
    return nil, fmt.Errorf("only single digits, commas, dashes and spaces allowed")
  }

  // convert to []strings, no checks on purpose
  pages := strings.ReplaceAll(input, " ", "")
  pagesStr := strings.Split(pages, ",")
  for _, p := range pagesStr {
    if p == "" { continue }
    pagesArr = append(pagesArr, p)
  }
  return pagesArr, nil
}

func reportRes(r report) bool {
  if r.err != nil {
    packMsg(r.msgRow, r.msgLabel, r.msgFail + r.err.Error())
    return false
  }
  packMsg(r.msgRow, r.msgLabel, r.msgSuccess + r.result)
  return true
}

// goHome destroys frame, unbinds key and calls home view
func goHome (buttonsArr buttonDefs, win **tk.FrameWidget) {
    buttonsArr.unbind()
    tk.Destroy(*win)
    *win = nil
    showHome()
  }
