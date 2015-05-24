package sftpd

const (
	ssh_FXP_INIT           = 1
	ssh_FXP_VERSION        = 2
	ssh_FXP_OPEN           = 3
	ssh_FXP_CLOSE          = 4
	ssh_FXP_READ           = 5
	ssh_FXP_WRITE          = 6
	ssh_FXP_LSTAT          = 7
	ssh_FXP_FSTAT          = 8
	ssh_FXP_SETSTAT        = 9
	ssh_FXP_FSETSTAT       = 10
	ssh_FXP_OPENDIR        = 11
	ssh_FXP_READDIR        = 12
	ssh_FXP_REMOVE         = 13
	ssh_FXP_MKDIR          = 14
	ssh_FXP_RMDIR          = 15
	ssh_FXP_REALPATH       = 16
	ssh_FXP_STAT           = 17
	ssh_FXP_RENAME         = 18
	ssh_FXP_READLINK       = 19
	ssh_FXP_SYMLINK        = 20
	ssh_FXP_STATUS         = 101
	ssh_FXP_HANDLE         = 102
	ssh_FXP_DATA           = 103
	ssh_FXP_NAME           = 104
	ssh_FXP_ATTRS          = 105
	ssh_FXP_EXTENDED       = 200
	ssh_FXP_EXTENDED_REPLY = 201
)

const (
	ssh_FX_OK                = 0
	ssh_FX_EOF               = 1
	ssh_FX_NO_SUCH_FILE      = 2
	ssh_FX_PERMISSION_DENIED = 3
	ssh_FX_FAILURE           = 4
	ssh_FX_BAD_MESSAGE       = 5
	ssh_FX_NO_CONNECTION     = 6
	ssh_FX_CONNECTION_LOST   = 7
	ssh_FX_OP_UNSUPPORTED    = 8
)

const (
	ssh_FILEXFER_ATTR_SIZE        = 0x00000001
	ssh_FILEXFER_ATTR_UIDGID      = 0x00000002
	ssh_FILEXFER_ATTR_PERMISSIONS = 0x00000004
	ssh_FILEXFER_ATTR_ACMODTIME   = 0x00000008
	ssh_FILEXFER_ATTR_EXTENDED    = 0x80000000
)
