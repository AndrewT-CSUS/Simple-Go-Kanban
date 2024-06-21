package main

import (
	"errors"
	"fmt"
	"log"
	"github.com/jroimartin/gocui"
)

const(
    notStartedColumn = "Not Started"
    startedColumn = "Started"
    doneColumn = "Done"
    logBox = "Log Box"
    progressPopup = "Progress Popup"
    regressPopup = "Regress Popup"
    hintText = "Hint Text"
    newItemPopup = "New Item Popup"
    deletePopup = "Delete Popup"
)

var (       //state
	viewArr = []string{notStartedColumn, startedColumn, doneColumn}
	activeCol  = 0
    prevActiveCol = notStartedColumn
    cursorPositions = []int{0, 0, 0}
    showLog = false
    edited = false
)

var queuedItem string

func activeColAsConst() (string, error){
    switch activeCol{
        case 0: 
            return notStartedColumn, nil
        case 1: 
            return startedColumn, nil
        case 2: 
            return doneColumn, nil
        default:
            return "", errors.New("Invalid Column")
    }
}

var kanban Kanban

func main(){
    err := kanban.assembleKanbanObject()
    if(err != nil){}

   	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	
    g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
    

	g.SetManagerFunc(layout)
    buildKeymaps(g)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func buildKeymaps(g *gocui.Gui) error {
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
        return err
	}

    if err := g.SetKeybinding(progressPopup, gocui.KeyEsc, gocui.ModNone, clearPopups); err != nil {
        return err
	}
    
    if err := g.SetKeybinding(progressPopup, 'q', gocui.ModNone, clearPopups); err != nil {
        return err
	}
    
    if err := g.SetKeybinding(progressPopup, 'n', gocui.ModNone, clearPopups); err != nil {
        return err
	}

    if err := g.SetKeybinding(progressPopup, 'y', gocui.ModNone, confimProgressPopup); err != nil {
        return err
	}
    
    if err := g.SetKeybinding(regressPopup, gocui.KeyEsc, gocui.ModNone, clearPopups); err != nil {
        return err
	}
    
    if err := g.SetKeybinding(regressPopup, 'q', gocui.ModNone, clearPopups); err != nil {
        return err
	}
    
    if err := g.SetKeybinding(regressPopup, 'n', gocui.ModNone, clearPopups); err != nil {
        return err
	}

    if err := g.SetKeybinding(regressPopup, 'y', gocui.ModNone, confirmRegressPopup); err != nil {
        return err
	}

    if err := g.SetKeybinding("", '~', gocui.ModNone, toggleLog); err != nil {
        return err
	}
    
    if err := g.SetKeybinding(logBox, 'c', gocui.ModNone, logClear); err != nil {
        return err
	}

    if err := g.SetKeybinding(newItemPopup, gocui.KeyEnter, gocui.ModNone, confirmNewItemPopup); err != nil {
        return err
	}

    if err := g.SetKeybinding(newItemPopup, gocui.KeyCtrlQ, gocui.ModNone, clearPopups); err != nil {
        return err
	}

    if err := g.SetKeybinding(deletePopup, gocui.KeyEsc, gocui.ModNone, clearPopups); err != nil {
        return err
	}
    
    if err := g.SetKeybinding(deletePopup, 'q', gocui.ModNone, clearPopups); err != nil {
        return err
	}
    
    if err := g.SetKeybinding(deletePopup, 'n', gocui.ModNone, clearPopups); err != nil {
        return err
	}

    if err := g.SetKeybinding(deletePopup, 'y', gocui.ModNone, confirmDeletePopup); err != nil {
        return err
	}


    for _, element := range viewArr {
        buildAllColumKeymaps(g, element)
    }

    return nil
}

func buildAllColumKeymaps(g *gocui.Gui, column string) error {

    if err := g.SetKeybinding(column, 'q', gocui.ModNone, quit); err != nil {
        return err
	}

	if err := g.SetKeybinding(column, gocui.KeyArrowDown, gocui.ModNone, cursorUp); err != nil {
        return err
	}
	
    if err := g.SetKeybinding(column, 'j', gocui.ModNone, cursorUp); err != nil {
        return err
	}

	if err := g.SetKeybinding(column, gocui.KeyArrowUp, gocui.ModNone, cursorDown); err != nil {
        return err
	}
	
    if err := g.SetKeybinding(column, 'k', gocui.ModNone, cursorDown); err != nil {
        return err
	}

	if err := g.SetKeybinding(column, gocui.KeyArrowRight, gocui.ModNone, nextView); err != nil {
        return err
	}

    if err := g.SetKeybinding(column, 'l', gocui.ModNone, nextView); err != nil {
        return err
	}

	if err := g.SetKeybinding(column, gocui.KeyArrowLeft, gocui.ModNone, prevView); err != nil {
        return err
	}

	if err := g.SetKeybinding(column, 'h', gocui.ModNone, prevView); err != nil {
        return err
	}
	
    if err := g.SetKeybinding(column, 'f', gocui.ModNone, progressKanbanItemPopup); err != nil {
        return err
	}
	
    if err := g.SetKeybinding(column, 'g', gocui.ModNone, regressKanbanItemPopup); err != nil {
        return err
	}
    
    if err := g.SetKeybinding(column, 'L', gocui.ModNone, progressKanbanItemPopup); err != nil {
        return err
	}
	
    if err := g.SetKeybinding(column, 'H', gocui.ModNone, regressKanbanItemPopup); err != nil {
        return err
	}

    if err := g.SetKeybinding(column, 'J', gocui.ModNone, moveKanbanItemDown); err != nil {
        return err
	}
    
    if err := g.SetKeybinding(column, 'K', gocui.ModNone, moveKanbanItemUp); err != nil {
        return err
	}

    if err := g.SetKeybinding(column, 'n', gocui.ModNone, addKanbanItem); err != nil {
        return err
	}

    if err := g.SetKeybinding(column, 'd', gocui.ModNone, showDeletePopup); err != nil {
        return err
	}

    return nil
}


func layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    colWidth := int(maxX/3)
    colStart := 0
    if v, err := g.SetView(notStartedColumn, 0, 0, maxX/3-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
        
        v.Frame = true
        v.Title = v.Name() 
        v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	
        if _, err = g.SetCurrentView(notStartedColumn); err != nil {
			return err
		}
	}

    colStart = colStart+colWidth
	
	if v, err := g.SetView(startedColumn, maxX/3-1, 0, (maxX/3-1)*2, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

        v.Frame = true
        v.Title = v.Name() 
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		
	}

    colStart = colStart+colWidth

	if v, err := g.SetView(doneColumn, (maxX/3-1)*2, 0, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

        v.Frame = true
        v.Title = v.Name() 
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

	}

	if v, err := g.SetView(logBox, 0, maxY/2, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

        v.Frame = true
        v.Title = v.Name() 
		v.FgColor = gocui.ColorCyan
        _, err =  g.SetViewOnBottom(logBox)
	}


	if v, err := g.SetView(hintText, 0, maxY-2, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

        v.Frame = false
        _, err =  g.SetViewOnBottom(hintText)
        fmt.Fprintf(v, "[h] - move left  [j] - move down  [k] - move up  [l] - move right  |  [H], [g] - move item back  [J] - move item down  [K] - move item up  [L], [f] - move item forward |  [q] - quit")

	}

    //TODO: See if I can find a more optimal place to put this. No need to redraw the text each time
    if err := populateColumns(g); err != nil {
        return err
    }

	return nil

}

func nextView(g *gocui.Gui, view *gocui.View) error {
	nextIndex := activeCol + 1
    if(nextIndex > 2){
        nextIndex = 0
    }
	name := viewArr[nextIndex]

    logWriter(g, "Going from view", view.Name(), "to", name)

    newView, err := g.SetCurrentView(name)
    if err != nil {
		return err
	}
    
    newView.Highlight = true
    newView.SelBgColor = gocui.ColorGreen
    newView.SelFgColor = gocui.ColorBlack
    view.Highlight = false
    view.SelBgColor = gocui.ColorDefault
	view.SelFgColor = gocui.ColorDefault

	activeCol = nextIndex
	return nil
}

func prevView(g *gocui.Gui, view *gocui.View) error {
	nextIndex := activeCol - 1
    if(nextIndex < 0){
        nextIndex = len(viewArr)-1
    }
	name := viewArr[nextIndex]

    newView, err := g.SetCurrentView(name)
    if err != nil {
		return err
	}

    logWriter(g, "Going from view", view.Name(), "to", name)
    
    newView.Highlight = true
    newView.SelBgColor = gocui.ColorGreen
    newView.SelFgColor = gocui.ColorBlack
    view.Highlight = false
    view.SelBgColor = gocui.ColorDefault
	view.SelFgColor = gocui.ColorDefault

	activeCol = nextIndex
	return nil
    
}


func cursorUp(g *gocui.Gui, v *gocui.View) error {
    rc, err := incCursor()
    if(err != nil){
        return err
    }else if(rc){
        return nil
    }
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		if err := v.SetCursor(cx, cy+1); err != nil {
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
    }

    logWriter(g, "Cursor moved to item", cursorPositions[activeCol]+1)
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
    rc, err := decCursor()
    if(err != nil){
        return err
    }else if(rc){
        return nil
    }
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
        }

	}
    

    logWriter(g, "Cursor moved to item", cursorPositions[activeCol]+1)
	return nil
}

func incCursor() (bool, error){ 
    var colSize int;
    switch activeCol{
    case 0:
        colSize = len(kanban.NotStarted)
    case 1:
        colSize = len(kanban.Started)
    case 2:
        colSize = len(kanban.Done)
    }

    if(cursorPositions[activeCol] >= colSize-1){
        return true,  nil
    }else if(cursorPositions[activeCol] >= colSize){
        return false, errors.New("Cursor out of bounds")
    }

    cursorPositions[activeCol]++

    return false, nil
}

func decCursor() (bool, error){
    // logWriter(g, "DEC")
    if(cursorPositions[activeCol] < 0){
        return false, errors.New("Cursor out of bounds")
    }else if(cursorPositions[activeCol] == 0){
        return true, nil
    }

    cursorPositions[activeCol] -= 1
    
    return false, nil
}

func moveKanbanItemUp(g *gocui.Gui, v *gocui.View) error{
    if err:= kanban.MoveKanbanItemUp(activeCol, cursorPositions[activeCol]); err != nil{
        logWriter(g, err)
        if(err == ErrInvalidMoveUp){
            //TODO: New popup: error message
            return nil
        }else{
            return err
        }
    }
    return cursorDown(g, v)

}

func moveKanbanItemDown(g *gocui.Gui, v *gocui.View) error{
    if err := kanban.MoveKanbanItemDown(activeCol, cursorPositions[activeCol]); err != nil{
        logWriter(g, err)
        if(err == ErrInvalidMoveDown){
            //TODO: New popup: error message
            return nil
        }else{
            return err
        }
    }
    return cursorUp(g, v)
}

func addKanbanItem(g *gocui.Gui, v *gocui.View) error{
    prevActiveCol = g.CurrentView().Name()
    if err := showNewItemPopup(g); err != nil{
        return err
    }

    return nil
}

func showNewItemPopup(g *gocui.Gui) error{
	maxX, maxY := g.Size()
	if v, err := g.SetView(newItemPopup, maxX/2-30, maxY/2-2, maxX/2+30, maxY/2+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
        v.Editable = true
        v.Frame = true
        v.Title = "New Item Name:"
		if _, err := g.SetCurrentView(newItemPopup); err != nil {
			return err
		}
	}
	return nil
}

func confirmNewItemPopup(g *gocui.Gui, v *gocui.View) error{
    if _, err := g.SetCurrentView(newItemPopup); err != nil {
        return err
    }
    queuedItem := v.Buffer()
    if err := clearPopups(g, v); err != nil{
        return err
    }
    if (len(queuedItem) != 0){   
        if err := kanban.AddKanbanItem(activeCol, cursorPositions[activeCol], queuedItem[:len(queuedItem)-1]); err != nil{
            logWriter(g, err)
            return err
        }else{
           queuedItem = ""
           edited = true 
        }
    }
    return nil
}



func deleteKanbanItem(g *gocui.Gui, v *gocui.View) error{
    if err := kanban.RemoveKanbanItem(activeCol, cursorPositions[activeCol]); err != nil{
        logWriter(g, err)
        if(err == ErrInvalidIndex){
            //TODO: New popup: error message
            return nil
        }else{
            return err
        }
    }
    edited = true 
    temp, err := kanban.getColumnByNumber(activeCol)
    if(err != nil){
        return err
    }
    if(cursorPositions[activeCol] == len(*temp)){
        return cursorDown(g, v)
    }else{
        return nil
    }
}

func progressKanbanItemPopup(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error
    prevActiveCol = v.Name()
    logWriter(g, "Progressing, prep return to " + prevActiveCol)

	_, cy := v.Cursor()
    
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView(progressPopup, maxX/2-30, maxY/2-2, maxX/2+30, maxY/2+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
        //TODO: DYNAMIC MESSAGE BASED OFF prevActive
		fmt.Fprintln(v, "Move " + l + " forward?")
		fmt.Fprintln(v, "[y] - confirm  [n] - cancel")
		if _, err := g.SetCurrentView(progressPopup); err != nil {
			return err
		}
	}
	return nil
}


func regressKanbanItemPopup(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error
    prevActiveCol = v.Name()
    logWriter(g, "Reverting, prep return to " + prevActiveCol)

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView(regressPopup, maxX/2-30, maxY/2-2, maxX/2+30, maxY/2+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
        //TODO: DYNAMIC MESSAGE BASED OFF prevActive
		fmt.Fprintln(v, "Move " + l + " back?")
		fmt.Fprintln(v, "[y] - confirm  [n] - cancel")
		if _, err := g.SetCurrentView(regressPopup); err != nil {
			return err
		}
	}
	return nil
}

func toggleLog(g *gocui.Gui, _ *gocui.View) (error){ 
    var err error
    if(showLog){
        _, err =  g.SetViewOnBottom(logBox)
        g.SelFgColor = gocui.ColorGreen
        showLog = false
		if _, err := g.SetCurrentView(prevActiveCol); err != nil {
			return err
		}
    }else{
        _, err = g.SetViewOnTop(logBox)
        g.SelFgColor = gocui.ColorCyan
        showLog = true
        if prevActiveCol, err = activeColAsConst(); err != nil {
            log.Panic(err)
            return err
        }
        if _, err := g.SetCurrentView(logBox); err != nil {
			return err
		}
    }
    return err
}

func confimProgressPopup(g *gocui.Gui, v *gocui.View) error {
    if err := kanban.ProgressKanbanItem(activeCol, cursorPositions[activeCol]); err != nil{
        if(err == ErrInvalidProgress){
            //TODO: New popup: error message
        }else{
            return err
        }
    }
    edited = true
    return clearPopups(g, v)
}

func confirmRegressPopup(g *gocui.Gui, v *gocui.View) error {
    if err := kanban.RegressKanbanItem(activeCol, cursorPositions[activeCol]); err != nil{
        if(err == ErrInvalidRegress){
            //TODO: New popup: error message
        }else{
            return err
        }
    }
    edited = true
    clearPopups(g, v)
	return nil
}

func showDeletePopup(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error
	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	maxX, maxY := g.Size()
	if v, err := g.SetView(deletePopup, maxX/2-30, maxY/2-2, maxX/2+30, maxY/2+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
        v.Frame = true
        v.Title = "Confirm"
        fmt.Fprintln(v, "Are you sure you want to delete " + l + "?")
		fmt.Fprintln(v, "[y] - confirm  [n] - cancel")
		if _, err := g.SetCurrentView(deletePopup); err != nil {
			return err
		}
	}
    return nil
}

func confirmDeletePopup(g *gocui.Gui, v *gocui.View) error {
    if err := deleteKanbanItem(g, v); err != nil{
        return err
    }
    edited = true
    clearPopups(g, v)
	return nil
}

func clearPopups(g *gocui.Gui, _ *gocui.View) error {
	if err := g.DeleteView(progressPopup); err != nil && err.Error() != gocui.ErrUnknownView.Error(){
		return err
	}
	if err := g.DeleteView(regressPopup); err != nil && err.Error() != gocui.ErrUnknownView.Error(){
		return err
	}
	if err := g.DeleteView(newItemPopup); err != nil && err.Error() != gocui.ErrUnknownView.Error(){
		return err
	}
	if err := g.DeleteView(deletePopup); err != nil && err.Error() != gocui.ErrUnknownView.Error(){
		return err
	}
	if _, err := g.SetCurrentView(prevActiveCol); err != nil {
		return err
	}
	return nil
}


func quit(g *gocui.Gui, v *gocui.View) error {
    if(edited){
        //TODO: Ask user if they want to save changes
        if err := kanban.Save(); err != nil{
            return err;
        }else{
	        return gocui.ErrQuit
        }
    }
	return gocui.ErrQuit
}


func populateColumns(g *gocui.Gui) error{
    // logWriter(g, kanban)
    if nsc, err := g.View(notStartedColumn); err != nil {
        return errors.Join(err, errors.New(notStartedColumn))
    }else{
        nsc.Clear()
        for _, element := range kanban.NotStarted {
            fmt.Fprintf(nsc, "%s\n", element)
        }
    }
    
    if sc, err := g.View(startedColumn); err != nil {
        return errors.Join(err, errors.New(startedColumn))
    }else{
        sc.Clear()
        for _, element := range kanban.Started {
            fmt.Fprintf(sc, "%s\n", element)
        }
    }

    if dc, err := g.View(doneColumn); err != nil {
        return errors.Join(err, errors.New(doneColumn))
    }else{
        dc.Clear()
        for _, element := range kanban.Done {
            fmt.Fprintf(dc, "%s\n", element)
        }
    }

    return nil
}

func logWriter(g *gocui.Gui, a ...any) (err error){
	out, err := g.View(logBox)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, a...)
    return nil
}

func logClear(g *gocui.Gui, v *gocui.View) (err error){
	out, err := g.View(logBox)
	if err != nil {
		return err
	}
    out.Clear()
    return nil
}
