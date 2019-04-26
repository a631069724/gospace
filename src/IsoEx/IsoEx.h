#ifndef __ISO_EX_H__
#define __ISO_EX_H__

#define MAX_ISO_DATA	(1024 * 1)

#define MIN_ISO_LEN  10
#define FIELD_MAX_SIZE   512

typedef struct {
		short len;
			unsigned char Xtype;
} IsoTabEx;

typedef struct {
		short bitf;
			short len;
				short off;
} IsoField;

typedef struct {
		unsigned char dbuf[MAX_ISO_DATA];
			short		off;
#define BCDMSG	0
#define ASCMSG	1
				int		msgtype;
#define BCDBIT	0
#define ASCBIT	1
					int		bittype;
#define BCDVAR	0
#define ASCVAR	1
#define HEXVAR	2
						int		vartype;
							char		msgid[5];
								IsoField	f[128];
									IsoTabEx	*deftab;

} IsoEx;

void EBcd2Asc(char *outstr,char *instr,int lenth);
void Asc2EBcd(char *outstr,char *instr,int lenth);
void AtoE(char *s,int len);
void EtoA(char *s,int len);
void Asc2Bcd(unsigned char *bcd, unsigned char *asc, int len, int r_align);
void Bcd2Asc(unsigned char *asc, unsigned char *bcd, int len, int r_align);
int InitIsoEx( IsoEx *Iso, int MsgType, int BitType,
					int VarType, IsoTabEx *DefIso);
int Str2IsoEx( IsoEx *Iso, unsigned char *dStr, int inLen );
int Iso2StrEx(IsoEx *Iso, unsigned char *dStr, int CanUseLen );
int SetBitEx(IsoEx *Iso, int nBitNo, char *sData, int len);
int GetBitEx(IsoEx *Iso, int nBitNo, char *rData , int MaxLen);
int ExportIso( char *TabName, IsoTabEx *IsoDefEx);

#endif

