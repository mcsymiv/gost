package driver

import (
	"fmt"
)

type InputSource string

const (
	NullInput    InputSource = "null"
	KeyInput     InputSource = "key"
	PointerInput InputSource = "pointer"
	WheelInput   InputSource = "wheel"
)

// PointerType is the type of pointer used by StorePointerActions.
// There are 3 different types according to the WC3 implementation.
type PointerType string

const (
	MousePointer PointerType = "mouse"
	PenPointer   PointerType = "pen"
	TouchPointer PointerType = "touch"
)

// PointerMoveOrigin controls how the offset for
// the pointer move action is calculated.
type PointerMoveOrigin string

const (
	// FromViewport calculates the offset from the viewport at 0,0.
	FromViewport PointerMoveOrigin = "viewport"
	// FromPointer calculates the offset from the current pointer position.
	FromPointer PointerMoveOrigin = "pointer"
)

// KeyAction represents an activity involving a keyboard key.
type KeyAction map[string]interface{}

// PointerAction represents an activity involving a pointer.
type PointerAction map[string]interface{}

// Actions stores KeyActions and PointerActions for later execution.
type Actions []map[string]interface{}

// KeyActions
// Resigters KeyActions
func (w *WebDriver) KeyActions(inputId string, actions ...KeyAction) {
	rawActions := []map[string]interface{}{}
	keyActions := []map[string]interface{}{}

	for _, action := range actions {
		rawActions = append(rawActions, action)
	}

	keyActions = append(keyActions, map[string]interface{}{
		"type":    KeyInput,
		"id":      inputId,
		"actions": rawActions,
	})

	// err := w.WebClient.Action(key, string(action), w.SessionId)
	// if err != nil {
	// 	panic(fmt.Sprintf("error on action: %v", err))
	// }
	//
	// err = w.WebClient.ReleaseAction(w.SessionId)
	// if err != nil {
	// 	panic(fmt.Sprintf("error on release action: %v", err))
	// }
}

// PointerActions
// Registers Pointer Actions
func (wd *WebDriver) PointerActions(inputId string, actions ...PointerAction) {
	rawActions := []map[string]interface{}{}
	pointerActions := []map[string]interface{}{}

	for _, action := range actions {
		rawActions = append(rawActions, action)
	}

	pointerActions = append(pointerActions, map[string]interface{}{
		"type":       PointerInput,
		"id":         inputId,
		"parameters": map[string]string{"pointerType": string(PointerInput)},
		"actions":    rawActions,
	})
}

type ActionType string

const (
	KeyUpAction   ActionType = "keyUp"
	KeyDownAction ActionType = "keyDown"

	PointerDown   ActionType = "pointerDown"
	PointerUp     ActionType = "pointerUp"
	PointerMove   ActionType = "pointerMove"
	PointerCancel ActionType = "pointerCancel"
)

func KeyDown(key string) KeyAction {
	return KeyAction{
		"type":  KeyDownAction,
		"value": key,
	}
}

func (w *WebDriver) Action(key string, action ActionType) {
	err := w.WebClient.Action(key, string(action), w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on action: %v", err))
	}

	err = w.WebClient.ReleaseAction(w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on release action: %v", err))
	}
}

// ReleaseAction
// Causes events to be fired
// as if the state was released by an explicit series of actions.
// It also clears all the internal state of the virtual devices.
func (w *WebDriver) ReleaseAction() {
	err := w.WebClient.ReleaseAction(w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on action: %v", err))
	}

}

func (w *WebDriver) Keys(key string) {
	err := w.WebClient.Action(key, string(KeyDownAction), w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on action: %v", err))
	}

	err = w.WebClient.ReleaseAction(w.SessionId)
	if err != nil {
		panic(fmt.Sprintf("error on release action: %v", err))
	}
}

