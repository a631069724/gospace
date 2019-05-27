#ifndef	__T_LOG_H__
#define	__T_LOG_H__

#define	SYS	0
#define	MOD	1
#define	ERR	2
#define	INF	3
#define	DEB	4
#define	ASC	5
#define	HEX	6
#define	PUB	7
#define	MAIL	8

#define mERR	__FILE__,__LINE__,ERR
#define mINF	__FILE__,__LINE__,INF
#define mDEB	__FILE__,__LINE__,DEB

void	tLog_Init();
int	tLog_DebugOk();
void	tLog_PrintAsc(char *pFile, int iLine, int iLineType, char iLogPos, char* pFormat, ...);
void	tLog_PrintHex();
void	tLog_PrintIso();
int     tLog_Monitor();
int	tLog_File();
void	DoAsc();
void	DoHex();
void	DoIso();
int	LogIt();

#endif
