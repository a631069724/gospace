#ifndef	__T_SHM_H__
#define	__T_SHM_H__

int	tShm_IsReload();
int	tShm_InitReloadNum();
int	tShm_GetReloadNum();
int	tShm_InitShmReloadNum();
int	tShm_GetShmReloadNum();
int	tShm_GetReloadFlag();
int	tShm_SetReloadFlag(int,char);
int	tShm_DebugAll();
int	tShm_SetDebugAll(char);
int	tShm_SetDebugAllInvalid();
int	tShm_DebugOne(char*);
int	tShm_SetDebugOne(char*,char);
int	tShm_DebugBeLog(char*);
int	tShm_UpdatePid(char*, int ,char);
int	tShm_SelectPid();
int	tShm_PidPause();
int	tShm_SetHostLs(int);
int	tShm_GetHostLs();
int	IsNull(char*);
int	PidPause(char*);

#endif
