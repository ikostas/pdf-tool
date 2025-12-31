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

// ShowHome draws main menu
func ShowHome() {
  var buttonsArr, bottomBtnsArr ButtonDefs
  var merge, split, rotate, remove, extract, closeApp, unBindDestroy func()
  var msgLabel *tk.TLabelWidget

  home := tk.App.Frame()
  t := Title{
    wmTitle: "PDF Tool -- main menu",
    title: "Home",
    tipString: "Choose operation",
    isMainMenu: true,
    win: home,
    msgLabel: &msgLabel,
  }

  _, _, btnRow := MakeTitle(t)
  br1 := home.Frame()

  closeApp = func() {
    tk.Destroy(tk.App)
  }
  unBindDestroy = func() {
    buttonsArr.UnBind()
    bottomBtnsArr.UnBind()
    tk.Destroy(home)
    home = nil
  }
  merge = func() {
    unBindDestroy()
    MergePdf()
  }
  split = func() {
    unBindDestroy()
    SplitPdf()
  }
  rotate = func() {
    unBindDestroy()
    RotatePdf()
  }
  remove = func() {
    unBindDestroy()
    RemovePdfPages()
  }
  extract = func() {
    unBindDestroy()
    ExtractPdf()
  }

  buttonsArr = ButtonDefs{
    // todo: unbind everything before calling other functions
    {"Merge PDFs", "icons/icons8-merge-24.png", "left", 0, merge, "<Alt-m>"},
    {"Split PDF", "icons/icons8-split-24.png", "left", 0, split, "<Alt-s>"},
    {"Rotate PDF pages", "icons/icons8-rotate-24.png", "left", 0, rotate, "<Alt-r>"},
    {"Remove PDF pages", "icons/icons8-remove-24.png", "left", 3, remove, "<Alt-o>"},
    {"Extract PDF pages", "icons/icons8-tweezers-24.png", "left", 1, extract, "<Alt-x>"},
  }
  bottomBtnsArr = ButtonDefs{
    {"About", "icons/icons8-about-24.png", "left", 0, aboutMessageBox, "<Alt-a>"},
    {"Exit", "icons/icons8-exit-24.png", "left", 0, closeApp, "<Alt-e>"},
  }
  _ = buttonsArr.CreateButtons(br1)
  _ = bottomBtnsArr.CreateButtons(btnRow)
  tk.Pack(br1, tk.Fill("x"))
  PackBottomBtns(btnRow)
  buttonsArr.SetHotkeys()
  bottomBtnsArr.SetHotkeys()
}
