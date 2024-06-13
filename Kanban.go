package main

import (
	"encoding/json"
	"errors"
	"os"
)

var (
	//ErrInvalidColumn is returned when progress/regress item is given an invalid column number 
	ErrInvalidColumn = errors.New("Invalid Kanban Column")

    //ErrInvalidRegress is returned when trying to access an index ouside of the range of the columns
    ErrInvalidIndex = errors.New("Invalid Kanban Item Index")

	//ErrInvalidProgress is returned when trying to move an item past "Done" 
	ErrInvalidProgress = errors.New("Can't progress past 'Done'")

	//ErrInvalidRegress is returned when trying to move an item before "Not Started" 
	ErrInvalidRegress = errors.New("Can't regress to before 'Not Started'")
	
    //ErrInvalidProgress is returned when trying to move an item past "Done" 
	ErrInvalidMoveUp = errors.New("Can't move first item in list up")

	//ErrInvalidRegress is returned when trying to move an item before "Not Started" 
	ErrInvalidMoveDown = errors.New("Can't move final item in list down")
)

const (
    notStartedColIndex = iota
    startedColIndex 
    doneColIndex
    )

type Kanban struct {
    NotStarted []string `json:"notStarted"`
    Started []string    `json:"started"`
    Done []string       `json:"done"`
}

func (k *Kanban) assembleKanbanObject() error{
    file, err := os.ReadFile("kanban.json")
    if(err != nil){
        if(err == os.ErrExist){
            if _, err := os.Create("kanban.json"); err != nil{
                return err
            }else{
                return k.assembleKanbanObject() //TODO: Double check that this is a good idea 
            }
        }
        return err
    }

    err = json.Unmarshal(file, &k)
    if(err != nil){
        return err
    }

    return nil
}

func (k *Kanban) Save() error{
    out, err := json.MarshalIndent(&k, "", "    ")
    if(err != nil){
        return err
    }

    file, err := os.Create("kanban.json")
    if(err != nil){
        return err
    }
    if _, err := file.Write(out); err != nil{
        return err
    }
    
    return nil
}

func (k *Kanban) getColumnByNumber(column int) (*[]string, error){
    switch column{
        case notStartedColIndex:
            return &k.NotStarted, nil
        case startedColIndex:
            return &k.Started, nil
        case doneColIndex:
            return &k.Done, nil
        default:
            return nil, ErrInvalidColumn 
    }
}

func (k *Kanban) ProgressKanbanItem(column int , index int) error{
    if(index < 0){
       return ErrInvalidIndex
    }

    var item string
    switch column{
        case notStartedColIndex:
            item = k.NotStarted[index]
            k.NotStarted = append(k.NotStarted[:index], k.NotStarted[index+1:]...) 
            k.Started = append(k.Started, item)
        case startedColIndex:
            item = k.Started[index]
            k.Started = append(k.Started[:index], k.Started[index+1:]...) 
            k.Done = append(k.Done, item)
        case doneColIndex:
            return ErrInvalidProgress
        default:
            return ErrInvalidColumn 
    }

    return nil
}


func (k *Kanban) RegressKanbanItem(column int , index int) error{
    if(index < 0){
       return ErrInvalidIndex
    }

    var item string
    switch column{
        case notStartedColIndex:
            return ErrInvalidRegress
        case startedColIndex:
            item = k.Started[index]
            k.Started = append(k.Started[:index], k.Started[index+1:]...) 
            k.NotStarted = append(k.NotStarted, item)
        case doneColIndex:
            if(index >= len(k.Done)){
                return ErrInvalidIndex
            }
            item = k.Done[index]
            k.Done = append(k.Done[:index], k.Done[index+1:]...) 
            k.Started = append(k.Started, item)
        default:
            return ErrInvalidColumn
    }

    return nil
}

func (k* Kanban) MoveKanbanItemUp(column int, index int) error{
    arr, err := k.getColumnByNumber(column)
    if(err != nil){
        return err
    }

    if err:= moveUpIndexChecker(index, len(*arr)); err != nil{
        return err
    }

    item := (*arr)[index]
    (*arr)[index] = (*arr)[index-1] 
    (*arr)[index-1] = item

    return nil
}


func (k* Kanban) MoveKanbanItemDown(column int, index int) error{
    arr, err := k.getColumnByNumber(column)
    if(err != nil){
        return err
    }

    if err:= moveDownIndexChecker(index, len(*arr)); err != nil{
        return err
    }

    item := (*arr)[index+1]
    (*arr)[index+1] = (*arr)[index] 
    (*arr)[index] = item

    return nil
}

func (k* Kanban) AddKanbanItem(column int, index int, item string) error{
    arr, err := k.getColumnByNumber(column)
    if(err != nil){
        return err
    }

    if(index == len(*arr)){
        *arr = append(*arr, item)
    }else{
        //TODO: this was producting odd results... see if you can replicate
        *arr = append((*arr)[:index+1], (*arr)[index:]...) 
        (*arr)[index] = item
        
        //Try this if ^ still gives issues
        // *arr = append(*arr, "") 
        // copy((*arr)[index+1:], (*arr)[index:])
        // (*arr)[index] = item 
        
    }
    

    return nil
}

func (k* Kanban) RemoveKanbanItem(column int, index int) error{
    arr, err := k.getColumnByNumber(column)
    if(err != nil){
        return err
    }

    if err:= removeIndexChecker(index, len(*arr)); err != nil{
        return err
    }

    if(index == len(*arr)-1){
        *arr = (*arr)[:index]
    }else{
        *arr = append((*arr)[:index], (*arr)[index+1:]...) 
    }

    return nil
}

func moveUpIndexChecker(index int, length int) error{
    if(index < 0 || index >= length) {
        return ErrInvalidIndex
    }else if (index == 0) {
        return ErrInvalidMoveUp
    }

    return nil
}

func moveDownIndexChecker(index int, length int) error{
    if(index < 0 || index >= length) {
        return ErrInvalidIndex
    }else if (index == length-1) {
        return ErrInvalidMoveDown
    }

    return nil
}


func removeIndexChecker(index int, length int) error{
    if(index < 0 || index >= length) {
        return ErrInvalidIndex
    }

    return nil
}
