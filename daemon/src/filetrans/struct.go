package filetrans

type Command struct {
	Command string // 发送的命令
	Args    string // 所带的Args.
}

type localServerFile struct {
	Name    string
	Mode    string
	IsDir   bool
	Size    int64
	ModTIme int64
}
