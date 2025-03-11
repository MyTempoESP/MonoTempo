#ifndef _CF_API_H_
#define _CF_API_H_

#include <stdint.h>
#include <stdio.h>
#include <wchar.h>
#include <sys/types.h>
#include <sys/time.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <strings.h>
#include <string.h>
#include <stdlib.h>
#include <stdbool.h>
#include <errno.h>       
#include <dirent.h>       
#include <time.h>         
#include <semaphore.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <iostream>
#include <sstream>
#include <vector>
#include <time.h>
#include <termios.h>
#include <signal.h>
#include <pthread.h>
#include <assert.h>

struct PARA
{
	unsigned char RFIDPRO;
	unsigned short STRATFREI;
	unsigned short STRATFRED;
	unsigned short STEPFRE;
	unsigned char CN;
	unsigned char POWER;
	unsigned char ANTENNA;
	unsigned char REGION;
	unsigned char RESERVED;
};

typedef struct
{
	unsigned char firmVersion[32];
	unsigned char hardVersion[32];
	unsigned char SN[12];
	unsigned char PARAS[12];
}DeviceInfo;

typedef struct
{
	unsigned char DevicehardVersion[32];
	unsigned char DevicefirmVersion[32];
	unsigned char DeviceSN[12];
	unsigned char hardVersion[32];
	unsigned char firmVersion[32];
	unsigned char SN[12];
}DeviceFullInfo;

typedef struct
{
	unsigned char DEVICEARRD;
	unsigned char RFIDPRO;
	unsigned char WORKMODE;
	unsigned char INTERFACE;
	unsigned char BAUDRATE;
	unsigned char WGSET;
	unsigned char ANT;
	unsigned char REGION;
	unsigned char STRATFREI[2];
	unsigned char STRATFRED[2];
	unsigned char STEPFRE[2];
	unsigned char CN;
	unsigned char RFIDPOWER;
	unsigned char INVENTORYAREA;
	unsigned char QVALUE;
	unsigned char SESSION;
	unsigned char ACSADDR;
	unsigned char ACSDATALEN;
	unsigned char FILTERTIME;
	unsigned char TRIGGLETIME;
	unsigned char BUZZERTIME;
	unsigned char INTENERLTIME;
}DevicePara;

typedef struct
{
	unsigned char CodeEn;
	unsigned char Code[4];
	unsigned char MaskEn;
	unsigned char StartAdd;
	unsigned char MaskLen;
	unsigned char MaskData[12];
	unsigned char MaskCondition;
}PermissonPara;

typedef struct
{
	unsigned char CodeEn;
	unsigned char Code[4];
	unsigned char MaskEn;
	unsigned char StartAdd;
	unsigned char MaskLen;
	unsigned char MaskData[31];
	unsigned char MaskCondition;
}LongPermissonPara;

typedef struct
{
	unsigned char KCEn;
	unsigned char RelayTime;
	unsigned char KCPowerEn;
	unsigned char TriggleMode;
	unsigned char BufferEn;
	unsigned char ProtocolEn;
	unsigned char ProtocolType;
	unsigned char ProtocolFormat[10];
}GpioPara;

typedef struct
{
	short BasciRssi;
	unsigned char  AntDelta[16];
}RssiPara;

typedef struct
{
	unsigned char wifiEn;
	unsigned char SSID[32];
	unsigned char PASSWORD[64];
	unsigned char IP[4];
	unsigned char PORT[2];
}WiFiPara;

typedef struct
{
	unsigned char IP[4];
	unsigned char MAC[6];
	unsigned char PORT[2];
	unsigned char NetMask[4];
	unsigned char Gateway[4];
}NetInfo;

typedef struct
{
	unsigned char Enable;
	unsigned char IP[4];
	unsigned char PORT[2];
	unsigned char HeartTime;
}RemoteNetInfo;

typedef struct
{
	unsigned char region;
	unsigned short StartFreq;
	unsigned short StopFreq;
	unsigned short StepFreq;
	unsigned char cnt;
}FreqInfo;

typedef struct
{
	unsigned short addr;
	unsigned char  val;
}RFIcRegs;

typedef struct
{
	unsigned char tc;
	unsigned char blf;
	unsigned char miller;
	unsigned char trext;
	unsigned char modu;
}GBRFParam;

typedef struct
{
	unsigned char target;
	unsigned char action;
	unsigned char memBank;
	unsigned short maskPtr;
	unsigned char maskLen;
	unsigned char maskData[255];
}GBSortParam;

typedef struct
{
	unsigned char condition;
	unsigned char session;
	unsigned char target;
}QueryParam;

typedef struct
{
	unsigned short NO;
	short rssi;
	unsigned char antenna;
	unsigned char channel;
	unsigned char crc[2];
	unsigned char pc[2];
	unsigned char codeLen;
	unsigned char code[255];
}TagInfo;

typedef struct
{
	unsigned char tagStatus;
	unsigned char antenna;
	unsigned char crc[2];
	unsigned char pc[2];
	unsigned char codeLen;
	unsigned char code[255];
}TagResp;

typedef struct
{
	float tari;
	float rtcal;
	float trcal;
	unsigned char dr;
	unsigned char miller;
	unsigned char trext;
	unsigned char modu;
}ISORFParam;

typedef struct
{
	unsigned char resv;
	unsigned char trucate;
	unsigned char target;
	unsigned char action;
	unsigned char membank;
	unsigned short ptr;
	unsigned char len;
	unsigned char mask[240];
}ISOSelectParam;

typedef struct
{
	unsigned char sel;
	unsigned char session;
	unsigned char target;
}ISOQueryParam;

typedef struct
{
	unsigned char readlock;
	unsigned char membank;
	unsigned short blockPtr;
	unsigned char blockRange;
	unsigned char mask[247];
}ISOPermalockParam;

typedef struct
{
	unsigned long blf;
	unsigned char miller;
	unsigned char trext;
	unsigned short rxDelay;
	unsigned short rxLen;
	unsigned long rxNum;
	unsigned char autoMode;
}CP_Sensi_Prm_Typ;

typedef struct
{
	unsigned long frame_err;
	unsigned long frame_total;
	unsigned long blf;
	unsigned char miller;
	unsigned char trext;
	unsigned char freqOffset;
}CP_Sensi_Result_Typ;

typedef struct
{
	unsigned char i_start;
	unsigned char i_stop;
	unsigned char q_start;
	unsigned char q_stop;
}IQ_Axial_Typ;

typedef struct
{
	unsigned char i_origin;
	unsigned char q_origin;
	unsigned char size;
	unsigned char step;
}JSC_AUTO_SCAN_PRM_Typ;

typedef struct
{
	unsigned char number;
	unsigned char dataLen;
	unsigned char data[253];
}JSC_Data_Typ;

typedef struct
{
	unsigned char option;
	unsigned short addr;
	unsigned short val;
	unsigned short interval;
}Read_Write_Reg_Cmd_Item_Typ;

typedef struct
{
	unsigned char items;
	unsigned short regs;
	unsigned short val[126];
}Read_Regs_Result_Typ;

typedef struct
{
	unsigned char status;
	unsigned long time;
}Int_Status_Item_Typ;

typedef struct
{
	unsigned short addr;
	unsigned char value;
	unsigned long time;
}CR_Log_Item_Typ;

typedef struct
{
	unsigned char target;
	unsigned char trucate;
	unsigned char action;
	unsigned char membank;
	unsigned short m_ptr;
	unsigned char len;
	unsigned char mask[31];
}SelectSortParam;

//========================Dynamic library external error code=================================
#define STAT_OK								0x00000000
#define STAT_PORT_HANDLE_ERR				0xFFFFFF01    // Handle error or input serial port parameter error
#define STAT_PORT_OPEN_FAILED				0xFFFFFF02	  // Failed to open the serial port
#define STAT_DLL_INNER_FAILED				0xFFFFFF03    // Dynamic library internal error
#define STAT_CMD_PARAM_ERR					0xFFFFFF04    // The parameter value is incorrect or out of range, or the module does not support the parameter value
#define STAT_CMD_SERIAL_NUM_EXIT			0xFFFFFF05    // Serial number already exists
#define STAT_CMD_INNER_ERR					0xFFFFFF06    // Command execution failed due to an internal module error
#define STAT_CMD_INVENTORY_STOP				0xFFFFFF07    // The inventory has not reached the label or has ended
#define STAT_CMD_TAG_NO_RESP				0xFFFFFF08    // Label response timeout
#define STAT_CMD_DECODE_TAG_DATA_FAIL		0xFFFFFF09    // Demodulation tag data error
#define STAT_CMD_CODE_OVERFLOW				0xFFFFFF0A    // Label data exceeds the maximum transmission length of the serial port
#define STAT_CMD_AUTH_FAIL					0xFFFFFF0B    // Authentication failed
#define STAT_CMD_PWD_ERR					0xFFFFFF0C    // password error
#define STAT_CMD_SAM_NO_RESP				0xFFFFFF0D    // SAM card not responding
#define STAT_CMD_SAM_CMD_FAIL				0xFFFFFF0E    // PSAM card command execution failed
#define STAT_CMD_RESP_FORMAT_ERR			0xFFFFFF0F    // Reader/writer response format error
#define STAT_CMD_HAS_MORE_DATA				0xFFFFFF10    // The command was executed successfully, but there are still subsequent data that have not been transferred
#define STAT_CMD_BUF_OVERFLOW				0xFFFFFF11    // Incoming cache too small, data overflow
#define STAT_CMD_COMM_TIMEOUT				0xFFFFFF12    // Timeout waiting for reader response
#define STAT_CMD_COMM_WR_FAILED				0xFFFFFF13    // Failed to write data to the serial port
#define STAT_CMD_COMM_RD_FAILED				0xFFFFFF14    // Failed to read serial port data
#define STAT_CMD_NOMORE_DATA				0xFFFFFF15    // No more data
#define STAT_DLL_UNCONNECT       			0xFFFFFF16    // The network connection has not been established
#define STAT_DLL_DISCONNECT       			0xFFFFFF17    // The network has been disconnected
#define STAT_CMD_RESP_CRC_ERR				0xFFFFFF18    // Reader/writer response CRC check error
#define STAT_CMD_IAP_CRC_ERR				0xFFFFFF21    // Download program CRC verification error
#define STAT_CMD_DOWMLOAD_ERR				0xFFFFFF22    // Download data error
#define STAT_CMD_DOWM_NONE_ERR				0xFFFFFF23    // User zone download incomplete

//========================Label Status Code=================================
#define STAT_GB_TAG_LOW_POWER				0xFFFFFF40    // Insufficient power supply to the tag, the tag does not have enough energy to complete the operation
#define STAT_GB_TAG_OPR_LIMIT				0xFFFFFF41    // Insufficient label operation permissions, unauthorized access
#define STAT_GB_TAG_MEM_OVF					0xFFFFFF42    // The label operation storage area overflowed or the target storage area does not exist
#define STAT_GB_TAG_MEM_LCK					0xFFFFFF43    // Label storage area locking, performing write or erase operations on storage areas that are locked as non writable, and performing read operations on storage areas that are locked as non readable
#define STAT_GB_TAG_PWD_ERR					0xFFFFFF44    // Label operation password error, access command password error
#define STAT_GB_TAG_AUTH_FAIL				0xFFFFFF45    // Label authentication failed, failed authentication
#define STAT_GB_TAG_UNKNW_ERR				0xFFFFFF46    // Unknown error in label operation, an undetermined error occurred
#define STAT_ISO_TAG_OTHER_ERR				0xFFFFFF50    // Other errors
#define STAT_ISO_TAG_NOT_SUPPORT			0xFFFFFF51    // The tag does not support this parameter
#define STAT_ISO_TAG_OPR_LIMIT				0xFFFFFF52    // Insufficient permissions
#define STAT_ISO_TAG_MEM_OVF				0xFFFFFF53    // Storage Area Overflow
#define STAT_ISO_TAG_MEM_LCK				0xFFFFFF54    // Store lock
#define STAT_ISO_TAG_CRYPTO_ERR				0xFFFFFF55    // Label encryption suite error
#define STAT_ISO_TAG_NOT_ENCAP				0xFFFFFF56    // Empty port command not encapsulated
#define STAT_ISO_TAG_RESP_OVF				0xFFFFFF57    // Label response buffer overflow
#define STAT_ISO_TAG_SEC_TIMEOUT			0xFFFFFF58    // The tag is in a security timeout state
#define STAT_ISO_TAG_LOW_POWER				0xFFFFFF59    // Insufficient power supply
#define STAT_ISO_TAG_UNKNW_ERR				0xFFFFFF5A    // Label response unknown error
#define STAT_ISO_TAG_SENSOR_CFG				0xFFFFFF5B    // Sensor timed task configuration exceeds upper limit
#define STAT_ISO_TAG_TAG_BUSY				0xFFFFFF5C    // Label busy
#define STAT_ISO_TAG_MEASU_NOT_SUPPORT		0xFFFFFF5D    // The sensor does not support this measurement type

#ifdef __cplusplus
extern "C" {
#endif

	/// <summary>
	/// Connecting the serial port
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="ComPort"></param>
	/// <param name="Baudrate"></param>
	/// <returns></returns>
	int OpenDevice(int64_t* hComm, char* pcCom, int iBaudRate);
	/// <summary>
	/// Connection network port
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="strIP">ip</param>
	/// <param name="wPort">prot</param>
	/// <param name="timeoutMs"></param>
	/// <returns></returns>
	int OpenNetConnection(int64_t* hComm, char* strIP, unsigned short wPort, long timeoutMs);
	/// <summary>
	/// Close the reader connection
	/// </summary>
	/// <param name="hComm"></param>
	/// <returns>0x00 </returns>
	int CloseDevice(int64_t hComm);
	/// <summary>
	/// Get Usb Count
	/// </summary>
	/// <param name=""></param>
	/// <returns></returns>
	int CFHid_GetUsbCount();
	/// <summary>
	/// Get Usb Info
	/// </summary>
	/// <param name="index"></param>
	/// <param name="pucDeviceInfo"></param>
	/// <returns></returns>
	int CFHid_GetUsbInfo(unsigned short index, char* pucDeviceInfo);
	/// <summary>
	///  Hid Usb Connection
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="index"></param>
	/// <returns></returns>
	int OpenHidConnection(int64_t* hComm, unsigned short index);
	/// <summary>
	/// Get Reader/Writer Module Information Command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="devInfo"></param>
	/// <returns></returns>
	int GetInfo(int64_t hComm, DeviceInfo* devInfo);
	/// <summary>
	/// Command to obtain the information of reader-writer all-in-one machine
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="devInfo"></param>
	/// <returns>0x00 </returns>
	int GetDeviceInfo(int64_t hComm, DeviceFullInfo* devInfo);
	/// <summary>
	/// Get Device Parameters Command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="devInfo"></param>
	/// <returns>0x00 </returns>
	int GetDevicePara(int64_t hComm, DevicePara* devInfo);
	/// <summary>
	/// Set Device Parameters Command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="devInfo"></param>
	/// <returns>0x00 </returns>
	int SetDevicePara(int64_t hComm, DevicePara devInfo);
	/// <summary>
	/// Obtain card reading permission parameters
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="PermissonPara"></param>
	/// <returns>0x00 </returns>
	int GetLongPermissonPara(int64_t hComm, LongPermissonPara* PermissonPara);
	/// <summary>
	/// Set card reading permission parameters
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="PermissonPara"></param>
	/// <returns>0x00 </returns>
	int SetLongPermissonPara(int64_t hComm, LongPermissonPara PermissonPara);
	/// <summary>
	/// Obtain card reading permission parameters
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="PermissonPara"></param>
	/// <returns>0x00 </returns>
	int GetPermissonPara(int64_t hComm, PermissonPara* PermissonPara);
	/// <summary>
	/// Set card reading permission parameters
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="PermissonPara">Permisson</param>
	/// <returns>0x00 </returns>
	int SetPermissonPara(int64_t hComm, PermissonPara PermissonPara);
	/// <summary>
	/// Get input and output parameters
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="GpioPara"></param>
	/// <returns>0x00 </returns>
	int GetGpioPara(int64_t hComm, GpioPara* GpioPara);
	/// <summary>
	/// Set input and output parameters
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="GpioPara">Gpio</param>
	/// <returns>0x00 </returns>
	int SetGpioPara(int64_t hComm, GpioPara GpioPara);
	/// <summary>
	/// Get remote network parameters
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="type"></param>
	/// <returns>0x00 </returns>
	int GetNetInfo(int64_t hComm, NetInfo* type);
	/// <summary>
	/// Set remote network parameters
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="type">Net</param>
	/// <returns>0x00 </returns>
	int SetNetInfo(int64_t hComm, NetInfo type);
	/// <summary>
	/// Get WIFI parameters
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="WiFiPara"></param>
	/// <returns></returns>
	int GetwifiPara(int64_t hComm, WiFiPara* WiFiPara);
	/// <summary>
	/// Set WIFI parameters
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="PermissonPara"></param>
	/// <returns></returns>
	int SetwifiPara(int64_t hComm, WiFiPara WiFiPara);
	/// <summary>
	/// Acquire power command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="power"></param>
	/// <param name="reserved"></param>
	/// <returns>0x00 </returns>
	int GetRFPower(int64_t hComm, unsigned char* power, unsigned char* reserved);
	/// <summary>
	/// Set power command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="power"></param>
	/// <param name="reserved"></param>
	/// <returns></returns>
	int SetRFPower(int64_t hComm, unsigned char power, unsigned char reserved);
	/// <summary>
	/// Obtain operating frequency command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="freqInfo"></param>
	/// <returns></returns>
	int GetFreq(int64_t hComm, FreqInfo* freqInfo);
	/// <summary>
	/// Set operating frequency command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="freqInfo"></param>
	/// <returns></returns>
	int SetFreq(int64_t hComm, const FreqInfo* freqInfo);
	/// <summary>
	/// Obtain antenna configuration commands
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="antenna"></param>
	/// <returns>0x00 </returns>
	int GetAntenna(int64_t hComm, unsigned char* antenna);
	/// <summary>
	/// Set Antenna Configuration Command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="antenna">Antenna</param>
	/// <returns>0x00 </returns>
	int SetAntenna(int64_t hComm, unsigned char* antenna);
	/// <summary>
	/// Get Protocol Type Command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="type"></param>
	/// <returns></returns>
	int GetRFIDType(int64_t hComm, unsigned char* type);
	/// <summary>
	/// Set Protocol Type Command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="type"></param>
	/// <returns></returns>
	int SetRFIDType(int64_t hComm, unsigned char type);
	/// <summary>
	/// Get remote network parameters
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="type"></param>
	/// <returns></returns>
	int GetRemoteNetInfo(int64_t hComm, RemoteNetInfo* type);
	/// <summary>
	/// Set remote network parameters
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="type"></param>
	/// <returns></returns>
	int SetRemoteNetInfo(int64_t hComm, RemoteNetInfo type);
	/// <summary>
	/// Get Temperature Threshold Command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="pTempCur"></param>
	/// <param name="pTempLimit"></param>
	/// <returns></returns>
	int GetTemperature(int64_t hComm, unsigned char* pTempCur, unsigned char* pTempLimit);
	/// <summary>
	/// Set Temperature Threshold Command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="tempLimit"></param>
	/// <param name="resv"></param>
	/// <returns></returns>
	int SetTemperature(int64_t hComm, unsigned char tempLimit, unsigned char resv);
	/// <summary>
	/// Device Restore Factory Settings Command
	/// </summary>
	/// <param name="hComm"></param>
	/// <returns>0x00 \B1\EDʾ\B3ɹ\A6</returns>
	int RebootDevice(int64_t hComm);
	/// <summary>
	/// Set relay release
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="time"></param>
	/// <returns></returns>
	int Release_Relay(int64_t hComm, unsigned char time);
	/// <summary>
	/// Set relay closed
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="time"></param>
	/// <returns></returns>
	int Close_Relay(int64_t hComm, unsigned char time);
	/// <summary>
	/// Inquiry command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="btInvCount"></param>
	/// <param name="dwInvParam"></param>
	/// <returns>0x00 </returns>
	int InventoryContinue(int64_t hComm, unsigned char btInvCount, unsigned long dwInvParam);
	/// <summary>
	/// Obtain label information and return it in TagInfo format
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="tag_info"></param>
	/// <param name="timeout"></param>
	/// <returns>0x00 </returns>
	int GetTagUii(int64_t hComm, TagInfo* tag_info, unsigned short timeout);
	/// <summary>
	/// Stop counting
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="timeout"></param>
	/// <returns>0x00 </returns>
	int InventoryStop(int64_t hComm, unsigned short timeout);
	/// <summary>
	/// Read Data Command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="option"></param>
	/// <param name="accPwd"></param>
	/// <param name="memBank"></param>
	/// <param name="wordPtr"></param>
	/// <param name="wordCount"></param>
	/// <returns></returns>
	int ReadTag(int64_t hComm, unsigned char option, unsigned char* accPwd, unsigned char memBank, unsigned short wordPtr, unsigned char wordCount);
	/// <summary>
	/// Get read command response command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="resp"></param>
	/// <param name="wordCount"></param>
	/// <param name="readData"></param>
	/// <param name="timeout"></param>
	/// <returns></returns>
	int GetReadTagResp(int64_t hComm, TagResp* resp, unsigned char* wordCount, unsigned char* readData, unsigned short timeout);
	/// <summary>
	/// Get tag response command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="cmd"></param>
	/// <param name="resp"></param>
	/// <param name="timeout"></param>
	/// <returns></returns>
	int GetTagResp(int64_t hComm, unsigned short cmd, TagResp* resp, unsigned short timeout);
	/// <summary>
	/// Block Write Command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="option"></param>
	/// <param name="accPwd"></param>
	/// <param name="memBank"></param>
	/// <param name="wordPtr"></param>
	/// <param name="wordCount"></param>
	/// <param name="writeData"></param>
	/// <returns></returns>
	int WriteTag(int64_t hComm, unsigned char option, unsigned char* accPwd, unsigned char memBank, unsigned short wordPtr, unsigned char wordCount, unsigned char* writeData);
	/// <summary>
	/// Lock label command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="accPwd"></param>
	/// <param name="erea"></param>
	/// <param name="action"></param>
	/// <returns></returns>
	int LockTag(int64_t hComm, unsigned char* accPwd, unsigned char erea, unsigned char action);
	/// <summary>
	/// Inactivation tag command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="accPwd"></param>
	/// <returns></returns>
	int KillTag(int64_t hComm, unsigned char* accPwd);
	/// <summary>
	/// Select Label Command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="maskPtr"></param>
	/// <param name="maskBits"></param>
	/// <param name="mask"></param>
	/// <returns></returns>
	int SetSelectMask(int64_t hComm, unsigned short maskPtr, unsigned char maskBits, unsigned char* mask);
	/// <summary>
	/// Get Q value command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="pqVal"></param>
	/// <param name="reserved"></param>
	/// <returns></returns>
	int GetCoilPRM(int64_t hComm, unsigned char* pqVal, unsigned char* reserved);
	/// <summary>
	/// Set Q value command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="qVal"></param>
	/// <param name="reserved"></param>
	/// <returns></returns>
	int SetCoilPRM(int64_t hComm, unsigned char qVal, unsigned char reserved);
	/// <summary>
	/// Get Count Tag Selset Command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="proto"></param>
	/// <param name="param"></param>
	/// <returns></returns>
	int SelectOrSortGet(int64_t hComm, unsigned char proto, SelectSortParam* param);
	/// <summary>
	/// Set inventory label Selset command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="proto"></param>
	/// <param name="param"></param>
	/// <returns></returns>
	int SelectOrSortSet(int64_t hComm, unsigned char proto, SelectSortParam* param);
	/// <summary>
	/// Query command to obtain inventory tag
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="proto"></param>
	/// <param name="param"></param>
	/// <returns></returns>
	int QueryCfgGet(int64_t hComm, unsigned char proto, QueryParam* param);
	/// <summary>
	/// Set Inventory Tag Query Command
	/// </summary>
	/// <param name="hComm"></param>
	/// <param name="proto"></param>
	/// <param name="param"></param>
	/// <returns></returns>
	int QueryCfgSet(int64_t hComm, unsigned char proto, QueryParam* param);

#ifdef __cplusplus
}
#endif

#endif

