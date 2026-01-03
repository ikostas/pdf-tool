package main
/* todo:
- look for more operations in pdfcpu -- add functions? */

import (
  tk "modernc.org/tk9.0"
  _ "modernc.org/tk9.0/themes/azure"
)

func main() {
  tk.ActivateTheme("azure light")
  showHome()
  tk.App.Center().Wait()
}
