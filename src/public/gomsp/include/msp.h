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

/**********************ÇëÇóÀàĞÍ*******************************/
#define iATTACHNUM	10000	/*°´ÓÊÏäºÅMSP_AttachÇëÇó	*/
#define iDETACH		20000	/*°´ÓÊÏäºÅMSP_DetachÇëÇó	*/
#define iATTACHNAME	11000	/*°´ÓÊÏäÃûMSP_AttachÇëÇó	*/
#define iATTACHREQ	12000	/*ÉêÇëÓÊÏäMSP_AttachÇëÇó	*/
#define iPUTMSG		30000	/*MSP_Put_MsgÇëÇó*/
#define iGETMSG		40000	/*MSP_Get_MsgÇëÇó*/
#define iSENDMSG	50000	/*·¢ËÍÓÊ¼şÇëÇó*/
#define iDEALSUCC	51000	/*ÓÊ¼ş´¦Àí³É¹¦*/
#define iDEALFAIL	52000	/*ÓÊ¼ş´¦Àí³É¹¦*/
#define iRESEND		53000	/*¶ÏµãĞø´«ÇëÇó*/
#define iDELMAIL	54000   /*É¾³ıÓÊ¼şÇëÇó*/
#define iLOGIN		60000	/*Ç©µ½ÇëÇó*/
#define iLOGOUT		61000	/*Ç©ÍËÇëÇó*/
#define iMSGSTATUS	64000	/*ÏûÏ¢×´Ì¬*/

/********************MSP_KNL·µ»ØÂë*************************/
#define	MSPSUCCESS		20000   /*³É¹¦*/
#define	MSPFAIL		20001   /*Ê§°Ü*/
#define ERR_FINDMB      20002	/*²éÕÒÓÊÏäÊ§°Ü*/
#define NOMB		20003   /*ÎŞ´ËÓÊÏä*/
#define WPROCFAIL       20004   /*Ğ´½ø³Ì¿Õ¼äÊ§°Ü*/
#define ERR_ENGROSS	20005   /*ÓÊÏäÒÑ±»¶ÀÕ¼*/
#define ERR_FINDEPTMB   20006   /*²éÕÒ¿ÕÓÊÏä¿Õ¼äÊ§°Ü*/
#define DELPROCFAIL     20007   /*É¾³ı½ø³Ì¿Õ¼äÊ§°Ü*/
#define NOTATTACH	20008   /*Ô­ÓÊÏäÃ»ÓĞattach*/
#define NOTXGROUP	20009   /*·Ç½»²æ×é*/
#define WMBFAIL		20010   /*Ğ´ÓÊÏäÊ§°Ü*/
#define RMBFAIL		20011   /*¶ÁÓÊÏäÊ§°Ü*/
#define DELMAILFAIL	20012   /*É¾³ıÓÊ¼şÊ§°Ü*/
#define NOMAIL		20013   /*ÓÊÏäÖĞÃ»ÓĞÓÊ¼ş*/
#define DELMBFAIL	20014   /*ÎŞ´ËÓÊÏä*/
#define NOSPACE		20015   /*¹²ÏíÄÚ´æÖĞÃ»ÓĞÓÊ¼şÍ·¿Õ¼ä*/
#define LINKMBMEMERR	20016   /*Á¬½ÓÓÊÏä¹²ÏíÄÚ´æÊ§°Ü*/
#define NOPROC		20017   /*ÎŞ´Ë½ø³Ì*/
#define NAMEERR         20018   /*ÎÄ¼şÃû´íÎó*/
#define PARERR          20019   /*²ÎÊı´íÎó*/
#define BUFSMALL        20020   /*´æ·Å±¨ÎÄµÄ¿Õ¼äÌ«Ğ¡*/
#define NOTFILE         20021   /*ÎÄ¼ş²»´æÔÚ»òÃ»ÓĞÈ¨ÏŞ*/
#define READQUEERR      20022   /*¶ÁÏûÏ¢¶ÓÁĞÊ§°Ü*/
#define WRITEQUEERR     20023   /*¶ÁÏûÏ¢¶ÓÁĞÊ§°Ü*/
#define SHORTSPACE      20024   /*ÓÃ»§»º´æÇøÌ«Ğ¡*/
#define READTIMEOUT	20025   /*¶ÁÓÊÏä³¬Ê±*/
#define STRNULL		20026   /*¶ÁÓÊÏä³¬Ê±*/
/*********************ÆäËü**************************/
#define COMPRESSSIZE	8192    /*ĞèÒªÑ¹ËõµÄ±¨ÎÄ´óĞ¡*/

#define MAXMBNUM	1000         /*×î´óÓÊÏäÊı*/
#define	LICENSESS	16			
#define MAXPROCNUM	65535        /*×î´ó½ø³ÌÊı*/
#define PROCTIMEOUT     3600	     /*ÓÃ»§½ø³ÌµÄ³¬Ê±Ê±¼ä*/
#define SHARE           '0'	     /*ÓÊÏä¹²Ïí*/
#define ENGROSS         '1'          /*ÓÊÏä¶ÀÕ¼*/
#define UNRESERVE       '0'	     /*ÓÊ¼ş¶Áºó²»±£Áô*/
#define RESERVE         '1'	     /*ÓÊ¼ş¶Áºó±£Áô*/
#define INSHM		0    	     /*±¨ÎÄÔÚ¹²ÏíÄÚ´æÖĞ*/
#define INFILE		1	     /*±¨ÎÄÔÚÎÄ¼şÖĞ*/
#define CHECKPROC	1	     /*ÊÇ·ñ¼ì²é½ø³ÌAttach*/
#define NOCHECKPROC	0	     /*²»ĞèÒª½ø³Ìcheck*/
#define FILENAMELEN      512         /*ÎÄ¼şÃû³¤¶È*/
#define	REALTIME	1	     /*ÊµÊ±²»¿É¿¿*/
#define	NOREALTIME	0	     /*¿É¿¿·ÇÊµÊ±*/
#define TSPREALTIME	2	     /*½»»»Æ½Ì¨×¨ÓÃÀàĞÍ*/
#define TMPMB		0	     /*ÁÙÊ±ÓÊÏä*/
#define LONGMB		1	     /*ÓÀ¾ÃÓÊÏä*/
#define SENDMAIL 	1
#define CHECKLINK  	2
#define LINKSUCC 	3
#define LINKFAIL 	4
#define KEYSIZE         8            /*ÃÜÔ¿³¤¶È*/
#define MACSIZE         8            /*MAC³¤¶È*/
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
/************************Msg_Send´íÎóÂë*************************/
#define MSND_FILE_ERR		30001   /*  ÎÄ¼ş¶ÁĞ´³ö´í */
#define MSND_SHM_ERR 		30002   /*  ¹²ÏíÄÚ´æ³ö´í */
#define MSND_MSGQ_ERR 		30003   /* ÏûÏ¢¶ÓÁĞ³ö´í */
#define MSND_SOCK_ERR 		30004   /*  Socket³ö´í */
#define MSND_TIMEOUT 		30005   /* ³¬Ê±³ö´í    */
#define MSND_OTHER_ERR 		30006   /* ÆäËü´íÎó */

/************************Msg_RCV´íÎóÂë**********************/
#define MSRV_FILE_ERR 		40001   /*  ÎÄ¼ş¶ÁĞ´³ö´í */
#define MSRV_SHM_ERR 		40002   /*  ¹²ÏíÄÚ´æ³ö´í */
#define MSRV_MSGQ_ERR 		40003   /* ÏûÏ¢¶ÓÁĞ³ö´í */
#define MSRV_SOCK_ERR 		40004   /*  Socket³ö´í */
#define MSRV_TIMEOUT 		40005   /* ³¬Ê±³ö´í */
#define MSRV_OTHER_ERR 		40006   /* ÆäËü´íÎó */

/****************************¿Í»§¶ËÅäÖÃ½á¹¹****************************/
typedef struct
{
	char	HostIp[20];		/*ÓÊÏä·şÎñÆ÷TcpµØÖ·*/
	char	HostPort[5+1];		/*ÓÊÏä·şÎñÆ÷Port¿Ú*/
	char	ErrLogPath[40];		/*´íÎó¼ÇÂ¼Ä¿Â¼*/
	char	SysLogPath[40];		/*ÏµÍ³ÈÕÖ¾Ä¿Â¼*/
	char	LogNum[5+1];		/*Log×î´óĞĞÊı*/
	char	FrmSize[8];		/*Ã¿Ö¡±¨ÎÄµÄ´óĞ¡(²»°üº¬Êı¾İ±¨ÎÄÍ·)*/
}CLIENTCFG ;

/****************************¹«¹²ĞÅÏ¢ÅäÖÃ½á¹¹****************************/
typedef struct {
	char	MaxMail[10];		/*MSP¹ÜÀí×î´óÓÊ¼şÊı*/
	char	MaxApp[5+1];  		/*×î´ó¹ÜÀí½ø³ÌÊı*/
	char	TmpMbNum[5+1];		/*×î´ó¿ÉÉêÇëÁÙÊ±ÓÊÏäÊı*/
	char    MailShmNum[8];		/*Õû¸öMail´æ·Å¹²ÏíÄÚ´æ´óĞ¡¿éÊı*/
	char	LogNum[5];		/*Log×î´óĞĞÊı*/
	char	BlockSize[8];		/*¿éµÄ´óĞ¡*/
}SHARECFG;

/****************************½ÚµãÅäÖÃ½á¹¹****************************/
typedef  struct {
	char	NodeName[10];           /*½ÚµãÃû*/
	short	iNodeId;		/*½ÚµãºÅ*/
	char	NodeIp[20];		/*½Úµãip µØÖ·*/
	short	iNodePort;		/*½ÚµãPort¿Ú*/
	short	iFrmSize;		/*Ã¿Ö¡±¨ÎÄµÄ´óĞ¡(²»°üº¬Êı¾İ±¨ÎÄÍ·)*/
    	short   iTimeOut;             	/*³¬Ê±Ê±¼ä*/
    	char    NodeList[185];          /*×ª·¢½ÚµãidÁĞ±í*/
	int	iSocketFd;
}NODECFG;

/****************************ÓÊÏäÅäÖÃ½á¹¹****************************/
struct MAILBOXCFG{
	char	MailBoxId[5+1];		/*ÏäºÅ*/
	char	MailBoxName[10];	/*ÓÊÏäÃû*/
	char	MbGrpId[5+1];		/*ËùÊô×éºÅ*/
	char	MaxMail[5+1];		/*ÓÊÏä×î´óÓÊ¼şÊı*/
	char	CpsFlag[1+1];		/*ÊÇ·ñÖ§³ÖÑ¹Ëõ*/
	char	CryptFlag[1+1];		/*ÊÇ·ñÖ§³Ö¼ÓÃÜ*/
	char	ShareFlag[1+1];		/*ÊÇ·ñ¶ÀÕ¼±êÖ¾*/
	struct MAILBOXCFG	*pre;
	struct MAILBOXCFG	*next;
};

/****************************¹«¹²ÃèÊöÇø****************************/
typedef	struct{
	unsigned int   iMaxStoreNum;       /*ÓÊ¼ş´æ´¢Çø¹²ÏíÄÚ´æ¿éÊı*/
	unsigned int   iBlockSize;
/*	unsigned int   iCurUsdNum;-del by hqz	ÓÊ¼ş´æ´¢Çøµ±Ç°Õ¼ÓÃÄÚ´æ¿Õ¼ä*/
	unsigned int   iMaxMbNum;	    /*×î´ó¹ÜÀíÓÊÏäÊı*/
	unsigned int   iMaxMailNum;	    /*×î´ó¹ÜÀíÓÊ¼şÊı*/
	unsigned short iMaxAttachNum;	    /*×î´óÁ¬½Ó½ø³ÌÊı*/
	unsigned short iCurMbNum;	    /*µ±Ç°ÓÊÏäÊı*/
	unsigned int   iCurMailNum;	    /*µ±Ç°¹ÜÀíÓÊ¼şÊı*/
	unsigned short iCurAttachNum;	    /*µ±Ç°Áª½Ó½ø³ÌÊı*/
	unsigned short iInQueFlag;          /*Èë¿ÚÏûÏ¢¶ÓÁĞÊÇÂ·ñÓĞ±¨ÎÄ±êÖ¾*/
	int            iPubShmId;           /*¹«¹²ÃèÊöÇø¹²ÏíÄÚ´æId*/
	int            iPubSemId;           /*¹«¹²ÃèÊöÇøĞÅºÅµÆId*/
	int            iMailShmId;          /*ÓÊÏäºÍ½Úµã¹²ÏíÄÚ´æID*/
	int	       iMailSemId;	    /*ÓÊÏäÃèÊöÇøĞÅºÅµÆId*/
	int            iBodySemId;          /*ÓÊ¼şĞÅºÅµÆId*/
	int            iToProc;             /*½ø³ÌÃèÊöÇøÆ«ÒÆ*/
	int            iToMailHead;         /*ÓÊ¼şÍ·ÃèÊöÇøÆ«ÒÆ*/
	int            iToMailBody;         /*ÓÊ¼şÌåÃèÊöÇøÆ«ÒÆ*/
	int 	       iInQueId;	    /*ºËĞÄÈë¿ÚÏûÏ¢¶ÓÁĞId*/
	int	       iOutQueId;	    /*ºËĞÄ³ö¿ÚÏûÏ¢¶ÓÁĞId*/
	unsigned int   iReqFileSer;         /*ÎÄ¼şĞòºÅ*/
	unsigned int   iFstEmptyMH;	    /*¿ÕÓÊ¼şÍ·Ê×Î»ÖÃ*/
	/* -del by hqz 20060703
	unsigned int   iFstEmptyMB;	    ¿ÕÓÊ¼şÌåÊ×Î»ÖÃ*/
	unsigned int   iFstEmptyProc;	    /*¿Õ½ø³ÌÁ´Ê×Î»ÖÃ*/
	unsigned int   iFstUsdProc;	    /*ÒÑÓÃ½ø³ÌÊ×Î»ÖÃ*/
	long	       iLogNum;             /*Log×î´óĞĞÊı*/
}Msp_Pub_Tab;
/****************************ÓÊÏäÃèÊöÇø****************************/
typedef struct{
	char	       cMbDscFlag;      /*ÓÊÏäÊ¹ÓÃ±êÖ¾*/
	char	       cPermFlag;	/*ÓÊÏäÊÇ·ñÓÀ¾ÃÓÊÏä:1:ÓÀ¾Ã£»0:ÁÙÊ±*/
	unsigned short iMbId;	        /*ÓÊÏäºÅ*/
	char	       sMbName[10];     /*ÓÊÏäÃû	*/
	unsigned short iMbGrpId;	/*ÓÊÏäµÄ×éºÅ*/
	unsigned int   iFstMail;	/*±¾ÓÊÏäµÚÒ»·âÓÊ¼şÎ»ÖÃ*/
	unsigned short iProcNum;        /*±¾ÓÊÏäÁª½Ó½ø³ÌÊı*/
	unsigned int   iMaxMailNum;	/*ÓÊÏä×î´óÓÊ¼şÊı*/
	unsigned int   iMailNum;	/*ÓÊÏäµ±Ç°ÓÊ¼şÊı*/
	unsigned int   iMailSize;	/*ÓÊÏäµ±Ç°ÓÊ¼şÕ¼ÓÃ¿Õ¼ä*/
	unsigned int   iRcvMailNum;	/*ÓÊÏä½ÓÊÕÓÊ¼şÊı*/
	unsigned int   iRcvMailSize;	/*ÓÊÏä½ÓÊÕÊı¾İ´óĞ¡*/
	unsigned int   iSndMailNum;	/*ÓÊÏä·¢ËÍÓÊ¼şÊı*/
	unsigned int   iSndMailSize;	/*ÓÊÏä·¢ËÍÊı¾İ´óĞ¡*/
	char           cMbOpenMode;     /*ÓÊÏä´ò¿ªÄ£Ê½£¨¹²ÏíÓë¶ÀÕ¼)*/
	char	       cCompressFlag;	/*ÓÊ¼ş´æ´¢Êı¾İÊÇ·ñÑ¹Ëõ*/
	char	       cRecryptFlag;	/*ÓÊ¼ş´æ´¢Êı¾İÊÇ·ñ¼ÓÃÜ*/
}Msp_Mb_Tab;
/****************************ÓÊ¼şÃèÊöÇø****************************/
typedef struct {
	unsigned short iSrcMbId;	/*Ô´ÓÊÏäºÅ*/
	unsigned short iSrcGrpId;	/*Ô´ÓÊÏäµÄ×éºÅ*/
	unsigned short iDstMbId;	/*Ä¿±êÓÊÏäºÅ*/
	unsigned short iDstMbGrpId;	/*Ä¿±êÓÊÏäµÄ×éºÅ*/
	unsigned int   iMailNum;        /*ÓÊ¼şºÅ*/
	unsigned short ipriority;	/*ÓÊ¼şPriority*/
	unsigned short iclass;		/*ÓÊ¼şClass*/
	int	       iType;		/*ÓÊ¼şType*/
	char	       cStorageFlag;	/*ÓÊ¼ş´æ·ÅÎ»ÖÃ*/
	unsigned int   iMailSize;	/*ÓÊ¼ş´óĞ¡*/
	unsigned long   iMailBeginPosi;	/*ÓÊ¼ş´æ·ÅÆğÊ¼Î»ÖÃ*/
	unsigned int   iMailPosi;	/*±¾ÓÊ¼şÎ»ÖÃ*/
	unsigned int   iNextMail;	/*ÏÂÒ»ÓÊ¼şÎ»ÖÃ*/
}Msp_Mail_Tab;

/****************************½ø³ÌÃèÊöÇø****************************/
typedef struct  {
	int	       iSeriNo;		/*ĞòºÅ*/
	int	       iUserPid;	/*½ø³ÌºÅ*/
	char           sUserName[15+1];	/*ÓÃ»§Ãû*/
	unsigned short iProcBlockFlag;	/*±¾½ø³Ì×èÈû±êÖ¾*/
	unsigned short iMbIdPosi;      	/*ÓÊÏäÎ»ÖÃ*/
	long	       lBeginTime;	/*½ø³ÌÉÏ´Î²Ù×÷Ê±¼ä*/
	unsigned short iNextPosi;	/*ÏÂÒ»½ø³ÌÎ»ÖÃ*/
}Msp_Proc_Tab;

/****************************½Úµã±í****************************/
typedef struct {
	char	       sNodeName[7+1];	/*½ÚµãÃû*/
	unsigned short iNodeId;	 	/*½ÚµãºÅ*/
	char	       sNodeTcpip[16];	/*½ÚµãIP Address*/
	char	       sNodePort[5+1];	/*½ÚµãTCP port*/
	unsigned short iCurLinkNum;	/*Á´Â·Êı*/
	short          iTimeOut;        /*³¬Ê±Ê±¼ä*/
	short	       iFrmNum;		/*´°¿Ú°üº¬Ö¡Êı*/
	short	       iFrmSize;	/*Ã¿Ö¡±¨ÎÄµÄ´óĞ¡(²»°üº¬Êı¾İ±¨ÎÄÍ·)*/
    	int            iNodePid;        /*½ø³ÌºÅ1*/
}Msp_Node_Tab;

/****************************SafÓëSendµÄ½Ó¿Ú****************************/
typedef struct{
	unsigned short iReqType;	     /*ÇëÇóÀàĞÍ*/
	unsigned int   iMailPosi;	     /*ÓÊ¼şÎ»ÖÃ*/
	unsigned int   iMailNo;		     /*ÓÊ¼şĞòºÅ*/
}Msp_Saf_Tab;

/****************************ÏûÏ¢´«µİ¿ØÖÆÍ·****************************/
typedef struct{
	int		iType;		     /*ÓÊ¼şType*/
	int		iResp_code;          /* ·µ»ØÂë*/
	unsigned short	iReqType;	     /*ÇëÇóÀàĞÍ*/
	unsigned short	iSrcMbId;	     /*Ô´ÓÊÏäºÅ*/
    	unsigned short  iSrcGrpId;           /*Ô­ÓÊÏä×éºÅ*/
    	unsigned short	iSrcBusId;	     /*Ô´ÓÊÏäbusºÅ*/
	unsigned short	iDstMbId;	     /*Ä¿µÄÓÊÏäºÅ*/
	unsigned short	iDstMbGrpId;	     /*Ä¿µÄÓÊÏä×éºÅ*/
	unsigned short	iPriority;	     /*ÓÊ¼şÓÅÏÈ¼¶±ğ*/
	unsigned short	iClass;		     /*ÓÊ¼şClass*/
	unsigned short  iTimeOut;            /*³¬Ê±Ê±¼ä*/
	char		sSrcMbName[20];	     /*Ô´ÓÊÏäÃû*/
	unsigned char	sMacCheck[KEYSIZE];  /*MacĞ£Ñé*/
	char            cFlag;               /*ÊÇ·ñ¶ÀÕ¼¡¢ÊµÊ±¡¢¶Áºó±£Áô*/
}Msp_Msg_Head;
/****************************ÏûÏ¢Ìå****************************/
typedef struct {
	unsigned char FileName[10];
	unsigned int  iMsgLength;
}Msg_Request;
typedef struct {
	unsigned int	iCurDataPosi;	     /*µ±Ç°ÓÊ¼şÎ»ÖÃ*/
	unsigned int    iCurSize;            /*±¾Ö¡Êı¾İµÄ´óĞ¡*/
	char		cEndFlag;	     /*ÓÊ¼ş´«µİ±êÖ¾*/
	unsigned char	sMacCheck[MACSIZE];  /*MacĞ£Ñé*/
}Msg_Body;
typedef struct {
	unsigned int  iRetPosi;
}Msg_Return;

/*************************ÏûÏ¢½á¹¹*********************************/
typedef struct{
	long	mtype;
	unsigned char	mtext[512];
}msg;

/*****************************È«¾Ö±äÁ¿*****************************/
Msp_Pub_Tab  		*pgPubDsc;	            /*¹«¹²ÃèÊöÇø½á¹¹±äÁ¿Ö¸Õë*/
unsigned short		igBusId;	            /*±¾BusºÅ*/
unsigned char 		*pgMspmsg;	            /*´«µİ±¨ÎÄ½á¹¹Ö¸Õë*/
unsigned char 		*pgMbHead;	            /*ÓÊÏäÆğÊ¼µØÖ·*/
unsigned char 		*pgMailHead;	            /*ÓÊ¼şÍ·ÆğÊ¼µØÖ·*/
unsigned char 		*pgProcHead;	            /*½ø³ÌÆğÊ¼µØÖ·*/
unsigned char 		*pgMailBody;	            /*ÓÊÏäÌåÆğÊ¼µØÖ·*/
unsigned char 		*pgNodeHead;	            /*½áµãÆğÊ¼µØÖ·*/
char    		cgShutdownFlag;               /*¹Ø±Õ±êÖ¾*/

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
	int			/*  Ì×½Ó×ÖÃèÊö·û		*/
);

/* add by yl 20060704 */
/*****************************º¯ÊıÉùÃ÷*****************************/
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
        int             sockfd,         /*  Ì×½Ó×ÖÃèÊö·û        */
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
char		*ptrFile		/*  ÎÄ¼şÈ«Â·¾¶Ãû	*/
);

int FileExist
(
char		*ptrFile		/*  ÎÄ¼şÈ«Â·¾¶Ãû	*/
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
sigset_t	*pOldSet	/*  Ô­ÓĞĞÅºÅÆÁ±Î×´Ì¬	*/
);

int UnblockAllSig
(
sigset_t	*pOldSet	/*  Ô­ÓĞĞÅºÅÆÁ±Î×´Ì¬	*/
);

int BlockSig
(
int			Sig			/*  ĞÅºÅ	*/
);

int UnblockSig
(
int			Sig			/*  ĞÅºÅ	*/
);

int ProcSig
(
int			Sig			/*  ĞÅºÅ	*/
);

int SigDirect
(
int			Sig,		/*  ĞÅºÅ			*/
void(* func)(int),		/*  ĞÅºÅ´¦Àíº¯Êı	*/
struct sigaction	*poldact
);

int SigDefault
(
int		Sig				/*  ĞÅºÅ	*/
);

int TestSig
(
int		Sig		/*  ĞÅºÅ	*/
);

int ClearTimer(void);

int SetTimer
(
long	sec,			/*  Ãë			*/
long	usec,			/*  Î¢Ãë		*/
char	cCycFlag		/*  Ñ­»·±êÖ¾	*/
);

int ReadMailDsc
(
char		*ptrFile,		/*  ÎÄ¼şÈ«Â·¾¶Ãû	*/
Msp_Mail_Tab	*pMailDsc		/*  ÓÊ¼şÃèÊöµØÖ·	*/
);

void ClearMsg
(
int		MsgQueID,	/*  ÏûÏ¢¶ÓÁĞ±êÊ¶·û			*/
long	lType		/*  ËùÒªÇå¿ÕÏûÏ¢µÄTypeÖµ	*/
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

