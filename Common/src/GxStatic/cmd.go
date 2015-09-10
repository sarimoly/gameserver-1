package GxStatic

const (
	//login
	CmdRegister = 1000
	CmdLogin    = 1001

	//gate
	CmdHeartbeat   = 2000
	CmdGetRoleList = 2001
	CmdSelectRole  = 2002
	CmdCreateRole  = 2003

	//server
	CmdServerConnectGate = 3000
)

type ServerStatus int

const (
	ServerStatusHot ServerStatus = iota
	ServerStatusNew              //1
	ServerStatusMaintain
)
