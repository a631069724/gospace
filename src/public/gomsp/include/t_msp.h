#ifndef	__T_MSP_H__
#define	__T_MSP_H__

#define	MSPGETTIMEOUT	1	

#define	SKPMID	11		/*��SKPM�����Ĵ���ģ��*/
#define	SFPMID	12		/*��SFPM��ʧ�ܴ���ģ��*/
#define	SSSMID	13		/*��SSSM�����ܴ���ģ��*/
#define	SAFMID	14		/*��SAFM���洢ת��ģ��*/
#define SRSMID	15		/*��SRSM����ѯ����ģ��*/
#define SMCIID	16		/*��SMCI������˿��ƽӿ�*/
#define SMMIID	17		/*��SMMI������˼�ؽӿ�*/
#define SMMMID	18		/*��SMMM����ط���ģ��*/
#define TESTID	20		/*��TEST������ģ��*/

int	tMsp_Attach(unsigned short uId);
int	tMsp_Detach();
int	tMsp_Put(char *sMsg, int iLen, unsigned short uDstId);
int	tMsp_Get(char *sMsg, int *iLen, unsigned short *uSrcId, int iTimeOut);

#endif
