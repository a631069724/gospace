#ifndef	__T_MSP_H__
#define	__T_MSP_H__

#define	MSPGETTIMEOUT	1	

#define	SKPMID	11		/*－SKPM－核心处理模块*/
#define	SFPMID	12		/*－SFPM－失败处理模块*/
#define	SSSMID	13		/*－SSSM－加密处理模块*/
#define	SAFMID	14		/*－SAFM－存储转发模块*/
#define SRSMID	15		/*－SRSM－轮询服务模块*/
#define SMCIID	16		/*－SMCI－管理端控制接口*/
#define SMMIID	17		/*－SMMI－管理端监控接口*/
#define SMMMID	18		/*－SMMM－监控服务模块*/
#define TESTID	20		/*－TEST－测试模块*/

int	tMsp_Attach(unsigned short uId);
int	tMsp_Detach();
int	tMsp_Put(char *sMsg, int iLen, unsigned short uDstId);
int	tMsp_Get(char *sMsg, int *iLen, unsigned short *uSrcId, int iTimeOut);

#endif
