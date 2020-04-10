package main

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"path/filepath"
	"time"
)

var bookmarkMenuItems []*gtk.MenuItem

type Bookmark struct {
	Path       string
	Page       uint
	TotalPages uint
	Added      time.Time
}

func (gui *GUI) AddBookmark() {
	if !gui.Loaded() {
		gui.ShowError("Cannot add a bookmark before an archive is opened.")
		return
	}

	defer gui.RebuildBookmarksMenu()

	for i := range gui.Config.Bookmarks {
		b := &gui.Config.Bookmarks[i]
		if b.Path == gui.State.ArchivePath {
			b.Page = uint(gui.State.ArchivePos + 1)
			b.TotalPages = uint(gui.State.Archive.Len())
			b.Added = time.Now()
			return
		}
	}

	gui.Config.Bookmarks = append(gui.Config.Bookmarks, Bookmark{
		Path:       gui.State.ArchivePath,
		TotalPages: uint(gui.State.Archive.Len()),
		Page:       uint(gui.State.ArchivePos + 1),
		Added:      time.Now(),
	})
}

func (gui *GUI) RebuildBookmarksMenu() {
	for i := range bookmarkMenuItems {
		gui.MenuBookmarks.Remove(bookmarkMenuItems[i])
		bookmarkMenuItems[i].Destroy()
	}

	bookmarkMenuItems = nil
	gc()

	for i := range gui.Config.Bookmarks {
		bookmark := &gui.Config.Bookmarks[i]
		base := filepath.Base(bookmark.Path)
		label := fmt.Sprintf("%s (%d/%d)", base, bookmark.Page, bookmark.TotalPages)
		bookmarkMenuItem, err := gtk.MenuItemNewWithLabel(label)
		if err != nil {
			gui.ShowError(err.Error())
			return
		}
		bookmarkMenuItem.Connect("activate", func() {
			if gui.State.ArchivePath != bookmark.Path {
				gui.LoadArchive(bookmark.Path)
			}
			gui.SetPage(int(bookmark.Page) - 1)
		})
		bookmarkMenuItems = append(bookmarkMenuItems, bookmarkMenuItem)
		gui.MenuBookmarks.Append(bookmarkMenuItem)
	}

	gui.MenuBookmarks.ShowAll()
}
