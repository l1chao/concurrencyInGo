package main

// 定义状态类型
type State string

// 定义所有可能状态
const (
	Idle      State = "idle"
	Selected  State = "selected"
	Paid      State = "paid"
	Dispensed State = "dispensed"
)

// 定义事件类型
type Event string

// 定义所有可能事件
const (
	SelectItem   Event = "select_item"
	InsertCoin   Event = "insert_coin"
	Dispense     Event = "dispense"
	ReturnChange Event = "return_change"
)

// 状态转移规则（关键：明确定义合法转换）
var stateTransitions = map[State]map[Event]State{
	Idle: {
		SelectItem: Selected,
	},
	Selected: {
		InsertCoin:   Paid,
		ReturnChange: Idle, // 取消购买
	},
	Paid: {
		Dispense: Dispensed,
	},
	Dispensed: {
		ReturnChange: Idle, // 完成交易
	},
}

// 状态机实现
type VendingMachine struct {
	currentState State
}

func (vm *VendingMachine) Transition(event Event) {
	// 检查当前状态是否有该事件的转移规则
	if nextState, ok := stateTransitions[vm.currentState][event]; ok {
		vm.executeSideEffects(vm.currentState, event, nextState)
		vm.currentState = nextState
	} else {
		println("非法操作：状态", vm.currentState, "不允许事件", event)
	}
}

// 处理状态转换时的副作用（如出货、退币等）
func (vm *VendingMachine) executeSideEffects(current State, event Event, next State) {
	switch {
	case current == Paid && next == Dispensed:
		println("出货中...")
	case next == Idle:
		println("重置机器...")
	}
}

// 使用示例
func main() {
	vm := &VendingMachine{currentState: Idle}

	vm.Transition(SelectItem)   // 正常：Idle → Selected
	vm.Transition(InsertCoin)   // 正常：Selected → Paid
	vm.Transition(SelectItem)   // 非法：Paid状态不允许选商品
	vm.Transition(Dispense)     // 正常：Paid → Dispensed
	vm.Transition(ReturnChange) // 正常：Dispensed → Idle
}
