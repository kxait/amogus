package pvm

/*
*	Data packing styles for pvm_initsend()
 */

type DataPackingStyle int64

const (
	DataDefault DataPackingStyle = 0
	DataRaw     DataPackingStyle = 2
	DataInPlace DataPackingStyle = 2
	DataFoo     DataPackingStyle = DataDefault
	DataTrace   DataPackingStyle = 4
)

/*
*	pvm_spawn options
 */

type SpawnOptions int64

const (
	TaskDefault SpawnOptions = 0
	TaskHost    SpawnOptions = 1 /* specify host */
	TaskArch    SpawnOptions = 2 /* specify architecture */
	TaskDebug   SpawnOptions = 4 /* start task in debugger */
	TaskTrace   SpawnOptions = 8 /* process generates trace data */
	/* for MPP ports */
	MppFront  SpawnOptions = 16 /* spawn task on service node */
	HostCompl SpawnOptions = 32 /* complement host set */
	/* for parent-less spawning */
	NoSpawnParent SpawnOptions = 64
)

/*
*	pvm_notify kinds
 */
type NotifyKind int64

const (
	TaskExit    NotifyKind = 1 /* on task exit */
	HostDelete  NotifyKind = 2 /* on host fail/delete */
	HostAdd     NotifyKind = 3 /* on host startup */
	RouteAdd    NotifyKind = 4 /* new task-task route opened */
	RouteDelete NotifyKind = 5 /* task-task route closed */

	/* flags combined with notify kind */

	NotifyCancel NotifyKind = 256 /* cancel (complete immediately) notifies */
)

/*
*	for pvm_setopt and pvm_getopt
 */

type LibpvmOption int64

const (
	Route             LibpvmOption = 1  /* routing policy */
	DontRoute         LibpvmOption = 1  /* don't allow direct task-task links */
	AllowDirect       LibpvmOption = 2  /* allow direct links, but don't request */
	RouteDirect       LibpvmOption = 3  /* request direct links */
	DebugMask         LibpvmOption = 2  /* debugmask */
	AutoErr           LibpvmOption = 3  /* auto error reporting */
	OutputTid         LibpvmOption = 4  /* stdout destination for children */
	OutputCode        LibpvmOption = 5  /* stdout message tag */
	TraceTid          LibpvmOption = 6  /* trace destination for children */
	TraceCode         LibpvmOption = 7  /* trace message tag */
	TraceBuffer       LibpvmOption = 8  /* trace buffering for children */
	TraceOptions      LibpvmOption = 9  /* trace options for children */
	TraceFull         LibpvmOption = 1  /* do full trace events */
	TraceTime         LibpvmOption = 2  /* only do PVM routine timings */
	TraceCount        LibpvmOption = 3  /* only do PVM routine profiling */
	FragSize          LibpvmOption = 10 /* message fragment size */
	ResvTids          LibpvmOption = 11 /* allow reserved message tids and codes */
	SelfOutputTid     LibpvmOption = 12 /* stdout destination for task */
	SelfOutputCode    LibpvmOption = 13 /* stdout message tag */
	SelfTraceTid      LibpvmOption = 14 /* trace destination for task */
	SelfTraceCode     LibpvmOption = 15 /* trace message tag */
	SelfTraceBuffer   LibpvmOption = 16 /* trace buffering for task */
	SelfTraceOptions  LibpvmOption = 17 /* trace options for task */
	ShowTids          LibpvmOption = 18 /* pvm_catchout prints task ids with output */
	PollType          LibpvmOption = 19 /* shared memory wait method */
	PollConstant      LibpvmOption = 1
	PollSleep         LibpvmOption = 2
	PollTime          LibpvmOption = 20 /* time before sleep if PvmPollSleep */
	OutputContext     LibpvmOption = 21 /* stdout message context */
	TraceContext      LibpvmOption = 22 /* trace message context */
	SelfOutputContext LibpvmOption = 23 /* stdout message context */
	SelfTraceContext  LibpvmOption = 24 /* trace message context */
	NoReset           LibpvmOption = 25 /* do not kill task on reset */
)

/*
*	for pvm_[sg]ettmask
 */

type TraceMask int64

const (
	TaskSelf  TraceMask = 0 /* this task */
	TaskChild TraceMask = 1 /* (future) child tasks */
)

/*
*	Need to have PvmBaseContext defined
 */

const BaseContext = 0

/*
*	for message mailbox operations: pvm_putinfo and pvm_recvinfo
 */

type MailboxFlag int64

const (
	MboxDefault MailboxFlag = 0 /* put: single locked instance */
	/* recv: 1st entry */
	/* start w/index=0 */
	MboxPersistent    MailboxFlag = 1  /* entry remains after owner exit */
	MboxMultiInstance MailboxFlag = 2  /* multiple entries in class */
	MboxOverWritable  MailboxFlag = 4  /* can write over this entry */
	MboxFirstAvail    MailboxFlag = 8  /* select 1st index >= specified */
	MboxReadAndDelete MailboxFlag = 16 /* atomic read / delete */
	/* requires read & delete rights */
	MboxWaitForInfo MailboxFlag = 32  /* for blocking recvinfo */
	MboxMaxFlag     MailboxFlag = 512 /* maximum mbox flag bit value */

	MboxDirectIndexShift MailboxFlag = 10 /* log2(PvmMboxMaxFlag) + 1 */
)

/*
*	pre-defined system message mailbox classes
 */

type MailboxClass string

const (
	NORESETCLASS MailboxClass = "###_PVM_NO_RESET_###"

	HOSTERCLASS MailboxClass = "###_PVM_HOSTER_###"

	TASKERCLASS MailboxClass = "###_PVM_TASKER_###"

	TRACERCLASS MailboxClass = "###_PVM_TRACER_###"

	RMCLASS MailboxClass = "###_PVM_RM_###"
)

/*
*	Libpvm error codes
 */

type ErrorCode int64

const (
	Ok           ErrorCode = 0   /* Success */
	BadParam     ErrorCode = -2  /* Bad parameter */
	Mismatch     ErrorCode = -3  /* Parameter mismatch */
	Overflow     ErrorCode = -4  /* Value too large */
	NoData       ErrorCode = -5  /* End of buffer */
	NoHost       ErrorCode = -6  /* No such host */
	NoFile       ErrorCode = -7  /* No such file */
	Denied       ErrorCode = -8  /* Permission denied */
	NoMem        ErrorCode = -10 /* Malloc failed */
	BadMsg       ErrorCode = -12 /* Can't decode message */
	SysErr       ErrorCode = -14 /* Can't contact local daemon */
	NoBuf        ErrorCode = -15 /* No current buffer */
	NoSuchBuf    ErrorCode = -16 /* No such buffer */
	NullGroup    ErrorCode = -17 /* Null group name */
	DupGroup     ErrorCode = -18 /* Already in group */
	NoGroup      ErrorCode = -19 /* No such group */
	NotInGroup   ErrorCode = -20 /* Not in group */
	NoInst       ErrorCode = -21 /* No such instance */
	HostFail     ErrorCode = -22 /* Host failed */
	NoParent     ErrorCode = -23 /* No parent task */
	NotImpl      ErrorCode = -24 /* Not implemented */
	DSysErr      ErrorCode = -25 /* Pvmd system error */
	BadVersion   ErrorCode = -26 /* Version mismatch */
	OutOfRes     ErrorCode = -27 /* Out of resources */
	DupHost      ErrorCode = -28 /* Duplicate host */
	CantStart    ErrorCode = -29 /* Can't start pvmd */
	Already      ErrorCode = -30 /* Already in progress */
	NoTask       ErrorCode = -31 /* No such task */
	NotFound     ErrorCode = -32 /* Not Found */
	Exists       ErrorCode = -33 /* Already exists */
	HostrNMstr   ErrorCode = -34 /* Hoster run on non-master host */
	ParentNotSet ErrorCode = -35 /* Spawning parent set PvmNoSpawnParent */
	IPLoopback   ErrorCode = -36 /* Master Host's IP is Loopback */
)

/*
*	crusty error constants from 3.3, now redefined...
 */
const (
	NoEntry  ErrorCode = NotFound /* No such entry */
	DupEntry ErrorCode = Denied   /* Duplicate entry */
)

/*
*	Data types for pvm_reduce(), pvm_psend(), pvm_precv()
 */

type DataType int64

const (
	STR    DataType = 0  /* string */
	BYTE   DataType = 1  /* byte */
	SHORT  DataType = 2  /* short */
	INT    DataType = 3  /* int */
	FLOAT  DataType = 4  /* real */
	CPLX   DataType = 5  /* complex */
	DOUBLE DataType = 6  /* double */
	DCPLX  DataType = 7  /* double complex */
	LONG   DataType = 8  /* long integer */
	USHORT DataType = 9  /* unsigned short int */
	UINT   DataType = 10 /* unsigned int */
	ULONG  DataType = 11 /* unsigned long int */
)
