package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type menuItem struct {
	displayString  string
	displayStrings []string
	onPress        func() error
}

// every item in a list context needs an ID
func (i *menuItem) ID() string {
	if i.displayString != "" {
		return i.displayString
	}

	return strings.Join(i.displayStrings, "-")
}

// list panel functions

func (gui *Gui) handleMenuSelect() error {
	return nil
}

// specific functions

func (gui *Gui) renderMenuOptions() error {
	optionsMap := map[string]string{
		gui.getKeyDisplay("universal.return"): gui.Tr.SLocalize("close"),
		fmt.Sprintf("%s %s", gui.getKeyDisplay("universal.prevItem"), gui.getKeyDisplay("universal.nextItem")): gui.Tr.SLocalize("navigate"),
		gui.getKeyDisplay("universal.select"): gui.Tr.SLocalize("execute"),
	}
	return gui.renderOptionsMap(optionsMap)
}

func (gui *Gui) handleMenuClose(g *gocui.Gui, v *gocui.View) error {
	err := g.DeleteView("menu")
	if err != nil {
		return err
	}
	return gui.returnFromContext()
}

type createMenuOptions struct {
	showCancel bool
}

func (gui *Gui) createMenu(title string, items []*menuItem, createMenuOptions createMenuOptions) error {
	if createMenuOptions.showCancel {
		// this is mutative but I'm okay with that for now
		items = append(items, &menuItem{
			displayStrings: []string{gui.Tr.SLocalize("cancel")},
			onPress: func() error {
				return nil
			},
		})
	}

	gui.State.MenuItems = items

	stringArrays := make([][]string, len(items))
	for i, item := range items {
		if item.displayStrings == nil {
			stringArrays[i] = []string{item.displayString}
		} else {
			stringArrays[i] = item.displayStrings
		}
	}

	list := utils.RenderDisplayStrings(stringArrays)

	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(false, list)
	menuView, _ := gui.g.SetView("menu", x0, y0, x1, y1, 0)
	menuView.Title = title
	menuView.FgColor = theme.GocuiDefaultTextColor
	menuView.ContainsList = true
	menuView.Clear()
	menuView.SetOnSelectItem(gui.onSelectItemWrapper(func(selectedLine int) error {
		return nil
	}))
	fmt.Fprint(menuView, list)
	gui.State.Panels.Menu.SelectedLine = 0

	gui.g.Update(func(g *gocui.Gui) error {
		return gui.switchContext(gui.Contexts.Menu.Context)
	})
	return nil
}

func (gui *Gui) onMenuPress() error {
	selectedLine := gui.State.Panels.Menu.SelectedLine
	if err := gui.State.MenuItems[selectedLine].onPress(); err != nil {
		return err
	}

	return gui.returnFromContext()
}
