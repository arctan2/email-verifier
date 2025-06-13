package verifier

type verifierManager struct {
	running map[int64]*Verifier
}

func (vm *verifierManager) Add(fileId int64, v *Verifier) {
	vm.running[fileId] = v
}

func (vm *verifierManager) Remove(fileId int64) {
	delete(vm.running, fileId)
}

func (vm *verifierManager) Get(fileId int64) (*Verifier) {
	v, ok := vm.running[fileId]
	if !ok {
		return nil
	}
	return v
}

var VerifierManager verifierManager = verifierManager{running: make(map[int64]*Verifier)}
