#include	<stdio.h>
#include	<sys/types.h>
#include	<sys/socket.h>
#include	<netinet/in.h>
#include	<string.h>
#include	<arpa/inet.h>
#include	<stdlib.h>
#include	<unistd.h>
#include	<errno.h>
#include	<pwd.h>
#include	<grp.h>
#include	<sys/stat.h>
#include	<sys/ipc.h>
#include	<sys/sem.h>
#include	<sys/msg.h>
#include	<sys/shm.h>
#include 	<fcntl.h>
#include 	<sys/file.h>
#include 	<sys/time.h>
#include	<sys/signal.h>
#include	<signal.h>
#include	<setjmp.h>
#include  	<dirent.h>
#include  	<limits.h>

/*
#include  	<varargs.h>
--del by hqz for LINUX 20060908
*/
#include <stdarg.h> /*add by wangbo*/
#include	<time.h>
#include	<sys/timeb.h>
#include	<netdb.h>

/**********************��������*******************************/
#define iATTACHNUM	10000	/*�������MSP_Attach����	*/
#define iDETACH		20000	/*�������MSP_Detach����	*/
#define iATTACHNAME	11000	/*��������MSP_Attach����	*/
#define iATTACHREQ	12000	/*��������MSP_Attach����	*/
#define iPUTMSG		30000	/*MSP_Put_Msg����*/
#define iGETMSG		40000	/*MSP_Get_Msg����*/
#define iSENDMSG	50000	/*�����ʼ�����*/
#define iDEALSUCC	51000	/*�ʼ�����ɹ�*/
#define iDEALFAIL	52000	/*�ʼ�����ɹ�*/
#define iRESEND		53000	/*�ϵ���������*/
#define iDELMAIL	54000   /*ɾ���ʼ�����*/
#define iLOGIN		60000	/*ǩ������*/
#define iLOGOUT		61000	/*ǩ������*/
#define iMSGSTATUS	64000	/*��Ϣ״̬*/

/********************MSP_KNL������*************************/
#define	MSPSUCCESS		20000   /*�ɹ�*/
#define	MSPFAIL		20001   /*ʧ��*/
#define ERR_FINDMB      20002	/*��������ʧ��*/
#define NOMB		20003   /*�޴�����*/
#define WPROCFAIL       20004   /*д���̿ռ�ʧ��*/
#define ERR_ENGROSS	20005   /*�����ѱ���ռ*/
#define ERR_FINDEPTMB   20006   /*���ҿ�����ռ�ʧ��*/
#define DELPROCFAIL     20007   /*ɾ�����̿ռ�ʧ��*/
#define NOTATTACH	20008   /*ԭ����û��attach*/
#define NOTXGROUP	20009   /*�ǽ�����*/
#define WMBFAIL		20010   /*д����ʧ��*/
#define RMBFAIL		20011   /*������ʧ��*/
#define DELMAILFAIL	20012   /*ɾ���ʼ�ʧ��*/
#define NOMAIL		20013   /*������û���ʼ�*/
#define DELMBFAIL	20014   /*�޴�����*/
#define NOSPACE		20015   /*�����ڴ���û���ʼ�ͷ�ռ�*/
#define LINKMBMEMERR	20016   /*�������乲���ڴ�ʧ��*/
#define NOPROC		20017   /*�޴˽���*/
#define NAMEERR         20018   /*�ļ�������*/
#define PARERR          20019   /*��������*/
#define BUFSMALL        20020   /*��ű��ĵĿռ�̫С*/
#define NOTFILE         20021   /*�ļ������ڻ�û��Ȩ��*/
#define READQUEERR      20022   /*����Ϣ����ʧ��*/
#define WRITEQUEERR     20023   /*����Ϣ����ʧ��*/
#define SHORTSPACE      20024   /*�û�������̫С*/
#define READTIMEOUT	20025   /*�����䳬ʱ*/
#define STRNULL		20026   /*�����䳬ʱ*/
/*********************����**************************/
#define COMPRESSSIZE	8192    /*��Ҫѹ���ı��Ĵ�С*/

#define MAXMBNUM	1000         /*���������*/
#define	LICENSESS	16			
#define MAXPROCNUM	65535        /*��������*/
#define PROCTIMEOUT     3600	     /*�û����̵ĳ�ʱʱ��*/
#define SHARE           '0'	     /*���乲��*/
#define ENGROSS         '1'          /*�����ռ*/
#define UNRESERVE       '0'	     /*�ʼ����󲻱���*/
#define RESERVE         '1'	     /*�ʼ�������*/
#define INSHM		0    	     /*�����ڹ����ڴ���*/
#define INFILE		1	     /*�������ļ���*/
#define CHECKPROC	1	     /*�Ƿ������Attach*/
#define NOCHECKPROC	0	     /*����Ҫ����check*/
#define FILENAMELEN      512         /*�ļ�������*/
#define	REALTIME	1	     /*ʵʱ���ɿ�*/
#define	NOREALTIME	0	     /*�ɿ���ʵʱ*/
#define TSPREALTIME	2	     /*����ƽ̨ר������*/
#define TMPMB		0	     /*��ʱ����*/
#define LONGMB		1	     /*��������*/
#define SENDMAIL 	1
#define CHECKLINK  	2
#define LINKSUCC 	3
#define LINKFAIL 	4
#define KEYSIZE         8            /*��Կ����*/
#define MACSIZE         8            /*MAC����*/
#define	MAXBUFLEN	8192
#define ERROR		-1

#ifndef MSPTRUE
	#define MSPTRUE        1
#endif
#ifndef MSPFALSE
	#define MSPFALSE       0
#endif


#ifndef strmove
#define strmove( s1, s2, nmove ) ( (strncpy(s1,s2,nmove) ), (s1[nmove] = '\0') )
#endif
/************************Msg_Send������*************************/
#define MSND_FILE_ERR		30001   /*  �ļ���д���� */
#define MSND_SHM_ERR 		30002   /*  �����ڴ���� */
#define MSND_MSGQ_ERR 		30003   /* ��Ϣ���г��� */
#define MSND_SOCK_ERR 		30004   /*  Socket���� */
#define MSND_TIMEOUT 		30005   /* ��ʱ����    */
#define MSND_OTHER_ERR 		30006   /* �������� */

/************************Msg_RCV������**********************/
#define MSRV_FILE_ERR 		40001   /*  �ļ���д���� */
#define MSRV_SHM_ERR 		40002   /*  �����ڴ���� */
#define MSRV_MSGQ_ERR 		40003   /* ��Ϣ���г��� */
#define MSRV_SOCK_ERR 		40004   /*  Socket���� */
#define MSRV_TIMEOUT 		40005   /* ��ʱ���� */
#define MSRV_OTHER_ERR 		40006   /* �������� */

/****************************�ͻ������ýṹ****************************/
typedef struct
{
	char	HostIp[20];		/*���������Tcp��ַ*/
	char	HostPort[5+1];		/*���������Port��*/
	char	ErrLogPath[40];		/*�����¼Ŀ¼*/
	char	SysLogPath[40];		/*ϵͳ��־Ŀ¼*/
	char	LogNum[5+1];		/*Log�������*/
	char	FrmSize[8];		/*ÿ֡���ĵĴ�С(���������ݱ���ͷ)*/
}CLIENTCFG ;

/****************************������Ϣ���ýṹ****************************/
typedef struct {
	char	MaxMail[10];		/*MSP��������ʼ���*/
	char	MaxApp[5+1];  		/*�����������*/
	char	TmpMbNum[5+1];		/*����������ʱ������*/
	char    MailShmNum[8];		/*����Mail��Ź����ڴ��С����*/
	char	LogNum[5];		/*Log�������*/
	char	BlockSize[8];		/*��Ĵ�С*/
}SHARECFG;

/****************************�ڵ����ýṹ****************************/
typedef  struct {
	char	NodeName[10];           /*�ڵ���*/
	short	iNodeId;		/*�ڵ��*/
	char	NodeIp[20];		/*�ڵ�ip ��ַ*/
	short	iNodePort;		/*�ڵ�Port��*/
	short	iFrmSize;		/*ÿ֡���ĵĴ�С(���������ݱ���ͷ)*/
    	short   iTimeOut;             	/*��ʱʱ��*/
    	char    NodeList[185];          /*ת���ڵ�id�б�*/
	int	iSocketFd;
}NODECFG;

/****************************�������ýṹ****************************/
struct MAILBOXCFG{
	char	MailBoxId[5+1];		/*���*/
	char	MailBoxName[10];	/*������*/
	char	MbGrpId[5+1];		/*�������*/
	char	MaxMail[5+1];		/*��������ʼ���*/
	char	CpsFlag[1+1];		/*�Ƿ�֧��ѹ��*/
	char	CryptFlag[1+1];		/*�Ƿ�֧�ּ���*/
	char	ShareFlag[1+1];		/*�Ƿ��ռ��־*/
	struct MAILBOXCFG	*pre;
	struct MAILBOXCFG	*next;
};

/****************************����������****************************/
typedef	struct{
	unsigned int   iMaxStoreNum;       /*�ʼ��洢�������ڴ����*/
	unsigned int   iBlockSize;
/*	unsigned int   iCurUsdNum;-del by hqz	�ʼ��洢����ǰռ���ڴ�ռ�*/
	unsigned int   iMaxMbNum;	    /*������������*/
	unsigned int   iMaxMailNum;	    /*�������ʼ���*/
	unsigned short iMaxAttachNum;	    /*������ӽ�����*/
	unsigned short iCurMbNum;	    /*��ǰ������*/
	unsigned int   iCurMailNum;	    /*��ǰ�����ʼ���*/
	unsigned short iCurAttachNum;	    /*��ǰ���ӽ�����*/
	unsigned short iInQueFlag;          /*�����Ϣ������·��б��ı�־*/
	int            iPubShmId;           /*���������������ڴ�Id*/
	int            iPubSemId;           /*�����������źŵ�Id*/
	int            iMailShmId;          /*����ͽڵ㹲���ڴ�ID*/
	int	       iMailSemId;	    /*�����������źŵ�Id*/
	int            iBodySemId;          /*�ʼ��źŵ�Id*/
	int            iToProc;             /*����������ƫ��*/
	int            iToMailHead;         /*�ʼ�ͷ������ƫ��*/
	int            iToMailBody;         /*�ʼ���������ƫ��*/
	int 	       iInQueId;	    /*���������Ϣ����Id*/
	int	       iOutQueId;	    /*���ĳ�����Ϣ����Id*/
	unsigned int   iReqFileSer;         /*�ļ����*/
	unsigned int   iFstEmptyMH;	    /*���ʼ�ͷ��λ��*/
	/* -del by hqz 20060703
	unsigned int   iFstEmptyMB;	    ���ʼ�����λ��*/
	unsigned int   iFstEmptyProc;	    /*�ս�������λ��*/
	unsigned int   iFstUsdProc;	    /*���ý�����λ��*/
	long	       iLogNum;             /*Log�������*/
}Msp_Pub_Tab;
/****************************����������****************************/
typedef struct{
	char	       cMbDscFlag;      /*����ʹ�ñ�־*/
	char	       cPermFlag;	/*�����Ƿ���������:1:���ã�0:��ʱ*/
	unsigned short iMbId;	        /*�����*/
	char	       sMbName[10];     /*������	*/
	unsigned short iMbGrpId;	/*��������*/
	unsigned int   iFstMail;	/*�������һ���ʼ�λ��*/
	unsigned short iProcNum;        /*���������ӽ�����*/
	unsigned int   iMaxMailNum;	/*��������ʼ���*/
	unsigned int   iMailNum;	/*���䵱ǰ�ʼ���*/
	unsigned int   iMailSize;	/*���䵱ǰ�ʼ�ռ�ÿռ�*/
	unsigned int   iRcvMailNum;	/*��������ʼ���*/
	unsigned int   iRcvMailSize;	/*����������ݴ�С*/
	unsigned int   iSndMailNum;	/*���䷢���ʼ���*/
	unsigned int   iSndMailSize;	/*���䷢�����ݴ�С*/
	char           cMbOpenMode;     /*�����ģʽ���������ռ)*/
	char	       cCompressFlag;	/*�ʼ��洢�����Ƿ�ѹ��*/
	char	       cRecryptFlag;	/*�ʼ��洢�����Ƿ����*/
}Msp_Mb_Tab;
/****************************�ʼ�������****************************/
typedef struct {
	unsigned short iSrcMbId;	/*Դ�����*/
	unsigned short iSrcGrpId;	/*Դ��������*/
	unsigned short iDstMbId;	/*Ŀ�������*/
	unsigned short iDstMbGrpId;	/*Ŀ����������*/
	unsigned int   iMailNum;        /*�ʼ���*/
	unsigned short ipriority;	/*�ʼ�Priority*/
	unsigned short iclass;		/*�ʼ�Class*/
	int	       iType;		/*�ʼ�Type*/
	char	       cStorageFlag;	/*�ʼ����λ��*/
	unsigned int   iMailSize;	/*�ʼ���С*/
	unsigned long   iMailBeginPosi;	/*�ʼ������ʼλ��*/
	unsigned int   iMailPosi;	/*���ʼ�λ��*/
	unsigned int   iNextMail;	/*��һ�ʼ�λ��*/
}Msp_Mail_Tab;

/****************************����������****************************/
typedef struct  {
	int	       iSeriNo;		/*���*/
	int	       iUserPid;	/*���̺�*/
	char           sUserName[15+1];	/*�û���*/
	unsigned short iProcBlockFlag;	/*������������־*/
	unsigned short iMbIdPosi;      	/*����λ��*/
	long	       lBeginTime;	/*�����ϴβ���ʱ��*/
	unsigned short iNextPosi;	/*��һ����λ��*/
}Msp_Proc_Tab;

/****************************�ڵ��****************************/
typedef struct {
	char	       sNodeName[7+1];	/*�ڵ���*/
	unsigned short iNodeId;	 	/*�ڵ��*/
	char	       sNodeTcpip[16];	/*�ڵ�IP Address*/
	char	       sNodePort[5+1];	/*�ڵ�TCP port*/
	unsigned short iCurLinkNum;	/*��·��*/
	short          iTimeOut;        /*��ʱʱ��*/
	short	       iFrmNum;		/*���ڰ���֡��*/
	short	       iFrmSize;	/*ÿ֡���ĵĴ�С(���������ݱ���ͷ)*/
    	int            iNodePid;        /*���̺�1*/
}Msp_Node_Tab;

/****************************Saf��Send�Ľӿ�****************************/
typedef struct{
	unsigned short iReqType;	     /*��������*/
	unsigned int   iMailPosi;	     /*�ʼ�λ��*/
	unsigned int   iMailNo;		     /*�ʼ����*/
}Msp_Saf_Tab;

/****************************��Ϣ���ݿ���ͷ****************************/
typedef struct{
	int		iType;		     /*�ʼ�Type*/
	int		iResp_code;          /* ������*/
	unsigned short	iReqType;	     /*��������*/
	unsigned short	iSrcMbId;	     /*Դ�����*/
    	unsigned short  iSrcGrpId;           /*ԭ�������*/
    	unsigned short	iSrcBusId;	     /*Դ����bus��*/
	unsigned short	iDstMbId;	     /*Ŀ�������*/
	unsigned short	iDstMbGrpId;	     /*Ŀ���������*/
	unsigned short	iPriority;	     /*�ʼ����ȼ���*/
	unsigned short	iClass;		     /*�ʼ�Class*/
	unsigned short  iTimeOut;            /*��ʱʱ��*/
	char		sSrcMbName[20];	     /*Դ������*/
	unsigned char	sMacCheck[KEYSIZE];  /*MacУ��*/
	char            cFlag;               /*�Ƿ��ռ��ʵʱ��������*/
}Msp_Msg_Head;
/****************************��Ϣ��****************************/
typedef struct {
	unsigned char FileName[10];
	unsigned int  iMsgLength;
}Msg_Request;
typedef struct {
	unsigned int	iCurDataPosi;	     /*��ǰ�ʼ�λ��*/
	unsigned int    iCurSize;            /*��֡���ݵĴ�С*/
	char		cEndFlag;	     /*�ʼ����ݱ�־*/
	unsigned char	sMacCheck[MACSIZE];  /*MacУ��*/
}Msg_Body;
typedef struct {
	unsigned int  iRetPosi;
}Msg_Return;

/*************************��Ϣ�ṹ*********************************/
typedef struct{
	long	mtype;
	unsigned char	mtext[512];
}msg;

/*****************************ȫ�ֱ���*****************************/
Msp_Pub_Tab  		*pgPubDsc;	            /*�����������ṹ����ָ��*/
unsigned short		igBusId;	            /*��Bus��*/
unsigned char 		*pgMspmsg;	            /*���ݱ��Ľṹָ��*/
unsigned char 		*pgMbHead;	            /*������ʼ��ַ*/
unsigned char 		*pgMailHead;	            /*�ʼ�ͷ��ʼ��ַ*/
unsigned char 		*pgProcHead;	            /*������ʼ��ַ*/
unsigned char 		*pgMailBody;	            /*��������ʼ��ַ*/
unsigned char 		*pgNodeHead;	            /*�����ʼ��ַ*/
char    		cgShutdownFlag;               /*�رձ�־*/

SHARECFG        	shacfg;
struct MAILBOXCFG 	mailboxcfg;

unsigned int		igBlockSize;
unsigned short		igNodeId;
unsigned short		igSrcMbId,igSrcGrpId,igSrcMailBoxPosi,igProcPosi;
pid_t			igUsrPid;
msg			mspmsg;
unsigned int		igSendMsgLen;
char 			agDebugfile[256];
char			sgMspPath[256];

int Msp_Cls_Child
(
	int			/*  �׽���������		*/
);

/* add by yl 20060704 */
/*****************************��������*****************************/
short Msp_attach(
  unsigned short *iSrcMbId, 
  char 		*MailBoxName, 
  unsigned short iSrcGrpId,
  unsigned short iBusId,
  unsigned short iAttachType
);

short Msp_Put_Msg(
char 		*sMsg, 
int		iMsgSize,
char            cPutType,	
unsigned short  iDstMbId, 
unsigned short  iDstMbGrpId, 
unsigned short	iPriority,
unsigned short	iClass,
int		iType,
char            cStoreFlag
);

short Msp_Get_Gate(
char 		*sMsg, 
unsigned int	*iUsrLen,	
unsigned short  *iSrcMbId, 
unsigned short  *iSrcGrpId, 
unsigned short	*iPriority,
unsigned short	*iClass,
int		*iType,
char		*cStoreFlag,
char            cReserve,
short  		iTimeOut,
Msp_Msg_Head	*stpHead
);

short Msp_Get_Msg
(
char 		*sMsg, 
unsigned int	*iUsrLen,	
unsigned short  *iSrcMbId, 
unsigned short  *iSrcGrpId, 
unsigned short	*iPriority,
unsigned short	*iClass,
int		*iType,
char		*cStoreFlag,
char            cReserve,
short  		iTimeOut
);

short Msp_Detach();

unsigned short ReadMailHead
(
	unsigned int 	*ilFstMail, 
	Msp_Msg_Head	*stpHead,
	int 		*iUsrLen,
	unsigned int	*ilFstBkId,
	char		*slFileName,
	char		*cStorageFlag,
	char		cReserve
);

unsigned short DelMailHead(unsigned int iPosi);

unsigned short ReadMemList
( 
	char 		*msg,
	int 		msgsize, 
	unsigned int 	iFstbkid,
	char		cStorageFlag
);

unsigned short ReadFileList
( 
	 char 	*msg,
	int 	msgsize, 
	char	*sFileName,
	char	cStorageFlag
);

unsigned short WriteMemList
( 
	char *msg, 
	int msgsize, 
	unsigned int *iFstbkid
);

unsigned short DeleteMemList( unsigned int iFstbkid);

unsigned short	BlockRead(unsigned short ilProcPosi,short iTimeOut);

short Msp_attach_Gate(
	int		*iSocketFd,
	short		iFrmSize,
	char		*NodeIp,
	short		iNodePort,
	unsigned short *iSrcMbId, 
	unsigned short iSrcGrpId,
	unsigned short iBusId
);

short Msp_Put_Gate(
int		iSocketFd,
short		iFrmSize,
Msp_Msg_Head	*stpHead,
char 		*sMsg, 
int		iMsgSize,
char            cStoreFlag
);

void Msp_Init(SHARECFG *shacfg, struct MAILBOXCFG *mailboxcfg);

unsigned short CreateMbShm(key_t key, int iMbShmSize,int iMailNum);

unsigned short  CreateIPC( );

unsigned short DeleteMbSpace();

unsigned short DeleteIPC();

void LoadOldMail();

unsigned short 	ReadFileSpace(char *sFilename, Msp_Mail_Tab *pTmp);

int MspCheckSocket(int sock, int tosec, int tousec); 

int SendSocketReturn( int sockfd, Msg_Return *plReturn);

int RecvSocketReturn( int sockfd, Msg_Return *plReturn);

int RecvHead
(
        int             sockfd,         /*  �׽���������        */
        Msp_Msg_Head    *plMsgHead
);

int RecvSocketRequest(int sockfd,Msg_Request *plRequest);

int SendSocketRequest(int sockfd,Msg_Request *plRequest);

int RecvSocketBuf(int sockfd,unsigned char *strbuf,Msg_Body *plbody);

int SendSocketBuf(int sockfd,unsigned char *strbuf);

int SendHead(int sockfd,Msp_Msg_Head *plMsgHead);

short Msg_Put(
int		ilsocketfd,
char 		*sMsg, 
int		iMsgSize,
unsigned short	ilFrmSize,
char            cStoreFlag
);

short	Msg_Recv
(
	int		ilsocketfd,
	unsigned char	*plCommArea,
	unsigned char 	*cStorageFlag,
	unsigned int	*iMsgSize
);

void t_catch();

void *getshm(void);

int  MSPGetTime(char *date);

int MSPP(int ilsemid);

int MSPV(int ilsemid);

void Init();

int TakeCount(char Flag);

int Msp_CliCfg_Load(char *filename, CLIENTCFG  *clientcfg);

int Msp_SrvCfg_Load(
char       *filename,
SHARECFG   *sharecfg,
struct MAILBOXCFG *mailboxcfg
);

int Msp_Item_Get(FILE *plFile, char *strlItemName, char *strlItemVar);

int Msp_Str_Trim(char *pstrlBuf);

long FileLen
(
char		*ptrFile		/*  �ļ�ȫ·����	*/
);

int FileExist
(
char		*ptrFile		/*  �ļ�ȫ·����	*/
);

short CompressData(unsigned char *InBuf,
unsigned char *OutBuf, short InBufLength);

void pack_init();

short lookup_ct( short code, unsigned char thischar);

short putcode( short code,unsigned char * buf, short bufPosition);

short UnCompressData(unsigned char *InBuf, 
unsigned char *OutBuf,short InBufLength);

void unpack_init();

void insert_ct( short code, short oldcode);

short putx(short code,unsigned char *buf,short bufPosition);


void CreatMAC
(
void	*pData,
int	DataSize,
void	*pKey,
int	KeySize,
void	*pMAC,
int	MACSize
);

int CheckMAC
(
void	*pData,
int	DataSize,
void	*pKey,
int	KeySize,
void	*pMAC,
int	MACSize
);

void * CreatKey
(
void	*pKey,
int		KeySize	
);

int BlockAllSig
(
sigset_t	*pOldSet	/*  ԭ���ź�����״̬	*/
);

int UnblockAllSig
(
sigset_t	*pOldSet	/*  ԭ���ź�����״̬	*/
);

int BlockSig
(
int			Sig			/*  �ź�	*/
);

int UnblockSig
(
int			Sig			/*  �ź�	*/
);

int ProcSig
(
int			Sig			/*  �ź�	*/
);

int SigDirect
(
int			Sig,		/*  �ź�			*/
void(* func)(int),		/*  �źŴ�����	*/
struct sigaction	*poldact
);

int SigDefault
(
int		Sig				/*  �ź�	*/
);

int TestSig
(
int		Sig		/*  �ź�	*/
);

int ClearTimer(void);

int SetTimer
(
long	sec,			/*  ��			*/
long	usec,			/*  ΢��		*/
char	cCycFlag		/*  ѭ����־	*/
);

int ReadMailDsc
(
char		*ptrFile,		/*  �ļ�ȫ·����	*/
Msp_Mail_Tab	*pMailDsc		/*  �ʼ�������ַ	*/
);

void ClearMsg
(
int		MsgQueID,	/*  ��Ϣ���б�ʶ��			*/
long	lType		/*  ��Ҫ�����Ϣ��Typeֵ	*/
);

unsigned short ReqFileSer(unsigned int *iFileNum);

unsigned short MSPSendMsg(int iMsgId,  msg  *SendMsg, int iLen);

unsigned short MSPRcvMsg(int iMsgId,msg *RcvMsg, int *iLen,int iTimeOut);

int  GetTime_mm(char *date);

void daemon_init(void);

void AddSubProc(unsigned short ilPosi,unsigned short ilFlag);

unsigned short CreateShm(
	key_t shmkey, 
	int *shmid, 
	unsigned int shmsize 
);

unsigned short DeleteShm(int shmid);

unsigned short CreatePubDscSpace(key_t iShmKey);

unsigned short DeletePubDscSpace();

unsigned short DeleteMsg(int msgid);

unsigned short CreateMsg(key_t msgkey, int *msgid);

unsigned short DeleteSem(int semid);

unsigned short CreateSem(key_t semkey,int *semid);

unsigned short FindMbBynum(
	 unsigned short mbid,
	unsigned short mbgrpid,
	 unsigned short *mbdsc_id,
	char	cPutType
) ;

unsigned short FindMbByname(
	char * mbname, 
	unsigned short *mbdsc_id ,
	char cPutType
);

unsigned short FindEmptyMb(unsigned short *mbdsc_id);

unsigned short WriteMbDsc(struct MAILBOXCFG * mailbox,unsigned short *mbdscid);

unsigned short WriteProcdscSpace
(
	char *pidname, 
	int pid,
	unsigned short ilMbPosi,
	unsigned short *ProcPosi
);

unsigned short UpdateProc
(
	unsigned short iProcPosi, 
	int pid
);

unsigned short DeleteProcdsc(int pid);

int ProcessMessage(short iMbId);

