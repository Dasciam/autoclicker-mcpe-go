package application

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/dasciam/autoclicker-mcpe-go/interfaces"
	"github.com/dasciam/autoclicker-mcpe-go/platform"
	"github.com/samber/lo"
	"image/color"
	"math/rand"
	"os"
	signal2 "os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"
)

type updateCallback struct {
	element fyne.CanvasObject
	fn      func(object fyne.CanvasObject)
}

type App struct {
	min, max atomic.Int32

	matcher atomic.Pointer[WindowMatcher]

	matchInput atomic.Pointer[string]
	matchIndex atomic.Int32

	matcherError atomic.Pointer[string]

	window fyne.Window

	display platform.Display

	updateCallbacks []updateCallback

	iteration chan time.Time
	close     chan struct{}
}

func New() (*App, error) {
	a := &App{
		close:     make(chan struct{}),
		iteration: make(chan time.Time),
	}

	a.min.Store(12)
	a.max.Store(16)

	minecraft := "Minecraft"
	a.matchInput.Store(&minecraft)
	a.matchIndex.Store(2)

	a.updateMatchSettings()
	return a, nil
}

func (a *App) Run() (err error) {
	a.display, err = interfaces.Open()
	if err != nil {
		return
	}

	a.window = a.deployWindow()

	return a.run()
}

func (a *App) Close() error {
	close(a.close)

	return a.display.Close()
}

func (a *App) run() error {
	go a.catchSignal()
	go a.setupClicker()

	// TODO: If there is ever a headless version,
	//  we need to run a.doMainAppCycle in the main goroutine
	go func() {
		err := a.doMainAppCycle()
		if err != nil {
			panic(err)
		}
	}()

	a.doWindowAppCycle()
	return nil
}

func (a *App) doMainAppCycle() error {
	appTicker := time.NewTicker(time.Second / 40)

	defer appTicker.Stop()

	for {
		select {
		case <-a.iteration:
			pointer, err := a.display.Pointer()
			if err != nil {
				return err
			}
			window, err := a.display.WindowFocus()
			if err != nil {
				return err
			}
			title := window.Title()

			matcher := a.matcher.Load()
			if matcher == nil {
				return errors.New("no matcher set")
			}

			if (*matcher).Match(title) && pointer.LoadMask(platform.FlagLMB) {
				width, _ := window.Size()
				wWidth := width / 2
				if runtime.GOOS == "windows" || wWidth+5 > pointer.X && wWidth-5 < pointer.X {
					err := window.EmulateMouseClick()
					if err != nil {
						return err
					}
				}
			}
		case <-appTicker.C:
			for _, cb := range a.updateCallbacks {
				cb.fn(cb.element)
			}
		case <-a.close:
			return nil
		}
	}
}

func (a *App) catchSignal() {
	signal := make(chan os.Signal, 1)

	signal2.Notify(signal, syscall.SIGINT)

	<-signal

	_ = a.Close()
}

func (a *App) setupClicker() {
	a.iteration = make(chan time.Time)

	var target int32 = 1

	ticker := time.NewTicker(time.Second / 2)
	defer ticker.Stop()

	for {
		select {
		case <-a.close:
			return
		case <-ticker.C:
			minCps := a.min.Load()
			maxCps := a.max.Load()

			atomic.StoreInt32(
				&target,
				minCps+rand.Int31n(max(maxCps-minCps, 1)),
			)
		default:
			a.iteration <- time.Now()
			time.Sleep(time.Second / time.Duration(target))
		}
	}
}

func (a *App) doWindowAppCycle() {
	a.window.SetMaster()

	a.window.SetOnClosed(func() {
		_ = a.Close()
	})

	a.window.ShowAndRun()
}

func (a *App) deployWindow() fyne.Window {
	fApp := app.New()

	w := fApp.NewWindow("Auto-Clicker")

	minSlider := widget.NewSlider(1, 20)
	minSlider.Value = float64(a.min.Load())
	maxSlider := widget.NewSlider(1, 20)
	maxSlider.Value = float64(a.max.Load())

	maxSlider.OnChanged = func(f float64) {
		s := int32(f)
		if s < a.min.Load() {
			minSlider.SetValue(f)
		}
		a.max.Store(s)
	}
	minSlider.OnChanged = func(f float64) {
		s := int32(f)
		if s > a.max.Load() {
			maxSlider.SetValue(f)
		}
		a.min.Store(s)
	}

	minText := canvas.NewText("", color.White)
	addGenericUpdateCallback(a, minText, func(object *canvas.Text) {
		object.Text = fmt.Sprintf("Min CPS: %d", a.min.Load())
	})

	maxText := canvas.NewText("", color.White)
	addGenericUpdateCallback(a, maxText, func(object *canvas.Text) {
		object.Text = fmt.Sprintf("Max CPS: %d", a.max.Load())
	})

	matcherErrorText := canvas.NewText("", color.RGBA{R: 0xff, A: 0xff})
	addGenericUpdateCallback(a, matcherErrorText, func(object *canvas.Text) {
		text := a.matcherError.Load()
		if text == nil {
			return
		}
		object.Text = fmt.Sprint("Error: ", *text)
	})

	content := container.NewBorder(
		container.New(layout.NewVBoxLayout(),
			// Target minimum CPS.
			minText,
			minSlider,
			// Target maximum CPS.
			maxText,
			maxSlider,
			container.New(layout.NewHBoxLayout(),
				canvas.NewText("Window checking: ", color.White),
				canvasObject(func() fyne.CanvasObject {
					v := widget.NewSelect(lo.Map(windowMatchStrategies[:], func(v WindowMatchStrategy, _ int) string {
						return v.String()
					}), nil)

					v.OnChanged = func(s string) {
						index := lo.IndexOf(v.Options, s)
						if index == -1 {
							return
						}
						a.matchIndex.Store(int32(index))

						a.updateMatchSettings()
					}

					v.SetSelectedIndex(int(a.matchIndex.Load()))
					return v
				}),
			),
			container.NewGridWithRows(3,
				canvas.NewText("Window name: ", color.White),
				canvasObject(func() fyne.CanvasObject {
					entry := widget.NewEntry()

					entry.Text = *a.matchInput.Load()

					entry.OnChanged = func(s string) {
						a.matchInput.Store(&s)

						a.updateMatchSettings()
					}
					return entry
				}),
				matcherErrorText,
			),
		),
		nil,
		nil,
		nil,
	)

	w.SetContent(content)

	w.Resize(fyne.NewSize(500, 400))
	w.CenterOnScreen()
	return w
}

func (a *App) updateMatchSettings() {
	matcher, err := windowMatchStrategies[a.matchIndex.Load()].New(*a.matchInput.Load())
	if err != nil {
		str := err.Error()
		a.matcherError.Store(&str)
		return
	}
	a.matcher.Store(
		&matcher,
	)
}

func (a *App) addUpdateCallback(object fyne.CanvasObject, fn func(object fyne.CanvasObject)) {
	a.updateCallbacks = append(a.updateCallbacks, updateCallback{
		element: object,
		fn:      fn,
	})
}

func addGenericUpdateCallback[T fyne.CanvasObject](a *App, v T, fn func(object T)) {
	a.addUpdateCallback(v, func(object fyne.CanvasObject) {
		fn(object.(T))
	})
}

func canvasObject(v func() fyne.CanvasObject) fyne.CanvasObject {
	return v()
}
