package gogp

type gogpProcessStep int

func (me gogpProcessStep) IsReverse() bool {
	return me >= gogpStepREQUIRE && me <= gogpStepREVERSE
}

func (me gogpProcessStep) String() (s string) {
	switch me {
	case gogpStepREQUIRE:
		s = "Step=[1RequireReplace]"
	case gogpStepREVERSE:
		s = "Step=[2ReverseWork]"
	case gogpStepPRODUCE:
		s = "Step=[3NormalProduce]"
	default:
		s = "Step=Unknown"
	}
	return
}

const (
	gogpStepREQUIRE gogpProcessStep = iota + 1 // require replace in fake go file
	gogpStepREVERSE                            // gen gp file from fake go file
	gogpStepPRODUCE                            // gen go file from gp file
)
