package views

import (
    "gioui.org/io/key"
    "gioui.org/layout"
    "gioui.org/unit"
    "gioui.org/widget"
    "github.com/ul-gaul/go-basestation/ui/components"
    "github.com/ul-gaul/go-basestation/ui/types"
    tabs2 "github.com/ul-gaul/go-basestation/ui/views/tabs"
    "time"
)

var _ types.IView = (*mainView)(nil)

type mainView struct {
    generalTab *tabs2.GeneralTab
    motorTab   *tabs2.MotorTab
    tabBar     *components.TabLayout
    btnMenu    widget.Clickable
}

func (m *mainView) Keypress(ev key.Event) {
    // Switch Tab (ctrl+tab // ctrl+shift+tab)
    if ev.Name == key.NameTab && ev.Modifiers.Contain(key.ModCtrl) && !ev.Modifiers.Contain(key.ModAlt) {
        i := m.tabBar.CurrentTab() + 1
        if ev.Modifiers.Contain(key.ModShift) {
            i = m.tabBar.CurrentTab() - 1
        }
        
        n := len(m.tabBar.Tabs)
        if i < 0 {
            m.tabBar.SetCurrentTab(i%n + n)
        } else {
            m.tabBar.SetCurrentTab(i % n)
        }
    }
}

func (m *mainView) Tick(delta time.Duration) {
    switch m.tabBar.CurrentTab() {
    case 0:
        m.generalTab.Tick(delta)
    case 1:
        m.motorTab.Tick(delta)
    }
}

func (m *mainView) Draw(gtx layout.Context) layout.Dimensions {
    return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return layout.Dimensions{} // TODO Menu button
                }),
                layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
                layout.Rigid(m.tabBar.TabBar))
        }),
        layout.Flexed(1, m.tabBar.Content))
}
