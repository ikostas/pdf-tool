package main 

import (
  tk "modernc.org/tk9.0"
)

func aboutMessageBox() {
  tk.MessageBox(tk.Title("About"), tk.Type("ok"), tk.Msg(`PDF Tool
Copyright (C) 2025 Konstantin Ovchinnikov k@kovchinnikov.info

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

You should have received a copy of the GNU General Public License along with this program. If not, see https://www.gnu.org/licenses/.

Source code: https://github.com/ikostas/pdftool

Icons by icons8.com`), )
}

// showHome draws main menu
func showHome() {
  var buttonsArr, bottomBtnsArr buttonDefs
  var merge, split, rotate, remove, extract, closeApp, unBindDestroy func()
  var msgLabel *tk.TLabelWidget

  home := tk.App.Frame()
  t := title{
    wmTitle: "PDF Tool -- main menu",
    title: "Home",
    tipString: "Choose operation",
    isMainMenu: true,
    win: home,
    msgLabel: &msgLabel,
  }

  r := makeTitle(t)
  _, _, btnRow := r.ir, r.mr, r.br
  br1 := home.Frame()

  closeApp = func() {
    tk.Destroy(tk.App)
  }
  unBindDestroy = func() {
    buttonsArr.unbind()
    bottomBtnsArr.unbind()
    tk.Destroy(home)
    home = nil
  }
  merge = func() {
    unBindDestroy()
    mergePdf()
  }
  split = func() {
    unBindDestroy()
    splitPdf()
  }
  rotate = func() {
    unBindDestroy()
    rotatePdf()
  }
  remove = func() {
    unBindDestroy()
    removePdfPages()
  }
  extract = func() {
    unBindDestroy()
    extractPdf()
  }

  buttonsArr = buttonDefs{
    {"Merge PDFs", "icons/icons8-merge-24.png", "left", 0, merge, "<Alt-m>"},
    {"Split PDF", "icons/icons8-split-24.png", "left", 0, split, "<Alt-s>"},
    {"Rotate PDF pages", "icons/icons8-rotate-24.png", "left", 0, rotate, "<Alt-r>"},
    {"Remove PDF pages", "icons/icons8-remove-24.png", "left", 3, remove, "<Alt-o>"},
    {"Extract PDF pages", "icons/icons8-tweezers-24.png", "left", 1, extract, "<Alt-x>"},
  }
  bottomBtnsArr = buttonDefs{
    {"About", "icons/icons8-about-24.png", "left", 0, aboutMessageBox, "<Alt-a>"},
    {"Exit", "icons/icons8-exit-24.png", "left", 0, closeApp, "<Alt-e>"},
  }
  _ = buttonsArr.createButtons(br1)
  _ = bottomBtnsArr.createButtons(btnRow)
  tk.Pack(br1, tk.Fill("x"))
  packBottomBtns(btnRow)
  buttonsArr.setHotkeys()
  bottomBtnsArr.setHotkeys()
}
