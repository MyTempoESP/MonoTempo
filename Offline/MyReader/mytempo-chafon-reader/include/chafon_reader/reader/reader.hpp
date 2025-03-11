#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include "../config/config.h"
#include "../tag/tag_monitor.hpp"
#include <thread>

/*
 * ~@ NAME Rodrigo Monteiro Junior
 * ~@ LINK(TsukiGva2) https://github.com/TsukiGva2
 * ~@ DATE 19/06/2024
 *
 * ~@ TITLE Device Parameters:
 *
 * Get: GetDevicePara(DEVICE_HANDLE, DevicePara* DEST)
 * Set: SetDevicePara(DEVICE_HANDLE, DevicePara SRC) // yes, by value
 *
 * DevicePara:
                unsigned char DEVICEARRD; // DEVICEADDR (correction)

                unsigned char RFIDPRO;
                unsigned char WORKMODE;
                unsigned char INTERFACE;
                unsigned char BAUDRATE;
                unsigned char WGSET;
                unsigned char ANT;
                unsigned char REGION;

        // Frequency {{
                unsigned char STRATFREI[2];
                unsigned char STRATFRED[2];
                unsigned char STEPFRE[2];
                unsigned char CN;
        // }}

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

 * SetDevicePara workflow example:
 *			...
        GetDevicePara(DEVICE_HANDLE, &param);

        param.PARAMETER = VALUE; // repeat as needed

        assert(SetDevicePara(DEVICE_HANDLE, param), 0);
 *			...
 *
 * THE FREQUENCY PARAMETER:
 *
 * 	FreqInfo freq:
 * 		uint8		region;
 * 		uint16		StartFreq;
 * 		uint16		StopFreq;
 * 		uint16		StepFreq;
 * 		uint8		cnt;
 *
 *	StartFreq can be modified directly or using the
 *	param.STARTFREI property:
 *
 *	(correction: translation issues turned STARTFRE<I|D>
 *	             into *STRAT*FRE<I|D>) // FIXME
 *
 *	the param.STARTFRE<I|D> properties's type is:
 *		Vec<uint8>
 *	they are supposed to be mapped from left to right
 *	as follows:
 *		uint16 <- ([0] SHIFT_BY 8)
 *		 	   BINARY_OR
 *		          ([1])
 *
 * 	~@ ASSIGNMENT :: Vec<uint8> -> uint16
 * 	freq.StartFreq =
 *	  (param.
 *	    STARTFREI
 *	             [0]) << 8
 *	        |
 *	  (param.
 *	    STARTFREI
 *	             [1]);
 *
 * 	StopFreq can also be modified in the same way,
 * 	instead using the param.STARTFRED property:
 *
 * 	~@ ASSIGNMENT :: Vec<uint8> -> uint16
 * 	freq.StopFreq =
 *	  (param.
 *	    STARTFRED
 *	             [0]) << 8
 *	        |
 *	  (param.
 *	    STARTFRED
 *	             [1]);
 *
 *	Finally, we have the StepFreq, which also follows that logic.
 *	Only with a different param property:
 *
 * 	~@ ASSIGNMENT :: Vec<uint8> -> uint16
 * 	freq.StepFreq =
 *	  (param.
 *	    STEPFRE
 *	           [0]) << 8
 *	        |
 *	  (param.
 *	    STEPFRE
 *	           [1]);
 *
 *	The last FreqInfo field is the count (cnt), it is region specific
 *	like the last parameters, but it's easily set with a simple assignment
 *
 *	freq.cnt = param.CN
 *
 *	CONVERTING FREQUENCIES TO SOMETHING SOMEWHAT HUMAN
 *
 *	we can convert the given frequencies to readable values
 *	with the following expressions
 *
 *	start = (float)freq.StartFreq + (float)freq.StopFreq / 1000
 *	end   = start + ((float)freq.StepFreq / 1000 * (float)freq.cnt)
 *
 * 	let's break this down:
 *
 * 	let StartFreq = freq.StartFreq,
 * 	    StopFreq  = ... .StopFreq,
 * 	    StepFreq  = ... .StepFreq
 * 	and
 * 	    Steps     = freq.cnt
 *
 * 	                    StopFreq
 * 	start = StartFreq + --------
 * 	                      1000
 *
 *	end = start + steps
 *	where
 *	             StepFreq
 *	     steps = -------- x Steps
 *	               1000
 *
 * 	~@ DEF_REFERENCE YAPPING
 *	FREQUENCIES, Final commentary:
 *
 *		The frequency stuff is really confusing,
 *		tho the only thing i didn't really get is
 *		the StopFreq field, which is assigned to at
 *		the reference function with:
 *
 *		  freq.StopFreq =
 *		    (param->STARTFRED[0] << 8)
 *		           |
 *		    (param->STARTFRED[1])
 *
 *		The name *START*FRED imply that this is
 *		somehow related to StartFreq.
 *		With the 'D' at the end probably meaning
 *		DECIMAL, which leads to the conclusion
 *		that StartFreq is a decimal number broken
 *		down into two UINT16s
 *
 *		The conversion formula provided furthers
 *		this theory even more so, considering:
 *
 * 	        	            StopFreq
 * 		start = StartFreq + --------
 * 	        	              1000
 *
 * 		it can't get clearer than this, so
 * 		consider STARTFREI as the integer part
 * 		and STARTFRED as the decimal part, to
 * 		form the number:
 *
 * 			~@ ASSIGNMENT :: uint16 -> uint16 -> float
 * 			frequency = STARTFREI.STARTFRED
 *
 * 		if you can gather further understanding on the topic,
 * 		please send an email to tsukigva@gmail.com
 *
 *
 * OUTRO:
 *  FREQUENCY_TABLE:
 *    ----------------
 *    | USA
 *    |   freq  : 902.750 MHz
 *    |   step  : 500
 *    |   count : 50
 *    ----------------
 *    | Korea
 *    |   freq  : 917.100 MHz
 *    |   step  : 200
 *    |   count : 15
 *    ----------------
 *    | Europe
 *    |   freq  : 865.100 MHz
 *    |   step  : 200
 *    |   count : 15
 *    ----------------
 *    | Japan
 *    |   freq  : 952.200 MHz
 *    |   step  : 200
 *    |   count : 8
 *    ----------------
 *    | Malaysia
 *    |   freq  : 919.500 MHz
 *    |   step  : 500
 *    |   count : 7
 *    ----------------
 *    | Europe3
 *    |   freq  : 865.700 MHz
 *    |   step  : 600
 *    |   count : 4
 *    ----------------
 *    | China_1
 *    |   freq  : 840.125 MHz
 *    |   step  : 250
 *    |   count : 20
 *    ----------------
 *    | China_2
 *    |   freq  : 920.125 MHz
 *    |   step  : 250
 *    |   count : 20
 *    ----------------
 */

/*
static void getFrequencyParam(DevicePara *param, float *start, float *end) {
 FreqInfo freq;
 freq.region = param->REGION;

 unsigned short start_freq;
 unsigned short stop_freq;

 start_freq = (param->STRATFREI[0] << 8) | param->STRATFREI[1];
 stop_freq = (param->STRATFRED[0] << 8) |
             param->STRATFRED[1]; // ~@ FIXME REFERENCE YAPPING

 freq.StartFreq = start_freq;
 freq.StopFreq = stop_freq;

 freq.StepFreq = (param->STEPFRE[0] << 8) | param->STEPFRE[1];
 freq.cnt = param->CN;

 *start = (float)start_freq + (float)stop_freq / 1000;

 float steps = (float)freq.StepFreq / 1000 * freq.cnt;
 *end = *start + steps;
}
*/

class Reader {
  int64_t handle;

public:
#ifdef __PYTHON_READER_MODULE__
  void read_n(int n);
#else
  void createMonitor(void);
  void joinMonitor(void);
#endif

  void SetFrequency();

  STATUS close(void);
  STATUS openAndFetchParams(DevicePara *param);
  STATUS startInventory(void);
  STATUS stopInventory(void);

  Reader();
  ~Reader();

private:
  TagMonitor tm;

};

#endif
