package gui

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/cfanatic/go-expense/account"
	"github.com/cfanatic/go-expense/database"
	"github.com/cfanatic/go-expense/datasheet"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	NORMAL = iota
	BOLD
	UNDERLINE
)

var (
	FILTER = []string{""}
)

type Gui struct {
	widgets.QWidget

	qapp    *widgets.QApplication
	lapp    *widgets.QVBoxLayout
	lbutton *widgets.QHBoxLayout

	twidget  *widgets.QTabWidget
	twinfo   *widgets.QWidget
	twlabel  *widgets.QWidget
	twmonth  *widgets.QWidget
	twyear   *widgets.QWidget
	twlinfo  *widgets.QHBoxLayout
	twllabel *widgets.QHBoxLayout
	twlmonth *widgets.QVBoxLayout
	twlyear  *widgets.QHBoxLayout

	tview *widgets.QTableView
	tlist *gui.QStandardItemModel

	ds    *datasheet.Datasheet
	db    *database.Database
	dlist []database.Content

	_ func() `constructor:"init"`

	_ func() `slot:"connect"`
}

func (w *Gui) init() {
	geometry := widgets.QApplication_Desktop().AvailableGeometry(0)
	w.SetWindowTitle("go-expensegui")
	w.SetGeometry2(0, 0, int(float32(geometry.Width())/2.2), geometry.Height()/2)
	w.Move2((geometry.Width()-w.Width())/2, (geometry.Height()-w.Height())/2)

	w.lapp = widgets.NewQVBoxLayout()
	w.lbutton = widgets.NewQHBoxLayout()

	w.twidget = widgets.NewQTabWidget(nil)
	w.twinfo = widgets.NewQWidget(nil, core.Qt__Widget)
	w.twlabel = widgets.NewQWidget(nil, core.Qt__Widget)
	w.twmonth = widgets.NewQWidget(nil, core.Qt__Widget)
	w.twyear = widgets.NewQWidget(nil, core.Qt__Widget)
	w.twlinfo = widgets.NewQHBoxLayout()
	w.twllabel = widgets.NewQHBoxLayout()
	w.twlmonth = widgets.NewQVBoxLayout()
	w.twlyear = widgets.NewQHBoxLayout()

	w.twidget.AddTab(w.twinfo, "Info")
	w.twidget.AddTab(w.twlabel, "Label")
	w.twidget.AddTab(w.twmonth, "Month")
	w.twidget.AddTab(w.twyear, "Year")
	w.twidget.SetTabEnabled(1, false)
	w.twidget.SetTabEnabled(2, false)
	w.twidget.SetTabEnabled(3, false)

	w.twmonth.SetLayout(w.twlmonth)

	w.tview = widgets.NewQTableView(nil)
	w.tlist = gui.NewQStandardItemModel(nil)

	w.tview.SetModel(w.tlist)
	w.twllabel.AddWidget(w.tview, 0, 0)
	w.twlabel.SetLayout(w.twllabel)

	blabel := widgets.NewQPushButton2("Label", nil)
	banalyze := widgets.NewQPushButton2("Analyze", nil)
	bquit := widgets.NewQPushButton2("Quit", nil)

	w.lbutton.AddWidget(blabel, 0, 0)
	w.lbutton.AddWidget(banalyze, 0, 0)
	w.lbutton.AddWidget(bquit, 0, 0)
	w.lapp.AddWidget(w.twidget, 0, 0)
	w.lapp.AddLayout(w.lbutton, 0)
	w.SetLayout(w.lapp)

	w.twidget.ConnectTabBarClicked(w.analyze)
	w.tlist.ConnectItemChanged(w.update)
	blabel.ConnectClicked(w.label)
	bquit.ConnectClicked(func(bool) { w.qapp.Exit(0) })
	w.ConnectKeyPressEvent(w.keypressevent)
}

func (w *Gui) label(bool) {
	var export []datasheet.Content
	var err error

	if f := widgets.QFileDialog_GetOpenFileNames(nil, "Open datasheet", core.QDir_HomePath(), "*.xlsx", "", 0); len(f) > 0 {
		w.ds = datasheet.New(f[0], account.TAB)
		w.db = database.New(account.ADDRESS, account.NAME, account.COLLECT)
		if w.ds.Err != nil {
			widgets.QMessageBox_Critical(nil,
				"Cannot open datasheet",
				"Error during datasheet initialization attempt!",
				widgets.QMessageBox__Default,
				widgets.QMessageBox__Default,
			)
			return
		}
		if w.db.Err != nil {
			widgets.QMessageBox_Critical(nil,
				"Cannot open database",
				"Error during database initialization attempt!",
				widgets.QMessageBox__Default,
				widgets.QMessageBox__Default,
			)
			return
		}
	} else {
		return
	}

	if export, err = w.ds.Content(); err == nil {
		w.twidget.SetTabEnabled(1, true)
		w.twidget.SetTabEnabled(2, true)
		w.twidget.SetCurrentIndex(1)
	} else {
		widgets.QMessageBox_Critical(nil,
			"Cannot import datasheet",
			"Error during datasheet import!",
			widgets.QMessageBox__Default,
			widgets.QMessageBox__Default,
		)
		return
	}

	if w.tlist.RowCount(core.NewQModelIndex()) > 0 {
		w.tlist.RemoveRows(0, w.tlist.RowCount(core.NewQModelIndex()), core.NewQModelIndex())
		w.dlist = []database.Content{}
	}

	for _, trans := range export {
		items := []*gui.QStandardItem{}
		if item, err := w.db.Document("payee", trans.Payee); err == nil {
			items = append(items,
				gui.NewQStandardItem2(trans.Date),
				gui.NewQStandardItem2(trans.Payee),
				gui.NewQStandardItem2(fmt.Sprintf("%.2f", trans.Amount)),
				gui.NewQStandardItem2(item.Label),
			)
		} else {
			items = append(items,
				gui.NewQStandardItem2(trans.Date),
				gui.NewQStandardItem2(trans.Payee),
				gui.NewQStandardItem2(fmt.Sprintf("%.2f", trans.Amount)),
				gui.NewQStandardItem2(""),
			)
		}
		date, _ := time.Parse("01-02-06", trans.Date)
		items[0].SetEditable(false)
		items[0].SetSelectable(false)
		items[1].SetToolTip(trans.Description)
		w.tlist.AppendRow(items)
		w.dlist = append(w.dlist, database.Content{Date: date, Payee: trans.Payee, Amount: trans.Amount})
	}

	width := float32(w.Geometry().Width())
	w.tlist.SetHorizontalHeaderLabels([]string{"Date", "Payee", "Amount", "Label"})
	w.tview.HorizontalHeader().SetSectionResizeMode(widgets.QHeaderView__Interactive)
	w.tview.HorizontalHeader().SetStretchLastSection(true)
	w.tview.HorizontalHeader().ResizeSection(0, int(width/10))
	w.tview.HorizontalHeader().ResizeSection(1, int(width/2))
	w.tview.HorizontalHeader().ResizeSection(2, int(width/10))
	w.tview.HorizontalHeader().ResizeSection(3, int(width/5))
	w.tview.VerticalHeader().SetSectionResizeMode(widgets.QHeaderView__Stretch)

	w.save(true)
}

func (w *Gui) analyze(index int) {
	if w.twidget.TabText(index) != "Month" {
		return
	}

	row := func(style int, items ...string) *widgets.QWidget {
		widget := widgets.NewQWidget(nil, 0)
		layout := widgets.NewQHBoxLayout()
		for _, item := range items {
			label := widgets.NewQLabel2(item, nil, core.Qt__Widget)
			switch font := label.Font(); style {
			case BOLD:
				font.SetBold(true)
				label.SetFont(font)
			case UNDERLINE:
				font.SetUnderline(true)
				label.SetFont(font)
			default:
			}
			layout.AddWidget(label, 0, 0)
		}
		widget.SetLayout(layout)
		return widget
	}

	sarea := widgets.NewQScrollArea(nil)
	swidget := widgets.NewQWidget(nil, 0)
	slayout := widgets.NewQVBoxLayout2(swidget)

	dateIn := w.dlist[len(w.dlist)-1].Date
	dateOut := w.dlist[0].Date
	labels, _ := w.db.Labels("label")
	cnt := make(map[string]int)
	res := make(map[string]float64)
	tot := float64(0.0)

	// analyze
	for _, label := range labels {
		content, _ := w.db.Content("label", label.(string), dateIn.Format("01-02-06"), dateOut.Format("01-02-06"))
		cnt[label.(string)] = len(content)
		for _, trans := range content {
			res[trans.Label] = res[trans.Label] + (-1.0 * float64(trans.Amount))
			tot = tot + (-1.0 * float64(trans.Amount))
		}
	}
	// filter
	for _, label := range FILTER {
		tot = tot - res[label]
		delete(res, label)
	}
	// sort
	cats := make([]string, 0, len(res))
	for cat := range res {
		cats = append(cats, cat)
	}
	sort.Strings(cats)

	slayout.AddWidget(row(UNDERLINE, fmt.Sprintf("%s %d", dateIn.Month(), dateIn.Year()), "Euro / Month", "Percent / Month", "Count / Month"), 0, 0)
	for _, cat := range cats {
		slayout.AddWidget(row(NORMAL,
			fmt.Sprintf("%s", cat),
			fmt.Sprintf("%.f €", res[cat]),
			fmt.Sprintf("%.f %%", math.Round(res[cat]/tot*100)),
			fmt.Sprintf("%d", cnt[cat])), 0, 0)
	}
	spacer := widgets.NewQSpacerItem(0, 0, widgets.QSizePolicy__Minimum, widgets.QSizePolicy__Expanding)
	slayout.AddSpacerItem(spacer)

	sarea.SetWidget(swidget)
	sarea.SetWidgetResizable(true)

	if w.twlmonth.Count() > 0 {
		tmp := w.twlmonth.ItemAt(0).Widget()
		tmp.Hide()
		w.twlmonth.RemoveWidget(tmp)
	}
	w.twlmonth.InsertWidget(0, sarea, 0, 0)
}

func (w *Gui) save(bool) {
	for row := 0; row < w.tlist.RowCount(core.NewQModelIndex()); row++ {
		trans := []string{}
		for col := 0; col < w.tlist.ColumnCount(core.NewQModelIndex()); col++ {
			index := w.tlist.Index(row, col, core.NewQModelIndex())
			data := w.tlist.Data(index, int(core.Qt__DisplayRole))
			trans = append(trans, data.ToString())
		}
		if doc := w.document(trans); len(doc.Label) > 0 {
			w.db.Save(doc)
		}
	}
}

func (w *Gui) update(item *gui.QStandardItem) {
	var trans []string

	trans = make([]string, 0)
	for col := 0; col < w.tlist.ColumnCount(core.NewQModelIndex()); col++ {
		index := w.tlist.Index(item.Row(), col, core.NewQModelIndex())
		data := w.tlist.Data(index, int(core.Qt__DisplayRole))
		trans = append(trans, data.ToString())
	}

	dsold := datasheet.Content{
		Date:   w.dlist[item.Row()].Date.Format("01-02-06"),
		Payee:  w.dlist[item.Row()].Payee,
		Amount: w.dlist[item.Row()].Amount,
	}
	dsnew := datasheet.Content{
		Date:   w.document(trans).Date.Format("01-02-06"),
		Payee:  w.document(trans).Payee,
		Amount: w.document(trans).Amount,
	}
	dbold := w.dlist[item.Row()]
	dbnew := w.document(trans)

	w.ds.Update(dsold, dsnew)
	w.db.Update(dbold, dbnew)

	w.dlist[item.Row()] = w.document(trans)
}

func (w *Gui) document(trans []string) database.Content {
	date, _ := time.Parse("01-02-06", trans[0])
	payee := trans[1]
	amount, _ := strconv.ParseFloat(trans[2], 32)
	label := trans[3]
	return database.Content{Date: date, Payee: payee, Amount: float32(amount), Label: label}
}

func (w *Gui) keypressevent(e *gui.QKeyEvent) {
	if e.Key() == int(core.Qt__Key_Escape) {
		w.qapp.Exit(0)
	}
}

func (w *Gui) InitWith(qapp *widgets.QApplication) {
	w.qapp = qapp
}
